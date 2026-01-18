package file

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// dedupeCmd represents the dedupe command
var dedupeCmd = &cobra.Command{
	Use:   "dedupe [path]",
	Short: "Find and remove duplicate files",
	Long: `Find duplicate files by hash and optionally remove them.

Examples:
  devkit file dedupe ./downloads --by hash
  devkit file dedupe ./photos --by name --action delete --dry-run`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDedupe,
}

func init() {
	fileCmd.AddCommand(dedupeCmd)

	dedupeCmd.Flags().StringP("by", "b", "hash", "Comparison method: hash, name")
	dedupeCmd.Flags().StringP("action", "a", "list", "Action: list, delete")
	dedupeCmd.Flags().BoolP("dry-run", "d", false, "Show what would be deleted without making changes")
	dedupeCmd.Flags().BoolP("recursive", "r", false, "Search recursively")
	dedupeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runDedupe(cmd *cobra.Command, args []string) error {
	searchPath := "."
	if len(args) > 0 {
		searchPath = args[0]
	}

	by, _ := cmd.Flags().GetString("by")
	action, _ := cmd.Flags().GetString("action")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	recursive, _ := cmd.Flags().GetBool("recursive")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	fileMap := make(map[string][]string)

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

		var key string
		if by == "hash" {
			hash, err := calculateFileHash(path)
			if err != nil {
				return nil
			}
			key = hash
		} else {
			key = info.Name()
		}

		fileMap[key] = append(fileMap[key], path)
		return nil
	})

	if err != nil {
		return fmt.Errorf("dedupe error: %w", err)
	}

	var duplicates []map[string]interface{}
	var toDelete []string

	for key, files := range fileMap {
		if len(files) > 1 {
			// Keep first, mark others for deletion
			keep := files[0]
			dups := files[1:]

			duplicate := map[string]interface{}{
				"key":  key,
				"keep": keep,
				"duplicates": dups,
				"count": len(dups),
			}
			duplicates = append(duplicates, duplicate)

			if action == "delete" {
				toDelete = append(toDelete, dups...)
			}
		}
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"path":      searchPath,
			"method":    by,
			"duplicates": duplicates,
			"count":     len(duplicates),
			"to_delete": len(toDelete),
			"dry_run":   dryRun,
		})
	} else {
		for _, dup := range duplicates {
			fmt.Printf("\nDuplicate group (key: %s):\n", dup["key"])
			fmt.Printf("  Keep: %s\n", dup["keep"])
			for _, file := range dup["duplicates"].([]string) {
				fmt.Printf("  Duplicate: %s\n", file)
			}
		}

		if action == "delete" && len(toDelete) > 0 {
			if dryRun {
				fmt.Println("\nDRY RUN - Would delete:")
			} else {
				fmt.Println("\nDeleting duplicates:")
			}
			for _, file := range toDelete {
				if !dryRun {
					if err := os.Remove(file); err != nil {
						fmt.Printf("  Failed to delete %s: %v\n", file, err)
					} else {
						fmt.Printf("  Deleted: %s\n", file)
					}
				} else {
					fmt.Printf("  Would delete: %s\n", file)
				}
			}
		}

		fmt.Printf("\nFound %d duplicate groups\n", len(duplicates))
	}

	return nil
}

func calculateFileHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
