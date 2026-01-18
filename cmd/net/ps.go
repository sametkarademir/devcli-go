package net

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Process list and management",
	Long: `List and manage running processes.

Examples:
  devkit net ps
  devkit net ps --sort cpu
  devkit net ps --filter "go"`,
	RunE: runPS,
}

func init() {
	netCmd.AddCommand(psCmd)

	psCmd.Flags().StringP("sort", "s", "cpu", "Sort by: cpu, mem, pid")
	psCmd.Flags().StringP("filter", "f", "", "Filter processes by name")
	psCmd.Flags().IntP("limit", "n", 20, "Limit number of processes")
	psCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json, table")
}

func runPS(cmd *cobra.Command, args []string) error {
	sortBy, _ := cmd.Flags().GetString("sort")
	filter, _ := cmd.Flags().GetString("filter")
	limit, _ := cmd.Flags().GetInt("limit")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var procList []map[string]interface{}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if filter != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(filter)) {
			continue
		}

		cpuPercent, _ := p.CPUPercent()
		memInfo, err := p.MemoryInfo()
		if err != nil || memInfo == nil {
			continue
		}
		pid := p.Pid

		proc := map[string]interface{}{
			"pid":        pid,
			"name":       name,
			"cpu":        fmt.Sprintf("%.1f%%", cpuPercent),
			"memory":     formatBytesPS(memInfo.RSS),
			"cpu_value":  cpuPercent,
			"mem_value":  memInfo.RSS,
		}
		procList = append(procList, proc)

		if len(procList) >= limit {
			break
		}
	}

	// Sort
	switch sortBy {
	case "cpu":
		// Sort by CPU (already sorted by gopsutil)
	case "mem":
		// Sort by memory
		for i := 0; i < len(procList); i++ {
			for j := i + 1; j < len(procList); j++ {
				if procList[i]["mem_value"].(uint64) < procList[j]["mem_value"].(uint64) {
					procList[i], procList[j] = procList[j], procList[i]
				}
			}
		}
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"processes": procList,
			"count":    len(procList),
		})
	} else if format == output.FormatTable {
		fmt.Printf("%-8s %-30s %10s %12s\n", "PID", "NAME", "CPU", "MEMORY")
		fmt.Println(strings.Repeat("-", 65))
		for _, proc := range procList {
			fmt.Printf("%-8d %-30s %10s %12s\n",
				proc["pid"], truncate(proc["name"].(string), 30), proc["cpu"], proc["memory"])
		}
	} else {
		for _, proc := range procList {
			fmt.Printf("PID: %d, Name: %s, CPU: %s, Memory: %s\n",
				proc["pid"], proc["name"], proc["cpu"], proc["memory"])
		}
	}

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatBytesPS(bytes uint64) string {
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
