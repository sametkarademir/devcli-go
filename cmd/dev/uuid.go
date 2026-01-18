package dev

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// uuidCmd represents the uuid command
var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate UUID (v4 or v7)",
	Long: `Generate UUID (Universally Unique Identifier) values.

Supported versions:
  - v4: Random UUID (default)
  - v7: Time-ordered UUID

Examples:
  devkit dev uuid                    # Generate UUID v4
  devkit dev uuid --version 7        # Generate UUID v7
  devkit dev uuid --count 5          # Generate 5 UUIDs
  devkit dev uuid --version 7 --count 3 --output json`,
	RunE: runUUID,
}

func init() {
	devCmd.AddCommand(uuidCmd)

	// Flag definitions
	uuidCmd.Flags().Int("version", 4, "UUID version (4 or 7)")
	uuidCmd.Flags().IntP("count", "c", 1, "Number of UUIDs to generate")
	uuidCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runUUID(cmd *cobra.Command, args []string) error {
	version, _ := cmd.Flags().GetInt("version")
	count, _ := cmd.Flags().GetInt("count")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if count < 1 {
		return fmt.Errorf("count must be at least 1")
	}

	if count > 1000 {
		return fmt.Errorf("count cannot exceed 1000")
	}

	var uuids []string
	var err error

	switch version {
	case 4:
		uuids, err = generateUUIDv4(count)
	case 7:
		uuids, err = generateUUIDv7(count)
	default:
		return fmt.Errorf("unsupported UUID version: %d (supported: 4, 7)", version)
	}

	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Prepare result based on format
	if format == output.FormatJSON {
		result := map[string]interface{}{
			"version": version,
			"count":   count,
			"uuids":   uuids,
		}
		output.PrintSuccess(format, result)
	} else {
		// Plain format - print each UUID on a new line
		if count == 1 {
			output.PrintSuccess(format, uuids[0])
		} else {
			output.PrintSuccess(format, uuids)
		}
	}

	return nil
}

func generateUUIDv4(count int) ([]string, error) {
	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		id := uuid.New()
		uuids[i] = id.String()
	}
	return uuids, nil
}

func generateUUIDv7(count int) ([]string, error) {
	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		uuids[i] = id.String()
	}
	return uuids, nil
}
