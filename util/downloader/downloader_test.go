package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloader(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		name string
		size int
	}{
		{size: _defaultNumberOfChunks * 1024},
		{size: (_defaultNumberOfChunks * 1024) - 1},
		{size: (_defaultNumberOfChunks * 1024) + 1},
		{size: ts.RandIntMax(64 * 1000 * 1000)},
		{size: 67850301},
	}

	for _, tt := range tests {
		fileContent := ts.RandBytes(tt.size)

		fileURL := "/testfile.zip"
		expectedSHA256 := sha256.Sum256(fileContent)
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
				w.Header().Set("Content-Length", fmt.Sprintf("%d", tt.size))
				_, _ = w.Write(fileContent)

				return
			}

			// Expecting "bytes=start-end"
			var start, end int64
			_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
			if err != nil {
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)

				return
			}
			if end >= int64(tt.size) {
				end = int64(tt.size - 1)
			}
			if start > end || start >= int64(tt.size) {
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)

				return
			}

			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, tt.size))
			w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write(fileContent[start : end+1])
		}))

		filePath := util.TempFilePath()

		downloader := New(
			server.URL+fileURL, filePath, expectedSHA256Hex,
			WithCustomClient(server.Client()),
			WithStatsCallback(printDownloaderStats),
			WithMaxRetries(1),
		)

		err := downloader.Download(t.Context())
		require.NoError(t, err)

		assert.Equal(t, "application/octet-stream", downloader.FileType())
	}
}

func printDownloaderStats(sts Stats) {
	if !sts.Completed {
		fmt.Printf("Downloaded: %d / %d (%.2f%%)\n", sts.Downloaded, sts.TotalSize, sts.Percent)
	}
}
