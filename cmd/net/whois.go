package net

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// whoisCmd represents the whois command
var whoisCmd = &cobra.Command{
	Use:   "whois [domain]",
	Short: "Domain whois lookup",
	Long: `Perform whois lookup for a domain.

Examples:
  devkit net whois example.com
  devkit net whois google.com`,
	Args: cobra.ExactArgs(1),
	RunE: runWhois,
}

func init() {
	netCmd.AddCommand(whoisCmd)

	whoisCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runWhois(cmd *cobra.Command, args []string) error {
	domain := args[0]
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	// Use whois server (simplified implementation)
	whoisServer := "whois.iana.org"
	conn, err := net.DialTimeout("tcp", whoisServer+":43", 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to whois server: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	fmt.Fprintf(conn, "%s\r\n", domain)

	var response strings.Builder
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}
		response.Write(buffer[:n])
	}

	whoisData := response.String()

	result := map[string]interface{}{
		"domain": domain,
		"data":   whoisData,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Print(whoisData)
	}

	return nil
}
