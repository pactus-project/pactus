package downloader

import (
	"context"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloader(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	fileSize := ts.RandIntMax(10 * 1024 * 1024)
	fileContent := ts.RandBytes(fileSize)

	fileURL := "/testfile.zip"
	expectedSHA256 := sha3.Sum256(fileContent)
	expectedSHA256Hex := hex.EncodeToString(expectedSHA256[:])

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != fileURL {
			http.NotFound(w, r)
			return
		}

		// Parse Range header
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			// Full file request (used only for HEAD or first chunk)
			w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
			w.Write(fileContent)
			return
		}

		// Expecting "bytes=start-end"
		var start, end int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		if end >= int64(fileSize) {
			end = int64(fileSize - 1)
		}
		if start > end || start >= int64(fileSize) {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
		w.WriteHeader(http.StatusPartialContent)
		w.Write(fileContent[start : end+1])
	}))

	filePath := util.TempFilePath()

	defer func() {
		require.NoError(t, os.RemoveAll("./testdata"))
	}()

	downloader := New(
		server.URL+fileURL, filePath, expectedSHA256Hex,
		WithCustomClient(server.Client()),
		WithStatsCallback(printDownloaderStats),
	)

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	err := downloader.Download(ctx)
	require.NoError(t, err)

	assert.Equal(t, "application/octet-stream", downloader.FileType())
}

func printDownloaderStats(sts Stats) {
	if !sts.Completed {
		fmt.Printf("Downloaded: %d / %d (%.2f%%)\n", sts.Downloaded, sts.TotalSize, sts.Percent)
	}
}
