package net

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// httpCmd represents the http command group
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "HTTP request operations",
	Long: `Send HTTP requests (GET, POST, PUT, DELETE).

Examples:
  devkit net http get https://api.example.com/users
  devkit net http post https://api.example.com/users --data '{"name":"John"}'
  devkit net http get https://api.example.com --header "Authorization: Bearer token"`,
}

// httpGetCmd represents the get subcommand
var httpGetCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Send GET request",
	RunE:  runHTTPGet,
}

// httpPostCmd represents the post subcommand
var httpPostCmd = &cobra.Command{
	Use:   "post [url]",
	Short: "Send POST request",
	RunE:  runHTTPPost,
}

// httpPutCmd represents the put subcommand
var httpPutCmd = &cobra.Command{
	Use:   "put [url]",
	Short: "Send PUT request",
	RunE:  runHTTPPut,
}

// httpDeleteCmd represents the delete subcommand
var httpDeleteCmd = &cobra.Command{
	Use:   "delete [url]",
	Short: "Send DELETE request",
	RunE:  runHTTPDelete,
}

func init() {
	netCmd.AddCommand(httpCmd)
	httpCmd.AddCommand(httpGetCmd)
	httpCmd.AddCommand(httpPostCmd)
	httpCmd.AddCommand(httpPutCmd)
	httpCmd.AddCommand(httpDeleteCmd)

	// Common flags
	for _, cmd := range []*cobra.Command{httpGetCmd, httpPostCmd, httpPutCmd, httpDeleteCmd} {
		cmd.Flags().StringSliceP("header", "H", []string{}, "HTTP headers (key:value)")
		cmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
	}

	httpPostCmd.Flags().StringP("data", "d", "", "Request body data")
	httpPutCmd.Flags().StringP("data", "d", "", "Request body data")
}

func runHTTPGet(cmd *cobra.Command, args []string) error {
	return runHTTPRequest(cmd, args, "GET", "")
}

func runHTTPPost(cmd *cobra.Command, args []string) error {
	data, _ := cmd.Flags().GetString("data")
	return runHTTPRequest(cmd, args, "POST", data)
}

func runHTTPPut(cmd *cobra.Command, args []string) error {
	data, _ := cmd.Flags().GetString("data")
	return runHTTPRequest(cmd, args, "PUT", data)
}

func runHTTPDelete(cmd *cobra.Command, args []string) error {
	return runHTTPRequest(cmd, args, "DELETE", "")
}

func runHTTPRequest(cmd *cobra.Command, args []string, method, body string) error {
	if len(args) == 0 {
		return fmt.Errorf("URL required")
	}

	url := args[0]
	headers, _ := cmd.Flags().GetStringSlice("header")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	if body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	result := map[string]interface{}{
		"method":      method,
		"url":         url,
		"status_code": resp.StatusCode,
		"status":      resp.Status,
		"headers":     resp.Header,
		"body":        string(respBody),
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, result)
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Response:\n%s\n", string(respBody))
	}

	return nil
}
