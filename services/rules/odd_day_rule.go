package rules

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelorm/receipt-processor/models"
)

// OddDayRule returns the rule for calculating points if the purchase day is odd
// Rule 6: 6 points if purchase day is odd
func OddDayRule() PointRule {
	return PointRule{
		Name:        "OddDayRule",
		Description: "6 points if purchase day is odd",
		Apply: func(ctx context.Context, r models.Receipt) int {
			purchaseDate, err := time.Parse("2006-01-02", string(r.PurchaseDate))
			if err != nil || purchaseDate.Day()%2 == 0 {
				return 0
			}
			return 6
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			if points == 0 {
				return ""
			}
			purchaseDate, _ := time.Parse("2006-01-02", string(r.PurchaseDate))
			return fmt.Sprintf("Rule 6: Added 6 points for purchase day %d being odd", purchaseDate.Day())
		},
	}
}