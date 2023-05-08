package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
)

// Recover recovers a wallet from mnemonic (seed phrase).
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passwordOpt := addPasswordOption(c)

		seedOpt := c.String(cli.StringOpt{
			Name: "s seed",
			Desc: "wallet seed phrase (mnemonic)",
		})
		c.Before = func() {}
		c.Action = func() {
			mnemonic := *seedOpt
			if mnemonic == "" {
				mnemonic = cmd.PromptInput("Seed")
			}
			password := *passwordOpt

			wallet, err := wallet.Create(*pathArg, mnemonic, password, 0)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintInfoMsg("Wallet recovered successfully at: %s", wallet.Path())
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
