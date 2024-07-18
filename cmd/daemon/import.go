package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	serverAddrOpt := importCmd.Flags().String("server-addr", "", "custom import server")

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

		serverAddr := cmd.SnapshotServer()

		if *serverAddrOpt != "" {
			serverAddr = *serverAddrOpt
		}

		snapshotURL := "mainnet/"

		switch gen.ChainType() {
		case genesis.Mainnet:
			snapshotURL = serverAddr + "mainnet/"
		case genesis.Testnet:
			snapshotURL = serverAddr + "testnet/"
		case genesis.Localnet:
			cmd.PrintErrorMsgf("Unsupported chain type: %s", gen.ChainType())

			return
		}

		metadata, err := cmd.GetSnapshotMetadata(c.Context(), snapshotURL)
		if err != nil {
			cmd.PrintErrorMsgf("Failed to get snapshot metadata: %s", err)

			return
		}

		sort.Slice(metadata, func(i, j int) bool {
			return i > j //nolint
		})

		cmd.PrintLine()

		for i, m := range metadata {
			fmt.Printf("%d. snapshot %s (%s)\n", //nolint
				i+1,
				parseDate(m.CreatedAt),
				util.FormatBytesToHumanReadable(uint64(m.TotalSize)))
		}

		cmd.PrintLine()

		var choice int
		fmt.Printf("Please select a snapshot [1-%d]: ", len(metadata)) //nolint
		_, err = fmt.Scanf("%d", &choice)
		cmd.FatalErrorCheck(err)

		if choice < 1 || choice > len(metadata) {
			cmd.PrintErrorMsgf("Invalid choice.")

			return
		}

		selected := metadata[choice-1]
		tmpDir := util.TempDirPath()
		extractPath := fmt.Sprintf("%s/data", tmpDir)

		err = os.MkdirAll(extractPath, 0o750)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()

		zipFileList := cmd.DownloadManager(
			c.Context(),
			&selected,
			snapshotURL,
			tmpDir,
			downloadProgressBar,
		)

		for _, zFile := range zipFileList {
			err := cmd.ExtractAndStoreFile(zFile, extractPath)
			cmd.FatalErrorCheck(err)
		}

		err = os.MkdirAll(filepath.Dir(conf.Store.StorePath()), 0o750)
		cmd.FatalErrorCheck(err)

		err = cmd.CopyAllFiles(extractPath, conf.Store.StorePath())
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

func downloadProgressBar(fileName string, totalSize, downloaded int64, percentage float64) {
	barWidth := 30
	completedWidth := int(float64(barWidth) * (percentage / 100))

	progressBar := fmt.Sprintf(
		"\r[%-*s] %3.0f%% %s (%s/ %s)",
		barWidth,
		strings.Repeat("=", completedWidth),
		percentage,
		fileName,
		util.FormatBytesToHumanReadable(uint64(downloaded)),
		util.FormatBytesToHumanReadable(uint64(totalSize)),
	)

	fmt.Print(progressBar) //nolint
}

func parseDate(dateString string) string {
	const layout = "2006-01-02T15:04:05.000000"

	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		return ""
	}

	return parsedTime.Format("2006-01-02")
}
