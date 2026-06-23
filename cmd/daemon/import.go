package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ezex-io/gopkg/signal"
	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/downloader"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/schollz/progressbar/v3"
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
		terminal.FatalErrorCheck(err)

		err = os.Chdir(workingDir)
		terminal.FatalErrorCheck(err)

		conf, gen, err := cmd.MakeConfig(workingDir)
		terminal.FatalErrorCheck(err)

		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock := flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		terminal.FatalErrorCheck(err)

		if !locked {
			terminal.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}

		terminal.PrintLine()

		snapshotURL := *serverAddrOpt

		importer, err := cmd.NewImporter(
			gen.ChainType(),
			snapshotURL,
			conf.Store.DataPath(),
		)
		terminal.FatalErrorCheck(err)

		metadata, err := importer.GetMetadata(cobra.Context())
		terminal.FatalErrorCheck(err)

		snapshots := make([]string, 0, len(metadata))

		for _, md := range metadata {
			item := fmt.Sprintf(
				"snapshot %s (%s)",
				md.CreatedAtTime().Format("2006-01-02"),
				util.FormatBytesToHumanReadable(md.Data.Size),
			)

			snapshots = append(snapshots, item)
		}

		terminal.PrintLine()

		choice := prompt.PromptSelect("Please select a snapshot", snapshots)

		selected := metadata[choice]

		signal.HandleInterrupt(func() {
			_ = fileLock.Unlock()
			_ = importer.Cleanup()
		})

		terminal.PrintLine()

		bar := terminal.ProgressBar(int64(selected.Data.Size), 30)
		err = importer.Download(cobra.Context(), &selected,
			func(stats downloader.Stats) {
				updateProgressBar(bar, selected.Data.Name, stats)
			})
		terminal.PrintLine()
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintLine()
		terminal.PrintInfoMsgf("📦 Extracting snapshot files...")

		err = importer.ExtractAndStoreFiles()
		terminal.FatalErrorCheck(err)

		terminal.PrintInfoMsgf("📁 Moving data to node directory...")
		err = importer.MoveStore()
		terminal.FatalErrorCheck(err)

		terminal.PrintInfoMsgf("🧹 Cleaning up temporary files...")
		err = importer.Cleanup()
		terminal.FatalErrorCheck(err)

		_ = fileLock.Unlock()

		terminal.PrintLine()
		terminal.PrintSuccessMsgf("✅ Node successfully imported pruned data!")
		terminal.PrintLine()
		terminal.PrintInfoMsgf("🚀 To start your node, run:")
		terminal.PrintInfoMsgBoldf("   %s start -w %s", cmd.PactusDaemonName(), workingDir)
	}
}

func updateProgressBar(bar *progressbar.ProgressBar, fileName string, stats downloader.Stats) {
	bar.Describe(fmt.Sprintf(
		"%s (%s/%s)",
		fileName,
		util.FormatBytesToHumanReadable(uint64(stats.Downloaded)),
		util.FormatBytesToHumanReadable(uint64(stats.TotalSize)),
	))
	// Ignore progress bar errors
	_ = bar.Set64(stats.Downloaded)

	if stats.Completed {
		_ = bar.Finish()
	}
}
