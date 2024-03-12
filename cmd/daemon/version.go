package main

import (
	"github.com/pactus-project/pactus/version"
	"github.com/spf13/cobra"
)

// Version prints the version of the Pactus node.
func buildVersionCmd(parentCmd *cobra.Command) {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "prints the Pactus version",
	}
	parentCmd.AddCommand(versionCmd)
	versionCmd.Run = func(c *cobra.Command, _ []string) {
		c.Printf("Pactus version: %v\n", version.Version())
	}
}
