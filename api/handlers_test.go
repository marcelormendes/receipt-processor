package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
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

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Add our JSON validation middleware to mimic production
	router.Use(JSONValidationMiddleware(1024 * 1024))

	store := storage.NewMemoryStorage()
	handler := NewReceiptHandler(store)

	router.POST("/receipts/process", handler.ProcessReceipt)
	router.GET("/receipts/:id/points", handler.GetPoints)

	return router
}

func TestProcessReceipt(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		receipt        any
		expectedStatus int
		expectID       bool
	}{
		{
			name: "Valid receipt - Target example",
			receipt: map[string]any{
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
			},
			expectedStatus: http.StatusOK,
			expectID:       true,
		},
		{
			name: "Valid receipt - M&M Corner Market example",
			receipt: map[string]any{
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
			},
			expectedStatus: http.StatusOK,
			expectID:       true,
		},
		{
			name:           "Invalid receipt - Empty",
			receipt:        map[string]any{},
			expectedStatus: http.StatusBadRequest,
			expectID:       false,
		},
		{
			name: "Invalid date format",
			receipt: map[string]any{
				"retailer":     "Shop",
				"purchaseDate": "01/01/2022", // Wrong format
				"purchaseTime": "13:01",
				"items": []map[string]any{
					{"shortDescription": "Item", "price": "5.00"},
				},
				"total": "5.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectID:       false,
		},
		{
			name: "Invalid time format",
			receipt: map[string]any{
				"retailer":     "Shop",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "1:01 PM", // Wrong format
				"items": []map[string]any{
					{"shortDescription": "Item", "price": "5.00"},
				},
				"total": "5.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectID:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tc.receipt)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d: %s", tc.expectedStatus, resp.Code, resp.Body.String())
			}

			if tc.expectID {
				var response models.ReceiptResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response.ID == "" {
					t.Error("Expected non-empty ID")
				}
			}
		})
	}
}

func TestGetPoints(t *testing.T) {
	router := setupRouter()

	// First, process a receipt to get an ID
	receipt := map[string]any{
		"retailer":     "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": []map[string]any{
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
			{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
		},
		"total": "18.74",
	}

	reqBody, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var processResp models.ReceiptResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &processResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	tests := []struct {
		name           string
		receiptID      string
		expectedStatus int
		expectPoints   bool
	}{
		{
			name:           "Valid receipt ID",
			receiptID:      processResp.ID,
			expectedStatus: http.StatusOK,
			expectPoints:   true,
		},
		{
			name:           "Invalid receipt ID",
			receiptID:      "nonexistent-id",
			expectedStatus: http.StatusNotFound,
			expectPoints:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/receipts/"+tc.receiptID+"/points", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d: %s", tc.expectedStatus, resp.Code, resp.Body.String())
			}

			if tc.expectPoints {
				var pointsResp models.PointsResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &pointsResp); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if pointsResp.Points <= 0 {
					t.Error("Expected positive points value")
				}
			}
		})
	}
}
