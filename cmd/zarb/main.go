package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd/zarb/key"
	"github.com/zarbchain/zarb-go/cmd/zarb/tx"
)

func zarb() *cli.Cli {
	app := cli.App("zarb", "Zarb blockchain node")

	app.Command("init", "Initialize the zarb blockchain", Init())
	app.Command("start", "Start the zarb blockchain", Start())
	app.Command("key", "Create zarb key file for signing messages", func(k *cli.Cmd) {
		k.Command("generate", "Generate a new key", key.Generate())
		k.Command("recover", "Recover a key from the seed", key.Recover())
		k.Command("inspect", "Inspect a key file", key.Inspect())
		k.Command("sign", "Sign a transaction or message with a key file", key.Sign())
		k.Command("verify", "Verify a signature", key.Verify())
		k.Command("change-auth", "Change the passphrase of a keyfile", key.ChangeAuth())
	})
	app.Command("tx", "Create raw transaction", func(k *cli.Cmd) {
		k.Command("bond", "Generate raw bond transaction", tx.BondTx())
		k.Command("send", "Generate raw send transaction", tx.SendTx())
	})
	app.Command("version", "Print the zarb version", Version())
	return app
}

func main() {
	if err := zarb().Run(os.Args); err != nil {
		panic(err)
	}
}
