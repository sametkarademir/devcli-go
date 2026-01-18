package errors

import "fmt"

// DevKitError represents a custom error type for DevKit
type DevKitError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *DevKitError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *DevKitError) Unwrap() error {
	return e.Err
}

// Predefined errors
var (
	ErrFileNotFound     = &DevKitError{Code: "FILE_NOT_FOUND", Message: "Dosya bulunamadı"}
	ErrInvalidInput     = &DevKitError{Code: "INVALID_INPUT", Message: "Geçersiz giriş"}
	ErrNetworkTimeout   = &DevKitError{Code: "NETWORK_TIMEOUT", Message: "Ağ zaman aşımı"}
	ErrPermissionDenied = &DevKitError{Code: "PERMISSION_DENIED", Message: "Erişim izni yok"}
)

// New creates a new DevKitError
func New(code, message string) *DevKitError {
	return &DevKitError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with a DevKitError
func Wrap(err error, code, message string) *DevKitError {
	return &DevKitError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
