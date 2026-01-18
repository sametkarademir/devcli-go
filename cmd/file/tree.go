package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree [directory]",
	Short: "Display directory structure as a tree",
	Long: `Display directory structure in a tree format.

Examples:
  devkit file tree .
  devkit file tree /path/to/directory
  devkit file tree . --depth 2`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTree,
}

func init() {
	fileCmd.AddCommand(treeCmd)

	treeCmd.Flags().IntP("depth", "d", -1, "Maximum depth to traverse (-1 for unlimited)")
	treeCmd.Flags().BoolP("all", "a", false, "Show hidden files")
	treeCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runTree(cmd *cobra.Command, args []string) error {
	depth, _ := cmd.Flags().GetInt("depth")
	showAll, _ := cmd.Flags().GetBool("all")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	var tree []string
	err := buildTree(root, "", 0, depth, showAll, &tree)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"root": root,
			"tree": tree,
		})
	} else {
		for _, line := range tree {
			fmt.Println(line)
		}
	}

	return nil
}

func buildTree(root, prefix string, level, maxDepth int, showAll bool, tree *[]string) error {
	if maxDepth >= 0 && level >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// Filter hidden files
	filtered := []os.DirEntry{}
	for _, entry := range entries {
		if showAll || !strings.HasPrefix(entry.Name(), ".") {
			filtered = append(filtered, entry)
		}
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1
		name := entry.Name()

		var connector string
		if isLast {
			connector = "└── "
			*tree = append(*tree, prefix+connector+name)
		} else {
			connector = "├── "
			*tree = append(*tree, prefix+connector+name)
		}

		if entry.IsDir() {
			var nextPrefix string
			if isLast {
				nextPrefix = prefix + "    "
			} else {
				nextPrefix = prefix + "│   "
			}
			buildTree(filepath.Join(root, name), nextPrefix, level+1, maxDepth, showAll, tree)
		}
	}

	return nil
}
