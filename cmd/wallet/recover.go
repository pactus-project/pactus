package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
)

// Recover recovers a wallet from mnemonic (seed phrase).
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() {}
		c.Action = func() {
			mnemonic := cmd.PromptInput("Seed")
			wallet, err := wallet.FromMnemonic(*pathOpt, mnemonic, "", 0)
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

// GetSeed prints the seed phrase (mnemonics).
func GetSeed() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*pathOpt, *offlineOpt)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet, *passOpt)
			mnemonic, err := wallet.Mnemonic(password)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
		}
	}
}
