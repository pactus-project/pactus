package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/pactus-project/pactus/util/logger"
)

const (
	_defaultConcurrencyPerChunk = 16
	_defaultMinSizeForChunk     = 1 << 20
	_defaultMaxRetries          = 3
)

type Downloader struct {
	client        *http.Client
	url           string
	filePath      string
	sha256Sum     string
	fileType      string
	fileName      string
	maxRetries    int
	cancel        context.CancelFunc
	statsCallback func(Stats)

	chunks []*chunk

	mu         sync.Mutex
	downloaded int64
}

type Stats struct {
	Downloaded int64
	TotalSize  int64
	Percent    float64
	Completed  bool
}

func New(url, filePath, sha256Sum string, opts ...Option) *Downloader {
	opt := defaultOptions()

	for _, o := range opts {
		o(opt)
	}

	return &Downloader{
		client:        opt.client,
		statsCallback: opt.statsCallBack,
		url:           url,
		filePath:      filePath,
		sha256Sum:     sha256Sum,
		chunks:        make([]*chunk, 0, _defaultConcurrencyPerChunk),
		maxRetries:    opt.maxRetries,
	}
}

func (d *Downloader) Start(ctx context.Context) {
	d.download(ctx)
}

func (d *Downloader) FileType() string {
	return d.fileType
}

func (d *Downloader) FileName() string {
	return d.fileName
}

func (d *Downloader) download(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	totalSize, err := d.getHeader(ctx)
	if err != nil {
		d.handleError(err)

		return
	}

	d.fileName = filepath.Base(d.filePath)
	if err := d.createDir(); err != nil {
		d.handleError(err)

		return
	}

	out, err := os.OpenFile(d.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		d.handleError(err)

		return
	}
	defer func() {
		_ = out.Close()
	}()

	d.updateStats(0, totalSize, false)

	var wg sync.WaitGroup
	for _, chk := range d.chunks {
		wg.Add(1)
		go func(chk *chunk) {
			defer wg.Done()
			if err := d.downloadChunk(ctx, out, chk, totalSize); err != nil {
				d.handleError(err)
			}
		}(chk)
	}

	wg.Wait()

	if ctx.Err() != nil {
		return
	}

	if err := d.finalizeDownload(); err != nil {
		d.handleError(err)

		return
	}
}

func (d *Downloader) getHeader(ctx context.Context) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, d.url, http.NoBody)
	if err != nil {
		return 0, &Error{Message: "failed to create new request for get header", Reason: err}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, &Error{Message: "failed to do request get header", Reason: err}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	d.fileType = resp.Header.Get("Content-Type")

	if resp.ContentLength > _defaultMinSizeForChunk {
		d.chunks = createChunks(resp.ContentLength, _defaultConcurrencyPerChunk)
	} else {
		d.chunks = append(d.chunks, &chunk{
			start: 0,
			end:   resp.ContentLength,
		})
	}

	return resp.ContentLength, nil
}

func (d *Downloader) createDir() error {
	dir := filepath.Dir(d.filePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return &Error{Message: "failed to create file path directory", Reason: err}
	}

	return nil
}

func (d *Downloader) downloadChunk(ctx context.Context, out *os.File, chk *chunk, totalSize int64) error {
	var err error
	for i := 0; i < d.maxRetries; i++ {
		err = d.downloadChunkWithContext(ctx, out, chk, totalSize)
		if err == nil || ctx.Err() != nil {
			return err
		}
	}

	return err
}

func (d *Downloader) downloadChunkWithContext(ctx context.Context, out *os.File, chk *chunk, totalSize int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, http.NoBody)
	if err != nil {
		return &Error{Message: "failed to create new request for download chunk", Reason: err}
	}

	req.Header.Set("Range", chk.rangeHeader())
	resp, err := d.client.Do(req)
	if err != nil {
		return &Error{Message: "failed to do request download chunk", Reason: err}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return &Error{
			Message: "response has invalid status code",
			Reason:  fmt.Errorf("got http response %s from %s: %w", resp.Status, d.url, err),
		}
	}

	buf := make([]byte, 32*1024) // 32KB buffer for reading the response body
	offset := chk.start
	for {
		count, err := resp.Body.Read(buf)
		if count > 0 {
			d.mu.Lock()
			for written := 0; written < count; {
				numBytes, err := out.WriteAt(buf[written:count], offset+int64(written))
				if err != nil {
					d.mu.Unlock()

					return &Error{Message: "failed write data into file", Reason: err}
				}
				written += numBytes
			}
			offset += int64(count)
			d.downloaded += int64(count)
			d.updateStats(d.downloaded, totalSize, false)
			d.mu.Unlock()
		}
		if err != nil {
			// if error is io.EOF stop write for loop response body.
			if errors.Is(err, io.EOF) {
				break
			}

			return &Error{Message: "error read body download chunk", Reason: err}
		}
	}

	return nil
}

func (d *Downloader) updateStats(downloaded, totalSize int64, completed bool) {
	if d.statsCallback != nil {
		if downloaded > totalSize {
			// In case of re-downloading a chunk...
			downloaded = totalSize
		}
		stats := Stats{
			Downloaded: downloaded,
			TotalSize:  totalSize,
			Percent:    float64(downloaded) / float64(totalSize) * 100,
		}

		if completed {
			stats.Completed = true
		}

		d.statsCallback(stats)
	}
}

func (d *Downloader) finalizeDownload() error {
	// Recalculate the hash by re-reading the entire file
	out, err := os.Open(d.filePath)
	if err != nil {
		return &Error{Message: "failed to open file", Reason: err}
	}
	defer func() {
		_ = out.Close()
	}()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, out); err != nil {
		return &Error{Message: "failed copy file data to hasher for calculate hash", Reason: err}
	}

	sum := hex.EncodeToString(hasher.Sum(nil))
	if sum != d.sha256Sum {
		return &Error{Message: "sha256 mismatch", Reason: err}
	}

	d.updateStats(0, 0, true)

	return nil
}

func (d *Downloader) handleError(err error) {
	logger.Error("failed to download", "error", err)
	if d.cancel != nil {
		d.cancel()
	}
}
