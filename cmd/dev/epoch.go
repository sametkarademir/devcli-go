package dev

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// epochCmd represents the epoch command
var epochCmd = &cobra.Command{
	Use:   "epoch [timestamp|date]",
	Short: "Convert between Unix timestamp and date",
	Long: `Convert between Unix timestamp and human-readable date.

Examples:
  devkit dev epoch 1699876543                    # Convert timestamp to date
  devkit dev epoch --to-unix "2024-01-15 10:30"  # Convert date to timestamp
  devkit dev epoch now                          # Current timestamp`,
	RunE: runEpoch,
}

func init() {
	devCmd.AddCommand(epochCmd)

	epochCmd.Flags().String("to-unix", "", "Convert date string to Unix timestamp")
	epochCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runEpoch(cmd *cobra.Command, args []string) error {
	toUnix, _ := cmd.Flags().GetString("to-unix")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var result map[string]interface{}

	if toUnix != "" {
		// Convert date to Unix timestamp
		layouts := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"2006-01-02",
			time.RFC822,
			time.RFC1123,
		}

		var t time.Time
		var err error
		parsed := false

		for _, layout := range layouts {
			t, err = time.Parse(layout, toUnix)
			if err == nil {
				parsed = true
				break
			}
		}

		if !parsed {
			return fmt.Errorf("failed to parse date: %s (supported formats: RFC3339, 2006-01-02 15:04:05, 2006-01-02)", toUnix)
		}

		unix := t.Unix()
		result = map[string]interface{}{
			"timestamp": unix,
			"date":      t.Format(time.RFC3339),
			"input":     toUnix,
		}
	} else if len(args) > 0 {
		input := args[0]

		if input == "now" {
			// Current timestamp
			now := time.Now()
			result = map[string]interface{}{
				"timestamp": now.Unix(),
				"date":      now.Format(time.RFC3339),
				"utc":       now.UTC().Format(time.RFC3339),
			}
		} else {
			// Convert timestamp to date
			timestamp, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid timestamp: %s", input)
			}

			t := time.Unix(timestamp, 0)
			result = map[string]interface{}{
				"timestamp": timestamp,
				"date":      t.Format(time.RFC3339),
				"utc":       t.UTC().Format(time.RFC3339),
				"unix":      timestamp,
			}
		}
	} else {
		return fmt.Errorf("timestamp or date not specified")
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		if timestamp, ok := result["timestamp"].(int64); ok {
			fmt.Printf("Timestamp: %d\n", timestamp)
		}
		if date, ok := result["date"].(string); ok {
			fmt.Printf("Date: %s\n", date)
		}
		if utc, ok := result["utc"].(string); ok {
			fmt.Printf("UTC: %s\n", utc)
		}
	}

	return nil
}
