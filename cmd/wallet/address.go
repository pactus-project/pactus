package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

/// AllAddresses lists all the wallet addresses
func AllAddresses() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
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
				cmd.PrintInfoMsg("%s %s", info.Address, label)
			}
		}
	}
}

/// NewAddress creates a new address
func NewAddress() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			label := cmd.PromptInput("Label: ")
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet)
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

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			balance, stake, err := wallet.Balance(*addrArg)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}
			cmd.PrintInfoMsg("balance: %v, stake: %v", balance, stake)
		}
	}
}

// GetPrivateKey returns the private key of an address
func GetPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet)
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

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet)
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
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			prv := cmd.PromptInput("Private Key: ")

			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			password := getPassword(wallet)
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

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			oldLabel := wallet.Label(*addrArg)
			newLabel := cmd.PromptInputWithSuggestion("Label: ", oldLabel)

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
