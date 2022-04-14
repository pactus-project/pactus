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
			mnemonic := wallet.GenerateMnemonic()
			wallet, err := wallet.FromMnemonic(*path, mnemonic, passphrase, 0)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			err = wallet.Save()
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintSuccessMsg("Wallet created successfully at: %s", wallet.Path())
			cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
			cmd.PrintWarnMsg("Please keep your seed in a safe place; if you lose it, you will not be able to restore your wallet.")
		}
	}
}
