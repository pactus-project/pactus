package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/version"
)

// Version prints the version of the Pactus node.
func Version() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.Pactus) }
		c.Action = func() {
			fmt.Println()
			cmd.PrintInfoMsg("Pactus version: %v", version.Version())
		}
	}
}
