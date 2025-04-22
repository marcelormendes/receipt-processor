package rules

import (
	"context"
	"fmt"
	"math"

	"github.com/marcelorm/receipt-processor/models"
)

// QuarterMultipleRule returns the rule for calculating points if the total is a multiple of 0.25
// Rule 3: 25 points if the total is a multiple of 0.25
func QuarterMultipleRule() PointRule {
	return PointRule{
		Name:        "QuarterMultipleRule",
		Description: "25 points if the total is a multiple of 0.25",
		Apply: func(ctx context.Context, r models.Receipt) int {
			total := float64(r.Total)
			if math.Mod(total*100, 25) == 0 {
				return 25
			}
			return 0
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			if points == 0 {
				return ""
			}
			return fmt.Sprintf("Rule 3: Added 25 points for total $%.2f being a multiple of 0.25", float64(r.Total))
		},
	}
}