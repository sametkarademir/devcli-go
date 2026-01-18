package utils

import (
	"io"
	"os"
	"strings"
)

// GetInput reads input from stdin, file, or arguments
func GetInput(cmdArgs []string, fileFlag string, stdinFlag bool) (string, error) {
	// Check stdin first
	if stdinFlag {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return "", err
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(string(bytes)), nil
		}
	}

	// Check file flag
	if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}

	// Check arguments
	if len(cmdArgs) > 0 {
		return cmdArgs[0], nil
	}

	return "", nil
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDir checks if a path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
