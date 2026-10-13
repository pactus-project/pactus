package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Checking a directory should never destroy the files of the user.
func TestIsValidDirPathKeepsExistingFiles(t *testing.T) {
	dir := t.TempDir()

	// The old implementation used this name to test if the directory is writable.
	userFile := filepath.Join(dir, "test")
	require.NoError(t, os.WriteFile(userFile, []byte("important data"), 0o600))

	assert.True(t, IsValidDirPath(dir))

	data, err := os.ReadFile(userFile)
	require.NoError(t, err, "the file of the user should still exist")
	assert.Equal(t, "important data", string(data))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "no temporary file should be left in the directory")
}

func TestIsValidDirPath(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		assert.True(t, IsValidDirPath(t.TempDir()))
	})

	t.Run("path is a file", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "file")
		require.NoError(t, os.WriteFile(file, []byte("data"), 0o600))

		assert.False(t, IsValidDirPath(file))
	})

	t.Run("directory doesn't exist yet", func(t *testing.T) {
		assert.True(t, IsValidDirPath(filepath.Join(t.TempDir(), "new")))
	})

	t.Run("directory can't be created", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "file")
		require.NoError(t, os.WriteFile(file, []byte("data"), 0o600))

		assert.False(t, IsValidDirPath(filepath.Join(file, "sub")))
	})
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	require.NoError(t, os.WriteFile(src, []byte("content"), 0o640))

	require.NoError(t, CopyFile(src, dst))

	data, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Equal(t, "content", string(data))

	srcInfo, err := os.Stat(src)
	require.NoError(t, err)
	dstInfo, err := os.Stat(dst)
	require.NoError(t, err)
	assert.Equal(t, srcInfo.Mode(), dstInfo.Mode(), "permissions should be preserved")

	t.Run("source doesn't exist", func(t *testing.T) {
		require.Error(t, CopyFile(filepath.Join(dir, "missing"), filepath.Join(dir, "other")))
	})

	t.Run("destination can't be created", func(t *testing.T) {
		require.Error(t, CopyFile(src, filepath.Join(dir, "missing-dir", "dst")))
	})
}

func TestCopyDirectory(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(t.TempDir(), "dst")
	require.NoError(t, os.MkdirAll(filepath.Join(src, "a", "b"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(src, "root.txt"), []byte("root"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(src, "a", "b", "deep.txt"), []byte("deep"), 0o600))
	require.NoError(t, os.MkdirAll(filepath.Join(src, "empty"), 0o750))

	require.NoError(t, CopyDirectory(src, dst))

	data, err := os.ReadFile(filepath.Join(dst, "root.txt"))
	require.NoError(t, err)
	assert.Equal(t, "root", string(data))
	data, err = os.ReadFile(filepath.Join(dst, "a", "b", "deep.txt"))
	require.NoError(t, err)
	assert.Equal(t, "deep", string(data))
	assert.True(t, IsDir(filepath.Join(dst, "empty")), "empty directories are copied")

	t.Run("source doesn't exist", func(t *testing.T) {
		require.Error(t, CopyDirectory(filepath.Join(t.TempDir(), "missing"), t.TempDir()))
	})
}

func TestMoveDirectoryFailures(t *testing.T) {
	t.Run("destination already exists", func(t *testing.T) {
		src := t.TempDir()
		dst := t.TempDir()

		require.Error(t, MoveDirectory(src, dst))
		assert.True(t, PathExists(src), "the source must be kept")
	})

	t.Run("source doesn't exist", func(t *testing.T) {
		require.Error(t, MoveDirectory(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "dst")))
	})
}

// SecureJoinPath is used to extract the archives and to find the wallets.
func TestSecureJoinPathPortable(t *testing.T) {
	baseDir := filepath.Join(t.TempDir(), "base")

	tests := []struct {
		name  string
		path  string
		valid bool
	}{
		{"file", "file.txt", true},
		{"file in a sub directory", filepath.Join("a", "b", "file.txt"), true},
		{"parent directory is cleaned inside the base", filepath.Join("a", "..", "file.txt"), true},
		{"parent directory", filepath.Join("..", "file.txt"), false},
		{"traversal after a directory", filepath.Join("a", "..", "..", "file.txt"), false},
		{"traversal to a directory with the same prefix", filepath.Join("..", "base-evil", "file.txt"), false},
		{"base directory itself", ".", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath, err := SecureJoinPath(baseDir, tt.path)
			if tt.valid {
				require.NoError(t, err)
				assert.Equal(t, filepath.Join(baseDir, tt.path), fullPath)
				assert.True(t, strings.HasPrefix(fullPath, baseDir+string(os.PathSeparator)))
			} else {
				require.Error(t, err)
				assert.Empty(t, fullPath)
			}
		})
	}
}

func TestFixedWriter(t *testing.T) {
	w := NewFixedWriter(5)

	n, err := w.Write([]byte("abc"))
	require.NoError(t, err)
	assert.Equal(t, 3, n)

	// A write that doesn't fit is refused entirely, and the content doesn't change.
	n, err = w.Write([]byte("def"))
	require.ErrorIs(t, err, io.ErrShortWrite)
	assert.Zero(t, n)

	n, err = w.Write([]byte("de"))
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte("abcde"), w.Bytes())

	_, err = w.Write([]byte("f"))
	require.ErrorIs(t, err, io.ErrShortWrite)
}

func TestFixedReader(t *testing.T) {
	r := NewFixedReader(5, []byte("abc"))

	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, []byte{'a', 'b', 'c', 0, 0}, data, "the missing bytes are zero")

	n, err := r.Read(make([]byte, 1))
	require.ErrorIs(t, err, io.EOF)
	assert.Zero(t, n)

	empty, err := io.ReadAll(NewFixedReader(3, nil))
	require.NoError(t, err)
	assert.Equal(t, []byte{0, 0, 0}, empty)
}
