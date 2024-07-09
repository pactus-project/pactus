package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/store"
	"github.com/spf13/cobra"
)

func buildPruneCmd(parentCmd *cobra.Command) {
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "prune old blocks and transactions from client",
		Long: "The prune command optimizes blockchain storage by removing outdated blocks and transactions, " +
			"freeing up disk space and enhancing client performance.",
	}
	parentCmd.AddCommand(pruneCmd)

	workingDirOpt := addWorkingDirOption(pruneCmd)

	pruneCmd.Run = func(_ *cobra.Command, _ []string) {
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

		conf, _, err := cmd.MakeConfig(workingDir)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintWarnMsgf("This command removes all the blocks and transactions up to %d days ago "+
			"and converts the node to prune mode.", conf.Store.RetentionDays)
		cmd.PrintLine()
		confirmed := cmd.PromptConfirm("Do you want to continue")
		if !confirmed {
			return
		}
		cmd.PrintLine()

		str, err := store.NewStore(conf.Store)
		cmd.FatalErrorCheck(err)

		prunedCount := uint32(0)
		skippedCount := uint32(0)
		totalCount := uint32(0)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-interrupt
			str.Close()
			_ = fileLock.Unlock()
		}()

		err = str.Prune(func(pruned, skipped, pruningHeight uint32) {
			prunedCount += pruned
			skippedCount += skipped

			if totalCount == 0 {
				totalCount = pruningHeight
			}

			pruningProgressBar(prunedCount, skippedCount, totalCount)
		})
		cmd.PrintLine()
		cmd.FatalErrorCheck(err)

		str.Close()
		_ = fileLock.Unlock()

		cmd.PrintLine()
		cmd.PrintInfoMsgf("✅ Your node successfully pruned and changed to prune mode.")
		cmd.PrintLine()
		cmd.PrintInfoMsgf("You can start the node by running this command:")
		cmd.PrintInfoMsgf("./pactus-daemon start -w %v", workingDir)
	}
}

func pruningProgressBar(prunedCount, skippedCount, totalCount uint32) {
	percentage := float64(prunedCount+skippedCount) / float64(totalCount) * 100
	if percentage > 100 {
		percentage = 100
	}

	barLength := 40
	filledLength := int(float64(barLength) * percentage / 100)

	bar := strings.Repeat("=", filledLength) + strings.Repeat(" ", barLength-filledLength)
	fmt.Printf("\r [%s] %.0f%% Pruned: %d | Skipped: %d", //nolint
		bar, percentage, prunedCount, skippedCount)
}
