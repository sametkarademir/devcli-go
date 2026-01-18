package utils

// ValidateNotEmpty validates that a string is not empty
func ValidateNotEmpty(s, fieldName string) error {
	if IsEmpty(s) {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
