package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	rperrors "github.com/marcelorm/receipt-processor/errors"
	"github.com/marcelorm/receipt-processor/models"
)

// JSONValidationMiddleware pre-validates JSON requests before they reach your handlers
// This can fail fast on malformed requests and save processing time
func JSONValidationMiddleware(maxBodySize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBodySize {
			c.JSON(http.StatusRequestEntityTooLarge, rperrors.APIError{
				Code:    rperrors.ErrRequestTooLarge,
				Message: rperrors.Error(rperrors.ErrRequestTooLarge),
			})
			c.Abort()
			return
		}
		// Only validate for POST /receipts/process
		if c.Request.Method == http.MethodPost && c.FullPath() == "/receipts/process" {
			var receipt models.Receipt
			if err := c.ShouldBindJSON(&receipt); err != nil {
				slog.Error("Invalid JSON in request", "error", err)
				c.JSON(http.StatusBadRequest, rperrors.APIError{
					Code:    rperrors.ErrInvalidJSON,
					Message: rperrors.Error(rperrors.ErrInvalidJSON),
				})
				c.Abort()
				return
			}
			if err := receipt.Validate(); err != nil {
				slog.Error("Invalid receipt data", "error", err)
				code := rperrors.ErrInvalidReceiptData
				switch {
				case err.Error() == "retailer is required":
					code = rperrors.ErrInvalidRetailer
				case err.Error() == "at least one item is required":
					code = rperrors.ErrMissingItems
				case err.Error() == "invalid date format":
					code = rperrors.ErrInvalidPurchaseDate
				case err.Error() == "invalid time format":
					code = rperrors.ErrInvalidPurchaseTime
				}
				c.JSON(http.StatusBadRequest, rperrors.APIError{
					Code:    code,
					Message: rperrors.Error(code),
				})
				c.Abort()
				return
			}
			c.Set("receipt", receipt)
		}
		c.Next()
	}
}

func SetupRouter(maxBodySize int64) *gin.Engine {
	r := gin.New()

	// Apply middleware at router group level
	v1 := r.Group("/receipts")
	v1.Use(JSONValidationMiddleware(maxBodySize))

	// Process endpoint
	v1.POST("/process", func(c *gin.Context) {
		var receipt models.Receipt
		if err := c.ShouldBindJSON(&receipt); err != nil {
			slog.Error("Invalid JSON in request", "error", err)
			c.JSON(http.StatusBadRequest, rperrors.APIError{
				Code:    rperrors.ErrInvalidJSON,
				Message: rperrors.Error(rperrors.ErrInvalidJSON),
			})
			c.Abort()
			return
		}

		// Validate the receipt schema
		if err := receipt.Validate(); err != nil {
			slog.Error("Invalid receipt data", "error", err)

			// Determine specific validation error code
			code := rperrors.ErrInvalidReceiptData

			switch {
			case err.Error() == "retailer is required":
				code = rperrors.ErrInvalidRetailer
			case err.Error() == "at least one item is required":
				code = rperrors.ErrMissingItems
			}

			c.JSON(http.StatusBadRequest, rperrors.APIError{
				Code:    code,
				Message: rperrors.Error(code),
			})
			c.Abort()
			return
		}

		// Continue processing
		c.Next()
	})

	return r
}
