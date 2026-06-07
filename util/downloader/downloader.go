package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/ezex-io/gopkg/util/retry"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

const (
	_defaultNumberOfChunks = 16
	_defaultMaxRetries     = 3
)

type Downloader struct {
	client        *http.Client
	url           string
	filePath      string
	sha3Sum       string
	fileType      string
	fileSize      int64
	maxRetries    int
	cancel        context.CancelFunc
	statsCallback func(Stats)

	mu         sync.Mutex
	downloaded int64
	completed  bool
}

type Stats struct {
	Downloaded int64
	TotalSize  int64
	Percent    float64
	Completed  bool
}

func New(url, filePath, sha3Sum string, opts ...Option) *Downloader {
	opt := defaultOptions()

	for _, o := range opts {
		o(opt)
	}

	return &Downloader{
		client:        opt.client,
		statsCallback: opt.statsCallBack,
		url:           url,
		filePath:      filePath,
		sha3Sum:       sha3Sum,
		maxRetries:    opt.maxRetries,
		downloaded:    0,
	}
}

func (d *Downloader) Download(ctx context.Context) error {
	if d.statsCallback != nil {
		scheduler.Every(1*time.Second).Do(ctx, d.reportStats)

		// Report for the first time.
		d.reportStats(ctx)
	}

	return d.download(ctx)
}

func (d *Downloader) reportStats(context.Context) {
	d.mu.Lock()
	stats := Stats{
		Downloaded: d.downloaded,
		TotalSize:  d.fileSize,
		Percent:    (float64(d.downloaded) / float64(d.fileSize)) * 100,
		Completed:  d.completed,
	}
	d.mu.Unlock()

	d.statsCallback(stats)
}

func (d *Downloader) FileType() string {
	return d.fileType
}

func (d *Downloader) download(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	err := d.parseHeaders(ctx)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(d.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	var wg sync.WaitGroup
	chunks := d.createChunks()
	for _, chk := range chunks {
		wg.Go(func() {
			if err := d.downloadChunkRetry(ctx, out, chk); err != nil {
				logger.Error("error on downloading a chunk", "error", err)
			}
		})
	}

	wg.Wait()

	if err := d.finalizeDownload(); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) parseHeaders(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, d.url, http.NoBody)
	if err != nil {
		return &Error{Message: "failed to create new request for get header", Reason: err}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return &Error{Message: "failed to do request get header", Reason: err}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	d.fileType = resp.Header.Get("Content-Type")
	d.fileSize = resp.ContentLength

	return nil
}

func (d *Downloader) createChunks() []*chunk {
	return createChunks(d.fileSize, _defaultNumberOfChunks)
}

func (d *Downloader) downloadChunkRetry(ctx context.Context, out *os.File, chk *chunk) error {
	return retry.ExecuteSync(ctx, func() error {
		err := d.downloadChunk(ctx, out, chk)
		if err != nil {
			logger.Warn("retry downloading a chuck, ", "error", err)
		}

		return err
	}, retry.WithSyncMaxRetries(d.maxRetries))
}

func (d *Downloader) downloadChunk(ctx context.Context, out *os.File, chk *chunk) error {
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

	bufSize := 32 * 1024 // 32KB buffer for reading the response body
	buf := make([]byte, bufSize)
	offset := chk.start
	for {
		remained := int64(chk.end+1) - offset
		if remained <= 0 {
			break
		}
		if int64(bufSize) > remained {
			buf = buf[:remained]
		}

		count, err := resp.Body.Read(buf)
		if err != nil {
			// if error is not io.EOF stop write for loop response body.
			if !errors.Is(err, io.EOF) {
				return &Error{Message: "error read body download chunk", Reason: err}
			}
		}
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
			d.mu.Unlock()
		}
	}

	return nil
}

func (d *Downloader) finalizeDownload() error {
	// Recalculate the hash by re-reading the entire file
	sum, err := util.CalculateChecksum(d.filePath)
	if err != nil {
		return &Error{Message: "unable to calculate downloaded checksum", Reason: err}
	}
	if sum != d.sha3Sum {
		return &Error{Message: "sha256 mismatch", Reason: fmt.Errorf("expected %s, got %s", d.sha3Sum, sum)}
	}

	d.completed = true

	return nil
}
