package dev

import (
	"fmt"
	"html"
	"io"
	"os"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// htmlCmd represents the html command group
var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "HTML entity encode/decode operations",
	Long: `Encode and decode HTML entities.

Examples:
  devkit dev html encode "hello <world>"
  devkit dev html decode "hello &lt;world&gt;"`,
}

// htmlEncodeCmd represents the encode subcommand
var htmlEncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "HTML entity encode a string",
	Long: `Encode special characters to HTML entities.

Examples:
  devkit dev html encode "hello <world>"
  devkit dev html encode --file input.txt`,
	RunE: runHTMLEncode,
}

// htmlDecodeCmd represents the decode subcommand
var htmlDecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "HTML entity decode a string",
	Long: `Decode HTML entities to special characters.

Examples:
  devkit dev html decode "hello &lt;world&gt;"
  devkit dev html decode --file encoded.txt`,
	RunE: runHTMLDecode,
}

func init() {
	devCmd.AddCommand(htmlCmd)
	htmlCmd.AddCommand(htmlEncodeCmd)
	htmlCmd.AddCommand(htmlDecodeCmd)

	// Flag definitions
	htmlEncodeCmd.Flags().StringP("file", "f", "", "Input file path")
	htmlEncodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	htmlEncodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	htmlDecodeCmd.Flags().StringP("file", "f", "", "Input file path")
	htmlDecodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	htmlDecodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runHTMLEncode(cmd *cobra.Command, args []string) error {
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var input string

	if stdinFlag {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("stdin error: %w", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin error: %w", err)
			}
			input = string(bytes)
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
		input = string(bytes)
	} else if len(args) > 0 {
		input = args[0]
	} else {
		return fmt.Errorf("input not specified")
	}

	encoded := html.EscapeString(input)

	if format == output.FormatJSON {
		result := map[string]interface{}{
			"encoded": encoded,
			"input":   input,
		}
		output.PrintSuccess(format, result)
	} else {
		output.PrintSuccess(format, encoded)
	}

	return nil
}

func runHTMLDecode(cmd *cobra.Command, args []string) error {
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var input string

	if stdinFlag {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("stdin error: %w", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin error: %w", err)
			}
			input = string(bytes)
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
		input = string(bytes)
	} else if len(args) > 0 {
		input = args[0]
	} else {
		return fmt.Errorf("input not specified")
	}

	decoded := html.UnescapeString(input)

	if format == output.FormatJSON {
		result := map[string]interface{}{
			"decoded": decoded,
			"input":   input,
		}
		output.PrintSuccess(format, result)
	} else {
		output.PrintSuccess(format, decoded)
	}

	return nil
}
