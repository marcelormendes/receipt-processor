package rules

import (
	"context"
	"fmt"
	"unicode"

	"github.com/marcelorm/receipt-processor/models"
)

// RetailerNameRule returns the rule for calculating points based on the retailer name
// Rule 1: One point for every alphanumeric character in the retailer name
func RetailerNameRule() PointRule {
	return PointRule{
		Name:        "RetailerNameRule",
		Description: "One point for every alphanumeric character in the retailer name",
		Apply: func(ctx context.Context, r models.Receipt) int {
			return countAlphanumericChars(r.Retailer)
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			return fmt.Sprintf("Rule 1: Added %d points for alphanumeric characters in retailer name '%s'", 
				points, r.Retailer)
		},
	}
}

// countAlphanumericChars counts alphanumeric characters in a string
func countAlphanumericChars(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return count
}