package net

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// interfacesCmd represents the interfaces command
var interfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "Network interfaces and connections",
	Long: `List network interfaces and their configurations.

Examples:
  devkit net interfaces
  devkit net interfaces --output json`,
	RunE: runInterfaces,
}

func init() {
	netCmd.AddCommand(interfacesCmd)

	interfacesCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runInterfaces(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get interfaces: %w", err)
	}

	var ifaceList []map[string]interface{}

	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		var ipAddrs []string
		for _, addr := range addrs {
			ipAddrs = append(ipAddrs, addr.String())
		}

		ifaceInfo := map[string]interface{}{
			"name":    iface.Name,
			"index":   iface.Index,
			"mtu":     iface.MTU,
			"flags":   iface.Flags.String(),
			"addresses": ipAddrs,
		}
		ifaceList = append(ifaceList, ifaceInfo)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"interfaces": ifaceList,
			"count":     len(ifaceList),
		})
	} else {
		for _, iface := range ifaceList {
			fmt.Printf("Interface: %s\n", iface["name"])
			fmt.Printf("  Index: %d, MTU: %d\n", iface["index"], iface["mtu"])
			fmt.Printf("  Flags: %s\n", iface["flags"])
			if addrs, ok := iface["addresses"].([]string); ok && len(addrs) > 0 {
				fmt.Printf("  Addresses:\n")
				for _, addr := range addrs {
					fmt.Printf("    %s\n", addr)
				}
			}
			fmt.Println()
		}
	}

	return nil
}
