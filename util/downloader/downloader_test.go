package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func cleanup(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}

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

	filePath := "./testdata/example_testfile.txt"

	defer func() {
		assert.NoError(t, os.RemoveAll("./testdata"))
	}()

	dl := New(server.URL+fileURL, filePath, expectedSHA256Hex, WithCustomClient(server.Client()))

	assrt := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	go func() {
		dl.Start(ctx)
	}()

	done := make(chan bool)

	go func() {
		for stat := range dl.Stats() {
			log.Printf("Downloaded: %d / %d (%.2f%%)\n", stat.Downloaded, stat.TotalSize, stat.Percent)
			assrt.True(stat.Downloaded <= stat.TotalSize, "Downloaded size should not exceed total size")
			assrt.True(stat.Percent <= 100, "Download percentage should not exceed 100")

			if stat.Completed {
				log.Println("Download completed successfully")
				assrt.Equal(float64(100), stat.Percent, "Download should be 100% complete")
				done <- true

				return
			}
		}
	}()

	go func() {
		for err := range dl.Errors() {
			assrt.Fail("Download encountered an error", err)
			done <- true

			return
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Minute):
		cancel()
		assrt.Fail("Download test timed out")
	}

	t.Log(dl.FileName())
	t.Log(dl.FileType())

	downloadedContent, err := os.ReadFile(filePath)
	assrt.NoError(err, "Failed to read the downloaded file")
	assrt.Equal(fileContent, downloadedContent, "Downloaded file content does not match expected content")

	assert.NoError(t, cleanup(dl.filePath))
}
