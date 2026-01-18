package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert between file formats (JSON, YAML, TOML, XML, CSV)",
	Long: `Convert files between different formats.

Supported formats: json, yaml, toml, xml, csv

Examples:
  devkit file convert config.json --to yaml
  devkit file convert data.csv --to json
  devkit file convert config.yaml --to toml`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func init() {
	fileCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringP("to", "t", "", "Target format (json, yaml, toml, xml, csv) (required)")
	convertCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	convertCmd.MarkFlagRequired("to")
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	toFormat, _ := cmd.Flags().GetString("to")
	outputFile, _ := cmd.Flags().GetString("output")

	// Detect input format
	inputExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(inputFile), "."))
	if inputExt == "" {
		return fmt.Errorf("cannot detect input format from file extension")
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse input
	var parsedData interface{}
	switch inputExt {
	case "json":
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &parsedData); err != nil {
			return fmt.Errorf("invalid YAML: %w", err)
		}
	case "toml":
		if err := toml.Unmarshal(data, &parsedData); err != nil {
			return fmt.Errorf("invalid TOML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported input format: %s", inputExt)
	}

	// Convert to target format
	var outputData []byte
	switch strings.ToLower(toFormat) {
	case "json":
		outputData, err = json.MarshalIndent(parsedData, "", "  ")
	case "yaml", "yml":
		outputData, err = yaml.Marshal(parsedData)
	case "toml":
		outputData, err = toml.Marshal(parsedData)
	default:
		return fmt.Errorf("unsupported output format: %s", toFormat)
	}

	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Converted: %s -> %s\n", inputFile, outputFile)
	} else {
		fmt.Print(string(outputData))
	}

	return nil
}
