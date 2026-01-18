package net

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// speedCmd represents the speed command
var speedCmd = &cobra.Command{
	Use:   "speed",
	Short: "Internet speed test",
	Long: `Test internet connection speed (download/upload).

Examples:
  devkit net speed
  devkit net speed --server-id 12345`,
	RunE: runSpeed,
}

func init() {
	netCmd.AddCommand(speedCmd)

	speedCmd.Flags().Int("server-id", 0, "Specific server ID to test")
	speedCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runSpeed(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	// Simple speed test using HTTP download
	testURL := "https://speed.cloudflare.com/__down?bytes=10000000" // 10MB

	fmt.Println("Testing download speed...")
	start := time.Now()
	resp, err := http.Get(testURL)
	if err != nil {
		return fmt.Errorf("speed test failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	duration := time.Since(start)
	bytesDownloaded := int64(10000000)
	mbps := float64(bytesDownloaded*8) / duration.Seconds() / 1000000

	// Ping test
	pingStart := time.Now()
	_, err = http.Get("https://www.google.com")
	pingDuration := time.Since(pingStart)

	result := map[string]interface{}{
		"download": fmt.Sprintf("%.2f Mbps", mbps),
		"ping":     fmt.Sprintf("%.2f ms", float64(pingDuration.Nanoseconds())/1000000),
		"note":     "Simplified speed test - for accurate results use dedicated speed test tools",
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("\nSpeed Test Results:\n")
		fmt.Printf("  Download: %s\n", result["download"])
		fmt.Printf("  Ping: %s\n", result["ping"])
		fmt.Printf("\nNote: %s\n", result["note"])
	}

	return nil
}
