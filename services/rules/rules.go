package rules

import (
	"context"

	"github.com/marcelorm/receipt-processor/models"
)

// PointRule defines a single point calculation rule
type PointRule struct {
	Name             string                                   // Name of the rule
	Description      string                                   // Description of the rule
	Apply            func(context.Context, models.Receipt) int // Function that applies the rule and returns points
	FormatLogMessage func(int, models.Receipt) string         // Function to generate log message
}

// GetAllRules returns all the point calculation rules in order
func GetAllRules() []PointRule {
	return []PointRule{
		RetailerNameRule(),
		RoundDollarRule(),
		QuarterMultipleRule(),
		ItemPairsRule(),
		ItemDescriptionLengthRule(),
		OddDayRule(),
		AfternoonTimeRule(),
	}
}