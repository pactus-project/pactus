package main

import (
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildStartCmd builds a sub-command to starts the Pactus blockchain node.
func buildStartCmd(parentCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the Pactus Blockchain node",
	}

	parentCmd.AddCommand(startCmd)

	workingDirOpt := startCmd.Flags().StringP("working-dir", "w", cmd.PactusDefaultHomeDir(),
		"the path to the working directory to load the wallet and node files")

	passwordOpt := startCmd.Flags().StringP("password", "p", "",
		"the wallet password")

	pprofOpt := startCmd.Flags().String("pprof", "",
		"pprof server address (for debugging)")

	startCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		// change working directory
		err := os.Chdir(workingDir)
		cmd.FatalErrorCheck(err)

		// Define the lock file path
		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock := flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		if err != nil {
			// handle unable to attempt to acquire lock
			cmd.FatalErrorCheck(err)
		}

		if !locked {
			cmd.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}

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

		passwordFetcher := func(wlt *wallet.Wallet) (string, bool) {
			if !wlt.IsEncrypted() {
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
			_ = fileLock.Unlock()
			node.Stop()
			cmd.PrintInfoMsgf("Exiting ...")
		})

		// run forever (the node will not be returned)
		select {}
	}
}
