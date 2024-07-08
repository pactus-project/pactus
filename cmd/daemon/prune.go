package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/util"
	"github.com/spf13/cobra"
)

func buildPruneCmd(parentCmd *cobra.Command) {
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "prune old blocks and transactions from client",
		Long: "The prune command optimizes blockchain storage by removing outdated blocks and transactions, " +
			"freeing up disk space and enhancing client performance. customize pruning criteria via flags to " +
			"manage storage effectively.",
	}
	parentCmd.AddCommand(pruneCmd)

	workingDirOpt := pruneCmd.Flags().StringP("working-dir", "w", cmd.PactusDefaultHomeDir(),
		"a path to the working directory of node files")

	pruneCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		if util.IsDirNotExistsOrEmpty(workingDir) {
			cmd.PrintErrorMsgf("The working directory is not exists: %s", workingDir)

			return
		}

		conf, _, err := cmd.MakeConfig(workingDir)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintWarnMsgf("Warning: This command removes all your blocks and transactions and changes " +
			"your node to a prune node. You cannot revert to a full node.")
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
		}()

		err = str.Prune(func(pruned, skipped, pruningHeight uint32) {
			prunedCount += pruned
			skippedCount += skipped

			if totalCount == 0 {
				totalCount = pruningHeight
			}

			pruningProgressBar(prunedCount, skippedCount, pruningHeight, totalCount)
		})
		cmd.PrintLine()
		cmd.FatalErrorCheck(err)

		str.Close()

		cmd.PrintLine()
		cmd.PrintInfoMsgf("✅ Your node successfully pruned and changed to prune client.")
		cmd.PrintLine()
		cmd.PrintInfoMsgf("You can start the node by running this command:")
		cmd.PrintInfoMsgf("./pactus-daemon start -w %v", workingDir)
	}
}

func pruningProgressBar(prunedCount, skippedCount, leftBlock, totalCount uint32) {
	percentage := float64(prunedCount+skippedCount) / float64(totalCount) * 100
	if percentage > 100 {
		percentage = 100
	}

	barLength := 40
	filledLength := int(float64(barLength) * percentage / 100)

	bar := strings.Repeat("=", filledLength) + strings.Repeat(" ", barLength-filledLength)
	fmt.Printf("\r [%s] %.0f%% Pruned: %d | Skipped: %d | Left Block: %d", //nolint
		bar, percentage, prunedCount, skippedCount, leftBlock)
}
