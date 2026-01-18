package dev

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// semverCmd represents the semver command group
var semverCmd = &cobra.Command{
	Use:   "semver",
	Short: "Semantic version operations",
	Long: `Compare and manipulate semantic versions.

Examples:
  devkit dev semver compare "1.2.3" "1.2.4"
  devkit dev semver bump major "1.2.3"
  devkit dev semver bump minor "1.2.3"
  devkit dev semver bump patch "1.2.3"`,
}

// semverCompareCmd represents the compare subcommand
var semverCompareCmd = &cobra.Command{
	Use:   "compare [version1] [version2]",
	Short: "Compare two semantic versions",
	Long: `Compare two semantic versions and show the result.

Examples:
  devkit dev semver compare "1.2.3" "1.2.4"
  devkit dev semver compare "2.0.0" "1.9.9"`,
	RunE: runSemverCompare,
}

// semverBumpCmd represents the bump subcommand
var semverBumpCmd = &cobra.Command{
	Use:   "bump [type] [version]",
	Short: "Bump a semantic version",
	Long: `Bump a semantic version by major, minor, or patch.

Types: major, minor, patch

Examples:
  devkit dev semver bump major "1.2.3"
  devkit dev semver bump minor "1.2.3"
  devkit dev semver bump patch "1.2.3"`,
	RunE: runSemverBump,
}

func init() {
	devCmd.AddCommand(semverCmd)
	semverCmd.AddCommand(semverCompareCmd)
	semverCmd.AddCommand(semverBumpCmd)

	semverCompareCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
	semverBumpCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runSemverCompare(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) < 2 {
		return fmt.Errorf("two versions required")
	}

	v1, err := semver.NewVersion(args[0])
	if err != nil {
		return fmt.Errorf("invalid version 1: %w", err)
	}

	v2, err := semver.NewVersion(args[1])
	if err != nil {
		return fmt.Errorf("invalid version 2: %w", err)
	}

	comparison := v1.Compare(v2)
	var resultStr string
	switch {
	case comparison < 0:
		resultStr = fmt.Sprintf("%s < %s", v1, v2)
	case comparison > 0:
		resultStr = fmt.Sprintf("%s > %s", v1, v2)
	default:
		resultStr = fmt.Sprintf("%s == %s", v1, v2)
	}

	result := map[string]interface{}{
		"version1":  v1.String(),
		"version2":  v2.String(),
		"comparison": comparison,
		"result":    resultStr,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		output.PrintSuccess(format, resultStr)
	}

	return nil
}

func runSemverBump(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) < 2 {
		return fmt.Errorf("bump type and version required")
	}

	bumpType := args[0]
	versionStr := args[1]

	v, err := semver.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	var bumped semver.Version
	switch bumpType {
	case "major":
		bumped = v.IncMajor()
	case "minor":
		bumped = v.IncMinor()
	case "patch":
		bumped = v.IncPatch()
	default:
		return fmt.Errorf("invalid bump type: %s (supported: major, minor, patch)", bumpType)
	}

	result := map[string]interface{}{
		"original": v.String(),
		"bumped":   bumped.String(),
		"type":     bumpType,
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		output.PrintSuccess(format, bumped.String())
	}

	return nil
}
