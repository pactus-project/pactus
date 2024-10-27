package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestDownloader(t *testing.T) {
	fileContent := []byte("This is a test file content")
	fileURL := "/testfile"
	expectedSHA256 := sha256.Sum256(fileContent)
	expectedSHA256Hex := hex.EncodeToString(expectedSHA256[:])

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fileURL {
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

	downloader := New(server.URL+fileURL, filePath, expectedSHA256Hex, WithCustomClient(server.Client()))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	go func() {
		downloader.Start(ctx)
	}()

	done := make(chan bool)

	go func() {
		for stat := range downloader.Stats() {
			log.Printf("Downloaded: %d / %d (%.2f%%)\n", stat.Downloaded, stat.TotalSize, stat.Percent)
			assert.True(t, stat.Downloaded <= stat.TotalSize, "Downloaded size should not exceed total size")
			assert.True(t, stat.Percent <= 100, "Download percentage should not exceed 100")

			if stat.Completed {
				log.Println("Download completed successfully")
				assert.Equal(t, float64(100), stat.Percent, "Download should be 100% complete")
				done <- true

				return
			}
		}
	}()

	go func() {
		for err := range downloader.Errors() {
			assert.Fail(t, "Download encountered an error", err)
			done <- true

			return
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Minute):
		cancel()
		assert.Fail(t, "Download test timed out")
	}

	t.Log(downloader.FileName())
	t.Log(downloader.FileType())

	downloadedContent, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read the downloaded file")
	assert.Equal(t, fileContent, downloadedContent, "Downloaded file content does not match expected content")
}
