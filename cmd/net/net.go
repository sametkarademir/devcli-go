package net

import (
	"github.com/spf13/cobra"
)

// netCmd represents the net command group
var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Network and system operations",
	Long: `Network and system operations for network analysis and system monitoring.

This command group includes utilities for:
- Port scanning and status checking
- DNS lookups
- IP information and geolocation
- HTTP requests
- Ping with statistics
- SSL certificate information
- Whois queries
- Internet speed testing
- System information
- Process management
- Disk usage analysis
- Network interfaces
- Open ports monitoring`,
}

// GetNetCmd returns the net command
func GetNetCmd() *cobra.Command {
	return netCmd
}

func init() {
	// This will be called when the package is imported
	// Commands will be added in their respective files
}
