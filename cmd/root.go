package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"devkit/cmd/dev"
	"devkit/cmd/file"
	"devkit/cmd/net"
	"devkit/pkg/version"
)

var (
	cfgFile string
	verbose bool
	quiet   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devkit",
	Short: "A versatile CLI toolkit for developers",
	Long: `DevKit is a comprehensive command-line toolkit designed for developers
and system administrators. It provides a wide range of utilities for
file operations, developer tools, network operations, data processing,
and time/date manipulation.

Built with Go, DevKit is distributed as a single binary and works
cross-platform.`,
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devkit.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (suppress non-error output)")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	// Add subcommands
	rootCmd.AddCommand(dev.GetDevCmd())
	rootCmd.AddCommand(file.GetFileCmd())
	rootCmd.AddCommand(net.GetNetCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".devkit" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devkit")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// GetVerbose returns the verbose flag value
func GetVerbose() bool {
	return verbose
}

// GetQuiet returns the quiet flag value
func GetQuiet() bool {
	return quiet
}
