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
