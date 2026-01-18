package dev

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// jwtCmd represents the jwt command group
var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT (JSON Web Token) operations",
	Long: `Decode and verify JWT tokens.

Examples:
  devkit dev jwt decode "eyJhbGciOiJIUzI1NiIs..."
  devkit dev jwt verify "eyJ..." --secret "my-secret-key"
  devkit dev jwt decode --file token.txt`,
}

// jwtDecodeCmd represents the decode subcommand
var jwtDecodeCmd = &cobra.Command{
	Use:   "decode [token]",
	Short: "Decode JWT token without verification",
	Long: `Decode a JWT token and display its header and payload.

This command does not verify the token signature, it only decodes it.

Examples:
  devkit dev jwt decode "eyJhbGciOiJIUzI1NiIs..."
  devkit dev jwt decode --file token.txt
  echo "eyJ..." | devkit dev jwt decode --stdin`,
	RunE: runJWTDecode,
}

// jwtVerifyCmd represents the verify subcommand
var jwtVerifyCmd = &cobra.Command{
	Use:   "verify [token]",
	Short: "Verify JWT token signature",
	Long: `Verify a JWT token's signature using a secret key.

Examples:
  devkit dev jwt verify "eyJ..." --secret "my-secret-key"
  devkit dev jwt verify --file token.txt --secret "my-secret-key"`,
	RunE: runJWTVerify,
}

func init() {
	devCmd.AddCommand(jwtCmd)
	jwtCmd.AddCommand(jwtDecodeCmd)
	jwtCmd.AddCommand(jwtVerifyCmd)

	// Flag definitions for decode
	jwtDecodeCmd.Flags().StringP("file", "f", "", "Input file path")
	jwtDecodeCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	jwtDecodeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")

	// Flag definitions for verify
	jwtVerifyCmd.Flags().StringP("file", "f", "", "Input file path")
	jwtVerifyCmd.Flags().BoolP("stdin", "s", false, "Read from stdin")
	jwtVerifyCmd.Flags().StringP("secret", "k", "", "Secret key for verification (required)")
	jwtVerifyCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
	jwtVerifyCmd.MarkFlagRequired("secret")
}

func runJWTDecode(cmd *cobra.Command, args []string) error {
	// Get input
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var tokenString string
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
			tokenString = strings.TrimSpace(string(bytes))
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
		tokenString = strings.TrimSpace(string(bytes))
	} else if len(args) > 0 {
		tokenString = args[0]
	} else {
		return fmt.Errorf("token not specified (use --file, --stdin, or provide as argument)")
	}

	// Parse token without verification
	// Split token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid token format: expected 3 parts separated by dots")
	}

	// Decode header (base64url)
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("failed to decode header: %w", err)
	}

	// Decode claims (base64url)
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode claims: %w", err)
	}

	var header map[string]interface{}
	var claims jwt.MapClaims

	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return fmt.Errorf("failed to parse header: %w", err)
	}

	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return fmt.Errorf("failed to parse claims: %w", err)
	}

	// Create a token object for compatibility
	token := &jwt.Token{
		Header: header,
		Claims: claims,
		Valid:  false, // Not verified
	}

	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to extract claims")
	}

	// Prepare result
	headerJSON, _ := json.MarshalIndent(header, "", "  ")
	claimsJSON, _ := json.MarshalIndent(claims, "", "  ")

	if format == output.FormatJSON {
		// For JSON output, parse the JSON strings back to objects
		var headerObj map[string]interface{}
		var claimsObj map[string]interface{}
		json.Unmarshal(headerJSON, &headerObj)
		json.Unmarshal(claimsJSON, &claimsObj)

		output.PrintSuccess(format, map[string]interface{}{
			"header":  headerObj,
			"claims":  claimsObj,
			"valid":   token.Valid,
			"expired": isExpired(claims),
		})
	} else {
		// Plain format
		fmt.Println("Header:")
		fmt.Println(string(headerJSON))
		fmt.Println("\nClaims:")
		fmt.Println(string(claimsJSON))
		if isExpired(claims) {
			fmt.Println("\n⚠ Token is expired")
		}
	}

	return nil
}

func runJWTVerify(cmd *cobra.Command, args []string) error {
	// Get input
	fileFlag, _ := cmd.Flags().GetString("file")
	stdinFlag, _ := cmd.Flags().GetBool("stdin")
	secret, _ := cmd.Flags().GetString("secret")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var tokenString string
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
			tokenString = strings.TrimSpace(string(bytes))
		} else {
			return fmt.Errorf("no data available from stdin")
		}
	} else if fileFlag != "" {
		bytes, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("read file error: %w", err)
		}
		tokenString = strings.TrimSpace(string(bytes))
	} else if len(args) > 0 {
		tokenString = args[0]
	} else {
		return fmt.Errorf("token not specified (use --file, --stdin, or provide as argument)")
	}

	if secret == "" {
		return fmt.Errorf("secret key is required (use --secret)")
	}

	// Parse and verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to extract claims")
	}

	// Prepare result
	result := map[string]interface{}{
		"valid":   token.Valid,
		"expired": isExpired(claims),
		"claims":  claims,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		if token.Valid {
			fmt.Println("✓ Token is valid")
			if isExpired(claims) {
				fmt.Println("⚠ Token is expired")
			}
		} else {
			fmt.Println("✗ Token is invalid")
		}
	}

	return nil
}

func isExpired(claims jwt.MapClaims) bool {
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		return time.Now().After(expTime)
	}
	return false
}
