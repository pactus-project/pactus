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

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/downloader"
)

const maxDecompressedSize = 10 << 20 // 10 MB

type DMStateFunc func(
	fileName string,
	totalSize, downloaded int64,
	percentage float64,
)

type Metadata struct {
	Name      string        `json:"name"`
	CreatedAt string        `json:"created_at"`
	Compress  string        `json:"compress"`
	Data      *SnapshotData `json:"data"`
}

type SnapshotData struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Sha  string `json:"sha"`
	Size uint64 `json:"size"`
}

type DownloadManager struct {
	snapshotURL  string
	tempDir      string
	storeDir     string
	dataFileName string
}

func NewDownloadManager(snapshotURL, tempDir, storeDir string) *DownloadManager {
	return &DownloadManager{
		snapshotURL: snapshotURL,
		tempDir:     tempDir,
		storeDir:    storeDir,
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
	done := make(chan struct{})
	dlLink, err := url.JoinPath(dl.snapshotURL, metadata.Data.Path)
	FatalErrorCheck(err)

	fileName := filepath.Base(dlLink)

	dl.dataFileName = fileName

	filePath := fmt.Sprintf("%s/%s", dl.tempDir, fileName)

	d := downloader.New(
		dlLink,
		filePath,
		metadata.Data.Sha,
	)

	d.Start(ctx)

	go func() {
		err := <-d.Errors()
		FatalErrorCheck(err)
	}()

	go func() {
		for state := range d.Stats() {
			stateFunc(fileName, state.TotalSize, state.Downloaded, state.Percent)
			if state.Completed {
				done <- struct{}{}
				close(done)

				return
			}
		}
	}()

	<-done
}

func (*DownloadManager) ParseTime(dateString string) time.Time {
	const layout = "2006-01-02T15:04:05.000000"

	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func (dl *DownloadManager) Cleanup() error {
	return os.RemoveAll(dl.tempDir)
}

func (dl *DownloadManager) ExtractAndStoreFiles() error {
	zipPath := filepath.Join(dl.tempDir, dl.dataFileName)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		if err := extractAndWriteFile(f, dl.tempDir); err != nil {
			return err
		}
	}

	return nil
}

func extractAndWriteFile(f *zip.File, destination string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in zip archive: %w", err)
	}
	defer func() {
		_ = rc.Close()
	}()

	fpath := fmt.Sprintf("%s/%s", destination, f.Name)
	if f.FileInfo().IsDir() {
		return util.Mkdir(fpath)
	}

	if err := util.Mkdir(filepath.Dir(fpath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	outFile, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	// Use a limited reader to prevent DoS attacks via decompression bomb
	lr := &io.LimitedReader{R: rc, N: maxDecompressedSize}
	written, err := io.Copy(outFile, lr)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	if written >= maxDecompressedSize {
		return fmt.Errorf("file exceeds maximum decompressed size limit: %s", fpath)
	}

	return nil
}
