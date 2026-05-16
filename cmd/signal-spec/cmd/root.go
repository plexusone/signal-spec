// Package cmd provides the signal-spec CLI commands.
package cmd

import (
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "signal-spec",
	Short:   "CLI for signal-spec operations",
	Version: Version,
	Long: `signal-spec is a CLI for working with operational signal specifications.

It provides commands for generating reports, validating signals, and
managing JSON schemas for the signal-spec data model.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}
