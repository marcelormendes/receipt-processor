package errors

// New creates a new AppError with the given code
func New(code ErrorCode, detail string) *AppError {
	message, ok := errorMap[code]
	if !ok {
		message = "Unknown error"
	}

	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// Wrap wraps an existing error with an error code
func Wrap(code ErrorCode, err error, detail string) *AppError {
	message, ok := errorMap[code]
	if !ok {
		message = "Unknown error"
	}

	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
		Err:     err,
	}
}

// IsCode checks if an error has a specific error code
func IsCode(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}

	return false
}

// GetCode returns the error code from an error or empty string if not an AppError
func GetCode(err error) ErrorCode {
	if err == nil {
		return ""
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}

	return ""
}

// Error returns the standard message for the given error code.
func Error(code ErrorCode) string {
    msg, ok := errorMap[code]
    if !ok {
        return "Unknown error"
    }
    return msg
}
