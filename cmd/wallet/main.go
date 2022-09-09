package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
)

var pathOpt *string
var offlineOpt *bool

func addPasswordOption(c *cli.Cmd) *string {
	return c.String(cli.StringOpt{
		Name:  "password",
		Desc:  "provide wallet password as a parameter instead of interactively",
		Value: "",
	})
}

func main() {
	app := cli.App("pactus-wallet", "Pactus wallet")

	pathOpt = app.String(cli.StringOpt{
		Name:  "w wallet file",
		Desc:  "a path to the wallet file",
		Value: cmd.PactusDefaultWalletPath(cmd.PactusHomeDir()),
	})

	offlineOpt = app.Bool(cli.BoolOpt{
		Name:  "offline",
		Desc:  "offline mode",
		Value: false,
	})

	app.Command("create", "Create a new wallet", Generate())
	app.Command("recover", "Recover waller from the seed phrase (mnemonic)", Recover())
	app.Command("seed", "Show secret seed phrase (mnemonic) that can be used to recover this wallet", GetSeed())
	app.Command("password", "Change wallet password", ChangePassword())
	app.Command("address", "Manage address book", func(k *cli.Cmd) {
		k.Command("new", "Creating a new address", NewAddress())
		k.Command("all", "Show all addresses", AllAddresses())
		k.Command("label", "Set label for the an address", SetLabel())
		k.Command("balance", "Show the balance of an address", Balance())
		k.Command("pub", "Show the public key of an address", PublicKey())
		k.Command("priv", "Show the private key of an address", PrivateKey())
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
