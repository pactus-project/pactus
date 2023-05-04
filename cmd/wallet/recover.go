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
			wallet, err := wallet.Create(*pathOpt, mnemonic, "", 0)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

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
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			password := getPassword(wallet, *passOpt)
			mnemonic, err := wallet.Mnemonic(password)
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
		}
	}
}
