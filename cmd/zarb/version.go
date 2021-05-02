package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/version"
)

// Version prints the version of the Zarb node
func Version() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			fmt.Println()
			cmd.PrintInfoMsg("Zarb version: %v", version.Version())
		}
	}
}
