package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrHeaderRequest      = errors.New("request header error")
	ErrSHA256Mismatch     = errors.New("sha256 mismatch")
	ErrCreateDir          = errors.New("create dir error")
	ErrOpenFile           = errors.New("open file error")
	ErrInvalidFilePath    = errors.New("file path is a directory, not a file")
	ErrGetFileInfo        = errors.New("get file info error")
	ErrCopyExistsFileData = errors.New("error copying existing file data")
	ErrDoRequest          = errors.New("error doing request")
	ErrFileWriting        = errors.New("error writing file")
	ErrNewRequest         = errors.New("error creating request")
	ErrOpenFileExists     = errors.New("error opening file exists")
)

type Downloader struct {
	url       string
	filePath  string
	sha256Sum string
	fileType  string
	fileName  string
	pause     chan bool
	resume    chan bool
	stop      chan bool
	statsCh   chan Stats
	errCh     chan error
	wg        sync.WaitGroup
}

type Stats struct {
	Downloaded int64
	TotalSize  int64
	Percent    float64
	Completed  bool
}

func NewDownloader(url, filePath, sha256Sum string) *Downloader {
	return &Downloader{
		url:       url,
		filePath:  filePath,
		sha256Sum: sha256Sum,
		pause:     make(chan bool),
		resume:    make(chan bool),
		stop:      make(chan bool),
		statsCh:   make(chan Stats),
		errCh:     make(chan error, 1),
	}
}

func (d *Downloader) Start() {
	d.wg.Add(1)
	go d.download()
}

func (d *Downloader) Pause() {
	d.pause <- true
}

func (d *Downloader) Resume() {
	d.resume <- true
}

func (d *Downloader) Stop() {
	d.stop <- true
	d.wg.Wait()
	close(d.statsCh)
	close(d.errCh)
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

func (d *Downloader) download() {
	defer d.wg.Done()

	resp, err := http.Head(d.url)
	if err != nil {
		d.handleError(ErrHeaderRequest)
		return
	}

	stats := Stats{
		TotalSize: resp.ContentLength,
	}

	d.fileType = resp.Header.Get("Content-Type")
	d.fileName = filepath.Base(d.filePath)

	dir := filepath.Dir(d.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		d.handleError(ErrCreateDir)
		return
	}

	fileInfo, err := os.Stat(d.filePath)
	if err == nil && fileInfo.IsDir() {
		d.handleError(ErrInvalidFilePath)
		return
	}

	out, err := os.OpenFile(d.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		d.handleError(ErrOpenFile)
		return
	}
	defer out.Close()

	fileInfo, err = out.Stat()
	if err != nil {
		d.handleError(ErrGetFileInfo)
		return
	}
	stats.Downloaded = fileInfo.Size()

	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		d.handleError(ErrNewRequest)
		return
	}
	if stats.Downloaded > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", stats.Downloaded))
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		d.handleError(ErrDoRequest)
		return
	}
	defer resp.Body.Close()

	buffer := make([]byte, 32*1024)
	hasher := sha256.New()

	// Update the hasher with the already downloaded part
	if stats.Downloaded > 0 {
		existingFile, err := os.Open(d.filePath)
		if err != nil {
			d.handleError(ErrOpenFileExists)
			return
		}
		defer existingFile.Close()
		if _, err := io.CopyN(hasher, existingFile, stats.Downloaded); err != nil {
			d.handleError(ErrCopyExistsFileData)
			return
		}
	}

	for {
		select {
		case <-d.pause:
			<-d.resume
		case <-d.stop:
			return
		default:
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				if _, err := out.Write(buffer[:n]); err != nil {
					d.handleError(ErrFileWriting)
					return
				}
				hasher.Write(buffer[:n])
				stats.Downloaded += int64(n)
				stats.Percent = float64(stats.Downloaded) / float64(stats.TotalSize) * 100
				d.statsCh <- stats
			}
			if err != nil {
				if err == io.EOF {
					stats.Completed = true
					sum := hex.EncodeToString(hasher.Sum(nil))
					if sum != d.sha256Sum {
						d.handleError(ErrSHA256Mismatch)
					} else {
						d.statsCh <- stats
					}
					return
				}
				d.handleError(fmt.Errorf("error reading response body: %v", err))
				return
			}
		}
	}
}

func (d *Downloader) handleError(err error) {
	select {
	case d.errCh <- err:
	default:
	}
	d.Stop()
}
