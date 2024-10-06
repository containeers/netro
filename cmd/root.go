/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "netro",
	Short: "Netro - A versatile networking and troubleshooting CLI tool",
	Long: `Netro is a command-line tool designed for developers and administrators 
to perform a variety of network-related operations. It is powered by Go's Cobra library, 
allowing you to perform DNS lookups, network diagnostics, container troubleshooting, 
and more through simple, powerful CLI commands.

Examples:

# Perform a DNS lookup for a domain:
netro dig example.com

# Display network interfaces and addresses:
netro ifconfig

# Perform a basic network diagnostic:
netro netstat
`,
	// The action when no subcommand is provided
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Netro! Use 'netro --help' to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This function is called by main.main() and sets the starting point for the CLI.
// It only needs to be called once to initiate the root command and its subcommands.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags are global and can be used with any subcommand of 'netro'.
	// Here you can define configuration-related flags for the entire application.
	// Example: configuration file support can be added.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netro.yaml)")

	// Local flags, specific to the root command itself (i.e., when no subcommands are provided).
	// The 'toggle' flag is an example of a boolean flag.
	rootCmd.Flags().BoolP("toggle", "t", false, "Enable or disable specific features in Netro")
}
