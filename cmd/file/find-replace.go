package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// findReplaceCmd represents the find-replace command
var findReplaceCmd = &cobra.Command{
	Use:   "find-replace [find] [replace] [path]",
	Short: "Find and replace text in multiple files",
	Long: `Find and replace text patterns in multiple files.

Examples:
  devkit file find-replace "old" "new" .
  devkit file find-replace "TODO" "DONE" ./src --recursive
  devkit file find-replace "error" "err" . --extensions "go,js" --dry-run`,
	Args: cobra.MinimumNArgs(2),
	RunE: runFindReplace,
}

func init() {
	fileCmd.AddCommand(findReplaceCmd)

	findReplaceCmd.Flags().StringP("path", "p", ".", "Path to search in")
	findReplaceCmd.Flags().BoolP("recursive", "r", false, "Search recursively")
	findReplaceCmd.Flags().String("extensions", "", "File extensions to search (comma-separated)")
	findReplaceCmd.Flags().String("ignore", "", "Directories to ignore (comma-separated)")
	findReplaceCmd.Flags().BoolP("regex", "e", false, "Use regex pattern")
	findReplaceCmd.Flags().BoolP("dry-run", "d", false, "Show what would be changed without making changes")
	findReplaceCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runFindReplace(cmd *cobra.Command, args []string) error {
	find := args[0]
	replace := args[1]
	searchPath := "."
	if len(args) > 2 {
		searchPath = args[2]
	}

	recursive, _ := cmd.Flags().GetBool("recursive")
	extensions, _ := cmd.Flags().GetString("extensions")
	ignore, _ := cmd.Flags().GetString("ignore")
	useRegex, _ := cmd.Flags().GetBool("regex")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	extList := strings.Split(extensions, ",")
	ignoreList := strings.Split(ignore, ",")

	var searchPattern *regexp.Regexp
	var err error

	if useRegex {
		searchPattern, err = regexp.Compile(find)
	} else {
		searchPattern, err = regexp.Compile(regexp.QuoteMeta(find))
	}

	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	var results []map[string]interface{}
	totalReplacements := 0

	err = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if !recursive && path != searchPath {
				return filepath.SkipDir
			}
			for _, ignoreDir := range ignoreList {
				if strings.Contains(path, strings.TrimSpace(ignoreDir)) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if extensions != "" {
			ext := filepath.Ext(path)
			found := false
			for _, e := range extList {
				if strings.TrimSpace(e) == strings.TrimPrefix(ext, ".") {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}

		scanner := bufio.NewScanner(file)
		var lines []string
		fileChanged := false
		lineNum := 0

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if searchPattern.MatchString(line) {
				newLine := searchPattern.ReplaceAllString(line, replace)
				if newLine != line {
					fileChanged = true
					totalReplacements++
				}
				lines = append(lines, newLine)
			} else {
				lines = append(lines, line)
			}
		}
		file.Close()

		if fileChanged {
			result := map[string]interface{}{
				"file":         path,
				"replacements": countReplacements(lines, searchPattern, find, replace),
			}
			results = append(results, result)

			if !dryRun {
				outputFile, err := os.Create(path)
				if err != nil {
					return fmt.Errorf("failed to write file %s: %w", path, err)
				}
				writer := bufio.NewWriter(outputFile)
				for _, line := range lines {
					fmt.Fprintln(writer, line)
				}
				writer.Flush()
				outputFile.Close()
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("find-replace error: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"find":      find,
			"replace":   replace,
			"path":     searchPath,
			"files":    results,
			"count":    len(results),
			"total":    totalReplacements,
			"dry_run":  dryRun,
		})
	} else {
		if dryRun {
			fmt.Println("DRY RUN - No files were modified")
		}
		for _, result := range results {
			fmt.Printf("Modified: %s (%d replacements)\n", result["file"], result["replacements"])
		}
		fmt.Printf("\nTotal: %d files, %d replacements\n", len(results), totalReplacements)
	}

	return nil
}

func countReplacements(lines []string, pattern *regexp.Regexp, find, replace string) int {
	count := 0
	for _, line := range lines {
		matches := pattern.FindAllString(line, -1)
		count += len(matches)
	}
	return count
}
