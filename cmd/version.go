/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the version number of the application, to be updated using build-time flags
var Version = "v0.0.1" // Default version

// BuildDate is the build date of the application, set at build time
var BuildDate = "undefined"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Netro",
	Long:  "All software has versions. This is Netro's version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Netro version: %s (built on %s)\n", Version, BuildDate)
	},
}

func init() {
	// Register the version command as a subcommand of the root command
	rootCmd.AddCommand(versionCmd)
}
