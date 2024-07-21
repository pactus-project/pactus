package main

import (
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
	serverAddrOpt := importCmd.Flags().String("server-addr", "https://download.pactus.org",
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
		extractPath := fmt.Sprintf("%s/data", tmpDir)

		err = os.MkdirAll(extractPath, 0o750)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()

		dl := cmd.NewDownloadManager(
			snapshotURL,
			extractPath,
			tmpDir,
			conf.Store.StorePath(),
		)

		metadata, err := dl.GetMetadata(c.Context())
		cmd.FatalErrorCheck(err)

		snapshots := make([]string, 0, len(metadata))

		for _, m := range metadata {
			item := fmt.Sprintf("snapshot %s (%s)",
				dl.ParseTime(m.CreatedAt).Format("2006-01-02"),
				util.FormatBytesToHumanReadable(m.TotalSize),
			)

			snapshots = append(snapshots, item)
		}

		cmd.PrintLine()

		choice := cmd.PromptSelect("Please select a snapshot", snapshots)

		selected := metadata[choice]

		cmd.PrintLine()

		dl.Download(
			c.Context(),
			&selected,
			downloadProgressBar,
		)

		err = dl.ExtractAndStoreFiles()
		cmd.FatalErrorCheck(err)

		err = os.MkdirAll(filepath.Dir(conf.Store.StorePath()), 0o750)
		cmd.FatalErrorCheck(err)

		err = dl.CopyAllFiles()
		cmd.FatalErrorCheck(err)

		err = os.RemoveAll(tmpDir)
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

func downloadProgressBar(fileName string, totalSize, downloaded int64, totalItem, downloadedItem int, _ float64) {
	bar := cmd.TerminalProgressBar(int(totalSize), 30, true)
	bar.Describe(fmt.Sprintf("%s (%d/%d)", fileName, downloadedItem, totalItem))
	err := bar.Add(int(downloaded))
	cmd.FatalErrorCheck(err)
}