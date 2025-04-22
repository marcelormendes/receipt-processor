package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	// Test creating a new error
	err := New(ErrInvalidReceiptData, "missing required fields")

	if err.Code != ErrInvalidReceiptData {
		t.Errorf("Expected code %s, got %s", ErrInvalidReceiptData, err.Code)
	}

	if err.Message != errorMap[ErrInvalidReceiptData] {
		t.Errorf("Expected message %s, got %s", errorMap[ErrInvalidReceiptData], err.Message)
	}

	if err.Detail != "missing required fields" {
		t.Errorf("Expected detail %s, got %s", "missing required fields", err.Detail)
	}

	if err.Err != nil {
		t.Errorf("Expected nil original error, got %v", err.Err)
	}

	// Test with unknown error code
	unknownCode := ErrorCode("UNKNOWN")
	err = New(unknownCode, "test detail")

	if err.Message != "Unknown error" {
		t.Errorf("Expected message 'Unknown error', got %s", err.Message)
	}
}

func TestWrap(t *testing.T) {
	// Create an original error
	originalErr := fmt.Errorf("database connection failed")

	// Wrap the error
	err := Wrap(ErrStorageFailure, originalErr, "while connecting to database")

	// Check properties
	if err.Code != ErrStorageFailure {
		t.Errorf("Expected code %s, got %s", ErrStorageFailure, err.Code)
	}

	if err.Message != errorMap[ErrStorageFailure] {
		t.Errorf("Expected message %s, got %s", errorMap[ErrStorageFailure], err.Message)
	}

	if err.Detail != "while connecting to database" {
		t.Errorf("Expected detail %s, got %s", "while connecting to database", err.Detail)
	}

	if err.Err != originalErr {
		t.Errorf("Original error not stored correctly")
	}

	// Test unwrap
	unwrapped := err.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Expected unwrapped error to be original error")
	}

	// Test with unknown error code
	unknownCode := ErrorCode("UNKNOWN")
	err = Wrap(unknownCode, originalErr, "test detail")

	if err.Message != "Unknown error" {
		t.Errorf("Expected message 'Unknown error', got %s", err.Message)
	}
}

func TestIsCode(t *testing.T) {
	// Create errors with different codes
	err1 := New(ErrInvalidReceiptData, "test detail")
	err2 := New(ErrReceiptNotFound, "test detail")

	// Test IsCode function
	if !IsCode(err1, ErrInvalidReceiptData) {
		t.Errorf("IsCode failed to identify correct error code")
	}

	if IsCode(err1, ErrReceiptNotFound) {
		t.Errorf("IsCode incorrectly matched wrong error code")
	}

	if IsCode(err2, ErrInvalidReceiptData) {
		t.Errorf("IsCode incorrectly matched wrong error code")
	}

	if !IsCode(err2, ErrReceiptNotFound) {
		t.Errorf("IsCode failed to identify correct error code")
	}

	// Test with nil error
	if IsCode(nil, ErrInvalidReceiptData) {
		t.Errorf("IsCode should return false for nil error")
	}

	// Test with non-AppError
	regularErr := fmt.Errorf("regular error")
	if IsCode(regularErr, ErrInvalidReceiptData) {
		t.Errorf("IsCode should return false for non-AppError")
	}
}

func TestGetCode(t *testing.T) {
	// Create an app error
	err := New(ErrInvalidReceiptData, "test detail")

	// Test GetCode function
	code := GetCode(err)
	if code != ErrInvalidReceiptData {
		t.Errorf("Expected code %s, got %s", ErrInvalidReceiptData, code)
	}

	// Test with nil error
	code = GetCode(nil)
	if code != "" {
		t.Errorf("Expected empty code for nil error, got %s", code)
	}

	// Test with non-AppError
	regularErr := fmt.Errorf("regular error")
	code = GetCode(regularErr)
	if code != "" {
		t.Errorf("Expected empty code for regular error, got %s", code)
	}
}
