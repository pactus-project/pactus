package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
)

// Start starts the pactus node.
func Start() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		workingDirOpt := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "A path to the working directory to read the wallet and node files",
			Value: cmd.PactusHomeDir(),
		})
		passwordOpt := c.String(cli.StringOpt{
			Name: "p password",
			Desc: "The wallet password",
		})
		pprofOpt := c.String(cli.StringOpt{
			Name: "pprof",
			Desc: "debug pprof server address(not recommended to expose to internet)",
		})

		c.LongDesc = "Starting the node from working directory"
		c.Before = func() { fmt.Println(cmd.Pactus) }
		c.Action = func() {
			workingDir, _ := filepath.Abs(*workingDirOpt)
			// change working directory
			err := os.Chdir(workingDir)
			cmd.FatalErrorCheck(err)

			if *pprofOpt != "" {
				cmd.PrintWarnMsg("Starting Debug pprof server on: http://%s/debug/pprof/\n", *pprofOpt)
				server := &http.Server{
					Addr:              *pprofOpt,
					ReadHeaderTimeout: 3 * time.Second,
				}
				go func() {
					err := server.ListenAndServe()
					cmd.FatalErrorCheck(err)
				}()
			}

			passwordFetcher := func(wallet *wallet.Wallet) (string, bool) {
				if !wallet.IsEncrypted() {
					return "", true
				}

				var password string
				if *passwordOpt != "" {
					password = *passwordOpt
				} else {
					password = cmd.PromptPassword("Wallet password", false)
				}
				return password, true
			}
			node, _, err := cmd.StartNode(
				workingDir, passwordFetcher)
			cmd.FatalErrorCheck(err)

			cmd.TrapSignal(func() {
				node.Stop()
				cmd.PrintInfoMsg("Exiting ...")
			})

			// run forever (the node will not be returned)
			select {}
		}
	}
}
