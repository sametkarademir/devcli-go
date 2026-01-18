package dev

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// base64Cmd represents the base64 command group
var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Base64 encode/decode operations (use 'base64 encode' or 'base64 decode')",
	Long: `Encode or decode base64 strings.

Subcommands:
  encode    Encode input to base64
  decode    Decode base64 string

Examples:
  devkit dev base64 encode "hello world"
  devkit dev base64 decode "aGVsbG8gd29ybGQ="
  devkit dev base64 encode --file ./image.png
  echo "test" | devkit dev base64 encode --stdin`,
}

// encodeCmd represents the encode subcommand
var encodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "Encode input to base64",
	Long: `Encode a string or file to base64.

Examples:
  devkit dev base64 encode "hello world"
  devkit dev base64 encode --file ./image.png
  echo "test" | devkit dev base64 encode --stdin`,
	RunE: runEncode,
}

// decodeCmd represents the decode subcommand
var decodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "Decode base64 string",
	Long: `Decode a base64 string.

Examples:
  devkit dev base64 decode "aGVsbG8gd29ybGQ="
  devkit dev base64 decode --file encoded.txt
  echo "aGVsbG8gd29ybGQ=" | devkit dev base64 decode --stdin`,
	RunE: runDecode,
}

func init() {
	devCmd.AddCommand(base64Cmd)
	base64Cmd.AddCommand(encodeCmd)
	base64Cmd.AddCommand(decodeCmd)

	// Flag definitions for encode
	encodeCmd.Flags().StringP("file", "f", "", "Input file path")
	encodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	encodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")

	// Flag definitions for decode
	decodeCmd.Flags().StringP("file", "f", "", "Input file path")
	decodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	decodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runEncode(cmd *cobra.Command, args []string) error {
	// Get input
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var input []byte
	var err error

	if stdinFlag {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("stdin error: %w", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin error: %w", err)
			}
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		input, err = os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
	} else if len(args) > 0 {
		input = []byte(args[0])
	} else {
		return fmt.Errorf("input not specified (use --file, --stdin, or provide as argument)")
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(input)

	// Prepare result based on format
	if format == output.FormatJSON {
		result := map[string]interface{}{
			"encoded": encoded,
			"input":   string(input),
		}
		output.PrintSuccess(format, result)
	} else {
		// Plain format - just print the encoded string
		output.PrintSuccess(format, encoded)
	}

	return nil
}

func runDecode(cmd *cobra.Command, args []string) error {
	// Get input
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var input string
	var err error

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
		return fmt.Errorf("input not specified (use --file, --stdin, or provide as argument)")
	}

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return fmt.Errorf("invalid base64 string: %w", err)
	}

	// Prepare result based on format
	if format == output.FormatJSON {
		result := map[string]interface{}{
			"decoded": string(decoded),
			"input":   input,
		}
		output.PrintSuccess(format, result)
	} else {
		// Plain format - just print the decoded string
		output.PrintSuccess(format, string(decoded))
	}

	return nil
}
