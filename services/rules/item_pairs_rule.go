package rules

import (
	"context"
	"fmt"

	"github.com/marcelorm/receipt-processor/models"
)

// ItemPairsRule returns the rule for calculating points based on item pairs
// Rule 4: 5 points for every two items
func ItemPairsRule() PointRule {
	return PointRule{
		Name:        "ItemPairsRule",
		Description: "5 points for every two items",
		Apply: func(ctx context.Context, r models.Receipt) int {
			return (len(r.Items) / 2) * 5
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			return fmt.Sprintf("Rule 4: Added %d points for %d pairs of items (%d items total)", 
				points, len(r.Items)/2, len(r.Items))
		},
	}
}