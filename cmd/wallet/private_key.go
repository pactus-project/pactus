package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/wallet"
)

// GetPrivateKey returns the private key of an address
func GetPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addressOpt := c.String(cli.StringOpt{
			Name: "a address",
			Desc: "Address string",
		})

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			passphrase := cmd.PromptPassphrase("Passphrase: ", false)

			fmt.Println()

			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			prv, err := wallet.PrivateKey(passphrase, *addressOpt)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintSuccessMsg("Private Key: \"%v\"", prv.String())
		}
	}
}
