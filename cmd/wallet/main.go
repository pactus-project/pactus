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
		Desc:  "A path to the wallet file",
		Value: cmd.ZarbWalletsDir() + "default_wallet",
	})

	app.Command("generate", "Generate a new key", Generate())
	app.Command("recover", "Recover waller from mnemonic (seed phrase) a new key", Recover())

	app.Run(os.Args)
}
