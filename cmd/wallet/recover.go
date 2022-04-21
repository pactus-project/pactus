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
			wallet, err := wallet.FromMnemonic(*path, mnemonic, "", "Recovered wallet", 0)
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
			cmd.PrintInfoMsg("Wallet recovered successfully at: %s", wallet.Path())
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
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet)
			mnemonic, err := wallet.Mnemonic(password)

			cmd.PrintLine()
			cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
		}
	}
}
