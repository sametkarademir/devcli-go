package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [pattern] [path]",
	Short: "Search for text in files",
	Long: `Search for text patterns in files with colored output and line numbers.

Examples:
  devkit file search "TODO" .
  devkit file search "function" ./src --recursive
  devkit file search "error" . --extensions "go,js" --ignore "node_modules"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

func init() {
	fileCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("path", "p", ".", "Path to search in")
	searchCmd.Flags().BoolP("recursive", "r", false, "Search recursively")
	searchCmd.Flags().String("extensions", "", "File extensions to search (comma-separated)")
	searchCmd.Flags().String("ignore", "", "Directories to ignore (comma-separated)")
	searchCmd.Flags().BoolP("case-sensitive", "c", false, "Case-sensitive search")
	searchCmd.Flags().BoolP("regex", "e", false, "Use regex pattern")
	searchCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runSearch(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	searchPath, _ := cmd.Flags().GetString("path")
	recursive, _ := cmd.Flags().GetBool("recursive")
	extensions, _ := cmd.Flags().GetString("extensions")
	ignore, _ := cmd.Flags().GetString("ignore")
	caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
	useRegex, _ := cmd.Flags().GetBool("regex")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	extList := strings.Split(extensions, ",")
	ignoreList := strings.Split(ignore, ",")

	var searchPattern *regexp.Regexp
	var err error

	if useRegex {
		if caseSensitive {
			searchPattern, err = regexp.Compile(pattern)
		} else {
			searchPattern, err = regexp.Compile("(?i)" + pattern)
		}
	} else {
		if caseSensitive {
			searchPattern, err = regexp.Compile(regexp.QuoteMeta(pattern))
		} else {
			searchPattern, err = regexp.Compile("(?i)" + regexp.QuoteMeta(pattern))
		}
	}

	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	var results []map[string]interface{}

	err = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if !recursive && path != searchPath {
				return filepath.SkipDir
			}
			// Check ignore list
			for _, ignoreDir := range ignoreList {
				if strings.Contains(path, strings.TrimSpace(ignoreDir)) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check extensions
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
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if searchPattern.MatchString(line) {
				matches := searchPattern.FindAllStringIndex(line, -1)
				highlighted := highlightMatches(line, matches)

				result := map[string]interface{}{
					"file":     path,
					"line":     lineNum,
					"content":  line,
					"highlighted": highlighted,
				}
				results = append(results, result)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("search error: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"pattern": pattern,
			"path":    searchPath,
			"results": results,
			"count":   len(results),
		})
	} else {
		if len(results) == 0 {
			fmt.Println("No matches found")
			return nil
		}

		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()

		currentFile := ""
		for _, result := range results {
			file := result["file"].(string)
			line := result["line"].(int)
			content := result["content"].(string)

			if file != currentFile {
				currentFile = file
				fmt.Printf("\n%s\n", blue(file))
			}
			fmt.Printf("  %s:%s %s\n", yellow(fmt.Sprintf("%d", line)), green("│"), content)
		}
		fmt.Printf("\nFound %d matches\n", len(results))
	}

	return nil
}

func highlightMatches(line string, matches [][]int) string {
	if len(matches) == 0 {
		return line
	}

	red := color.New(color.FgRed, color.Bold).SprintFunc()
	result := ""
	last := 0

	for _, match := range matches {
		result += line[last:match[0]]
		result += red(line[match[0]:match[1]])
		last = match[1]
	}
	result += line[last:]

	return result
}
