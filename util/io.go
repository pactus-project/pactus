package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func IsAbsPath(path string) bool {
	return filepath.IsAbs(path)
}

// MakeAbs returns the absolute form of path. If the working directory can't be
// determined, the cleaned relative path is returned, so the error shows up where
// the path is used instead of aborting the process here.
func MakeAbs(path string) string {
	return makeAbs(path, os.Getwd)
}

func makeAbs(path string, getwd func() (string, error)) string {
	if IsAbsPath(path) {
		return path
	}
	wd, err := getwd()
	if err != nil {
		return filepath.Clean(path)
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

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

// MakeTempDir creates a new temporary directory and returns its path.
func MakeTempDir() (string, error) {
	return os.MkdirTemp("", "pactus*")
}

// TempDirPath creates a new temporary directory and returns its path.
// It panics if the directory can't be created, so it is meant for tests;
// other code should use MakeTempDir.
func TempDirPath() string {
	path, err := MakeTempDir()
	if err != nil {
		panic(err)
	}

	return path
}

func TempFilePath() string {
	return filepath.Join(TempDirPath(), "file")
}

// IsDirEmpty checks if a directory is empty.
// It returns an error if the directory can't be read, for example when it doesn't exist.
func IsDirEmpty(path string) (bool, error) {
	return dirEmpty(path)
}

// IsDirNotExistsOrEmpty checks if the path exists and, if so, whether the directory is empty.
// A path that can't be inspected, for example because of its permissions, is not reported as missing.
func IsDirNotExistsOrEmpty(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return errors.Is(err, fs.ErrNotExist)
	}

	empty, err := dirEmpty(path)

	return err == nil && empty
}

func dirEmpty(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = file.Close()
	}()

	// read in ONLY one file
	_, err = file.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	return errors.Is(err, io.EOF), nil
}

// IsValidDirPath checks if a directory exists, or could be created, and is writable.
// It doesn't create anything: a missing directory is checked through its closest existing parent.
func IsValidDirPath(path string) bool {
	if strings.ContainsRune(path, 0) {
		return false
	}

	for current := path; ; {
		fi, err := os.Stat(current)
		if err == nil {
			return fi.IsDir() && isWritableDir(current)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return false
		}

		parent := filepath.Dir(current)
		if parent == current {
			return false
		}
		current = parent
	}
}

// isWritableDir creates a temporary file with a unique name to check that the directory is writable,
// without overwriting or removing the files of the user.
func isWritableDir(path string) bool {
	file, err := os.CreateTemp(path, ".write-check-*")
	if err != nil {
		return false
	}
	_ = file.Close()
	_ = os.Remove(file.Name())

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

	parentDir := filepath.Dir(dstDir)
	if err := Mkdir(parentDir); err != nil {
		return fmt.Errorf("failed to create parent directories for %s: %w", dstDir, err)
	}

	err := os.Rename(srcDir, dstDir)
	if err != nil {
		// To prevent invalid cross-device link perform a manual copy
		return moveByCopy(srcDir, dstDir)
	}

	return nil
}

// moveByCopy copies srcDir to dstDir and removes srcDir once the copy is complete.
// A partial copy is removed, and srcDir is kept, when the copy fails.
func moveByCopy(srcDir, dstDir string) error {
	if err := CopyDirectory(srcDir, dstDir); err != nil {
		_ = os.RemoveAll(dstDir)

		return fmt.Errorf("failed to move directory from %s to %s: %w", srcDir, dstDir, err)
	}

	// Remove source directory after successful copy
	if err := os.RemoveAll(srcDir); err != nil {
		return fmt.Errorf("failed to remove source directory %s: %w", srcDir, err)
	}

	return nil
}

func CopyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return CopyFile(path, dstPath)
	})
}

func CopyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = input.Close()
	}()

	info, err := input.Stat()
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if err := copyAndClose(out, input); err != nil {
		return err
	}

	// Preserve file permissions
	return os.Chmod(dst, info.Mode())
}

// copyAndClose copies the data to out, flushes it to the disk and closes it.
// Errors reported while flushing or closing are returned like write errors.
func copyAndClose(out interface {
	io.WriteCloser
	Sync() error
}, input io.Reader,
) error {
	if _, err := io.Copy(out, input); err != nil {
		_ = out.Close()

		return err
	}
	if err := out.Sync(); err != nil {
		_ = out.Close()

		return err
	}

	return out.Close()
}

// SecureJoinPath mitigates the "Zip Slip" vulnerability by sanitizing archive file paths.
// It ensures that the file path is contained within the specified base directory to prevent directory
// traversal attacks. For more details on the vulnerability, see https://snyk.io/research/zip-slip-vulnerability.
func SecureJoinPath(baseDir, filePath string) (fullPath string, err error) {
	fullPath = filepath.Join(baseDir, filePath)
	if strings.HasPrefix(fullPath, filepath.Clean(baseDir)+string(os.PathSeparator)) {
		return fullPath, nil
	}

	return "", fmt.Errorf("illegal file path: %s", filePath)
}

type listFilesConfig struct {
	excludeDirs bool
}

var defaultListFilesConfig = listFilesConfig{
	excludeDirs: false,
}

type ListFilesOption func(*listFilesConfig)

// ExcludeDirectories prevents entering subdirectories.
func ExcludeDirectories() ListFilesOption {
	return func(c *listFilesConfig) {
		c.excludeDirs = true
	}
}

// ListFilesInDir return list of files in directory.
func ListFilesInDir(dir string, opts ...ListFilesOption) ([]string, error) {
	cfg := defaultListFilesConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	files := make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if cfg.excludeDirs && entry.IsDir() {
			continue
		}

		fulPath := filepath.Join(dir, entry.Name())
		files = append(files, fulPath)
	}

	return files, nil
}

// LimitReaderClose returns a ReadCloser that reads from r
// but stops with EOF after n bytes.
// The returned ReadCloser always forwards Close to r.
// The underlying implementation is a *LimitedReadCloser.
func LimitReaderClose(r io.ReadCloser, n int64) io.ReadCloser {
	return &LimitedReadCloser{R: r, N: n}
}

// LimitedReadCloser reads from R but limits the amount of
// data returned to just N bytes. Each call to Read updates N
// to reflect the new amount remaining. Read returns EOF when
// N <= 0 or when the underlying R returns EOF.
// Close always calls R.Close.
type LimitedReadCloser struct {
	R io.ReadCloser // underlying reader with close support
	N int64         // max bytes remaining
}

func (l *LimitedReadCloser) Read(buf []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(buf)) > l.N {
		buf = buf[:l.N]
	}
	n, err = l.R.Read(buf)
	l.N -= int64(n)

	return n, err
}

func (l *LimitedReadCloser) Close() error {
	return l.R.Close()
}
