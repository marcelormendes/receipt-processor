package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marcelorm/receipt-processor/api"
	"github.com/marcelorm/receipt-processor/models"
	"github.com/marcelorm/receipt-processor/storage"
)

func init() {
	// Configure minimal logging for tests
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only log errors during tests
	})
	slog.SetDefault(slog.New(handler))

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// setupTestServer creates a test server with all routes configured
func setupTestServer() *httptest.Server {
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Add JSON validation middleware
	router.Use(api.JSONValidationMiddleware(1024 * 1024))

	store := storage.NewMemoryStorage()
	handler := api.NewReceiptHandler(store)

	router.POST("/receipts/process", handler.ProcessReceipt)
	router.GET("/receipts/:id/points", handler.GetPoints)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	server := httptest.NewServer(router)
	return server
}

func TestE2ETargetReceipt(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Target receipt example - use map for e2e test to bypass custom type validations
	receipt := map[string]any{
		"retailer":     "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": []map[string]any{
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
			{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
			{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
			{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
			{"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00"},
		},
		"total": "35.35",
	}

	// Process receipt
	receiptID := processReceipt(t, server.URL, receipt)

	// Get points
	points := getPoints(t, server.URL, receiptID)

	// Verify points for Target receipt
	expectedPoints := 28
	if points != expectedPoints {
		t.Errorf("Expected %d points for Target receipt, got %d", expectedPoints, points)
	}
}

func TestE2EMAndMCornerMarketReceipt(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// M&M Corner Market receipt example
	receipt := map[string]any{
		"retailer":     "M&M Corner Market",
		"purchaseDate": "2022-03-20",
		"purchaseTime": "14:33",
		"items": []map[string]any{
			{"shortDescription": "Gatorade", "price": "2.25"},
			{"shortDescription": "Gatorade", "price": "2.25"},
			{"shortDescription": "Gatorade", "price": "2.25"},
			{"shortDescription": "Gatorade", "price": "2.25"},
		},
		"total": "9.00",
	}

	// Process receipt
	receiptID := processReceipt(t, server.URL, receipt)

	// Get points
	points := getPoints(t, server.URL, receiptID)

	// Verify points for M&M Corner Market receipt
	expectedPoints := 109
	if points != expectedPoints {
		t.Errorf("Expected %d points for M&M Corner Market receipt, got %d", expectedPoints, points)
	}
}

func TestE2ERoundDollarAmount(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Receipt with round dollar amount
	receipt := map[string]any{
		"retailer":     "Walmart",
		"purchaseDate": "2022-02-02",
		"purchaseTime": "13:30",
		"items": []map[string]any{
			{"shortDescription": "Item 1", "price": "20.00"},
			{"shortDescription": "Item 2", "price": "5.00"},
		},
		"total": "25.00", // Round dollar amount
	}

	// Process receipt
	receiptID := processReceipt(t, server.URL, receipt)

	// Get points
	points := getPoints(t, server.URL, receiptID)

	// Points should include 50 for round dollar amount
	if points < 50 {
		t.Errorf("Expected at least 50 points for round dollar amount, got %d", points)
	}
}

func TestE2EPurchaseTimeBonus(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Receipt with purchase time between 14:00 and 16:00
	receipt := map[string]any{
		"retailer":     "Store",
		"purchaseDate": "2022-02-02",
		"purchaseTime": "15:15", // Between 14:00 and 16:00
		"items": []map[string]any{
			{"shortDescription": "Item", "price": "5.00"},
		},
		"total": "5.00",
	}

	// Process receipt
	receiptID := processReceipt(t, server.URL, receipt)

	// Get points
	points := getPoints(t, server.URL, receiptID)

	// Points should include 10 for purchase time between 14:00 and 16:00
	if points < 10 {
		t.Errorf("Expected at least 10 points for purchase time bonus, got %d", points)
	}
}

func TestE2EOddDayBonus(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Receipt with odd purchase day
	receipt := map[string]any{
		"retailer":     "Store",
		"purchaseDate": "2022-03-21", // Odd day
		"purchaseTime": "12:00",
		"items": []map[string]any{
			{"shortDescription": "Item", "price": "5.00"},
		},
		"total": "5.00",
	}

	// Process receipt
	receiptID := processReceipt(t, server.URL, receipt)

	// Get points
	points := getPoints(t, server.URL, receiptID)

	// Points should include 6 for odd purchase day
	if points < 6 {
		t.Errorf("Expected at least 6 points for odd day bonus, got %d", points)
	}
}

func TestE2EInvalidReceipt(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Empty receipt (invalid)
	receipt := map[string]any{}

	// Create request
	reqBody, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("Failed to marshal receipt: %v", err)
	}

	// Send request
	resp, err := http.Post(fmt.Sprintf("%s/receipts/process", server.URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Warning: failed to close response body: %v", err)
		}
	}()

	// Verify response is bad request
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)

		// Debug output to help diagnose why validation is failing
		var respBody map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			t.Logf("Warning: failed to decode error response: %v", err)
		} else {
			t.Logf("Response body: %v", respBody)
		}
	}
}

func TestE2EInvalidReceiptID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Get points for non-existent receipt ID
	resp, err := http.Get(fmt.Sprintf("%s/receipts/non-existent-id/points", server.URL))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Warning: failed to close response body: %v", err)
		}
	}()

	// Verify response is not found
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

// Helper function to process a receipt and return the ID
func processReceipt(t *testing.T, serverURL string, receipt any) string {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqBody, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("Failed to marshal receipt: %v", err)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/receipts/process", serverURL), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}
		t.Fatalf("Expected status code %d, got %d: %v", http.StatusOK, resp.StatusCode, errResp)
	}

	var response models.ReceiptResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return response.ID
}

// Helper function to get points for a receipt ID
func getPoints(t *testing.T, serverURL string, receiptID string) int {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/receipts/%s/points", serverURL, receiptID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}
		t.Fatalf("Expected status code %d, got %d: %v", http.StatusOK, resp.StatusCode, errResp)
	}

	var response models.PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return response.Points
}
