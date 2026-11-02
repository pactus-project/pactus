package util

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A file that is replaced must never be modified in place: a reader that
// still holds the old file (here, a hard link) keeps the old content.
func TestWriteFileReplacesTheFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "wallet")
	link := filepath.Join(dir, "wallet.link")

	require.NoError(t, WriteFile(file, []byte("old content")))
	if err := os.Link(file, link); err != nil {
		t.Skipf("hard links are not supported: %v", err)
	}

	require.NoError(t, WriteFile(file, []byte("new content")))

	data, err := os.ReadFile(file)
	require.NoError(t, err)
	assert.Equal(t, "new content", string(data))

	data, err = os.ReadFile(link)
	require.NoError(t, err)
	assert.Equal(t, "old content", string(data), "the old file should not be rewritten in place")
}

func TestWriteFileLeavesNoTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "config.toml")

	require.NoError(t, WriteFile(file, []byte("first")))
	require.NoError(t, WriteFile(file, []byte("second")))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "config.toml", entries[0].Name())
}

func TestWriteFileKeepsPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not supported")
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "file")

	require.NoError(t, WriteFile(file, []byte("data")))
	info, err := os.Stat(file)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(), "new files are private")

	require.NoError(t, os.Chmod(file, 0o640))
	require.NoError(t, WriteFile(file, []byte("other data")))
	info, err = os.Stat(file)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o640), info.Mode().Perm(), "the mode of an existing file is kept")
}

func TestWriteFileFollowsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	link := filepath.Join(dir, "link")

	require.NoError(t, os.WriteFile(target, []byte("old"), 0o600))
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links are not supported: %v", err)
	}

	require.NoError(t, WriteFile(link, []byte("new")))

	info, err := os.Lstat(link)
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink, "the link should still be a link")

	data, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "new", string(data))
}

type failingFile struct {
	bytes.Buffer
	syncErr  error
	closeErr error
	closed   bool
}

func (f *failingFile) Sync() error { return f.syncErr }

func (f *failingFile) Close() error {
	f.closed = true

	return f.closeErr
}

func TestCopyAndClose(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		out := &failingFile{}

		require.NoError(t, copyAndClose(out, strings.NewReader("data")))
		assert.Equal(t, "data", out.String())
		assert.True(t, out.closed)
	})

	t.Run("close error is returned", func(t *testing.T) {
		out := &failingFile{closeErr: errors.New("close failed")}

		require.ErrorIs(t, copyAndClose(out, strings.NewReader("data")), out.closeErr)
	})

	t.Run("sync error is returned and the file is closed", func(t *testing.T) {
		out := &failingFile{syncErr: errors.New("sync failed")}

		require.ErrorIs(t, copyAndClose(out, strings.NewReader("data")), out.syncErr)
		assert.True(t, out.closed)
	})
}

func TestMoveByCopyRemovesPartialDestination(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("an unreadable file can't be created")
	}

	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(t.TempDir(), "dst")
	require.NoError(t, os.MkdirAll(src, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(src, "b.txt"), []byte("b"), 0o000))

	require.Error(t, moveByCopy(src, dst))

	assert.False(t, PathExists(dst), "the partial copy should be removed")
	assert.True(t, PathExists(filepath.Join(src, "a.txt")), "the source should be kept")
	assert.True(t, PathExists(filepath.Join(src, "b.txt")), "the source should be kept")
}

func TestMoveByCopy(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(t.TempDir(), "dst")
	require.NoError(t, os.MkdirAll(filepath.Join(src, "sub"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(src, "sub", "a.txt"), []byte("a"), 0o600))

	require.NoError(t, moveByCopy(src, dst))

	assert.False(t, PathExists(src))
	data, err := os.ReadFile(filepath.Join(dst, "sub", "a.txt"))
	require.NoError(t, err)
	assert.Equal(t, "a", string(data))
}
