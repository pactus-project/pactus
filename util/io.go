package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func IsAbsPath(path string) bool {
	return filepath.IsAbs(path)
}

func MakeAbs(path string) string {
	if IsAbsPath(path) {
		return path
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Clean(filepath.Join(wd, path))
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

func Mkdir(path string) error {
	// create the directory
	if err := os.MkdirAll(path, 0o750); err != nil {
		return fmt.Errorf("could not create directory %s", path)
	}

	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func TempDirPath() string {
	path, err := os.MkdirTemp("", "pactus*")
	if err != nil {
		panic(err)
	}

	return path
}

func TempFilePath() string {
	return filepath.Join(TempDirPath(), "file")
}

// IsDirEmpty checks if a directory is empty.
func IsDirEmpty(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	// read in ONLY one file
	_, err = file.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	return errors.Is(err, io.EOF)
}

// IsDirNotExistsOrEmpty checks if the path exists and, if so, whether the directory is empty.
func IsDirNotExistsOrEmpty(path string) bool {
	if !PathExists(path) {
		return true
	}

	return IsDirEmpty(path)
}

func IsValidDirPath(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			if err := os.WriteFile(path+"/test", []byte{}, 0o600); err != nil {
				return false
			}
			_ = os.Remove(path + "/test")

			return true
		}

		return false
	}

	if err := Mkdir(path); err != nil {
		return false
	}
	_ = os.Remove(path)

	return true
}

// TODO: move these to a test suite

// FixedWriter implements the io.Writer interface and intentionally allows
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
func (w *FixedWriter) Write(data []byte) (int, error) {
	lenp := len(data)

	if w.pos+lenp > cap(w.b) {
		return 0, io.ErrShortWrite
	}

	w.pos += copy(w.b[w.pos:], data)

	return lenp, nil
}

// Bytes returns the bytes already written to the fixed writer.
func (w *FixedWriter) Bytes() []byte {
	return w.b
}

// NewFixedWriter returns a new io.Writer that will error once more bytes than
// the specified max have been written.
func NewFixedWriter(max int) *FixedWriter {
	b := make([]byte, max)
	fw := FixedWriter{b, 0}

	return &fw
}

// FixedReader implements the io.Reader interface and intentionally allows
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
	count, err := fr.iobuf.Read(p)
	if err != nil {
		return 0, err
	}

	fr.pos += count

	return count, nil
}

// NewFixedReader returns a new io.Reader that will error once more bytes than
// the specified max have been read.
func NewFixedReader(max int, data []byte) *FixedReader {
	buf := make([]byte, max)
	if data != nil {
		copy(buf, data)
	}

	iobuf := bytes.NewBuffer(buf)
	fr := FixedReader{buf, 0, iobuf}

	return &fr
}

// MoveDirectory moves a directory from srcDir to dstDir, including all its contents.
// If dstDir already exists and is not empty, it returns an error.
// If the parent directory of dstDir does not exist, it will be created.
func MoveDirectory(srcDir, dstDir string) error {
	if PathExists(dstDir) {
		return fmt.Errorf("destination directory %s already exists", dstDir)
	}

	// Get the parent directory of the destination directory
	parentDir := filepath.Dir(dstDir)
	if err := Mkdir(parentDir); err != nil {
		return fmt.Errorf("failed to create parent directories for %s: %w", dstDir, err)
	}

	if err := os.Rename(srcDir, dstDir); err != nil {
		return fmt.Errorf("failed to move directory from %s to %s: %w", srcDir, dstDir, err)
	}

	return nil
}

// SanitizeArchivePath mitigates the "Zip Slip" vulnerability by sanitizing archive file paths.
// It ensures that the file path is contained within the specified base directory to prevent directory
// traversal attacks. For more details on the vulnerability, see https://snyk.io/research/zip-slip-vulnerability.
func SanitizeArchivePath(baseDir, archivePath string) (fullPath string, err error) {
	fullPath = filepath.Join(baseDir, archivePath)
	if strings.HasPrefix(fullPath, filepath.Clean(baseDir)) {
		return fullPath, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", archivePath)
}

// ListFilesInDir return list of files in directory.
func ListFilesInDir(dir string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
