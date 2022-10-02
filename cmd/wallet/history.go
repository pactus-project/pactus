package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/hash"
)

func AddToHistory() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		txID := c.String(cli.StringArg{
			Name: "ID",
			Desc: "transaction id",
		})

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}
			id, err := hash.FromString(*txID)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}
			err = wallet.AddTransaction(id)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}
			err = wallet.Save()
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}
			cmd.PrintInfoMsg("Transaction added to wallet")
		}
	}
}

func ShowHistory() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addrArg := addAddressArg(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			history := wallet.GetHistory(*addrArg)
			for i, h := range history {
				cmd.PrintInfoMsg("%d %v %v\t%v",
					i+1, h.TxID, h.PayloadType, h.Amount)
			}
		}
	}
}
