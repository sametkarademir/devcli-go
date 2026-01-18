package dev

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// cronCmd represents the cron command group
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Cron expression operations",
	Long: `Parse and explain cron expressions.

Examples:
  devkit dev cron explain "0 9 * * 1-5"
  devkit dev cron next "0 9 * * 1-5" --count 5`,
}

// cronExplainCmd represents the explain subcommand
var cronExplainCmd = &cobra.Command{
	Use:   "explain [expression]",
	Short: "Explain a cron expression",
	Long: `Explain what a cron expression means in human-readable format.

Examples:
  devkit dev cron explain "0 9 * * 1-5"
  devkit dev cron explain "*/5 * * * *"`,
	RunE: runCronExplain,
}

// cronNextCmd represents the next subcommand
var cronNextCmd = &cobra.Command{
	Use:   "next [expression]",
	Short: "Show next execution times",
	Long: `Show the next execution times for a cron expression.

Examples:
  devkit dev cron next "0 9 * * 1-5"
  devkit dev cron next "*/5 * * * *" --count 10`,
	RunE: runCronNext,
}

func init() {
	devCmd.AddCommand(cronCmd)
	cronCmd.AddCommand(cronExplainCmd)
	cronCmd.AddCommand(cronNextCmd)

	cronExplainCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
	cronNextCmd.Flags().IntP("count", "c", 5, "Number of next executions to show")
	cronNextCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runCronExplain(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("cron expression not specified")
	}

	expr := args[0]
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	schedule, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// Parse parts
	parts := parseCronExpression(expr)
	explanation := explainCronExpression(parts)

	result := map[string]interface{}{
		"expression":  expr,
		"explanation": explanation,
		"valid":       true,
		"next_run":    schedule.Next(time.Now()).Format(time.RFC3339),
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Expression: %s\n", expr)
		fmt.Printf("Explanation: %s\n", explanation)
		fmt.Printf("Next run: %s\n", result["next_run"])
	}

	return nil
}

func runCronNext(cmd *cobra.Command, args []string) error {
	count, _ := cmd.Flags().GetInt("count")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("cron expression not specified")
	}

	expr := args[0]
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	schedule, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	now := time.Now()
	nextTimes := make([]string, 0, count)
	currentTime := now

	for i := 0; i < count; i++ {
		nextTime := schedule.Next(currentTime)
		nextTimes = append(nextTimes, nextTime.Format(time.RFC3339))
		currentTime = nextTime
	}

	result := map[string]interface{}{
		"expression": expr,
		"next_times": nextTimes,
		"count":      count,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Expression: %s\n", expr)
		fmt.Println("Next execution times:")
		for i, t := range nextTimes {
			fmt.Printf("  %d. %s\n", i+1, t)
		}
	}

	return nil
}

func parseCronExpression(expr string) []string {
	parts := make([]string, 5)
	
	// Simple parsing - split by space
	values := strings.Fields(expr)
	for i := 0; i < 5 && i < len(values); i++ {
		parts[i] = values[i]
	}
	
	return parts
}

func explainCronExpression(parts []string) string {
	if len(parts) < 5 {
		return "Invalid cron expression"
	}

	explanations := []string{
		fmt.Sprintf("Minute: %s", explainField(parts[0], 0, 59)),
		fmt.Sprintf("Hour: %s", explainField(parts[1], 0, 23)),
		fmt.Sprintf("Day of month: %s", explainField(parts[2], 1, 31)),
		fmt.Sprintf("Month: %s", explainField(parts[3], 1, 12)),
		fmt.Sprintf("Day of week: %s", explainField(parts[4], 0, 6)),
	}

	return strings.Join(explanations, ", ")
}

func explainField(field string, min, max int) string {
	if field == "*" {
		return "every value"
	}
	if strings.HasPrefix(field, "*/") {
		interval := strings.TrimPrefix(field, "*/")
		return fmt.Sprintf("every %s", interval)
	}
	if strings.Contains(field, "-") {
		return fmt.Sprintf("range %s", field)
	}
	if strings.Contains(field, ",") {
		return fmt.Sprintf("values %s", field)
	}
	return field
}
