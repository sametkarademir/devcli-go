package net

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// portCmd represents the port command group
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Port scanning and status checking",
	Long: `Scan ports and check port status.

Examples:
  devkit net port check 8080
  devkit net port scan localhost --range 1-1000
  devkit net port list --listening`,
}

// portCheckCmd represents the check subcommand
var portCheckCmd = &cobra.Command{
	Use:   "check [port]",
	Short: "Check if a port is open",
	Long: `Check if a specific port is open on a host.

Examples:
  devkit net port check 8080
  devkit net port check 8080 --host localhost`,
	Args: cobra.ExactArgs(1),
	RunE: runPortCheck,
}

// portScanCmd represents the scan subcommand
var portScanCmd = &cobra.Command{
	Use:   "scan [host]",
	Short: "Scan a range of ports",
	Long: `Scan a range of ports on a host.

Examples:
  devkit net port scan localhost --range 1-1000
  devkit net port scan 192.168.1.1 --range 80-443`,
	Args: cobra.ExactArgs(1),
	RunE: runPortScan,
}

// portListCmd represents the list subcommand
var portListCmd = &cobra.Command{
	Use:   "list",
	Short: "List listening ports",
	Long: `List all listening ports on the local system.

Examples:
  devkit net port list
  devkit net port list --output json`,
	RunE: runPortList,
}

func init() {
	netCmd.AddCommand(portCmd)
	portCmd.AddCommand(portCheckCmd)
	portCmd.AddCommand(portScanCmd)
	portCmd.AddCommand(portListCmd)

	portCheckCmd.Flags().String("host", "localhost", "Host to check")
	portCheckCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	portScanCmd.Flags().StringP("range", "r", "1-1000", "Port range to scan (e.g., 1-1000)")
	portScanCmd.Flags().IntP("timeout", "t", 1, "Timeout in seconds")
	portScanCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	portListCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runPortCheck(cmd *cobra.Command, args []string) error {
	portStr := args[0]
	host, _ := cmd.Flags().GetString("host")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port: %s", portStr)
	}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	
	isOpen := err == nil
	if conn != nil {
		conn.Close()
	}

	result := map[string]interface{}{
		"host":   host,
		"port":   port,
		"status": "open",
	}

	if !isOpen {
		result["status"] = "closed"
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		if isOpen {
			fmt.Printf("Port %d on %s is OPEN\n", port, host)
		} else {
			fmt.Printf("Port %d on %s is CLOSED\n", port, host)
		}
	}

	return nil
}

func runPortScan(cmd *cobra.Command, args []string) error {
	host := args[0]
	rangeStr, _ := cmd.Flags().GetString("range")
	timeout, _ := cmd.Flags().GetInt("timeout")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format: %s (expected: start-end)", rangeStr)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid start port: %s", parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid end port: %s", parts[1])
	}

	var openPorts []int
	for port := start; port <= end; port++ {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
		if err == nil {
			openPorts = append(openPorts, port)
			conn.Close()
		}
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"host":      host,
			"range":     rangeStr,
			"open_ports": openPorts,
			"count":     len(openPorts),
		})
	} else {
		if len(openPorts) == 0 {
			fmt.Printf("No open ports found in range %s on %s\n", rangeStr, host)
		} else {
			fmt.Printf("Open ports on %s:\n", host)
			for _, port := range openPorts {
				fmt.Printf("  %d\n", port)
			}
			fmt.Printf("\nTotal: %d open ports\n", len(openPorts))
		}
	}

	return nil
}

func runPortList(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	// List listening ports (simplified - would need platform-specific code for full implementation)
	result := map[string]interface{}{
		"message": "Port listing requires platform-specific implementation",
		"note":    "Use 'netstat' or 'lsof' for detailed port information",
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Println("Note: Full port listing requires platform-specific implementation.")
		fmt.Println("Use system commands like 'netstat -an' or 'lsof -i' for detailed information.")
	}

	return nil
}
