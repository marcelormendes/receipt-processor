package errors

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestAppError(t *testing.T) {
	// Test creating a new error
	err := &AppError{
		Code:    ErrInvalidReceiptData,
		Message: errorMap[ErrInvalidReceiptData],
		Detail:  "missing required fields",
	}

	if err.Code != ErrInvalidReceiptData {
		t.Errorf("Expected code %s, got %s", ErrInvalidReceiptData, err.Code)
	}

	if err.Message != errorMap[ErrInvalidReceiptData] {
		t.Errorf("Expected message %s, got %s", errorMap[ErrInvalidReceiptData], err.Message)
	}

	if err.Detail != "missing required fields" {
		t.Errorf("Expected detail %s, got %s", "missing required fields", err.Detail)
	}

	// Test error string
	expected := fmt.Sprintf("[%s] %s: %s", ErrInvalidReceiptData, errorMap[ErrInvalidReceiptData], "missing required fields")
	if err.Error() != expected {
		t.Errorf("Expected error string %s, got %s", expected, err.Error())
	}

	// Test error string without detail
	err.Detail = ""
	expected = fmt.Sprintf("[%s] %s", ErrInvalidReceiptData, errorMap[ErrInvalidReceiptData])
	if err.Error() != expected {
		t.Errorf("Expected error string %s, got %s", expected, err.Error())
	}
}

func TestAppErrorUnwrap(t *testing.T) {
	// Create an original error
	originalErr := fmt.Errorf("database connection failed")

	// Create an app error that wraps the original error
	err := &AppError{
		Code:    ErrStorageFailure,
		Message: errorMap[ErrStorageFailure],
		Detail:  "while connecting to database",
		Err:     originalErr,
	}

	// Test unwrap
	unwrapped := err.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Expected unwrapped error to be original error")
	}
}

func TestAppErrorLog(t *testing.T) {
	// This is mostly a compile-time test to ensure the Log method works
	// A more comprehensive test would use a logger mock
	err := &AppError{
		Code:    ErrInvalidReceiptData,
		Message: errorMap[ErrInvalidReceiptData],
		Detail:  "test detail",
		Err:     fmt.Errorf("original error"),
	}

	// Create a simple string builder to capture log output
	var logOutput strings.Builder

	// Log the error (this won't actually test capturing the output properly,
	// but at least ensures the method executes without panicking)
	err.Log(context.Background())

	// In a real test, we would verify the log output format
	if len(logOutput.String()) > 0 {
		// Just to use the variable
		t.Log(logOutput.String())
	}
}
