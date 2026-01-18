package utils

import "strings"

// TrimSpace trims whitespace from a string
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty checks if a string is empty or only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
