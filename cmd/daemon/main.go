package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("zarb-daemon", "Zarb daemon")

	app.Command("init", "Initialize the zarb blockchain", Init())
	app.Command("start", "Start the zarb blockchain", Start())
	app.Command("version", "Print the zarb version", Version())

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
