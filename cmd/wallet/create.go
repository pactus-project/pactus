package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

// Generate creates a new wallet
func Generate() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		testnetOpt := c.Bool(cli.BoolOpt{
			Name:  "testnet",
			Desc:  "creating wallet for testnet",
			Value: false,
		})

		c.Before = func() {}
		c.Action = func() {
			password := cmd.PromptPassword("Password", true)
			mnemonic := wallet.GenerateMnemonic()

			network := wallet.NetworkMainNet
			if *testnetOpt {
				network = wallet.NetworkTestNet
			}
			wallet, err := wallet.FromMnemonic(*path, mnemonic, password, network)
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

// ChangePassword updates the wallet password
func ChangePassword() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			oldPassword := getPassword(wallet, *passOpt)
			newPassword := cmd.PromptPassword("New Password", true)

			err = wallet.UpdatePassword(oldPassword, newPassword)
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
			cmd.PrintWarnMsg("Wallet password updated")
		}
	}
}
