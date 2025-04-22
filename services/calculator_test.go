package services

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/marcelorm/receipt-processor/models"
)

// BenchmarkCalculatePoints measures the performance of the point calculation
func BenchmarkCalculatePoints(b *testing.B) {
	// Save the original logger and replace with a no-op logger
	originalLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	
	// Restore the original logger when the benchmark is done
	defer slog.SetDefault(originalLogger)
	
	// Create a sample receipt for benchmarking
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: models.Date("2022-01-01"),
		PurchaseTime: models.Time("13:01"),
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
			{ShortDescription: "Emils Cheese Pizza", Price: 12.25},
			{ShortDescription: "Knorr Creamy Chicken", Price: 1.26},
			{ShortDescription: "Doritos Nacho Cheese", Price: 3.35},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: 12.00},
		},
		Total: 35.35,
	}

	ctx := context.Background()

	// Reset the timer for the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, _ = CalculatePoints(ctx, receipt)
	}
}

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int
	}{
		{
			name: "Target receipt example",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: models.Date("2022-01-01"),
				PurchaseTime: models.Time("13:01"),
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
					{ShortDescription: "Emils Cheese Pizza", Price: 12.25},
					{ShortDescription: "Knorr Creamy Chicken", Price: 1.26},
					{ShortDescription: "Doritos Nacho Cheese", Price: 3.35},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: 12.00},
				},
				Total: 35.35,
			},
			expected: 28,
		},
		{
			name: "M&M Corner Market receipt example",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: models.Date("2022-03-20"),
				PurchaseTime: models.Time("14:33"),
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: 2.25},
					{ShortDescription: "Gatorade", Price: 2.25},
					{ShortDescription: "Gatorade", Price: 2.25},
					{ShortDescription: "Gatorade", Price: 2.25},
				},
				Total: 9.00,
			},
			expected: 109,
		},
		// Additional edge cases
		{
			name: "Empty receipt with zero",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: models.Date("2022-01-01"),
				PurchaseTime: models.Time("12:00"),
				Items:        []models.Item{},
				Total:        0.00,
			},
			expected: 81, // 0 for retailer, 0 for items, 50 for round dollar, 25 for multiple of 0.25, 6 for odd day
		},
		{
			name: "Receipt with odd number of items",
			receipt: models.Receipt{
				Retailer:     "ABC",
				PurchaseDate: models.Date("2022-02-02"),
				PurchaseTime: models.Time("12:00"),
				Items: []models.Item{
					{ShortDescription: "Item 1", Price: 1.00},
					{ShortDescription: "Item 2", Price: 2.00},
					{ShortDescription: "Item 3", Price: 3.00},
				},
				Total: 6.00,
			},
			expected: 86, // 3 for retailer, 5 for 1 pair of items, 50 for round dollar amount, 25 for multiple of 0.25, 3 for item descriptions
		},
		{
			name: "Receipt with purchase time between 14:00 and 16:00",
			receipt: models.Receipt{
				Retailer:     "XYZ",
				PurchaseDate: models.Date("2022-02-01"),
				PurchaseTime: models.Time("15:30"),
				Items: []models.Item{
					{ShortDescription: "Item", Price: 1.00},
				},
				Total: 1.00,
			},
			expected: 94, // 3 for retailer, 0 pairs, 50 for round dollar, 25 for multiple of 0.25, 6 for odd day, 10 for time range
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create context for the test
			ctx := context.Background()

			// Call with context
			actual, err := CalculatePoints(ctx, tc.receipt)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Expected %d points, but got %d", tc.expected, actual)
			}
		})
	}

	// Test context cancellation
	t.Run("Context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		receipt := models.Receipt{
			Retailer:     "Test",
			PurchaseDate: models.Date("2022-01-01"),
			PurchaseTime: models.Time("12:00"),
			Items:        []models.Item{{ShortDescription: "Item", Price: 1.00}},
			Total:        1.00,
		}

		_, err := CalculatePoints(ctx, receipt)
		if err == nil {
			t.Error("Expected error due to context cancellation, got nil")
		}
	})
}
