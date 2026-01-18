package net

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping [host]",
	Short: "Ping a host with statistics",
	Long: `Ping a host and display statistics.

Examples:
  devkit net ping google.com
  devkit net ping 8.8.8.8 --count 10`,
	Args: cobra.ExactArgs(1),
	RunE: runPing,
}

func init() {
	netCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntP("count", "c", 4, "Number of ping packets")
	pingCmd.Flags().IntP("timeout", "t", 3, "Timeout in seconds")
	pingCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runPing(cmd *cobra.Command, args []string) error {
	host := args[0]
	count, _ := cmd.Flags().GetInt("count")
	timeout, _ := cmd.Flags().GetInt("timeout")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	var times []time.Duration
	var successCount int

	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", host), time.Duration(timeout)*time.Second)
		duration := time.Since(start)

		if err == nil {
			conn.Close()
			times = append(times, duration)
			successCount++
		}
	}

	if len(times) == 0 {
		return fmt.Errorf("all ping attempts failed")
	}

	var total time.Duration
	var min, max time.Duration = times[0], times[0]
	for _, t := range times {
		total += t
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}

	avg := total / time.Duration(len(times))
	loss := float64(count-successCount) / float64(count) * 100

	result := map[string]interface{}{
		"host":        host,
		"sent":        count,
		"received":    successCount,
		"loss":        fmt.Sprintf("%.1f%%", loss),
		"min":         min.String(),
		"max":         max.String(),
		"avg":         avg.String(),
		"times":       times,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Ping statistics for %s:\n", host)
		fmt.Printf("  Packets: Sent = %d, Received = %d, Lost = %d (%.1f%% loss)\n",
			count, successCount, count-successCount, loss)
		fmt.Printf("  Times: Min = %s, Max = %s, Avg = %s\n", min, max, avg)
	}

	return nil
}
