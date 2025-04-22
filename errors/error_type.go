package errors

import (
	"context"
	"fmt"
	"log/slog"
)

// AppError represents an application error with a standardized code and message
type AppError struct {
	Code    ErrorCode // Standardized error code
	Message string    // User-friendly error message
	Detail  string    // Detailed error information (for logging, not user-facing)
	Err     error     // Original error (if any)
}

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Log logs the error with the standard logger
func (e *AppError) Log(ctx context.Context) {
	// Create attribute list for structured logging
	attrs := []any{
		"error_code", e.Code,
		"error_message", e.Message,
	}

	if e.Detail != "" {
		attrs = append(attrs, "error_detail", e.Detail)
	}

	if e.Err != nil {
		attrs = append(attrs, "error_cause", e.Err.Error())
	}

	// Log with context
	slog.ErrorContext(ctx, "Application error occurred", attrs...)
}
