package net

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// sslCmd represents the ssl command group
var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "SSL certificate operations",
	Long: `Check SSL certificate information and validity.

Examples:
  devkit net ssl check google.com
  devkit net ssl expiry google.com`,
}

// sslCheckCmd represents the check subcommand
var sslCheckCmd = &cobra.Command{
	Use:   "check [host]",
	Short: "Check SSL certificate",
	Long: `Check SSL certificate information for a host.

Examples:
  devkit net ssl check google.com
  devkit net ssl check example.com:443`,
	Args: cobra.ExactArgs(1),
	RunE: runSSLCheck,
}

// sslExpiryCmd represents the expiry subcommand
var sslExpiryCmd = &cobra.Command{
	Use:   "expiry [host]",
	Short: "Check SSL certificate expiry",
	Long: `Check when SSL certificate expires.

Examples:
  devkit net ssl expiry google.com`,
	Args: cobra.ExactArgs(1),
	RunE: runSSLExpiry,
}

func init() {
	netCmd.AddCommand(sslCmd)
	sslCmd.AddCommand(sslCheckCmd)
	sslCmd.AddCommand(sslExpiryCmd)

	sslCheckCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
	sslExpiryCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runSSLCheck(cmd *cobra.Command, args []string) error {
	host := args[0]
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return fmt.Errorf("no certificate found")
	}

	cert := state.PeerCertificates[0]

	result := map[string]interface{}{
		"host":            host,
		"subject":         cert.Subject.String(),
		"issuer":          cert.Issuer.String(),
		"valid_from":      cert.NotBefore.Format(time.RFC3339),
		"valid_to":        cert.NotAfter.Format(time.RFC3339),
		"is_valid":        time.Now().Before(cert.NotAfter) && time.Now().After(cert.NotBefore),
		"days_remaining":  int(time.Until(cert.NotAfter).Hours() / 24),
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("SSL Certificate for %s:\n", host)
		fmt.Printf("  Subject: %s\n", result["subject"])
		fmt.Printf("  Issuer: %s\n", result["issuer"])
		fmt.Printf("  Valid From: %s\n", result["valid_from"])
		fmt.Printf("  Valid To: %s\n", result["valid_to"])
		fmt.Printf("  Days Remaining: %d\n", result["days_remaining"])
	}

	return nil
}

func runSSLExpiry(cmd *cobra.Command, args []string) error {
	host := args[0]
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return fmt.Errorf("no certificate found")
	}

	cert := state.PeerCertificates[0]
	daysRemaining := int(time.Until(cert.NotAfter).Hours() / 24)
	isExpired := time.Now().After(cert.NotAfter)

	result := map[string]interface{}{
		"host":           host,
		"expires":        cert.NotAfter.Format(time.RFC3339),
		"days_remaining": daysRemaining,
		"is_expired":     isExpired,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		if isExpired {
			fmt.Printf("Certificate for %s EXPIRED on %s\n", host, result["expires"])
		} else {
			fmt.Printf("Certificate for %s expires in %d days (%s)\n", host, daysRemaining, result["expires"])
		}
	}

	return nil
}
