package dev

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// envCmd represents the env command group
var envCmd = &cobra.Command{
	Use:   "env",
	Short: ".env file management",
	Long: `Manage .env files (read, set, unset variables).

Examples:
  devkit dev env get KEY --file .env
  devkit dev env set KEY=value --file .env
  devkit dev env unset KEY --file .env
  devkit dev env list --file .env`,
}

// envGetCmd represents the get subcommand
var envGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get value of an environment variable from .env file",
	Long: `Get the value of a specific key from .env file.

Examples:
  devkit dev env get DATABASE_URL --file .env
  devkit dev env get API_KEY`,
	RunE: runEnvGet,
}

// envSetCmd represents the set subcommand
var envSetCmd = &cobra.Command{
	Use:   "set [key=value]",
	Short: "Set an environment variable in .env file",
	Long: `Set or update an environment variable in .env file.

Examples:
  devkit dev env set DATABASE_URL=postgres://... --file .env
  devkit dev env set API_KEY=secret123`,
	RunE: runEnvSet,
}

// envUnsetCmd represents the unset subcommand
var envUnsetCmd = &cobra.Command{
	Use:   "unset [key]",
	Short: "Remove an environment variable from .env file",
	Long: `Remove an environment variable from .env file.

Examples:
  devkit dev env unset DATABASE_URL --file .env
  devkit dev env unset API_KEY`,
	RunE: runEnvUnset,
}

// envListCmd represents the list subcommand
var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment variables from .env file",
	Long: `List all environment variables from .env file.

Examples:
  devkit dev env list --file .env
  devkit dev env list`,
	RunE: runEnvList,
}

func init() {
	devCmd.AddCommand(envCmd)
	envCmd.AddCommand(envGetCmd)
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envUnsetCmd)
	envCmd.AddCommand(envListCmd)

	envGetCmd.Flags().StringP("file", "f", ".env", ".env file path")
	envGetCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	envSetCmd.Flags().StringP("file", "f", ".env", ".env file path")
	envSetCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	envUnsetCmd.Flags().StringP("file", "f", ".env", ".env file path")
	envUnsetCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")

	envListCmd.Flags().StringP("file", "f", ".env", ".env file path")
	envListCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func getEnvFilePath(cmd *cobra.Command) string {
	fileFlag, _ := cmd.Flags().GetString("file")
	if fileFlag == "" {
		return ".env"
	}
	return fileFlag
}

func readEnvFile(filePath string) (map[string]string, error) {
	env := make(map[string]string)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return env, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			env[key] = value
		}
	}

	return env, scanner.Err()
}

func writeEnvFile(filePath string, env map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range env {
		// Quote value if it contains spaces
		if strings.Contains(value, " ") {
			value = fmt.Sprintf(`"%s"`, value)
		}
		fmt.Fprintf(file, "%s=%s\n", key, value)
	}

	return nil
}

func runEnvGet(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("key not specified")
	}

	key := args[0]
	filePath := getEnvFilePath(cmd)

	env, err := readEnvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	value, exists := env[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	} else {
		output.PrintSuccess(format, value)
	}

	return nil
}

func runEnvSet(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("key=value not specified")
	}

	kv := args[0]
	parts := strings.SplitN(kv, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format: expected key=value")
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	filePath := getEnvFilePath(cmd)

	env, err := readEnvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	env[key] = value

	if err := writeEnvFile(filePath, env); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"key":   key,
			"value": value,
			"action": "set",
		})
	} else {
		output.PrintSuccess(format, fmt.Sprintf("Set %s=%s", key, value))
	}

	return nil
}

func runEnvUnset(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if len(args) == 0 {
		return fmt.Errorf("key not specified")
	}

	key := args[0]
	filePath := getEnvFilePath(cmd)

	env, err := readEnvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	if _, exists := env[key]; !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	delete(env, key)

	if err := writeEnvFile(filePath, env); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"key":    key,
			"action": "unset",
		})
	} else {
		output.PrintSuccess(format, fmt.Sprintf("Unset %s", key))
	}

	return nil
}

func runEnvList(cmd *cobra.Command, args []string) error {
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	filePath := getEnvFilePath(cmd)

	env, err := readEnvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, env)
	} else {
		if len(env) == 0 {
			fmt.Println("No environment variables found")
			return nil
		}
		for key, value := range env {
			fmt.Printf("%s=%s\n", key, value)
		}
	}

	return nil
}
