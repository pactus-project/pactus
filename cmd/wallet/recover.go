package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

/// Recover recovers a wallet from mnemonic (seed phrase)
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			mnemonic := cmd.PromptInput("Seed: ")
			w, err := wallet.FromMnemonic(*path, mnemonic, "", 0)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Wallet recovered successfully at: %s", w.Path())
			cmd.PrintWarnMsg("Never share your private key.")
			cmd.PrintWarnMsg("Don't forget to set a password for your wallet.")
		}
	}
}

/// GetSeed prints the seed phrase (mnemonics)
func GetSeed() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			mnemonic, err := w.Mnemonic(passphrase)

			cmd.PrintLine()
			cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
		}
	}
}
