package dev

import (
	"github.com/spf13/cobra"
)

// devCmd represents the dev command group
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Developer tools and utilities",
	Long: `Developer tools and utilities for common development tasks.

This command group includes utilities for:
- UUID/ULID generation
- Hash calculation
- Base64 encoding/decoding
- JWT operations
- Regex testing
- Color conversion
- URL encoding/decoding
- HTML entity operations
- JSON operations
- Random data generation
- And more...`,
}

// GetDevCmd returns the dev command
func GetDevCmd() *cobra.Command {
	return devCmd
}

func init() {
	// This will be called when the package is imported
	// Commands will be added in their respective files
}
