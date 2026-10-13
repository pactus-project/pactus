package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rangeServer serves `content` with support of the Range header.
// `truncate` shortens the body of the partial responses by the given number of bytes,
// while the headers are consistent with the shortened body.
func rangeServer(t *testing.T, content []byte, truncate int64) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		size := int64(len(content))
		if r.Method == http.MethodHead {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", size))

			return
		}

		var start, end int64
		if _, err := fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)

			return
		}
		end = min(end, size-1)
		if start > end {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)

			return
		}

		body := content[start : end+1]
		body = body[:max(int64(len(body))-truncate, 0)]

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)

	return server
}

func checksumOf(content []byte) string {
	sum := sha256.Sum256(content)

	return hex.EncodeToString(sum[:])
}

// downloadWithTimeout fails the test, instead of blocking it, if the download never ends.
func downloadWithTimeout(t *testing.T, d *Downloader) error {
	t.Helper()

	done := make(chan error, 1)
	go func() { done <- d.Download(t.Context()) }()

	select {
	case err := <-done:
		return err
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the download never ends")

		return nil
	}
}

func TestDownloaderChecksumMismatch(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	content := ts.RandBytes(_defaultNumberOfChunks * 1024)
	server := rangeServer(t, content, 0)

	filePath := util.TempFilePath()
	wrongChecksum := checksumOf(ts.RandBytes(100))
	d := New(server.URL, filePath, wrongChecksum, WithCustomClient(server.Client()), WithMaxRetries(1))

	err := downloadWithTimeout(t, d)

	var dlErr *Error
	require.ErrorAs(t, err, &dlErr)
	assert.Equal(t, "sha256 mismatch", dlErr.Message)
	require.ErrorContains(t, err, wrongChecksum, "the error should tell the expected checksum")
	require.ErrorContains(t, err, checksumOf(content), "the error should tell the checksum of the downloaded file")
}

// The body of the response is shorter than the requested range.
// This happens if the file on the server is replaced by a smaller one,
// or if the server closes the connection after a part of the data.
func TestDownloaderTruncatedBody(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	content := ts.RandBytes(_defaultNumberOfChunks * 1024)
	server := rangeServer(t, content, 100)

	d := New(server.URL, util.TempFilePath(), checksumOf(content),
		WithCustomClient(server.Client()), WithMaxRetries(1))

	err := downloadWithTimeout(t, d)

	require.Error(t, err, "a file with missing data is not a valid download")
}

func TestDownloaderServerFailures(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	content := ts.RandBytes(_defaultNumberOfChunks * 1024)

	t.Run("server is not reachable", func(t *testing.T) {
		server := rangeServer(t, content, 0)
		server.Close()

		d := New(server.URL, util.TempFilePath(), checksumOf(content),
			WithCustomClient(server.Client()), WithMaxRetries(1))

		var dlErr *Error
		require.ErrorAs(t, downloadWithTimeout(t, d), &dlErr)
		assert.Equal(t, "failed to do request get header", dlErr.Message)
	})

	t.Run("invalid url", func(t *testing.T) {
		d := New("://invalid-url", util.TempFilePath(), checksumOf(content), WithMaxRetries(1))

		var dlErr *Error
		require.ErrorAs(t, downloadWithTimeout(t, d), &dlErr)
		assert.Equal(t, "failed to create new request for get header", dlErr.Message)
	})

	t.Run("destination file can't be created", func(t *testing.T) {
		server := rangeServer(t, content, 0)
		missingDir := util.TempDirPath() + "/missing/file"

		d := New(server.URL, missingDir, checksumOf(content),
			WithCustomClient(server.Client()), WithMaxRetries(1))

		require.Error(t, downloadWithTimeout(t, d))
	})

	t.Run("server doesn't support range requests", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			_, _ = w.Write(content)
		}))
		defer server.Close()

		d := New(server.URL, util.TempFilePath(), checksumOf(content),
			WithCustomClient(server.Client()), WithMaxRetries(1))

		require.Error(t, downloadWithTimeout(t, d), "the file is not downloaded without range support")
	})
}

func TestDownloaderCancel(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	content := ts.RandBytes(_defaultNumberOfChunks * 1024)
	server := rangeServer(t, content, 0)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	d := New(server.URL, util.TempFilePath(), checksumOf(content),
		WithCustomClient(server.Client()), WithMaxRetries(1))

	err := d.Download(ctx)

	require.ErrorIs(t, err, context.Canceled)
}

func TestErrorMessages(t *testing.T) {
	cause := errors.New("connection reset")
	err := &Error{Message: "failed to do request", Reason: cause}

	assert.Equal(t, "failed to do request: connection reset", err.Error())
	assert.ErrorIs(t, err, cause)
}
