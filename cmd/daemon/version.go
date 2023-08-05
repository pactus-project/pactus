package main

import (
	"github.com/spf13/cobra"
	"github.com/pactus-project/pactus/version"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the Pactus version",
	Run: Version(),
}

// Version prints the version of the Pactus node.
func Version() func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		c.Printf("Pactus version: %v\n", version.Version())
	}
}
