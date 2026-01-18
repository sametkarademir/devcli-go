package file

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [file1] [file2]",
	Short: "Compare two files or directories",
	Long: `Compare two files or directories and show differences.

Examples:
  devkit file diff file1.txt file2.txt
  devkit file diff dir1/ dir2/`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	fileCmd.AddCommand(diffCmd)

	diffCmd.Flags().BoolP("unified", "u", false, "Show unified diff format")
	diffCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runDiff(cmd *cobra.Command, args []string) error {
	file1 := args[0]
	file2 := args[1]
	_, _ = cmd.Flags().GetBool("unified") // unified flag for future use
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	info1, err := os.Stat(file1)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", file1, err)
	}

	info2, err := os.Stat(file2)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", file2, err)
	}

	if info1.IsDir() || info2.IsDir() {
		return fmt.Errorf("directory diff not yet implemented")
	}

	lines1, err := readLines(file1)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1, err)
	}

	lines2, err := readLines(file2)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2, err)
	}

	diffs := computeDiff(lines1, lines2)

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"file1": file1,
			"file2": file2,
			"diffs": diffs,
		})
	} else {
		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Printf("%s %s\n", yellow("---"), file1)
		fmt.Printf("%s %s\n", yellow("+++"), file2)

		for _, diff := range diffs {
			switch diff.Type {
			case "removed":
				fmt.Printf("%s %s\n", red("-"), diff.Line)
			case "added":
				fmt.Printf("%s %s\n", green("+"), diff.Line)
			case "context":
				fmt.Printf("%s %s\n", " ", diff.Line)
			}
		}
	}

	return nil
}

type diffLine struct {
	Type string
	Line string
	Num  int
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func computeDiff(lines1, lines2 []string) []diffLine {
	var diffs []diffLine
	maxLen := len(lines1)
	if len(lines2) > maxLen {
		maxLen = len(lines2)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(lines1) {
			diffs = append(diffs, diffLine{Type: "added", Line: lines2[i], Num: i + 1})
		} else if i >= len(lines2) {
			diffs = append(diffs, diffLine{Type: "removed", Line: lines1[i], Num: i + 1})
		} else if lines1[i] != lines2[i] {
			diffs = append(diffs, diffLine{Type: "removed", Line: lines1[i], Num: i + 1})
			diffs = append(diffs, diffLine{Type: "added", Line: lines2[i], Num: i + 1})
		} else {
			diffs = append(diffs, diffLine{Type: "context", Line: lines1[i], Num: i + 1})
		}
	}

	return diffs
}
