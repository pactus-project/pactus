package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	ErrHeaderRequest      = errors.New("request header error")
	ErrSHA256Mismatch     = errors.New("sha256 mismatch")
	ErrCreateDir          = errors.New("create dir error")
	ErrInvalidFilePath    = errors.New("file path is a directory, not a file")
	ErrGetFileInfo        = errors.New("get file info error")
	ErrCopyExistsFileData = errors.New("error copying existing file data")
	ErrDoRequest          = errors.New("error doing request")
	ErrFileWriting        = errors.New("error writing file")
	ErrNewRequest         = errors.New("error creating request")
	ErrOpenFileExists     = errors.New("error opening existing file")
)

type Downloader struct {
	client    *http.Client
	url       string
	filePath  string
	sha256Sum string
	fileType  string
	fileName  string
	statsCh   chan Stats
	errCh     chan error
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
		client:    opt.client,
		url:       url,
		filePath:  filePath,
		sha256Sum: sha256Sum,
		statsCh:   make(chan Stats),
		errCh:     make(chan error, 1),
	}
}

func (d *Downloader) Start(ctx context.Context) {
	go d.download(ctx)
}

func (d *Downloader) Stats() <-chan Stats {
	return d.statsCh
}

func (d *Downloader) FileType() string {
	return d.fileType
}

func (d *Downloader) FileName() string {
	return d.fileName
}

func (d *Downloader) Errors() <-chan error {
	return d.errCh
}

func (d *Downloader) download(ctx context.Context) {
	stats, err := d.getHeader(ctx)
	if err != nil {
		d.handleError(err)

		return
	}

	d.fileName = filepath.Base(d.filePath)
	if err := d.createDir(); err != nil {
		d.handleError(err)

		return
	}

	out, err := d.openFile()
	if err != nil {
		d.handleError(err)

		return
	}
	defer func() {
		_ = out.Close()
	}()

	if err := d.validateExistingFile(out, &stats); err != nil {
		d.handleError(err)

		return
	}

	if err := d.downloadFile(ctx, out, &stats); err != nil {
		d.handleError(err)
	}
}

func (d *Downloader) getHeader(ctx context.Context) (Stats, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, d.url, http.NoBody)
	if err != nil {
		return Stats{}, ErrHeaderRequest
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return Stats{}, ErrHeaderRequest
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	d.fileType = resp.Header.Get("Content-Type")

	return Stats{
		TotalSize: resp.ContentLength,
	}, nil
}

func (d *Downloader) createDir() error {
	dir := filepath.Dir(d.filePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return ErrCreateDir
	}

	return nil
}

func (d *Downloader) openFile() (*os.File, error) {
	fileInfo, err := os.Stat(d.filePath)
	if err == nil && fileInfo.IsDir() {
		return nil, ErrInvalidFilePath
	}

	return os.OpenFile(d.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
}

func (*Downloader) validateExistingFile(out *os.File, stats *Stats) error {
	fileInfo, err := out.Stat()
	if err != nil {
		return ErrGetFileInfo
	}
	stats.Downloaded = fileInfo.Size()

	return nil
}

func (d *Downloader) downloadFile(ctx context.Context, out *os.File, stats *Stats) error {
	req, err := d.createRequest(ctx, stats.Downloaded)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return ErrDoRequest
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	buffer := make([]byte, 32*1024)
	hasher := sha256.New()

	if err := d.updateHasherWithExistingData(stats.Downloaded, hasher); err != nil {
		return err
	}

	return d.writeToFile(ctx, resp, out, buffer, hasher, stats)
}

func (d *Downloader) createRequest(ctx context.Context, downloaded int64) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, http.NoBody)
	if err != nil {
		return nil, ErrNewRequest
	}
	if downloaded > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", downloaded))
	}

	return req, nil
}

func (d *Downloader) updateHasherWithExistingData(downloaded int64, hasher io.Writer) error {
	if downloaded > 0 {
		existingFile, err := os.Open(d.filePath)
		if err != nil {
			return ErrOpenFileExists
		}
		defer func() {
			_ = existingFile.Close()
		}()

		if _, err := io.CopyN(hasher, existingFile, downloaded); err != nil {
			return ErrCopyExistsFileData
		}
	}

	return nil
}

func (d *Downloader) writeToFile(ctx context.Context, resp *http.Response, out *os.File, buffer []byte,
	hasher hash.Hash, stats *Stats,
) error {
	for {
		select {
		case <-ctx.Done():
			d.stop()

			return ctx.Err()
		default:
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				if _, err := out.Write(buffer[:n]); err != nil {
					return ErrFileWriting
				}

				if _, err := hasher.Write(buffer[:n]); err != nil {
					return ErrFileWriting
				}

				stats.Downloaded += int64(n)
				stats.Percent = float64(stats.Downloaded) / float64(stats.TotalSize) * 100
				d.statsCh <- *stats
			}
			if err != nil {
				if err == io.EOF {
					return d.finalizeDownload(hasher, stats)
				}

				return fmt.Errorf("error reading response body: %w", err)
			}
		}
	}
}

func (d *Downloader) finalizeDownload(hasher hash.Hash, stats *Stats) error {
	stats.Completed = true
	sum := hex.EncodeToString(hasher.Sum(nil))
	if sum != d.sha256Sum {
		return ErrSHA256Mismatch
	}
	d.statsCh <- *stats

	d.stop()

	return nil
}

func (d *Downloader) stop() {
	close(d.statsCh)
	close(d.errCh)
}

func (d *Downloader) handleError(err error) {
	d.errCh <- err
	d.stop()
}
