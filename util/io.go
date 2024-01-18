package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func IsAbsPath(p string) bool {
	return filepath.IsAbs(p)
}

func MakeAbs(p string) string {
	if IsAbsPath(p) {
		return p
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Clean(filepath.Join(wd, p))
}

func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func WriteFile(filename string, data []byte) error {
	// create directory
	if err := Mkdir(filepath.Dir(filename)); err != nil {
		return err
	}
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return fmt.Errorf("failed to write to %s: %w", filename, err)
	}

	return nil
}

func Mkdir(dir string) error {
	// create the directory
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("could not create directory %s", dir)
	}

	return nil
}

func PathExists(p string) bool {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func TempDirPath() string {
	p, err := os.MkdirTemp("", "pactus*")
	if err != nil {
		panic(err)
	}

	return p
}

func TempFilePath() string {
	return filepath.Join(TempDirPath(), "file")
}

func IsDirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	return errors.Is(err, io.EOF)
}

func IsDirNotExistsOrEmpty(name string) bool {
	if !PathExists(name) {
		return true
	}

	return IsDirEmpty(name)
}

func IsValidDirPath(fp string) bool {
	fi, err := os.Stat(fp)
	if err == nil {
		if fi.IsDir() {
			if err := os.WriteFile(fp+"/test", []byte{}, 0o600); err != nil {
				return false
			}
			os.Remove(fp + "/test")

			return true
		}

		return false
	}

	if err := Mkdir(fp); err != nil {
		return false
	}
	os.Remove(fp)

	return true
}

// TODO: move these to a test suite

// fixedWriter implements the io.Writer interface and intentionally allows
// testing of error paths by forcing short writes.
type FixedWriter struct {
	b   []byte
	pos int
}

// Write writes the contents of p to w. When the contents of p would cause
// the writer to exceed the maximum allowed size of the fixed writer,
// io.ErrShortWrite is returned and the writer is left unchanged.
//
// This satisfies the io.Writer interface.
func (w *FixedWriter) Write(p []byte) (int, error) {
	lenp := len(p)

	if w.pos+lenp > cap(w.b) {
		return 0, io.ErrShortWrite
	}

	w.pos += copy(w.b[w.pos:], p)

	return lenp, nil
}

// Bytes returns the bytes already written to the fixed writer.
func (w *FixedWriter) Bytes() []byte {
	return w.b
}

// newFixedWriter returns a new io.Writer that will error once more bytes than
// the specified max have been written.
func NewFixedWriter(max int) *FixedWriter {
	b := make([]byte, max)
	fw := FixedWriter{b, 0}

	return &fw
}

// fixedReader implements the io.Reader interface and intentionally allows
// testing of error paths by forcing short reads.
type FixedReader struct {
	buf   []byte
	pos   int
	iobuf *bytes.Buffer
}

// Read reads the next len(p) bytes from the fixed reader.  When the number of
// bytes read would exceed the maximum number of allowed bytes to be read from
// the fixed writer, an error is returned.
//
// This satisfies the io.Reader interface.
func (fr *FixedReader) Read(p []byte) (int, error) {
	n, err := fr.iobuf.Read(p)
	if err != nil {
		return 0, err
	}

	fr.pos += n

	return n, nil
}

// newFixedReader returns a new io.Reader that will error once more bytes than
// the specified max have been read.
func NewFixedReader(max int, buf []byte) *FixedReader {
	b := make([]byte, max)
	if buf != nil {
		copy(b, buf)
	}

	iobuf := bytes.NewBuffer(b)
	fr := FixedReader{b, 0, iobuf}

	return &fr
}
