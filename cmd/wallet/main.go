package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/wallet"
)

var pathArg *string
var offlineOpt *bool
var serverAddrOpt *string

func addPasswordOption(c *cli.Cmd) *string {
	return c.String(cli.StringOpt{
		Name: "p password",
		Desc: "the wallet password",
	})
}

func openWallet() (*wallet.Wallet, error) {
	if !*offlineOpt && *serverAddrOpt != "" {
		wallet, err := wallet.Open(*pathArg, true)
		if err != nil {
			return nil, err
		}

		err = wallet.Connect(*serverAddrOpt)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return wallet, err
	}
	wallet, err := wallet.Open(*pathArg, *offlineOpt)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func main() {
	app := cli.App("pactus-wallet", "Pactus wallet")

	pathArg = app.String(cli.StringArg{
		Name: "PATH",
		Desc: "the path to the wallet file",
	})

	offlineOpt = app.Bool(cli.BoolOpt{
		Name:  "offline",
		Desc:  "offline mode",
		Value: false,
	})

	serverAddrOpt = app.String(cli.StringOpt{
		Name: "server",
		Desc: "server gRPC address",
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
		k.Command("transfer", "Create, sign and publish a Transfer transaction", TransferTx())
		k.Command("unbond", "Create, sign and publish an unbond transaction", UnbondTx())
		k.Command("withdraw", "Create, sign and publish a Withdraw transaction", WithdrawTx())
	})
	app.Command("history", "Check the wallet history", func(k *cli.Cmd) {
		k.Command("add", "Add a transaction to the wallet history", AddToHistory())
		k.Command("get", "Show the transaction history of any address", ShowHistory())
	})

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
