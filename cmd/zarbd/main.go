package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func zarbd() *cli.Cli {
	app := cli.App("zarbd", "Zarb blockchain node(daemon)")

	app.Command("init", "Initialize the zarb blockchain", Init())
	app.Command("start", "Start the zarb blockchain", Start())
	app.Command("version", "Print the zarb version", Version())

	return app
}

func main() {
	if err := zarbd().Run(os.Args); err != nil {
		panic(err)
	}
}
