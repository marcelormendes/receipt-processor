package rules

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/marcelorm/receipt-processor/models"
)

// ItemDescriptionLengthRule returns the rule for calculating points based on item description length
// Rule 5: If trimmed item desc length % 3 == 0: points += ceil(price * 0.2)
func ItemDescriptionLengthRule() PointRule {
	return PointRule{
		Name:        "ItemDescriptionLengthRule",
		Description: "If trimmed item description length is a multiple of 3, add ceil(price * 0.2) points",
		Apply: func(ctx context.Context, r models.Receipt) int {
			return calculateItemDescriptionPoints(r)
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			if points == 0 {
				return ""
			}
			
			return fmt.Sprintf("Rule 5: Added %d points for items with description length multiple of 3", points)
		},
	}
}

// calculateItemDescriptionPoints calculates points based on item descriptions
func calculateItemDescriptionPoints(r models.Receipt) int {
	totalPoints := 0
	for _, item := range r.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc) > 0 && len(trimmedDesc)%3 == 0 {
			price := float64(item.Price)
			itemPoints := int(math.Ceil(price * 0.2))
			totalPoints += itemPoints
		}
	}
	return totalPoints
}