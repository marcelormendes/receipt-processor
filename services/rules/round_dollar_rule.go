package rules

import (
	"context"
	"fmt"
	"math"

	"github.com/marcelorm/receipt-processor/models"
)

// RoundDollarRule returns the rule for calculating points if the total is a round dollar amount
// Rule 2: 50 points if the total is a round dollar amount
func RoundDollarRule() PointRule {
	return PointRule{
		Name:        "RoundDollarRule",
		Description: "50 points if the total is a round dollar amount",
		Apply: func(ctx context.Context, r models.Receipt) int {
			total := float64(r.Total)
			if total == math.Floor(total) {
				return 50
			}
			return 0
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			if points == 0 {
				return ""
			}
			return fmt.Sprintf("Rule 2: Added 50 points for round dollar amount $%.2f", float64(r.Total))
		},
	}
}