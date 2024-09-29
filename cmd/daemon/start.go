package main

import (
	"os"
	"path/filepath"

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

	workingDirOpt := addWorkingDirOption(startCmd)

	passwordOpt := startCmd.Flags().StringP("password", "p", "",
		"the wallet password")

	debuggerOpt := startCmd.Flags().BoolP("debug", "d", false, "enable pprof debugger")

	startCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		// change working directory
		err := os.Chdir(workingDir)
		cmd.FatalErrorCheck(err)

		// Define the lock file path
		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock := flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		cmd.FatalErrorCheck(err)

		if !locked {
			cmd.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
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
			workingDir, passwordFetcher, *debuggerOpt)
		cmd.FatalErrorCheck(err)

		cmd.TrapSignal(func() {
			cmd.PrintInfoMsgf("Exiting...")

			_ = fileLock.Unlock()
			node.Stop()
		})

		// run forever (the node will not be returned)
		select {}
	}
}
