package utils

import (
	"strconv"
)

// StringToInt converts a string to int
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt64 converts a string to int64
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// IntToString converts an int to string
func IntToString(i int) string {
	return strconv.Itoa(i)
}
