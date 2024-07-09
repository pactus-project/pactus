package downloader

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func setup() *Downloader {
	fileURL := "https://github.com/pactus-project/Whitepaper/releases/latest/download/pactus_whitepaper.pdf"
	filePath := "./testdata/example.pdf"
	expectedSHA256 := "ea956128717b49669f29eeed116bc11b9bbdcd50f1df130e124ffd36afe71652"

	dl := NewDownloader(fileURL, filePath, expectedSHA256)
	return dl
}

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
	dl := setup()

	assrt := assert.New(t)
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dl.Start()
	}()

	go func() {
		select {
		case <-ctx.Done():
			dl.Stop()
			assrt.Fail("Download test timed out")
		}
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
		dl.Stop()
		assrt.Fail("Download test timed out")
	}

	assert.NoError(t, cleanup(dl.filePath))

	wg.Wait()
}

func TestDownloaderOperations(t *testing.T) {
	dl := setup()
	dl.Start()

	assrt := assert.New(t)

	go func() {
		for stat := range dl.Stats() {
			log.Printf("Downloaded: %d / %d (%.2f%%)\n", stat.Downloaded, stat.TotalSize, stat.Percent)
			assrt.True(stat.Downloaded <= stat.TotalSize, "Downloaded size should not exceed total size")
			assrt.True(stat.Percent <= 100, "Download percentage should not exceed 100")

			if stat.Completed {
				log.Println("Download completed successfully")
				assrt.Equal(float64(100), stat.Percent, "Download should be 100% complete")
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)
	dl.Pause()
	t.Log("Paused")
	time.Sleep(3 * time.Second)
	dl.Resume()
	t.Log("Resumed")

	time.Sleep(1 * time.Second)
	dl.Stop()

	assrt.NoError(cleanup(dl.filePath))
}
