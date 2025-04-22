package errors

import "testing"

func TestErrorCodes(t *testing.T) {
	// Test that all error codes have a corresponding message
	for code := range errorMap {
		if errorMap[code] == "" {
			t.Errorf("ErrorCode %s has an empty message", code)
		}
	}

	// Test error code prefix
	allCodes := []ErrorCode{
		// General errors
		ErrInternal, ErrInvalidJSON, ErrContextCancelled,

		// Validation errors
		ErrInvalidReceiptData, ErrInvalidRetailer, ErrInvalidPurchaseDate,
		ErrInvalidPurchaseTime, ErrInvalidTotal, ErrMissingItems,
		ErrInvalidItemData, ErrInvalidItemDescription, ErrInvalidItemPrice,

		// Storage errors
		ErrReceiptNotFound, ErrStorageFailure,

		// Calculation errors
		ErrCalculationFailed,
	}

	for _, code := range allCodes {
		if len(code) < 2 || code[:2] != "RP" {
			t.Errorf("ErrorCode %s does not start with the required prefix RP", code)
		}
	}
}
