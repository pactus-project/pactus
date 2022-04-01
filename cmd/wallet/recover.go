package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

// Generate creates a new wallet
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			mnemonic := cmd.PromptInput("Mnemonic: ")

			fmt.Println()

			_, err := wallet.RecoverWallet(*path, mnemonic)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

		}
	}
}
