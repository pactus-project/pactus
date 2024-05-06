package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util"
	"github.com/spf13/cobra"
)

type metadata struct {
	FileName string `json:"file_name"`
	Path     string `json:"path"`
	Hash     string `json:"hash"`
}

func buildExportCmd(root *cobra.Command) {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "make snapshot from last state of blockchain data",
	}
	root.AddCommand(exportCmd)
	workingDirOpt := exportCmd.Flags().StringP("working-dir", "w", cmd.PactusDefaultHomeDir(),
		"a path to the working directory node files")

	exportCmd.Run = func(_ *cobra.Command, _ []string) {
		log.Println("checking working directory...")

		workingDir, _ := filepath.Abs(*workingDirOpt)
		if util.IsDirNotExistsOrEmpty(workingDir) {
			cmd.PrintErrorMsgf("The working directory is empty: %s", workingDir)

			return
		}

		err := os.Chdir(workingDir)
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

		log.Println("work directory is ready for export.")

		srcDir := filepath.Join(workingDir, "data")

		log.Printf("checking old snapshot exists to remove...")

		snapshotDir, err := filepath.Abs(filepath.Join(workingDir, "snapshot"))
		if err != nil {
			cmd.FatalErrorCheck(err)
		}

		if !util.IsDirNotExistsOrEmpty(snapshotDir) {
			if err := os.RemoveAll(snapshotDir); err != nil {
				cmd.FatalErrorCheck(err)
			}
		}

		log.Println("creating snapshot from blockchain data...")

		md := make([]metadata, 0)

		if err := copyDir(srcDir, snapshotDir, &md); err != nil {
			cmd.PrintErrorMsgf(err.Error())

			return
		}

		log.Println("snapshot created")
		log.Println("create snapshot metadata...")

		mdFile, err := os.Create(snapshotDir + "/metadata.json")
		defer func() {
			if err := mdFile.Close(); err != nil {
				panic(err)
			}
		}()

		if err != nil {
			cmd.FatalErrorCheck(err)
		}

		b, err := json.MarshalIndent(&md, "", " ")
		if err != nil {
			cmd.PrintWarnMsgf("Failed to encode metadata got error %s", err.Error())

			return
		}

		_, err = mdFile.Write(b)
		if err != nil {
			cmd.PrintWarnMsgf("Failed to write metadata got error %s", err.Error())

			return
		}

		log.Println("Exporting Pactus blockchain data has been completed successfully.")
	}
}

func copyFile(base, src, dest string, md *[]metadata) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	currentFilePath, err := filepath.Rel(base, dest)
	if err != nil {
		return err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, destFile); err != nil {
		return err
	}

	sha256Hash := hex.EncodeToString(hash.Sum(nil))
	filename := filepath.Base(dest)

	*md = append(*md, metadata{
		FileName: filename,
		Hash:     sha256Hash,
		Path:     currentFilePath,
	})

	log.Printf("copying %s...\n", src)

	return nil
}

func copyDir(src, dest string, metadata *[]metadata) error {
	baseDir := filepath.Dir(dest)
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			if err := copyDir(srcPath, destPath, metadata); err != nil {
				return err
			}
		} else {
			if err := copyFile(baseDir, srcPath, destPath, metadata); err != nil {
				return err
			}
		}
	}

	return nil
}
