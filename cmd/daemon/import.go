package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/spf13/cobra"
)

func buildImportCmd(parentCmd *cobra.Command) {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "download and import pruned data",
	}
	parentCmd.AddCommand(importCmd)

	workingDirOpt := addWorkingDirOption(importCmd)
	serverAddrOpt := importCmd.Flags().String("server-addr", "https://snapshot.pactus.org",
		"import server address")

	importCmd.Run = func(c *cobra.Command, _ []string) {
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

		storeDir, _ := filepath.Abs(conf.Store.StorePath())
		if !util.IsDirNotExistsOrEmpty(storeDir) {
			cmd.PrintErrorMsgf("The data directory is not empty: %s", conf.Store.StorePath())

			return
		}

		snapshotURL := *serverAddrOpt

		switch gen.ChainType() {
		case genesis.Mainnet:
			snapshotURL += "/mainnet/"
		case genesis.Testnet:
			snapshotURL += "/testnet/"
		case genesis.Localnet:
			cmd.PrintErrorMsgf("Unsupported chain type: %s", gen.ChainType())

			return
		}

		tmpDir := util.TempDirPath()

		cmd.PrintLine()

		dm := cmd.NewDownloadManager(
			snapshotURL,
			tmpDir,
			conf.Store.StorePath(),
		)

		metadata, err := dm.GetMetadata(c.Context())
		cmd.FatalErrorCheck(err)

		snapshots := make([]string, 0, len(metadata))

		for _, m := range metadata {
			if m.Data == nil {
				cmd.FatalErrorCheck(errors.New("metadata is nil"))
			}

			item := fmt.Sprintf("snapshot %s (%s)",
				dm.ParseTime(m.CreatedAt).Format("2006-01-02"),
				util.FormatBytesToHumanReadable(m.Data.Size),
			)

			snapshots = append(snapshots, item)
		}

		cmd.PrintLine()

		choice := cmd.PromptSelect("Please select a snapshot", snapshots)

		selected := metadata[choice]

		cmd.TrapSignal(func() {
			_ = fileLock.Unlock()
			_ = dm.Cleanup()
		})

		cmd.PrintLine()

		dm.Download(
			c.Context(),
			&selected,
			downloadProgressBar,
		)

		cmd.PrintLine()
		cmd.PrintLine()
		cmd.PrintInfoMsgf("Extracting files...")

		err = dm.ExtractAndStoreFiles()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("Moving data...")
		err = util.MoveDirectory(filepath.Join(tmpDir, "data"), filepath.Join(workingDir, "data"))
		cmd.FatalErrorCheck(err)

		err = dm.Cleanup()
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

func downloadProgressBar(fileName string, totalSize, downloaded int64, _ float64) {
	bar := cmd.TerminalProgressBar(totalSize, 30)
	bar.Describe(fmt.Sprintf("%s (%s/%s)",
		fileName,
		util.FormatBytesToHumanReadable(uint64(downloaded)),
		util.FormatBytesToHumanReadable(uint64(totalSize)),
	))
	err := bar.Add64(downloaded)
	cmd.FatalErrorCheck(err)
}
