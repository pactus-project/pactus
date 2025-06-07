package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/downloader"
)

const DefaultSnapshotURL = "https://snapshot.pactus.org"

const maxDecompressedSize = 10 << 20 // 10 MB

type ImporterStateFunc func(
	fileName string,
	totalSize, downloaded int64,
	percentage float64,
)

type Metadata struct {
	Name      string       `json:"name"`
	CreatedAt string       `json:"created_at"`
	Compress  string       `json:"compress"`
	Data      SnapshotData `json:"data"`
}

type SnapshotData struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Sha  string `json:"sha"`
	Size uint64 `json:"size"`
}

func (md *Metadata) CreatedAtTime() time.Time {
	const layout = "2006-01-02T15:04:05.000000"

	parsedTime, err := time.Parse(layout, md.CreatedAt)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

// Importer downloads and imports the pruned data from a centralized server.
type Importer struct {
	snapshotURL  string
	tempDir      string
	storeDir     string
	dataFileName string
}

func NewImporter(chainType genesis.ChainType, snapshotURL, storeDir string) (*Importer, error) {
	if util.PathExists(storeDir) {
		return nil, fmt.Errorf("data directory is not empty: %s", storeDir)
	}

	switch chainType {
	case genesis.Mainnet:
		snapshotURL += "/mainnet/"
	case genesis.Testnet:
		snapshotURL += "/testnet/"
	case genesis.Localnet:
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}

	tempDir := util.TempDirPath()

	return &Importer{
		snapshotURL: snapshotURL,
		tempDir:     tempDir,
		storeDir:    storeDir,
	}, nil
}

func (i *Importer) GetMetadata(ctx context.Context) ([]Metadata, error) {
	cli := http.DefaultClient
	metaURL, err := url.JoinPath(i.snapshotURL, "metadata.json")
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

	sort.SliceStable(metadata, func(i, j int) bool {
		return metadata[i].CreatedAtTime().After(metadata[j].CreatedAtTime())
	})

	return metadata, nil
}

func (i *Importer) Download(ctx context.Context, metadata *Metadata,
	stateFunc ImporterStateFunc,
) error {
	done := make(chan error)
	defer close((done))

	dlLink, err := url.JoinPath(i.snapshotURL, metadata.Data.Path)
	if err != nil {
		return err
	}

	fileName := filepath.Base(dlLink)
	i.dataFileName = fileName
	filePath := fmt.Sprintf("%s/%s", i.tempDir, fileName)

	downloader := downloader.New(dlLink, filePath, metadata.Data.Sha)
	downloader.Start(ctx)

	go func() {
		for state := range downloader.Stats() {
			stateFunc(fileName, state.TotalSize, state.Downloaded, state.Percent)
			if state.Completed {
				log.Println("download completed")
				done <- nil

				return
			}
		}
	}()

	return <-done
}

func (i *Importer) Cleanup() error {
	return os.RemoveAll(i.tempDir)
}

func (i *Importer) ExtractAndStoreFiles() error {
	zipPath := filepath.Join(i.tempDir, i.dataFileName)
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	for _, file := range reader.File {
		if err := i.extractAndWriteFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (i *Importer) extractAndWriteFile(file *zip.File) error {
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in zip archive: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	fPath, err := util.SanitizeArchivePath(i.tempDir, file.Name)
	if err != nil {
		return fmt.Errorf("failed to make archive path: %w", err)
	}

	if file.FileInfo().IsDir() {
		return util.Mkdir(fPath)
	}

	if err := util.Mkdir(filepath.Dir(fPath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	outFile, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	// Use a limited reader to prevent DoS attacks via decompression bomb
	lr := &io.LimitedReader{R: reader, N: maxDecompressedSize}
	written, err := io.Copy(outFile, lr)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	if written >= maxDecompressedSize {
		return fmt.Errorf("file exceeds maximum decompressed size limit: %s", fPath)
	}

	return nil
}

func (i *Importer) MoveStore() error {
	return util.MoveDirectory(filepath.Join(i.tempDir, "data"), i.storeDir)
}
