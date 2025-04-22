package services

import (
	"context"
	"log/slog"

	rperrors "github.com/marcelorm/receipt-processor/errors"
	"github.com/marcelorm/receipt-processor/models"
	"github.com/marcelorm/receipt-processor/services/rules"
)

// CalculatePoints calculates the total points for a receipt according to the rules
func CalculatePoints(ctx context.Context, receipt models.Receipt) (int, error) {
	// Check if context is canceled before proceeding
	select {
	case <-ctx.Done():
		return 0, rperrors.Wrap(rperrors.ErrContextCancelled, ctx.Err(), "context cancelled during point calculation")
	default:
		// Continue with normal operation
	}

	slog.InfoContext(ctx, "Calculating points for receipt",
		"retailer", receipt.Retailer,
		"date", receipt.PurchaseDate,
		"time", receipt.PurchaseTime,
		"items_count", len(receipt.Items))

	// Get all rules
	pointRules := rules.GetAllRules()

	// Apply all rules and sum the points
	totalPoints := 0
	for i, rule := range pointRules {
		// Check if context is canceled before each rule evaluation
		select {
		case <-ctx.Done():
			return 0, rperrors.Wrap(rperrors.ErrContextCancelled, ctx.Err(), "context cancelled during rule evaluation")
		default:
			// Continue with normal operation
		}

		points := rule.Apply(ctx, receipt)
		if points > 0 {
			logMsg := rule.FormatLogMessage(points, receipt)
			if logMsg != "" {
				slog.InfoContext(ctx, logMsg, "rule_number", i+1, "points", points)
			}
			totalPoints += points
		}
	}

	slog.InfoContext(ctx, "Total points calculated",
		"retailer", receipt.Retailer,
		"total_points", totalPoints)

	return totalPoints, nil
}
