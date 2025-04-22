package models

import (
	"encoding/json"
	"testing"
)

func TestReceiptJSON(t *testing.T) {
	// Test marshaling and unmarshaling Receipt
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: Date("2022-01-01"),
		PurchaseTime: Time("13:01"),
		Items: []Item{
			{ShortDescription: "Item 1", Price: 5.99},
			{ShortDescription: "Item 2", Price: 10.00},
		},
		Total: 15.99,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("Failed to marshal receipt: %v", err)
	}
	
	// Unmarshal from JSON
	var unmarshaledReceipt Receipt
	if err := json.Unmarshal(jsonData, &unmarshaledReceipt); err != nil {
		t.Fatalf("Failed to unmarshal receipt: %v", err)
	}
	
	// Verify fields
	if unmarshaledReceipt.Retailer != receipt.Retailer {
		t.Errorf("Expected retailer '%s', got '%s'", receipt.Retailer, unmarshaledReceipt.Retailer)
	}
	
	if unmarshaledReceipt.PurchaseDate != receipt.PurchaseDate {
		t.Errorf("Expected purchase date '%s', got '%s'", receipt.PurchaseDate, unmarshaledReceipt.PurchaseDate)
	}
	
	if unmarshaledReceipt.PurchaseTime != receipt.PurchaseTime {
		t.Errorf("Expected purchase time '%s', got '%s'", receipt.PurchaseTime, unmarshaledReceipt.PurchaseTime)
	}
	
	if unmarshaledReceipt.Total != receipt.Total {
		t.Errorf("Expected total %f, got %f", float64(receipt.Total), float64(unmarshaledReceipt.Total))
	}
	
	if len(unmarshaledReceipt.Items) != len(receipt.Items) {
		t.Errorf("Expected %d items, got %d", len(receipt.Items), len(unmarshaledReceipt.Items))
	}
}

func TestReceiptResponseJSON(t *testing.T) {
	// Test marshaling and unmarshaling ReceiptResponse
	response := ReceiptResponse{
		ID: "test-id-123",
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal receipt response: %v", err)
	}
	
	// Unmarshal from JSON
	var unmarshaledResponse ReceiptResponse
	if err := json.Unmarshal(jsonData, &unmarshaledResponse); err != nil {
		t.Fatalf("Failed to unmarshal receipt response: %v", err)
	}
	
	// Verify fields
	if unmarshaledResponse.ID != response.ID {
		t.Errorf("Expected ID '%s', got '%s'", response.ID, unmarshaledResponse.ID)
	}
}

func TestPointsResponseJSON(t *testing.T) {
	// Test marshaling and unmarshaling PointsResponse
	response := PointsResponse{
		Points: 100,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal points response: %v", err)
	}
	
	// Unmarshal from JSON
	var unmarshaledResponse PointsResponse
	if err := json.Unmarshal(jsonData, &unmarshaledResponse); err != nil {
		t.Fatalf("Failed to unmarshal points response: %v", err)
	}
	
	// Verify fields
	if unmarshaledResponse.Points != response.Points {
		t.Errorf("Expected points %d, got %d", response.Points, unmarshaledResponse.Points)
	}
}