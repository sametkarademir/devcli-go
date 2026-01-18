package dev

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"devkit/internal/output"
	"devkit/internal/utils"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash [algorithm] [input]",
	Short: "Calculate hash of input",
	Long: `Calculate cryptographic hash of a string or file.

Supported algorithms: md5, sha1, sha256, sha512

Examples:
  devkit dev hash sha256 "hello world"
  devkit dev hash md5 --file /path/to/file
  echo "hello" | devkit dev hash sha256 --stdin`,
	Args: cobra.MinimumNArgs(1),
	ValidArgs: []string{"md5", "sha1", "sha256", "sha512"},
	RunE: runHash,
}

func init() {
	devCmd.AddCommand(hashCmd)

	// Flag definitions
	hashCmd.Flags().StringP("file", "f", "", "Input file path")
	hashCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	hashCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runHash(cmd *cobra.Command, args []string) error {
	algorithm := args[0]
	
	// Get input
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	
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
	} else if len(args) > 1 {
		input = args[1]
	} else {
		return fmt.Errorf("input not specified (use --file, --stdin, or provide as argument)")
	}

	// Calculate hash
	var hash string
	switch algorithm {
	case "md5":
		hash, err = calculateMD5(input)
	case "sha1":
		hash, err = calculateSHA1(input)
	case "sha256":
		hash, err = calculateSHA256(input)
	case "sha512":
		hash, err = calculateSHA512(input)
	default:
		return fmt.Errorf("unsupported algorithm: %s (supported: md5, sha1, sha256, sha512)", algorithm)
	}

	if err != nil {
		return err
	}

	// Get output format
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	// Prepare result based on format
	if format == output.FormatJSON {
		result := map[string]interface{}{
			"algorithm": algorithm,
			"hash":      hash,
			"input":     utils.TrimSpace(input),
		}
		output.PrintSuccess(format, result)
	} else {
		// Plain format - just print the hash
		output.PrintSuccess(format, hash)
	}

	return nil
}

func calculateMD5(input string) (string, error) {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func calculateSHA1(input string) (string, error) {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func calculateSHA256(input string) (string, error) {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func calculateSHA512(input string) (string, error) {
	hash := sha512.Sum512([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}
