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
	"time"

	"github.com/pactus-project/pactus/util/downloader"
)

const maxDecompressedSize = 10 << 20 // 10 MB

type DMStateFunc func(
	fileName string,
	totalSize, downloaded int64,
	totalItem, downloadedItem int,
	percentage float64,
)

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

type DownloadManager struct {
	snapshotURL string
	extractDir  string
	tempDir     string
	storeDir    string
	zipFileList []string
}

func NewDownloadManager(snapshotURL, extractDir, tempDir, storeDir string) *DownloadManager {
	return &DownloadManager{
		snapshotURL: snapshotURL,
		extractDir:  extractDir,
		tempDir:     tempDir,
		storeDir:    storeDir,
		zipFileList: make([]string, 0),
	}
}

func (dl *DownloadManager) GetMetadata(ctx context.Context) ([]Metadata, error) {
	cli := http.DefaultClient
	metaURL, err := url.JoinPath(dl.snapshotURL, "metadata.json")
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

func (dl *DownloadManager) Download(
	ctx context.Context,
	metadata *Metadata,
	stateFunc DMStateFunc,
) {
	downloadedItem := 0
	for _, data := range metadata.Data {
		done := make(chan struct{})
		dlLink, err := url.JoinPath(dl.snapshotURL, data.Path)
		FatalErrorCheck(err)

		fileName := filepath.Base(dlLink)

		filePath := fmt.Sprintf("%s/%s", dl.tempDir, fileName)

		d := downloader.New(
			dlLink,
			filePath,
			data.Sha,
		)

		d.Start(ctx)

		go func() {
			err := <-d.Errors()
			FatalErrorCheck(err)
		}()

		go func() {
			for state := range d.Stats() {
				stateFunc(fileName, state.TotalSize, state.Downloaded, len(metadata.Data), downloadedItem, state.Percent)
				if state.Completed {
					done <- struct{}{}
					close(done)

					return
				}
			}
		}()

		<-done
		dl.zipFileList = append(dl.zipFileList, filePath)
		downloadedItem++
	}
}

func (dl *DownloadManager) ExtractAndStoreFiles() error {
	for _, path := range dl.zipFileList {
		r, err := zip.OpenReader(path)
		if err != nil {
			return fmt.Errorf("failed to open zip file: %w", err)
		}

		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in zip archive: %w", err)
			}

			fpath := fmt.Sprintf("%s/%s", dl.extractDir, f.Name)

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
		_ = r.Close()
	}

	return nil
}

func (*DownloadManager) ParseTime(dateString string) time.Time {
	const layout = "2006-01-02T15:04:05.000000"

	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

// CopyAllFiles copies all files from srcDir to dstDir.
func (dl *DownloadManager) CopyAllFiles() error {
	return filepath.Walk(dl.extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		relativePath, err := filepath.Rel(dl.extractDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dl.storeDir, relativePath)

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
