package file

import (
	"github.com/spf13/cobra"
)

// fileCmd represents the file command group
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File and directory operations",
	Long: `File and directory operations for managing files and directories.

This command group includes utilities for:
- File search and find-replace
- Bulk rename operations
- Format conversion
- File comparison (diff)
- Duplicate file detection
- File watching
- Directory tree visualization
- File statistics
- And more...`,
}

// GetFileCmd returns the file command
func GetFileCmd() *cobra.Command {
	return fileCmd
}

func init() {
	// This will be called when the package is imported
	// Commands will be added in their respective files
}
