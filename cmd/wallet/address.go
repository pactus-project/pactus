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
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			addrs := w.Addresses()
			for addr, label := range addrs {
				cmd.PrintInfoMsg("%s %s", addr, label)
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
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			addr, err := w.NewAddress(passphrase, label)
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
		addrArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "address string",
		})

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			balance, stake, err := w.GetBalance(*addrArg)
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
		addrArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "address string",
		})

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			prv, err := w.PrivateKey(passphrase, *addrArg)
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
		addrArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "address string",
		})

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			pub, err := w.PublicKey(passphrase, *addrArg)
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

			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			err = w.ImportPrivateKey(passphrase, prv)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintSuccessMsg("Private Key imported")
		}
	}
}
