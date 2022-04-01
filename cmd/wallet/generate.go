package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

// Generate creates a new wallet
func Generate() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			passphrase := cmd.PromptPassphrase("Passphrase: ", true)

			fmt.Println()

			wallet, err := wallet.NewWallet(*path, passphrase)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			mnemonic := wallet.Mnemonic(passphrase)

			cmd.PrintSuccessMsg("mnemonic: \"%v\"", mnemonic)
		}
	}
}
