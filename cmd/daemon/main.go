package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("pactus-daemon", "Pactus daemon")

	app.Command("init", "Initialize the pactus blockchain", Init())
	app.Command("start", "Start the pactus blockchain", Start())
	app.Command("version", "Print the pactus version", Version())

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
