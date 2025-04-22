package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/marcelorm/receipt-processor/models"
)

func TestRetailerNameRule(t *testing.T) {
	ctx := context.Background()
	rule := RetailerNameRule()

	tests := []struct {
		name     string
		retailer string
		expected int
	}{
		{"Empty retailer", "", 0},
		{"Retailer with alphanumeric", "Target", 6},
		{"Retailer with spaces", "Target Store", 11},
		{"Retailer with special chars", "M&M Corner Market", 14},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			receipt := models.Receipt{Retailer: tc.retailer}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestRoundDollarRule(t *testing.T) {
	ctx := context.Background()
	rule := RoundDollarRule()

	tests := []struct {
		name     string
		total    float64
		expected int
	}{
		{"Round dollar amount", 100.00, 50},
		{"Round dollar zero", 0.00, 50},
		{"Not round dollar", 99.99, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			receipt := models.Receipt{Total: models.Price(tc.total)}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestQuarterMultipleRule(t *testing.T) {
	ctx := context.Background()
	rule := QuarterMultipleRule()

	tests := []struct {
		name     string
		total    float64
		expected int
	}{
		{"Multiple of 0.25 - #1", 10.00, 25},
		{"Multiple of 0.25 - #2", 10.25, 25},
		{"Multiple of 0.25 - #3", 10.50, 25},
		{"Multiple of 0.25 - #4", 10.75, 25},
		{"Not multiple of 0.25", 10.13, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			receipt := models.Receipt{Total: models.Price(tc.total)}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestItemPairsRule(t *testing.T) {
	ctx := context.Background()
	rule := ItemPairsRule()

	tests := []struct {
		name     string
		itemCount int
		expected int
	}{
		{"No items", 0, 0},
		{"One item", 1, 0},
		{"Two items", 2, 5},
		{"Three items", 3, 5},
		{"Four items", 4, 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a receipt with the specified number of items
			items := make([]models.Item, tc.itemCount)
			for i := 0; i < tc.itemCount; i++ {
				items[i] = models.Item{ShortDescription: "Item", Price: 1.00}
			}
			
			receipt := models.Receipt{Items: items}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestItemDescriptionLengthRule(t *testing.T) {
	ctx := context.Background()
	rule := ItemDescriptionLengthRule()

	// Debugging function to check divisions by 3
	checkDivisibility := func(s string) {
		trimmed := strings.TrimSpace(s)
		t.Logf("Description '%s': length=%d, divisible by 3=%v", 
			trimmed, len(trimmed), len(trimmed)%3 == 0)
	}

	// Check our test strings
	checkDivisibility("Item1")    // Length 5
	checkDivisibility("Item2")    // Length 5
	checkDivisibility("Item")     // Length 4

	tests := []struct {
		name        string
		descriptions []string
		prices      []float64
		expected    int
	}{
		{
			"No matching items",
			[]string{"Item1", "Item2"},  // Lengths 5 and 5, not divisible by 3
			[]float64{1.00, 2.00},
			0,
		},
		{
			"Length divisible by 3",
			[]string{"123", "123456"},
			[]float64{5.00, 10.00},
			3, // ceil(5.00 * 0.2) + ceil(10.00 * 0.2) = 1 + 2 = 3
		},
		{
			"With whitespace trimming",
			[]string{"   123   ", "Item"},
			[]float64{5.00, 1.00},
			1, // ceil(5.00 * 0.2) = 1
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create items with the specified descriptions and prices
			items := make([]models.Item, len(tc.descriptions))
			for i, desc := range tc.descriptions {
				items[i] = models.Item{
					ShortDescription: desc,
					Price:            models.Price(tc.prices[i]),
				}
			}
			
			receipt := models.Receipt{Items: items}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestOddDayRule(t *testing.T) {
	ctx := context.Background()
	rule := OddDayRule()

	tests := []struct {
		name     string
		date     string
		expected int
	}{
		{"Odd day", "2022-01-01", 6},   // January 1st is odd
		{"Even day", "2022-01-02", 0},  // January 2nd is even
		{"Invalid date", "invalid", 0}, // Invalid date should return 0 points
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			receipt := models.Receipt{PurchaseDate: models.Date(tc.date)}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestAfternoonTimeRule(t *testing.T) {
	ctx := context.Background()
	rule := AfternoonTimeRule()

	tests := []struct {
		name     string
		time     string
		expected int
	}{
		{"Before time window", "13:59", 0},  // Just before window
		{"Start of window", "14:00", 0},     // At start (exclusive)
		{"In time window", "15:00", 10},     // Inside window
		{"End of window", "16:00", 0},       // At end (exclusive)
		{"After time window", "16:01", 0},   // Just after window
		{"Invalid time", "invalid", 0},      // Invalid time should return 0 points
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			receipt := models.Receipt{PurchaseTime: models.Time(tc.time)}
			points := rule.Apply(ctx, receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestGetAllRules(t *testing.T) {
	rules := GetAllRules()
	
	// Verify we have 7 rules
	if len(rules) != 7 {
		t.Errorf("Expected 7 rules, got %d", len(rules))
	}
	
	// Verify the rule names
	expectedNames := []string{
		"RetailerNameRule",
		"RoundDollarRule",
		"QuarterMultipleRule",
		"ItemPairsRule",
		"ItemDescriptionLengthRule",
		"OddDayRule",
		"AfternoonTimeRule",
	}
	
	for i, rule := range rules {
		if rule.Name != expectedNames[i] {
			t.Errorf("Expected rule #%d to be %s, got %s", i+1, expectedNames[i], rule.Name)
		}
	}
}