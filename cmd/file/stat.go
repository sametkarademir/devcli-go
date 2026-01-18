package file

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// statCmd represents the stat command
var statCmd = &cobra.Command{
	Use:   "stat [file]",
	Short: "Display detailed file information",
	Long: `Display detailed information about a file or directory.

Examples:
  devkit file stat README.md
  devkit file stat /path/to/file
  devkit file stat .`,
	Args: cobra.ExactArgs(1),
	RunE: runStat,
}

func init() {
	fileCmd.AddCommand(statCmd)

	statCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runStat(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	path := args[0]
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	result := map[string]interface{}{
		"name":    info.Name(),
		"size":    info.Size(),
		"mode":    info.Mode().String(),
		"mod_time": info.ModTime().Format(time.RFC3339),
		"is_dir":  info.IsDir(),
	}

	if !info.IsDir() {
		result["size_human"] = formatSize(info.Size())
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Name: %s\n", result["name"])
		fmt.Printf("Size: %d bytes (%s)\n", result["size"], result["size_human"])
		fmt.Printf("Mode: %s\n", result["mode"])
		fmt.Printf("Modified: %s\n", result["mod_time"])
		fmt.Printf("Is Directory: %v\n", result["is_dir"])
	}

	return nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
