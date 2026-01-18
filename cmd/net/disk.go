package net

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// diskCmd represents the disk command
var diskCmd = &cobra.Command{
	Use:   "disk [path]",
	Short: "Disk usage analysis",
	Long: `Analyze disk usage for a path.

Examples:
  devkit net disk /
  devkit net disk . --top 10`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDisk,
}

func init() {
	netCmd.AddCommand(diskCmd)

	diskCmd.Flags().IntP("top", "t", 0, "Show top N largest directories")
	diskCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runDisk(cmd *cobra.Command, args []string) error {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	_, _ = cmd.Flags().GetInt("top") // topN for future use
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	usage, err := disk.Usage(path)
	if err != nil {
		return fmt.Errorf("failed to get disk usage: %w", err)
	}

	result := map[string]interface{}{
		"path":    path,
		"total":   formatBytesDisk(usage.Total),
		"used":    formatBytesDisk(usage.Used),
		"free":    formatBytesDisk(usage.Free),
		"percent": fmt.Sprintf("%.1f%%", usage.UsedPercent),
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else if format == output.FormatTable {
		fmt.Printf("%-20s %15s %15s %15s %10s\n", "PATH", "TOTAL", "USED", "FREE", "USED%")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-20s %15s %15s %15s %10s\n",
			path, result["total"], result["used"], result["free"], result["percent"])
	} else {
		fmt.Printf("Disk Usage for %s:\n", path)
		fmt.Printf("  Total: %s\n", result["total"])
		fmt.Printf("  Used: %s (%s)\n", result["used"], result["percent"])
		fmt.Printf("  Free: %s\n", result["free"])
	}

	return nil
}

func formatBytesDisk(bytes uint64) string {
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
