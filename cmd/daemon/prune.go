package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/signal"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/spf13/cobra"
)

func buildPruneCmd(parentCmd *cobra.Command) {
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "prune old blocks and transactions from the node",
		Long: "The prune command optimizes blockchain storage by removing outdated blocks and transactions, " +
			"freeing up disk space and enhancing node performance.",
	}
	parentCmd.AddCommand(pruneCmd)

	workingDirOpt := addWorkingDirOption(pruneCmd)

	pruneCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		// change working directory
		err := os.Chdir(workingDir)
		terminal.FatalErrorCheck(err)

		// Define the lock file path
		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock := flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		terminal.FatalErrorCheck(err)

		if !locked {
			terminal.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}

		conf, _, err := cmd.MakeConfig(workingDir)
		terminal.FatalErrorCheck(err)

		// Disable logger
		conf.Logger.Targets = []string{}
		logger.InitGlobalLogger(conf.Logger)

		terminal.PrintLine()
		terminal.PrintWarnMsgf("⚠️  PRUNE OPERATION WARNING")
		terminal.PrintWarnMsgf("   This will remove all blocks and transactions older than %d days", conf.Store.RetentionDays)
		terminal.PrintWarnMsgf("   and convert your node to prune mode.")
		terminal.PrintWarnMsgf("   This action cannot be undone.")
		terminal.PrintLine()
		confirmed := prompt.PromptConfirm("Are you sure you want to continue")
		if !confirmed {
			return
		}
		terminal.PrintLine()

		store, err := store.NewStore(conf.Store)
		terminal.FatalErrorCheck(err)

		prunedCount := uint32(0)
		skippedCount := uint32(0)
		totalCount := uint32(0)
		canceled := false
		closed := make(chan bool, 1)

		signal.HandleInterrupt(func() {
			canceled = true
			<-closed
		})

		err = store.Prune(func(pruned bool, pruningHeight uint32) bool {
			if pruned {
				prunedCount++
			} else {
				skippedCount++
			}

			if totalCount == 0 {
				totalCount = pruningHeight
			}

			pruningProgressBar(prunedCount, skippedCount, totalCount)

			return canceled
		})
		terminal.PrintLine()
		terminal.FatalErrorCheck(err)

		if canceled {
			terminal.PrintLine()
			terminal.PrintWarnMsgf("❌ Prune operation cancelled.")
			terminal.PrintLine()
		} else if prunedCount == 0 {
			terminal.PrintLine()
			terminal.PrintInfoMsgf("ℹ️  No pruning needed")
			terminal.PrintInfoMsgf("   Your node has not reached the retention period (%d days)", conf.Store.RetentionDays)
			terminal.PrintInfoMsgf("   or is already in prune mode.")
			terminal.PrintLine()
		} else {
			terminal.PrintLine()
			terminal.PrintSuccessMsgf("✅ Your node successfully pruned and changed to prune mode.")
			terminal.PrintLine()
			terminal.PrintInfoMsgf("🚀 To start your node, run:")
			terminal.PrintInfoMsgBoldf("   %s start -w %s", cmd.PactusDaemonName(), workingDir)
		}

		store.Close()
		_ = fileLock.Unlock()

		closed <- true
	}
}

func pruningProgressBar(prunedCount, skippedCount, totalCount uint32) {
	if (prunedCount+skippedCount)%1000 != 0 {
		return
	}

	bar := terminal.ProgressBar(int64(totalCount), 30)
	bar.Describe(fmt.Sprintf("Pruned: %d | Skipped: %d", prunedCount, skippedCount))
	err := bar.Add(int(prunedCount + skippedCount))
	terminal.FatalErrorCheck(err)
}
