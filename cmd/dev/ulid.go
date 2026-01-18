package dev

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// ulidCmd represents the ulid command
var ulidCmd = &cobra.Command{
	Use:   "ulid",
	Short: "Generate ULID (Universally Unique Lexicographically Sortable Identifier)",
	Long: `Generate ULID (Universally Unique Lexicographically Sortable Identifier) values.

ULID is a 26-character string that is:
- Universally unique
- Lexicographically sortable
- URL-safe
- Case-insensitive

Examples:
  devkit dev ulid                    # Generate a single ULID
  devkit dev ulid --count 5          # Generate 5 ULIDs
  devkit dev ulid --count 3 --output json`,
	RunE: runULID,
}

func init() {
	devCmd.AddCommand(ulidCmd)

	// Flag definitions
	ulidCmd.Flags().IntP("count", "c", 1, "Number of ULIDs to generate")
	ulidCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runULID(cmd *cobra.Command, args []string) error {
	count, _ := cmd.Flags().GetInt("count")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if count < 1 {
		return fmt.Errorf("count must be at least 1")
	}

	if count > 1000 {
		return fmt.Errorf("count cannot exceed 1000")
	}

	ulids, err := generateULID(count)
	if err != nil {
		return fmt.Errorf("failed to generate ULID: %w", err)
	}

	// Prepare result based on format
	if format == output.FormatJSON {
		result := map[string]interface{}{
			"count": count,
			"ulids": ulids,
		}
		output.PrintSuccess(format, result)
	} else {
		// Plain format - print each ULID on a new line
		if count == 1 {
			output.PrintSuccess(format, ulids[0])
		} else {
			output.PrintSuccess(format, ulids)
		}
	}

	return nil
}

func generateULID(count int) ([]string, error) {
	ulids := make([]string, count)
	entropy := ulid.DefaultEntropy()
	
	for i := 0; i < count; i++ {
		id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
		ulids[i] = id.String()
	}
	
	return ulids, nil
}
