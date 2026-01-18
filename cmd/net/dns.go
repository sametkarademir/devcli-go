package net

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// dnsCmd represents the dns command group
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "DNS lookup operations",
	Long: `Perform DNS lookups for various record types.

Examples:
  devkit net dns lookup google.com
  devkit net dns lookup google.com --type MX
  devkit net dns reverse 8.8.8.8`,
}

// dnsLookupCmd represents the lookup subcommand
var dnsLookupCmd = &cobra.Command{
	Use:   "lookup [domain]",
	Short: "Lookup DNS records",
	Long: `Lookup DNS records for a domain.

Record types: A, AAAA, MX, TXT, NS, CNAME

Examples:
  devkit net dns lookup google.com
  devkit net dns lookup google.com --type MX
  devkit net dns lookup google.com --type TXT`,
	Args: cobra.ExactArgs(1),
	RunE: runDNSLookup,
}

// dnsReverseCmd represents the reverse subcommand
var dnsReverseCmd = &cobra.Command{
	Use:   "reverse [ip]",
	Short: "Reverse DNS lookup",
	Long: `Perform reverse DNS lookup (PTR record).

Examples:
  devkit net dns reverse 8.8.8.8
  devkit net dns reverse 2001:4860:4860::8888`,
	Args: cobra.ExactArgs(1),
	RunE: runDNSReverse,
}

func init() {
	netCmd.AddCommand(dnsCmd)
	dnsCmd.AddCommand(dnsLookupCmd)
	dnsCmd.AddCommand(dnsReverseCmd)

	dnsLookupCmd.Flags().StringP("type", "t", "A", "DNS record type (A, AAAA, MX, TXT, NS, CNAME)")
	dnsLookupCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	dnsReverseCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runDNSLookup(cmd *cobra.Command, args []string) error {
	domain := args[0]
	recordType, _ := cmd.Flags().GetString("type")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var result map[string]interface{}
	var values []string

	switch strings.ToUpper(recordType) {
	case "A":
		ips, err := net.LookupIP(domain)
		if err != nil {
			return fmt.Errorf("DNS lookup failed: %w", err)
		}
		for _, ip := range ips {
			if ip.To4() != nil {
				values = append(values, ip.String())
			}
		}
	case "AAAA":
		ips, err := net.LookupIP(domain)
		if err != nil {
			return fmt.Errorf("DNS lookup failed: %w", err)
		}
		for _, ip := range ips {
			if ip.To4() == nil {
				values = append(values, ip.String())
			}
		}
	case "MX":
		mxRecords, err := net.LookupMX(domain)
		if err != nil {
			return fmt.Errorf("MX lookup failed: %w", err)
		}
		for _, mx := range mxRecords {
			values = append(values, fmt.Sprintf("%s (priority: %d)", mx.Host, mx.Pref))
		}
	case "TXT":
		txtRecords, err := net.LookupTXT(domain)
		if err != nil {
			return fmt.Errorf("TXT lookup failed: %w", err)
		}
		values = txtRecords
	case "NS":
		nsRecords, err := net.LookupNS(domain)
		if err != nil {
			return fmt.Errorf("NS lookup failed: %w", err)
		}
		for _, ns := range nsRecords {
			values = append(values, ns.Host)
		}
	case "CNAME":
		cname, err := net.LookupCNAME(domain)
		if err != nil {
			return fmt.Errorf("CNAME lookup failed: %w", err)
		}
		values = []string{cname}
	default:
		return fmt.Errorf("unsupported record type: %s (supported: A, AAAA, MX, TXT, NS, CNAME)", recordType)
	}

	result = map[string]interface{}{
		"domain": domain,
		"type":   recordType,
		"records": values,
		"count":  len(values),
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("DNS %s records for %s:\n", recordType, domain)
		for _, value := range values {
			fmt.Printf("  %s\n", value)
		}
	}

	return nil
}

func runDNSReverse(cmd *cobra.Command, args []string) error {
	ipStr := args[0]
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	names, err := net.LookupAddr(ipStr)
	if err != nil {
		return fmt.Errorf("reverse DNS lookup failed: %w", err)
	}

	result := map[string]interface{}{
		"ip":     ipStr,
		"names":  names,
		"count":  len(names),
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Reverse DNS for %s:\n", ipStr)
		for _, name := range names {
			fmt.Printf("  %s\n", name)
		}
	}

	return nil
}
