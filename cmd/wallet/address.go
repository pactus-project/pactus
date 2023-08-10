package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
)

// AllAddresses lists all the wallet addresses.
func AllAddresses() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		balanceOpt := c.Bool(cli.BoolOpt{
			Name:  "balance",
			Desc:  "show account balance",
			Value: false,
		})

		stakeOpt := c.Bool(cli.BoolOpt{
			Name:  "stake",
			Desc:  "show validator stake",
			Value: false,
		})

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			for i, info := range wallet.AddressLabels() {
				line := fmt.Sprintf("%v- %s\t", i+1, info.Address)

				if *balanceOpt {
					balance, _ := wallet.Balance(info.Address)
					line += fmt.Sprintf("%v\t", util.ChangeToCoin(balance))
				}

				if *stakeOpt {
					stake, _ := wallet.Stake(info.Address)
					line += fmt.Sprintf("%v\t", util.ChangeToCoin(stake))
				}

				line += info.Label
				if info.Imported {
					line += " (Imported)"
				}

				cmd.PrintInfoMsg(line)
			}
		}
	}
}

// NewAddress creates a new address.
func NewAddress() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() {}
		c.Action = func() {
			label := cmd.PromptInput("Label")
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			addr, err := wallet.DeriveNewAddress(label)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintInfoMsg("%s", addr)
		}
	}
}

// Balance shows the balance of an address.
func Balance() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			balance, err := wallet.Balance(*addrArg)
			cmd.FatalErrorCheck(err)

			stake, err := wallet.Stake(*addrArg)
			cmd.FatalErrorCheck(err)

			cmd.PrintInfoMsg("balance: %v\tstake: %v",
				util.ChangeToCoin(balance), util.ChangeToCoin(stake))
		}
	}
}

// PrivateKey returns the private key of an address.
func PrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			password := getPassword(wallet, *passOpt)
			prv, err := wallet.PrivateKey(password, *addrArg)
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintWarnMsg("Private Key: %v", prv)
		}
	}
}

// PublicKey returns the public key of an address.
func PublicKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			info := wallet.AddressInfo(*addrArg)
			if info == nil {
				cmd.PrintErrorMsg("Address not found")
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Public Key: %v", info.Pub.String())
			if !info.Imported {
				cmd.PrintInfoMsg("Path: %v", info.Path.String())
			}
		}
	}
}

// ImportPrivateKey imports a private key into the wallet.
func ImportPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			prvStr := cmd.PromptInput("Private Key")

			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			prv, err := bls.PrivateKeyFromString(prvStr)
			cmd.FatalErrorCheck(err)

			password := getPassword(wallet, *passOpt)
			err = wallet.ImportPrivateKey(password, prv)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintSuccessMsg("Private Key imported. Address: %v",
				prv.PublicKey().Address())
		}
	}
}

// SetLabel set label for the address.
func SetLabel() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Action = func() {
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			oldLabel := wallet.Label(*addrArg)
			newLabel := cmd.PromptInputWithSuggestion("Label", oldLabel)

			err = wallet.SetLabel(*addrArg, newLabel)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintSuccessMsg("Label set successfully")
		}
	}
}

func addAddressArg(c *cli.Cmd) *string {
	addrArg := c.String(cli.StringArg{
		Name: "ADDRESS",
		Desc: "address string",
	})

	return addrArg
}
