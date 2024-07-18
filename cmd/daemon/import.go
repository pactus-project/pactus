package main

import (
	"fmt"
	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func buildImportCmd(parentCmd *cobra.Command) {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "download and import pruned data",
	}
	parentCmd.AddCommand(importCmd)

	workingDirOpt := addWorkingDirOption(importCmd)

	importCmd.Run = func(c *cobra.Command, args []string) {
		workingDir, err := filepath.Abs(*workingDirOpt)
		cmd.FatalErrorCheck(err)

		err = os.Chdir(workingDir)
		cmd.FatalErrorCheck(err)

		conf, gen, err := cmd.MakeConfig(workingDir)
		cmd.FatalErrorCheck(err)

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

		cmd.PrintInfoMsgf("Checking data in %s exists...", workingDir)
		storeDir, _ := filepath.Abs(conf.Store.StorePath())
		if !util.IsDirNotExistsOrEmpty(storeDir) {
			cmd.PrintErrorMsgf("The data directory is not empty: %s", conf.Store.StorePath())

			return
		}

		snapshotUrl := ""

		switch gen.ChainType() {
		case genesis.Mainnet:
			snapshotUrl = cmd.SnapshotBaseUrl + "mainnet/"
		case genesis.Testnet:
			snapshotUrl = cmd.SnapshotBaseUrl + "testnet/"
		default:
			cmd.PrintErrorMsgf("Unsupported chain type: %s", gen.ChainType())

			return
		}

		cmd.PrintInfoMsgf("Getting snapshots metadata...")
		metadata, err := cmd.SnapshotMetadata(c.Context(), snapshotUrl)
		if err != nil {
			cmd.PrintErrorMsgf("Failed to get snapshot metadata: %s", err)

			return
		}

		sort.Slice(metadata, func(i, j int) bool {
			return i > j
		})

		cmd.PrintLine()

		for i, m := range metadata {
			fmt.Printf("%d. snapshot %s (%.2f MB)\n", i+1,
				m.CreatedAt,
				float64(m.TotalSize)/1024/1024)
		}

		cmd.PrintLine()

		var choice int
		fmt.Printf("Please select a snapshot [1-%d]: ", len(metadata))
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			cmd.PrintErrorMsgf("invalid input: %s", err)

			return
		}

		if choice < 1 || choice > len(metadata) {
			cmd.PrintErrorMsgf("Invalid choice.")

			return
		}

		selected := metadata[choice-1]
		tmpDir := util.TempDirPath()
		extractPath := fmt.Sprintf("%s/extracted", tmpDir)

		err = os.MkdirAll(extractPath, os.ModePerm)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()

		zipFileList := make([]string, 0)

		cmd.DownloadManager(
			c.Context(),
			selected,
			snapshotUrl,
			tmpDir,
			zipFileList,
			downloadProgressBar,
		)

		for _, zFile := range zipFileList {
			err := cmd.ExtractAndStoreFile(zFile, extractPath)
			cmd.FatalErrorCheck(err)
		}

		err = os.MkdirAll(filepath.Dir(conf.Store.StorePath()), os.ModePerm)
		cmd.FatalErrorCheck(err)

		err = cmd.CopyAllFiles(extractPath, conf.Store.StorePath())
		cmd.FatalErrorCheck(err)

		err = os.RemoveAll(tmpDir)
		cmd.FatalErrorCheck(err)

		_ = fileLock.Unlock()

	}
}

func downloadProgressBar(fileName string, totalSize, downloaded int64, percentage float64) {
	barWidth := 30
	completedWidth := int(float64(barWidth) * (percentage / 100))

	progressBar := fmt.Sprintf(
		"\r[%-*s] %3.0f%% %s (%.2f MB/ %.2f MB)",
		barWidth,
		strings.Repeat("=", completedWidth),
		percentage,
		fileName,
		float64(downloaded)/1024/1024,
		float64(totalSize)/1024/1024,
	)

	fmt.Print(progressBar)
}
