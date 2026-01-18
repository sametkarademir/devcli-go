package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatPlain OutputFormat = "plain"
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
)

// Result represents a command result
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Print prints the result in the specified format
func Print(format OutputFormat, result Result) {
	switch format {
	case FormatJSON:
		printJSON(result)
	case FormatTable:
		printTable(result)
	default:
		printPlain(result)
	}
}

// PrintSuccess prints a success result
func PrintSuccess(format OutputFormat, data interface{}) {
	Print(format, Result{
		Success: true,
		Data:    data,
	})
}

// PrintError prints an error result
func PrintError(format OutputFormat, err error) {
	errMsg := err.Error()
	Print(format, Result{
		Success: false,
		Error:   errMsg,
	})
}

// printJSON prints the result as JSON
func printJSON(result Result) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
	}
}

// printPlain prints the result as plain text
func printPlain(result Result) {
	if !result.Success {
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
		return
	}

	switch v := result.Data.(type) {
	case string:
		fmt.Println(v)
	case []string:
		for _, s := range v {
			fmt.Println(s)
		}
	case map[string]interface{}:
		for key, value := range v {
			fmt.Printf("%s: %v\n", key, value)
		}
	default:
		fmt.Printf("%v\n", v)
	}
}

// printTable prints the result as a table (basic implementation)
func printTable(result Result) {
	if !result.Success {
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
		return
	}

	// Basic table implementation
	// For more complex tables, we'll use tablewriter in specific commands
	printPlain(result)
}
