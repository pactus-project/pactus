package util

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIfPermissionsAreNotEnforced(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("file permissions are not enforced")
	}
}

// Checking a path is a read-only operation.
func TestIsValidDirPathCreatesNothing(t *testing.T) {
	dir := t.TempDir()

	assert.True(t, IsValidDirPath(filepath.Join(dir, "a", "b", "c")))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries, "no parent directory should be created")
}

func TestIsValidDirPathWithNulByte(t *testing.T) {
	assert.False(t, IsValidDirPath("a\x00b"))
}

func TestIsValidDirPathWithReadOnlyParent(t *testing.T) {
	skipIfPermissionsAreNotEnforced(t)

	dir := t.TempDir()
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	assert.False(t, IsValidDirPath(filepath.Join(dir, "a", "b")))
}

func TestIsDirNotExistsOrEmptyWithUnreadableDirectory(t *testing.T) {
	skipIfPermissionsAreNotEnforced(t)

	dir := filepath.Join(t.TempDir(), "unreadable")
	require.NoError(t, os.Mkdir(dir, 0o750))
	require.NoError(t, os.Chmod(dir, 0o000))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	assert.NotPanics(t, func() {
		assert.False(t, IsDirNotExistsOrEmpty(dir))
	})
}

func TestIsDirNotExistsOrEmptyWithUnsearchableParent(t *testing.T) {
	skipIfPermissionsAreNotEnforced(t)

	parent := t.TempDir()
	child := filepath.Join(parent, "child")
	require.NoError(t, os.Mkdir(child, 0o750))
	require.NoError(t, os.Chmod(parent, 0o000))
	t.Cleanup(func() { _ = os.Chmod(parent, 0o700) })

	assert.False(t, IsDirNotExistsOrEmpty(child), "a directory that can't be inspected is not a missing one")
}

func TestIsDirNotExistsOrEmpty(t *testing.T) {
	dir := t.TempDir()

	assert.True(t, IsDirNotExistsOrEmpty(filepath.Join(dir, "missing")))
	assert.True(t, IsDirNotExistsOrEmpty(dir))

	file := filepath.Join(dir, "file")
	require.NoError(t, os.WriteFile(file, []byte("data"), 0o600))
	assert.False(t, IsDirNotExistsOrEmpty(dir))
	assert.False(t, IsDirNotExistsOrEmpty(file))
}
