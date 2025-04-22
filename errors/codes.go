package errors

// ErrorCode represents a standardized error code for the application
type ErrorCode string

// Error codes for the Receipt Processor API
const (
	// General errors (0000-0099)
	ErrInternal         ErrorCode = "RP0001" // Internal server error
	ErrInvalidJSON      ErrorCode = "RP0002" // Invalid JSON in request body
	ErrContextCancelled ErrorCode = "RP0003" // Request context cancelled or timed out
	ErrRequestTooLarge  ErrorCode = "RP0004" // Request body too large
	ErrInvalidRequest   ErrorCode = "RP0005" // Invalid request

	// Validation errors (0100-0199)
	ErrInvalidReceiptData    ErrorCode = "RP0101" // Invalid or missing receipt data
	ErrInvalidRetailer       ErrorCode = "RP0102" // Invalid or missing retailer field
	ErrInvalidPurchaseDate   ErrorCode = "RP0103" // Invalid purchase date format
	ErrInvalidPurchaseTime   ErrorCode = "RP0104" // Invalid purchase time format
	ErrInvalidTotal          ErrorCode = "RP0105" // Invalid total amount
	ErrMissingItems          ErrorCode = "RP0106" // No items in receipt
	ErrInvalidItemData       ErrorCode = "RP0107" // Invalid item data
	ErrInvalidItemDescription ErrorCode = "RP0108" // Invalid item description
	ErrInvalidItemPrice      ErrorCode = "RP0109" // Invalid item price

	// Storage errors (0200-0299)
	ErrReceiptNotFound ErrorCode = "RP0201" // Receipt ID not found
	ErrStorageFailure  ErrorCode = "RP0202" // Failed to store receipt
	
	// Calculation errors (0300-0399)
	ErrCalculationFailed ErrorCode = "RP0301" // Failed to calculate points
)

// errorMap maps error codes to standard error messages
var errorMap = map[ErrorCode]string{
	// General errors
	ErrInternal:         "Internal server error",
	ErrInvalidJSON:      "Invalid JSON in request body",
	ErrContextCancelled: "Request cancelled or timed out",
	ErrRequestTooLarge:  "Request body too large",
	ErrInvalidRequest:   "Invalid request",

	// Validation errors
	ErrInvalidReceiptData:    "Invalid or missing receipt data",
	ErrInvalidRetailer:       "Invalid or missing retailer name",
	ErrInvalidPurchaseDate:   "Invalid purchase date format",
	ErrInvalidPurchaseTime:   "Invalid purchase time format",
	ErrInvalidTotal:          "Invalid total amount",
	ErrMissingItems:          "Receipt must contain at least one item",
	ErrInvalidItemData:       "Invalid item data",
	ErrInvalidItemDescription: "Invalid item description",
	ErrInvalidItemPrice:      "Invalid item price",

	// Storage errors
	ErrReceiptNotFound: "Receipt not found",
	ErrStorageFailure:  "Failed to store receipt data",
	
	// Calculation errors
	ErrCalculationFailed: "Failed to calculate receipt points",
}