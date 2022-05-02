package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

/// AllAddresses lists all the wallet addresses
func AllAddresses() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		balanceOpt := c.Bool(cli.BoolOpt{
			Name:  "balance",
			Desc:  "show balance",
			Value: false,
		})

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			for _, info := range wallet.AddressInfos() {
				label := info.Label
				if info.Imported {
					label += " (Imported)"
				}
				if *balanceOpt {
					balance, _ := wallet.Balance(info.Address)
					stake, _ := wallet.Stake(info.Address)
					cmd.PrintInfoMsg("%s\tbalance: %v\tstake: %v\t%s",
						info.Address, changeToCoin(balance),
						changeToCoin(stake), label)
				} else {
					cmd.PrintInfoMsg("%s\t%s", info.Address, label)

				}

			}
		}
	}
}

/// NewAddress creates a new address
func NewAddress() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			label := cmd.PromptInput("Label")
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet, *passOpt)
			addr, err := wallet.MakeNewAddress(password, label)
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
			cmd.PrintInfoMsg("%s", addr)
		}
	}
}

/// GetBalance shows the balance of an address
func GetBalance() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			balance, err := wallet.Balance(*addrArg)
			stake, err := wallet.Stake(*addrArg)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}
			cmd.PrintInfoMsg("%s\tbalance: %v\tstake: %v\t%s",
				changeToCoin(balance), changeToCoin(stake))
		}
	}
}

// GetPrivateKey returns the private key of an address
func GetPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet, *passOpt)
			prv, err := wallet.PrivateKey(password, *addrArg)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintWarnMsg("Private Key: \"%v\"", prv)
		}
	}
}

// GetPrivateKey returns the public key of an address
func GetPublicKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet, *passOpt)
			pub, err := wallet.PublicKey(password, *addrArg)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Public Key: \"%v\"", pub)
		}
	}
}

// ImportPrivateKey imports a private key into the wallet
func ImportPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			prv := cmd.PromptInput("Private Key")

			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet, *passOpt)
			err = wallet.ImportPrivateKey(password, prv)
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
			cmd.PrintSuccessMsg("Private Key imported")
		}
	}
}

// SetLabel set label for the address
func SetLabel() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			oldLabel := wallet.Label(*addrArg)
			newLabel := cmd.PromptInputWithSuggestion("Label", oldLabel)

			err = wallet.SetLabel(*addrArg, newLabel)
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
