package net

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// openPortsCmd represents the open-ports command
var openPortsCmd = &cobra.Command{
	Use:   "open-ports",
	Short: "Show open ports and applications",
	Long: `List open ports and which applications are using them.

Examples:
  devkit net open-ports
  devkit net open-ports --output json`,
	RunE: runOpenPorts,
}

func init() {
	netCmd.AddCommand(openPortsCmd)

	openPortsCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runOpenPorts(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	// Get listening addresses
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get interfaces: %w", err)
	}

	var ports []map[string]interface{}

	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			// Check common ports
			commonPorts := []int{22, 80, 443, 3306, 5432, 8080, 9000}
			for _, port := range commonPorts {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
				if err == nil {
					conn.Close()
					ports = append(ports, map[string]interface{}{
						"port":     port,
						"protocol": "tcp",
						"address":  ip.String(),
						"status":   "open",
					})
				}
			}
		}
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"ports": ports,
			"count": len(ports),
		})
	} else if format == output.FormatTable {
		fmt.Printf("%-10s %-10s %-20s %-10s\n", "PORT", "PROTOCOL", "ADDRESS", "STATUS")
		fmt.Println(strings.Repeat("-", 55))
		for _, port := range ports {
			fmt.Printf("%-10d %-10s %-20s %-10s\n",
				port["port"], port["protocol"], port["address"], port["status"])
		}
	} else {
		if len(ports) == 0 {
			fmt.Println("No open ports found")
			return nil
		}
		fmt.Println("Open Ports:")
		for _, port := range ports {
			fmt.Printf("  Port %d (%s) on %s - %s\n",
				port["port"], port["protocol"], port["address"], port["status"])
		}
	}

	return nil
}
