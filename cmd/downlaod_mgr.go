package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pactus-project/pactus/util/downloader"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	SnapshotBaseUrl = "https://data.pacviewer.com/files/"
)

type Metadata struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Compress  string `json:"compress"`
	TotalSize int    `json:"total_size"`
	Data      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Sha  string `json:"sha"`
		Size int    `json:"size"`
	} `json:"data"`
}

func SnapshotMetadata(ctx context.Context, snapshotUrl string) ([]Metadata, error) {
	cli := http.DefaultClient
	metaUrl, err := url.JoinPath(snapshotUrl, "metadata.json")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metaUrl, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	metadata := make([]Metadata, 0)
	return metadata, json.NewDecoder(resp.Body).Decode(&metadata)
}

func DownloadManager(
	ctx context.Context,
	metadata Metadata,
	baseUrl, tempPath string,
	zipFileListPath []string,
	stateFunc func(fileName string, totalSize, downloaded int64, percentage float64),
) {

	for _, data := range metadata.Data {
		done := make(chan struct{})
		dlLink, err := url.JoinPath(baseUrl, data.Path)
		FatalErrorCheck(err)

		fileName := filepath.Base(dlLink)

		filePath := fmt.Sprintf("%s/%s", tempPath, fileName)

		dl := downloader.New(
			dlLink,
			filePath,
			data.Sha,
		)

		dl.Start(ctx)

		go func() {
			for state := range dl.Stats() {
				stateFunc(fileName, state.TotalSize, state.Downloaded, state.Percent)
				if state.Completed {
					done <- struct{}{}
					close(done)
					return
				}
			}
		}()

		<-done
		zipFileListPath = append(zipFileListPath, filePath)
	}
}

func ExtractAndStoreFile(zipFilePath, extractPath string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip archive: %v", err)
		}

		fpath := filepath.Join(extractPath, f.Name)

		outFile, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to copy file contents: %v", err)
		}

		_ = rc.Close()
		_ = outFile.Close()
	}

	return nil
}

// CopyAllFiles copies all files from srcDir to dstDir
func CopyAllFiles(srcDir, dstDir string) error {
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		relativePath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relativePath)

		err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		err = copyFile(path, dstPath)
		if err != nil {
			return fmt.Errorf("failed to copy file from %s to %s: %v", path, dstPath, err)
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %v", err)
	}

	return nil
}
