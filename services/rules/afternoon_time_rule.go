package rules

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelorm/receipt-processor/models"
)

// AfternoonTimeRule returns the rule for calculating points if the purchase time is between 2-4 PM
// Rule 7: 10 points if purchase time is 14:00 < time < 16:00
func AfternoonTimeRule() PointRule {
	return PointRule{
		Name:        "AfternoonTimeRule",
		Description: "10 points if purchase time is between 14:00 and 16:00 (exclusive)",
		Apply: func(ctx context.Context, r models.Receipt) int {
			return calculateAfternoonTimePoints(r)
		},
		FormatLogMessage: func(points int, r models.Receipt) string {
			if points == 0 {
				return ""
			}
			return fmt.Sprintf("Rule 7: Added 10 points for purchase time %s being between 14:00 and 16:00",
				string(r.PurchaseTime))
		},
	}
}

// calculateAfternoonTimePoints calculates points for afternoon purchase time
func calculateAfternoonTimePoints(r models.Receipt) int {
	purchaseTime, err := time.Parse("15:04", string(r.PurchaseTime))
	if err != nil {
		return 0
	}

	hour, minute := purchaseTime.Hour(), purchaseTime.Minute()
	timeInMinutes := hour*60 + minute
	startTimeInMinutes, endTimeInMinutes := 14*60, 16*60

	if timeInMinutes <= startTimeInMinutes || timeInMinutes >= endTimeInMinutes {
		return 0
	}
	return 10
}
