package net

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP address information",
	Long: `Get IP address information (public/private) and geolocation.

Examples:
  devkit net ip                    # Public IP
  devkit net ip --local            # Local IP
  devkit net ip info 8.8.8.8       # IP information`,
	RunE: runIP,
}

func init() {
	netCmd.AddCommand(ipCmd)

	ipCmd.Flags().BoolP("local", "l", false, "Show local IP address")
	ipCmd.Flags().StringP("info", "i", "", "Get information about an IP address")
	ipCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runIP(cmd *cobra.Command, args []string) error {
	local, _ := cmd.Flags().GetBool("local")
	infoIP, _ := cmd.Flags().GetString("info")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if infoIP != "" {
		return showIPInfo(infoIP, format)
	}

	if local {
		return showLocalIP(format)
	}

	return showPublicIP(format)
}

func showPublicIP(format output.OutputFormat) error {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return fmt.Errorf("failed to get public IP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	ip := strings.TrimSpace(string(body))

	result := map[string]interface{}{
		"type": "public",
		"ip":   ip,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Public IP: %s\n", ip)
	}

	return nil
}

func showLocalIP(format output.OutputFormat) error {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return fmt.Errorf("failed to get local IP: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	result := map[string]interface{}{
		"type": "local",
		"ip":   ip,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Local IP: %s\n", ip)
	}

	return nil
}

func showIPInfo(ipStr string, format output.OutputFormat) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	result := map[string]interface{}{
		"ip":      ipStr,
		"version": "IPv4",
		"private": isPrivateIP(ip),
	}

	if ip.To4() == nil {
		result["version"] = "IPv6"
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("IP: %s\n", ipStr)
		fmt.Printf("Version: %s\n", result["version"])
		fmt.Printf("Private: %v\n", result["private"])
	}

	return nil
}

func isPrivateIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsPrivate()
}
