package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestDownloader(t *testing.T) {
	fileContent := make([]byte, 1*1024*1024) // 1 MB
	for i := range fileContent {
		fileContent[i] = byte(i % 256)
	}

	fileURL := "/testfile"
	expectedSHA256 := sha256.Sum256(fileContent)
	expectedSHA256Hex := hex.EncodeToString(expectedSHA256[:])

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fileURL {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))
			_, err := w.Write(fileContent)
			assert.NoError(t, err)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	filePath := util.TempFilePath()

	defer func() {
		assert.NoError(t, os.RemoveAll("./testdata"))
	}()

	downloader := New(server.URL+fileURL, filePath, expectedSHA256Hex,
		WithCustomClient(server.Client()),
		WithStatsCallback(printDownloaderStats),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	downloader.Start(ctx)

	t.Log(downloader.FileName())
	t.Log(downloader.FileType())

	downloadedContent, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read the downloaded file")
	assert.Equal(t, fileContent, downloadedContent, "Downloaded file content does not match expected content")
}

func printDownloaderStats(sts Stats) {
	if !sts.Completed {
		fmt.Printf("Downloaded: %d / %d (%.2f%%)\n", sts.Downloaded, sts.TotalSize, sts.Percent)
	}
}
