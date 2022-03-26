package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd/zarbcli/key"
	"github.com/zarbchain/zarb-go/cmd/zarbcli/tx"
)

func zarbcli() *cli.Cli {
	app := cli.App("zarbcli", "Zarb blockchain node (client)")

	app.Command("key", "Create zarb key file for signing messages", func(k *cli.Cmd) {
		k.Command("generate", "Generate a new key", key.Generate())
		k.Command("recover", "Recover a key from the seed", key.Recover())
		k.Command("inspect", "Inspect a key file", key.Inspect())
		k.Command("sign", "Sign a transaction or message with a key file", key.Sign())
		k.Command("verify", "Verify a signature", key.Verify())
		k.Command("change-auth", "Change the passphrase of a keyfile", key.ChangeAuth())
	})
	app.Command("tx", "Create, sign and publish a transaction", func(k *cli.Cmd) {
		k.Command("bond", "Create, sign and publish a bond transaction", tx.BondTx())
		k.Command("send", "Create, sign and publish a send transactio", tx.SendTx())
		k.Command("unbond", "Create, sign and publish an unbond transaction", tx.UnbondTx())
		k.Command("withdraw", "Create, sign and publish a withdraw transaction", tx.WithdrawTx())
	})
	return app
}

func main() {
	if err := zarbcli().Run(os.Args); err != nil {
		panic(err)
	}
}
