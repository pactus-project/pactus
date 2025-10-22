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
		cmd.FatalErrorCheck(err)

		if !locked {
			cmd.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}

		conf, _, err := cmd.MakeConfig(workingDir)
		cmd.FatalErrorCheck(err)

		// Disable logger
		conf.Logger.Targets = []string{}
		logger.InitGlobalLogger(conf.Logger)

		cmd.PrintLine()
		cmd.PrintWarnMsgf("This command removes all the blocks and transactions up to %d days ago "+
			"and converts the node to prune mode.", conf.Store.RetentionDays)
		cmd.PrintLine()
		confirmed := prompt.PromptConfirm("Do you want to continue")
		if !confirmed {
			return
		}
		cmd.PrintLine()

		store, err := store.NewStore(conf.Store)
		cmd.FatalErrorCheck(err)

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
		cmd.PrintLine()
		cmd.FatalErrorCheck(err)

		if canceled {
			cmd.PrintLine()
			cmd.PrintInfoMsgf("❌ The operation canceled.")
			cmd.PrintLine()
		} else if prunedCount == 0 {
			cmd.PrintLine()
			cmd.PrintInfoMsgf("⚠️ Your node is not passed the retention_days set in config or it's already a pruned node.")
			cmd.PrintLine()
			cmd.PrintInfoMsgf("Make sure you try to prune a node after retention_days specified in config.toml")
		} else {
			cmd.PrintLine()
			cmd.PrintInfoMsgf("✅ Your node successfully pruned and changed to prune mode.")
			cmd.PrintLine()
			cmd.PrintInfoMsgf("You can start the node by running this command:")
			cmd.PrintInfoMsgf("./pactus-daemon start -w %v", workingDir)
		}

		store.Close()
		_ = fileLock.Unlock()

		closed <- true
	}
}

func pruningProgressBar(prunedCount, skippedCount, totalCount uint32) {
	bar := cmd.TerminalProgressBar(int64(totalCount), 30)
	bar.Describe(fmt.Sprintf("Pruned: %d | Skipped: %d", prunedCount, skippedCount))
	err := bar.Add(int(prunedCount + skippedCount))
	cmd.FatalErrorCheck(err)
}
