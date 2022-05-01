package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
)

var path *string

func main() {
	app := cli.App("zarb-wallet", "Zarb wallet")

	path = app.String(cli.StringOpt{
		Name:  "w wallet file",
		Desc:  "a path to the wallet file",
		Value: cmd.ZarbDefaultWalletPath(cmd.ZarbHomeDir()),
	})

	app.Command("create", "Create a new wallet", Generate())
	app.Command("recover", "Recover waller from the seed phrase (mnemonic)", Recover())
	app.Command("seed", "Show secret seed phrase (mnemonic) that can be used to recover this wallet", GetSeed())
	app.Command("address", "Manage address book", func(k *cli.Cmd) {
		k.Command("new", "Creating a new address", NewAddress())
		k.Command("all", "Show all addresses", AllAddresses())
		k.Command("balance", "Show the balance of an address", GetBalance())
		k.Command("setlabel", "Set label for the an address", SetLabel())
		k.Command("pubkey", "Get public key of an address", GetPublicKey())
		k.Command("privkey", "Get private key of an address", GetPrivateKey())
		k.Command("import", "Import a private key into wallet", ImportPrivateKey())
	})
	app.Command("tx", "Create, sign and publish a transaction", func(k *cli.Cmd) {
		k.Command("bond", "Create, sign and publish a bond transaction", BondTx())
		k.Command("send", "Create, sign and publish a send transaction", SendTx())
		k.Command("unbond", "Create, sign and publish an unbond transaction", UnbondTx())
		k.Command("withdraw", "Create, sign and publish a withdraw transaction", WithdrawTx())
	})

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
