package dev

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// randomCmd represents the random command group
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate random data (string, number, password)",
	Long: `Generate random strings, numbers, and passwords.

Examples:
  devkit dev random string --length 32
  devkit dev random number --min 1 --max 100
  devkit dev random password --length 16`,
}

// randomStringCmd represents the string subcommand
var randomStringCmd = &cobra.Command{
	Use:   "string",
	Short: "Generate random string",
	Long: `Generate a random string of specified length.

Examples:
  devkit dev random string --length 32
  devkit dev random string --length 16 --charset "abc123"`,
	RunE: runRandomString,
}

// randomNumberCmd represents the number subcommand
var randomNumberCmd = &cobra.Command{
	Use:   "number",
	Short: "Generate random number",
	Long: `Generate a random number between min and max.

Examples:
  devkit dev random number --min 1 --max 100
  devkit dev random number --min 0 --max 1000`,
	RunE: runRandomNumber,
}

// randomPasswordCmd represents the password subcommand
var randomPasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "Generate random password",
	Long: `Generate a secure random password.

Examples:
  devkit dev random password --length 16
  devkit dev random password --length 20 --symbols`,
	RunE: runRandomPassword,
}

func init() {
	devCmd.AddCommand(randomCmd)
	randomCmd.AddCommand(randomStringCmd)
	randomCmd.AddCommand(randomNumberCmd)
	randomCmd.AddCommand(randomPasswordCmd)

	// Flag definitions
	randomStringCmd.Flags().IntP("length", "l", 16, "Length of the string")
	randomStringCmd.Flags().String("charset", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "Character set to use")
	randomStringCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	randomNumberCmd.Flags().IntP("min", "m", 0, "Minimum value")
	randomNumberCmd.Flags().IntP("max", "x", 100, "Maximum value")
	randomNumberCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	randomPasswordCmd.Flags().IntP("length", "l", 16, "Length of the password")
	randomPasswordCmd.Flags().BoolP("symbols", "s", false, "Include symbols")
	randomPasswordCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runRandomString(cmd *cobra.Command, args []string) error {
	length, _ := cmd.Flags().GetInt("length")
	charset, _ := cmd.Flags().GetString("charset")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if length < 1 {
		return fmt.Errorf("length must be at least 1")
	}

	if len(charset) == 0 {
		return fmt.Errorf("charset cannot be empty")
	}

	result := generateRandomString(length, charset)

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"string":  result,
			"length":  length,
			"charset": charset,
		})
	} else {
		output.PrintSuccess(format, result)
	}

	return nil
}

func runRandomNumber(cmd *cobra.Command, args []string) error {
	min, _ := cmd.Flags().GetInt("min")
	max, _ := cmd.Flags().GetInt("max")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if min >= max {
		return fmt.Errorf("min must be less than max")
	}

	result := generateRandomNumber(min, max)

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"number": result,
			"min":    min,
			"max":    max,
		})
	} else {
		output.PrintSuccess(format, result)
	}

	return nil
}

func runRandomPassword(cmd *cobra.Command, args []string) error {
	length, _ := cmd.Flags().GetInt("length")
	symbols, _ := cmd.Flags().GetBool("symbols")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if length < 1 {
		return fmt.Errorf("length must be at least 1")
	}

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if symbols {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	result := generateRandomString(length, charset)

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"password": result,
			"length":  length,
			"symbols": symbols,
		})
	} else {
		output.PrintSuccess(format, result)
	}

	return nil
}

func generateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func generateRandomNumber(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}
