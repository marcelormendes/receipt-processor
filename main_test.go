package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marcelorm/receipt-processor/api"
	"github.com/marcelorm/receipt-processor/storage"
)

// setupRouter creates a test router similar to the one in main
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Create a new in-memory receipt store
	store := storage.NewMemoryStorage()

	// Create a new receipt handler
	handler := api.NewReceiptHandler(store)

	// Create a new Gin router
	router := gin.Default()

	// Set up routes
	router.POST("/receipts/process", handler.ProcessReceipt)
	router.GET("/receipts/:id/points", handler.GetPoints)

	return router
}

func TestRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Process Receipt Route",
			method:         "POST",
			path:           "/receipts/process",
			expectedStatus: http.StatusBadRequest, // Expect bad request for empty body
		},
		{
			name:           "Get Points Route - Invalid ID",
			method:         "GET",
			path:           "/receipts/invalid-id/points",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unknown Route",
			method:         "GET",
			path:           "/unknown",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.Code)
			}
		})
	}
}
