package net

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// sysinfoCmd represents the sysinfo command
var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "System information",
	Long: `Display system information (CPU, RAM, disk, OS).

Examples:
  devkit net sysinfo
  devkit net sysinfo --cpu
  devkit net sysinfo --memory`,
	RunE: runSysinfo,
}

func init() {
	netCmd.AddCommand(sysinfoCmd)

	sysinfoCmd.Flags().Bool("cpu", false, "Show CPU information only")
	sysinfoCmd.Flags().Bool("memory", false, "Show memory information only")
	sysinfoCmd.Flags().Bool("disk", false, "Show disk information only")
	sysinfoCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runSysinfo(cmd *cobra.Command, args []string) error {
	cpuOnly, _ := cmd.Flags().GetBool("cpu")
	memOnly, _ := cmd.Flags().GetBool("memory")
	diskOnly, _ := cmd.Flags().GetBool("disk")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	result := make(map[string]interface{})

	if !memOnly && !diskOnly {
		// CPU Info
		cpuInfo, _ := cpu.Info()
		cpuPercent, _ := cpu.Percent(time.Second, false)
		result["cpu"] = map[string]interface{}{
			"cores":     runtime.NumCPU(),
			"model":     getCPUModel(cpuInfo),
			"usage":     fmt.Sprintf("%.1f%%", cpuPercent[0]),
		}
	}

	if !cpuOnly && !diskOnly {
		// Memory Info
		memInfo, _ := mem.VirtualMemory()
		result["memory"] = map[string]interface{}{
			"total":     formatBytes(memInfo.Total),
			"used":      formatBytes(memInfo.Used),
			"available": formatBytes(memInfo.Available),
			"percent":   fmt.Sprintf("%.1f%%", memInfo.UsedPercent),
		}
	}

	if !cpuOnly && !memOnly {
		// Disk Info
		diskInfo, _ := disk.Usage("/")
		result["disk"] = map[string]interface{}{
			"total":   formatBytes(diskInfo.Total),
			"used":    formatBytes(diskInfo.Used),
			"free":    formatBytes(diskInfo.Free),
			"percent": fmt.Sprintf("%.1f%%", diskInfo.UsedPercent),
		}
	}

	// OS Info
	hostInfo, _ := host.Info()
	result["os"] = map[string]interface{}{
		"platform": hostInfo.Platform,
		"family":   hostInfo.PlatformFamily,
		"version":  hostInfo.PlatformVersion,
		"hostname": hostInfo.Hostname,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		if cpu, ok := result["cpu"].(map[string]interface{}); ok {
			fmt.Printf("CPU: %s (%d cores) - Usage: %s\n", cpu["model"], cpu["cores"], cpu["usage"])
		}
		if mem, ok := result["memory"].(map[string]interface{}); ok {
			fmt.Printf("Memory: %s / %s (%s used)\n", mem["used"], mem["total"], mem["percent"])
		}
		if disk, ok := result["disk"].(map[string]interface{}); ok {
			fmt.Printf("Disk: %s / %s (%s used)\n", disk["used"], disk["total"], disk["percent"])
		}
		if os, ok := result["os"].(map[string]interface{}); ok {
			fmt.Printf("OS: %s %s (%s)\n", os["platform"], os["version"], os["hostname"])
		}
	}

	return nil
}

func getCPUModel(cpuInfo []cpu.InfoStat) string {
	if len(cpuInfo) > 0 {
		return cpuInfo[0].ModelName
	}
	return "Unknown"
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
