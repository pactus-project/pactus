package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("pactus-daemon", "Pactus daemon")

	app.Command("init", "Initialize the Pactus blockchain", Init())
	app.Command("start", "Start the Pactus blockchain", Start())
	app.Command("version", "Print the Pactus version", Version())

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
