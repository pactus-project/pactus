package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/downloader"
	"github.com/spf13/cobra"
)

func buildImportCmd(parentCmd *cobra.Command) {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "download and import pruned data",
	}
	parentCmd.AddCommand(importCmd)

	workingDirOpt := addWorkingDirOption(importCmd)
	serverAddrOpt := importCmd.Flags().String("server-addr", cmd.DefaultSnapshotURL,
		"import server address")

	importCmd.Run = func(cobra *cobra.Command, _ []string) {
		workingDir, err := filepath.Abs(*workingDirOpt)
		cmd.FatalErrorCheck(err)

		err = os.Chdir(workingDir)
		cmd.FatalErrorCheck(err)

		conf, gen, err := cmd.MakeConfig(workingDir)
		cmd.FatalErrorCheck(err)

		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock := flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		cmd.FatalErrorCheck(err)

		if !locked {
			cmd.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}

		cmd.PrintLine()

		snapshotURL := *serverAddrOpt
		importer, err := cmd.NewImporter(
			gen.ChainType(),
			snapshotURL,
			conf.Store.DataPath(),
		)
		cmd.FatalErrorCheck(err)

		metadata, err := importer.GetMetadata(cobra.Context())
		cmd.FatalErrorCheck(err)

		snapshots := make([]string, 0, len(metadata))

		for _, md := range metadata {
			item := fmt.Sprintf("snapshot %s (%s)",
				md.CreatedAtTime().Format("2006-01-02"),
				util.FormatBytesToHumanReadable(md.Data.Size),
			)

			snapshots = append(snapshots, item)
		}

		cmd.PrintLine()

		choice := cmd.PromptSelect("Please select a snapshot", snapshots)

		selected := metadata[choice]

		cmd.TrapSignal(func() {
			_ = fileLock.Unlock()
			_ = importer.Cleanup()
		})

		cmd.PrintLine()

		err = importer.Download(cobra.Context(), &selected, downloadProgressBar)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintLine()
		cmd.PrintInfoMsgf("Extracting files...")

		err = importer.ExtractAndStoreFiles()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("Moving data...")
		err = importer.MoveStore()
		cmd.FatalErrorCheck(err)

		err = importer.Cleanup()
		cmd.FatalErrorCheck(err)

		_ = fileLock.Unlock()

		cmd.PrintLine()
		cmd.PrintLine()
		cmd.PrintInfoMsgf("✅ Your node successfully imported prune data.")
		cmd.PrintLine()
		cmd.PrintInfoMsgf("You can start the node by running this command:")
		cmd.PrintInfoMsgf("./pactus-daemon start -w %v", workingDir)
	}
}

func downloadProgressBar(fileName string) func(stats downloader.Stats) {
	return func(stats downloader.Stats) {
		if !stats.Completed {
			bar := cmd.TerminalProgressBar(stats.TotalSize, 30)
			bar.Describe(fmt.Sprintf("%s (%s/%s)",
				fileName,
				util.FormatBytesToHumanReadable(uint64(stats.Downloaded)),
				util.FormatBytesToHumanReadable(uint64(stats.TotalSize)),
			))
			// Ignore progress bar errors
			_ = bar.Add64(stats.Downloaded)
		}
	}
}
