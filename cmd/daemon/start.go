package main

import (
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildStartCmd builds a sub-command to starts the Pactus blockchain node.
func buildStartCmd(parentCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Pactus blockchain node",
	}

	parentCmd.AddCommand(startCmd)

	workingDirOpt := startCmd.Flags().StringP("working-dir", "w",
		cmd.PactusHomeDir(), "A path to the working directory to read the wallet and node files")

	passwordOpt := startCmd.Flags().StringP("password", "p", "", "The wallet password")

	pprofOpt := startCmd.Flags().String("pprof", "", "debug pprof server address(not recommended to expose to internet)")

	startCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		// change working directory
		err := os.Chdir(workingDir)
		cmd.FatalErrorCheck(err)

		if *pprofOpt != "" {
			cmd.PrintWarnMsgf("Starting Debug pprof server on: http://%s/debug/pprof/", *pprofOpt)
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
			cmd.PrintInfoMsgf("Exiting ...")
		})

		// run forever (the node will not be returned)
		select {}
	}
}
