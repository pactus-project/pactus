package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pactus-project/pactus/util/downloader"
)

const maxDecompressedSize = 10 << 20 // 10 MB

type Metadata struct {
	Name      string          `json:"name"`
	CreatedAt string          `json:"created_at"`
	Compress  string          `json:"compress"`
	TotalSize uint64          `json:"total_size"`
	Data      []*SnapshotData `json:"data"`
}

type SnapshotData struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Sha  string `json:"sha"`
	Size uint64 `json:"size"`
}

func GetSnapshotMetadata(ctx context.Context, snapshotURL string) ([]Metadata, error) {
	cli := http.DefaultClient
	metaURL, err := url.JoinPath(snapshotURL, "metadata.json")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metaURL, http.NoBody)
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

	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(&metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func DownloadManager(
	ctx context.Context,
	metadata *Metadata,
	baseURL, tempPath string,
	stateFunc func(fileName string, totalSize, downloaded int64, percentage float64),
) []string {
	zipFileListPath := make([]string, 0)

	for _, data := range metadata.Data {
		done := make(chan struct{})
		dlLink, err := url.JoinPath(baseURL, data.Path)
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
			err := <-dl.Errors()
			FatalErrorCheck(err)
		}()

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

	return zipFileListPath
}

func ExtractAndStoreFile(zipFilePath, extractPath string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip archive: %w", err)
		}

		fpath := fmt.Sprintf("%s/%s", extractPath, f.Name)

		outFile, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		// fixed potential DoS vulnerability via decompression bomb
		lr := io.LimitedReader{R: rc, N: maxDecompressedSize}
		_, err = io.Copy(outFile, &lr)
		if err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		// check if the file size exceeds the limit
		if lr.N <= 0 {
			return fmt.Errorf("file exceeds maximum decompressed size limit: %s", fpath)
		}

		_ = rc.Close()
		_ = outFile.Close()
	}

	return nil
}

// CopyAllFiles copies all files from srcDir to dstDir.
func CopyAllFiles(srcDir, dstDir string) error {
	err := os.MkdirAll(dstDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
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

		err = os.MkdirAll(filepath.Dir(dstPath), 0o750)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		err = copyFile(path, dstPath)
		if err != nil {
			return fmt.Errorf("failed to copy file from %s to %s: %w", path, dstPath, err)
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = destinationFile.Close()
	}()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
