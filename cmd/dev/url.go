package dev

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// urlCmd represents the url command group
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "URL encode/decode and parse operations",
	Long: `Encode, decode, and parse URLs.

Examples:
  devkit dev url encode "hello world"
  devkit dev url decode "hello%20world"
  devkit dev url parse "https://example.com/path?key=value"`,
}

// urlEncodeCmd represents the encode subcommand
var urlEncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "URL encode a string",
	Long: `URL encode a string.

Examples:
  devkit dev url encode "hello world"
  devkit dev url encode --file input.txt
  echo "test" | devkit dev url encode --stdin`,
	RunE: runURLEncode,
}

// urlDecodeCmd represents the decode subcommand
var urlDecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "URL decode a string",
	Long: `URL decode a string.

Examples:
  devkit dev url decode "hello%20world"
  devkit dev url decode --file encoded.txt`,
	RunE: runURLDecode,
}

// urlParseCmd represents the parse subcommand
var urlParseCmd = &cobra.Command{
	Use:   "parse [url]",
	Short: "Parse a URL and display its components",
	Long: `Parse a URL and display its components (scheme, host, path, query, etc.).

Examples:
  devkit dev url parse "https://example.com/path?key=value"
  devkit dev url parse "https://user:pass@example.com:8080/path?key=value#fragment"`,
	RunE: runURLParse,
}

func init() {
	devCmd.AddCommand(urlCmd)
	urlCmd.AddCommand(urlEncodeCmd)
	urlCmd.AddCommand(urlDecodeCmd)
	urlCmd.AddCommand(urlParseCmd)

	// Flag definitions
	urlEncodeCmd.Flags().StringP("file", "f", "", "Input file path")
	urlEncodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	urlEncodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	urlDecodeCmd.Flags().StringP("file", "f", "", "Input file path")
	urlDecodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	urlDecodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	urlParseCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runURLEncode(cmd *cobra.Command, args []string) error {
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
			input = strings.TrimSpace(string(bytes))
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
		input = strings.TrimSpace(string(bytes))
	} else if len(args) > 0 {
		input = args[0]
	} else {
		return fmt.Errorf("input not specified")
	}

	encoded := url.QueryEscape(input)

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

func runURLDecode(cmd *cobra.Command, args []string) error {
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
			input = strings.TrimSpace(string(bytes))
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
		input = strings.TrimSpace(string(bytes))
	} else if len(args) > 0 {
		input = args[0]
	} else {
		return fmt.Errorf("input not specified")
	}

	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

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

func runURLParse(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("URL not specified")
	}

	u, err := url.Parse(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	result := map[string]interface{}{
		"scheme":   u.Scheme,
		"host":     u.Host,
		"path":     u.Path,
		"query":    u.RawQuery,
		"fragment": u.Fragment,
		"user":     u.User.String(),
	}

	// Parse query parameters
	if u.RawQuery != "" {
		queryParams := make(map[string]string)
		for key, values := range u.Query() {
			if len(values) > 0 {
				queryParams[key] = values[0]
			}
		}
		result["query_params"] = queryParams
	}

	output.PrintSuccess(format, result)
	return nil
}
