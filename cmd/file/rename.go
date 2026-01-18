package file

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename [pattern] [path]",
	Short: "Bulk rename files using patterns",
	Long: `Bulk rename files using regex, prefix, suffix, or case conversion.

Examples:
  devkit file rename --pattern "*.txt" --prefix "backup_" --path ./docs
  devkit file rename --pattern "IMG_*.jpg" --replace "IMG_" "photo_" --path ./images
  devkit file rename --pattern "*.txt" --case upper --path ./docs --dry-run`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRename,
}

func init() {
	fileCmd.AddCommand(renameCmd)

	renameCmd.Flags().StringP("pattern", "p", "*", "File pattern (glob)")
	renameCmd.Flags().StringP("path", "P", ".", "Path to search in")
	renameCmd.Flags().String("prefix", "", "Add prefix to filename")
	renameCmd.Flags().String("suffix", "", "Add suffix to filename (before extension)")
	renameCmd.Flags().String("replace", "", "Replace pattern (use with --with)")
	renameCmd.Flags().String("with", "", "Replacement text (use with --replace)")
	renameCmd.Flags().String("case", "", "Case conversion: lower, upper, title")
	renameCmd.Flags().BoolP("recursive", "r", false, "Search recursively")
	renameCmd.Flags().BoolP("dry-run", "d", false, "Show what would be renamed without making changes")
	renameCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runRename(cmd *cobra.Command, args []string) error {
	pattern, _ := cmd.Flags().GetString("pattern")
	searchPath, _ := cmd.Flags().GetString("path")
	prefix, _ := cmd.Flags().GetString("prefix")
	suffix, _ := cmd.Flags().GetString("suffix")
	replacePattern, _ := cmd.Flags().GetString("replace")
	replaceWith, _ := cmd.Flags().GetString("with")
	caseConv, _ := cmd.Flags().GetString("case")
	recursive, _ := cmd.Flags().GetBool("recursive")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) > 0 {
		searchPath = args[0]
	}

	var results []map[string]interface{}

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if !recursive && path != searchPath {
				return filepath.SkipDir
			}
			return nil
		}

		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}

		if !matched {
			return nil
		}

		dir := filepath.Dir(path)
		oldName := info.Name()
		ext := filepath.Ext(oldName)

		newName := oldName

		// Apply prefix
		if prefix != "" {
			newName = prefix + newName
		}

		// Apply suffix
		if suffix != "" {
			newName = strings.TrimSuffix(newName, ext) + suffix + ext
		}

		// Apply replace
		if replacePattern != "" && replaceWith != "" {
			re, err := regexp.Compile(replacePattern)
			if err != nil {
				return fmt.Errorf("invalid replace pattern: %w", err)
			}
			newName = re.ReplaceAllString(newName, replaceWith)
		}

		// Apply case conversion
		switch caseConv {
		case "lower":
			newName = strings.ToLower(newName)
		case "upper":
			newName = strings.ToUpper(newName)
		case "title":
			newName = strings.Title(strings.ToLower(newName))
		}

		if newName != oldName {
			newPath := filepath.Join(dir, newName)
			result := map[string]interface{}{
				"old": oldName,
				"new": newName,
				"path": path,
			}
			results = append(results, result)

			if !dryRun {
				if err := os.Rename(path, newPath); err != nil {
					return fmt.Errorf("failed to rename %s: %w", path, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("rename error: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"pattern": pattern,
			"path":    searchPath,
			"renames": results,
			"count":   len(results),
			"dry_run": dryRun,
		})
	} else {
		if dryRun {
			fmt.Println("DRY RUN - No files were renamed")
		}
		for _, result := range results {
			fmt.Printf("Rename: %s -> %s\n", result["old"], result["new"])
		}
		fmt.Printf("\nTotal: %d files renamed\n", len(results))
	}

	return nil
}
