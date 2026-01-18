package dev

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"devkit/internal/output"
)

// jsonCmd represents the json command group
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON operations (prettify, minify, validate, query)",
	Long: `JSON manipulation operations.

Examples:
  devkit dev json prettify '{"a":1,"b":2}'
  devkit dev json minify --file data.json
  devkit dev json validate --file data.json
  devkit dev json path '$.users[0].name' --file data.json`,
}

// jsonPrettifyCmd represents the prettify subcommand
var jsonPrettifyCmd = &cobra.Command{
	Use:   "prettify [json]",
	Short: "Prettify JSON string",
	Long: `Format JSON string with indentation.

Examples:
  devkit dev json prettify '{"a":1,"b":2}'
  devkit dev json prettify --file data.json`,
	RunE: runJSONPrettify,
}

// jsonMinifyCmd represents the minify subcommand
var jsonMinifyCmd = &cobra.Command{
	Use:   "minify [json]",
	Short: "Minify JSON string",
	Long: `Remove whitespace from JSON string.

Examples:
  devkit dev json minify '{"a": 1, "b": 2}'
  devkit dev json minify --file data.json`,
	RunE: runJSONMinify,
}

// jsonValidateCmd represents the validate subcommand
var jsonValidateCmd = &cobra.Command{
	Use:   "validate [json]",
	Short: "Validate JSON string",
	Long: `Check if a string is valid JSON.

Examples:
  devkit dev json validate '{"a":1}'
  devkit dev json validate --file data.json`,
	RunE: runJSONValidate,
}

// jsonPathCmd represents the path query subcommand
var jsonPathCmd = &cobra.Command{
	Use:   "path [query]",
	Short: "Query JSON using JSONPath",
	Long: `Query JSON data using JSONPath expression.

Examples:
  devkit dev json path '$.users[0].name' --file data.json
  devkit dev json path '$.items[*].id' --file data.json`,
	RunE: runJSONPath,
}

func init() {
	devCmd.AddCommand(jsonCmd)
	jsonCmd.AddCommand(jsonPrettifyCmd)
	jsonCmd.AddCommand(jsonMinifyCmd)
	jsonCmd.AddCommand(jsonValidateCmd)
	jsonCmd.AddCommand(jsonPathCmd)

	// Flag definitions
	jsonPrettifyCmd.Flags().StringP("file", "f", "", "Input file path")
	jsonPrettifyCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	jsonPrettifyCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	jsonMinifyCmd.Flags().StringP("file", "f", "", "Input file path")
	jsonMinifyCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	jsonMinifyCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	jsonValidateCmd.Flags().StringP("file", "f", "", "Input file path")
	jsonValidateCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	jsonValidateCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	jsonPathCmd.Flags().StringP("file", "f", "", "Input file path")
	jsonPathCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	jsonPathCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func getJSONInput(cmd *cobra.Command, args []string) (string, error) {
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")

	if stdinFlag {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return "", fmt.Errorf("stdin error: %w", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", fmt.Errorf("read stdin error: %w", err)
			}
			return string(bytes), nil
		} else {
			return "", fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return "", fmt.Errorf("read file error: %w", err)
		}
		return string(bytes), nil
	} else if len(args) > 0 {
		return args[0], nil
	}
	return "", fmt.Errorf("input not specified")
}

func runJSONPrettify(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	jsonInput, err := getJSONInput(cmd, args)
	if err != nil {
		return err
	}

	// Parse and prettify
	var data interface{}
	if err := json.Unmarshal([]byte(jsonInput), &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	prettified, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to prettify: %w", err)
	}

	result := string(pretty.Pretty(prettified))

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"prettified": result,
		})
	} else {
		output.PrintSuccess(format, result)
	}

	return nil
}

func runJSONMinify(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	jsonInput, err := getJSONInput(cmd, args)
	if err != nil {
		return err
	}

	// Parse and minify
	var data interface{}
	if err := json.Unmarshal([]byte(jsonInput), &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	minified, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to minify: %w", err)
	}

	result := string(pretty.Ugly(minified))

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"minified": result,
		})
	} else {
		output.PrintSuccess(format, result)
	}

	return nil
}

func runJSONValidate(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	jsonInput, err := getJSONInput(cmd, args)
	if err != nil {
		return err
	}

	// Validate JSON
	var data interface{}
	isValid := json.Unmarshal([]byte(jsonInput), &data) == nil

	if format == output.FormatJSON {
		result := map[string]interface{}{
			"valid": isValid,
		}
		if !isValid {
			result["error"] = "Invalid JSON format"
		}
		output.PrintSuccess(format, result)
	} else {
		if isValid {
			output.PrintSuccess(format, "✓ Valid JSON")
		} else {
			output.PrintError(format, fmt.Errorf("✗ Invalid JSON"))
		}
	}

	return nil
}

func runJSONPath(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("JSONPath query not specified")
	}

	query := args[0]
	jsonInput, err := getJSONInput(cmd, args[1:])
	if err != nil {
		return err
	}

	// Query JSON using gjson
	result := gjson.Get(jsonInput, query)

	if !result.Exists() {
		return fmt.Errorf("path not found: %s", query)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"query": query,
			"result": result.Value(),
		})
	} else {
		if result.IsArray() || result.IsObject() {
			// Pretty print for complex types
			prettyJSON, _ := json.MarshalIndent(result.Value(), "", "  ")
			output.PrintSuccess(format, string(prettyJSON))
		} else {
			output.PrintSuccess(format, result.String())
		}
	}

	return nil
}
