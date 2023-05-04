package main

import (
	"time"

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
			cmd.FatalErrorCheck(err)

			id, err := hash.FromString(*txID)
			cmd.FatalErrorCheck(err)

			err = wallet.AddTransaction(id)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

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
			cmd.FatalErrorCheck(err)

			history := wallet.GetHistory(*addrArg)
			for i, h := range history {
				if h.Time != nil {
					cmd.PrintInfoMsg("%d %v %v %v %s\t%v",
						i+1, h.Time.Format(time.RFC822), h.TxID, h.PayloadType, h.Desc, h.Amount)
				} else {
					cmd.PrintInfoMsg("%d              %v  %s\t%v",
						i+1, h.TxID, h.Desc, h.Amount)
				}
			}
		}
	}
}
