package main

import (
	"github.com/pactus-project/pactus/version"
	"github.com/spf13/cobra"
)

// Version prints the version of the Pactus node.
func buildVersionCmd(parentCmd *cobra.Command) {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the Pactus version",
	}
	parentCmd.AddCommand(versionCmd)
	versionCmd.Run = func(c *cobra.Command, args []string) {
		c.Printf("Pactus version: %v\n", version.Version())
	}
}
