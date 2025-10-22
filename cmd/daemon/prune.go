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
		terminal.PrintWarnMsgf("This command removes all the blocks and transactions up to %d days ago "+
			"and converts the node to prune mode.", conf.Store.RetentionDays)
		terminal.PrintLine()
		confirmed := prompt.PromptConfirm("Do you want to continue")
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
			terminal.PrintInfoMsgf("❌ The operation canceled.")
			terminal.PrintLine()
		} else if prunedCount == 0 {
			terminal.PrintLine()
			terminal.PrintInfoMsgf("⚠️ Your node is not passed the retention_days set in config or it's already a pruned node.")
			terminal.PrintLine()
			terminal.PrintInfoMsgf("Make sure you try to prune a node after retention_days specified in config.toml")
		} else {
			terminal.PrintLine()
			terminal.PrintInfoMsgf("✅ Your node successfully pruned and changed to prune mode.")
			terminal.PrintLine()
			terminal.PrintInfoMsgf("You can start the node by running this command:")
			terminal.PrintInfoMsgf("./pactus-daemon start -w %v", workingDir)
		}

		store.Close()
		_ = fileLock.Unlock()

		closed <- true
	}
}

func pruningProgressBar(prunedCount, skippedCount, totalCount uint32) {
	bar := terminal.ProgressBar(int64(totalCount), 30)
	bar.Describe(fmt.Sprintf("Pruned: %d | Skipped: %d", prunedCount, skippedCount))
	err := bar.Add(int(prunedCount + skippedCount))
	terminal.FatalErrorCheck(err)
}
