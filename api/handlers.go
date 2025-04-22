package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	rperrors "github.com/marcelorm/receipt-processor/errors"
	"github.com/marcelorm/receipt-processor/models"
	"github.com/marcelorm/receipt-processor/services"
	"github.com/marcelorm/receipt-processor/storage"
)

// ReceiptHandler handles receipt-related HTTP endpoints
type ReceiptHandler struct {
	store storage.ReceiptStorage
}

// NewReceiptHandler creates a new receipt handler
func NewReceiptHandler(store storage.ReceiptStorage) *ReceiptHandler {
	return &ReceiptHandler{
		store: store,
	}
}

// handleError standardizes error response handling across endpoints
func handleError(c *gin.Context, err error) {
	ctx := c.Request.Context()

	// Handle app errors with standard codes
	if appErr, ok := err.(*rperrors.AppError); ok {
		// Log the error
		appErr.Log(ctx)

		// Map error codes to HTTP status codes
		var status int
		switch appErr.Code {
		case rperrors.ErrInvalidJSON,
			rperrors.ErrInvalidReceiptData,
			rperrors.ErrInvalidRetailer,
			rperrors.ErrInvalidPurchaseDate,
			rperrors.ErrInvalidPurchaseTime,
			rperrors.ErrInvalidTotal,
			rperrors.ErrMissingItems,
			rperrors.ErrInvalidItemData,
			rperrors.ErrInvalidItemDescription,
			rperrors.ErrInvalidItemPrice:
			status = http.StatusBadRequest
		case rperrors.ErrReceiptNotFound:
			status = http.StatusNotFound
		case rperrors.ErrContextCancelled:
			status = http.StatusRequestTimeout
		default:
			status = http.StatusInternalServerError
		}

		// Return error response with code and message
		c.JSON(status, rperrors.APIError{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// Handle generic errors (should be avoided in production)
	slog.ErrorContext(ctx, "Unexpected non-application error", "error", err)
	c.JSON(http.StatusInternalServerError, rperrors.APIError{
		Code:    rperrors.ErrInternal,
		Message: "An unexpected error occurred",
	})
}

// ProcessReceipt handles the POST /receipts/process endpoint
func (h *ReceiptHandler) ProcessReceipt(c *gin.Context) {
	ctx := c.Request.Context()
	// Retrieve the validated receipt from context
	receiptVal, exists := c.Get("receipt")
	if !exists {
		handleError(c, rperrors.New(rperrors.ErrInvalidReceiptData, "receipt not found in context"))
		return
	}
	receipt, ok := receiptVal.(models.Receipt)
	if !ok {
		handleError(c, rperrors.New(rperrors.ErrInvalidReceiptData, "invalid receipt type in context"))
		return
	}

	// Calculate points for the receipt
	points, err := services.CalculatePoints(ctx, receipt)
	if err != nil {
		if rperrors.IsCode(err, rperrors.ErrContextCancelled) {
			handleError(c, err)
		} else {
			handleError(c, rperrors.Wrap(rperrors.ErrCalculationFailed, err,
				"error calculating points for receipt"))
		}
		return
	}

	// Store the points and get an ID
	id, err := h.store.SaveReceipt(ctx, points)
	if err != nil {
		if rperrors.IsCode(err, rperrors.ErrContextCancelled) {
			handleError(c, err)
		} else {
			handleError(c, rperrors.Wrap(rperrors.ErrStorageFailure, err,
				"unable to save receipt points"))
		}
		return
	}

	slog.InfoContext(ctx, "Receipt processed successfully",
		"id", id,
		"points", points,
		"retailer", receipt.Retailer)

	// Return the ID
	c.JSON(http.StatusOK, models.ReceiptResponse{ID: id})
}

// GetPoints handles the GET /receipts/{id}/points endpoint
func (h *ReceiptHandler) GetPoints(c *gin.Context) {
	// Get request context
	ctx := c.Request.Context()

	// Get the receipt ID from the URL
	id := c.Param("id")

	// Validate ID
	if id == "" {
		handleError(c, rperrors.New(rperrors.ErrInvalidReceiptData, "receipt ID is required"))
		return
	}

	slog.InfoContext(ctx, "Getting points for receipt", "id", id)

	// Look up the points for the ID
	points, err := h.store.GetPoints(ctx, id)
	if err != nil {
		// Use rperrors.IsCode for receipt not found
		if rperrors.IsCode(err, rperrors.ErrReceiptNotFound) {
			handleError(c, err)
			return
		}
		handleError(c, rperrors.Wrap(rperrors.ErrInternal, err, "error retrieving points for ID "+id))
		return
	}

	slog.InfoContext(ctx, "Points retrieved successfully", "id", id, "points", points)

	// Return the points
	c.JSON(http.StatusOK, models.PointsResponse{Points: points})
}
