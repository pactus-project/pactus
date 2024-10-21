package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const invalidDirName = "/invalid:path/\x00*folder?\\CON"

func TestWriteFile(t *testing.T) {
	p := TempDirPath()
	d := []byte("some-data")
	assert.NoError(t, WriteFile(p+"/d.dat", d))
	assert.NoError(t, WriteFile(p+"/another-folder/d.dat", d))
}

func TestEmptyPath(t *testing.T) {
	p := TempDirPath()
	assert.Equal(t, p, MakeAbs(p))
	assert.True(t, IsDirEmpty(p))

	f := TempFilePath()
	d := []byte("pactus")
	assert.NoError(t, WriteFile(f, d))
	o, err := ReadFile(f)
	assert.NoError(t, err)
	assert.Equal(t, d, o)
}

func TestAbsPath(t *testing.T) {
	abs := MakeAbs(".")
	assert.True(t, IsAbsPath(abs))
	assert.False(t, IsAbsPath("abs"))
	assert.False(t, IsDirEmpty(abs))
	assert.False(t, IsDirNotExistsOrEmpty(abs))
}

func TestTempDir(t *testing.T) {
	tmpDir := TempDirPath()

	assert.True(t, IsAbsPath(tmpDir))
	assert.True(t, IsDirEmpty(tmpDir))
	assert.True(t, PathExists(tmpDir))
	assert.True(t, IsDirNotExistsOrEmpty(tmpDir))
}

func TestTempFile(t *testing.T) {
	tmpFile := TempFilePath()

	assert.True(t, IsAbsPath(tmpFile))
	t.Run("Should panic because it doesn't exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		IsDirEmpty(tmpFile)
	})
	assert.False(t, PathExists(tmpFile))
	assert.True(t, IsDirNotExistsOrEmpty(tmpFile))
	assert.NoError(t, Mkdir(tmpFile))
	assert.True(t, IsDirNotExistsOrEmpty(tmpFile))
	assert.True(t, IsDirEmpty(tmpFile)) // no panic now
}

func isRoot() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		fmt.Println(err)
	}

	return i == 0
}

func TestIsValidPath(t *testing.T) {
	// To pass this tests inside docker
	if runtime.GOOS != "windows" && !isRoot() {
		assert.False(t, IsValidDirPath("/root"))
		assert.False(t, IsValidDirPath("/test"))
	}
	assert.False(t, IsValidDirPath(invalidDirName))
	assert.False(t, IsValidDirPath("./io_test.go"))
	assert.True(t, IsValidDirPath("/tmp"))
	assert.True(t, IsValidDirPath("/tmp/pactus"))
}

func TestMoveDirectory(t *testing.T) {
	t.Run("DestinationDirectoryExistsAndNotEmpty", func(t *testing.T) {
		srcDir := TempDirPath()
		dstDir := TempDirPath()

		err := MoveDirectory(srcDir, dstDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "destination directory")
	})

	t.Run("ParentDirectoryCreationFailure", func(t *testing.T) {
		srcDir := TempDirPath()

		err := MoveDirectory(srcDir, invalidDirName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create parent directories")
	})

	t.Run("SourceDirectoryRenameFailure", func(t *testing.T) {
		srcDir := TempDirPath()
		dstDir := TempDirPath()

		err := os.RemoveAll(dstDir)
		assert.NoError(t, err)

		// Remove the source directory to simulate the rename failure
		err = os.RemoveAll(srcDir)
		assert.NoError(t, err)

		err = MoveDirectory(srcDir, dstDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to move directory")
	})

	t.Run("MoveDirectorySuccess", func(t *testing.T) {
		// Create temporary directories
		srcDir := TempDirPath()
		dstDir := TempDirPath()
		defer func() { _ = os.RemoveAll(srcDir) }()
		defer func() { _ = os.RemoveAll(dstDir) }()

		// Create a subdirectory in the source directory
		subDir := filepath.Join(srcDir, "subdir")
		err := Mkdir(subDir)
		assert.NoError(t, err)

		// Create multiple files in the subdirectory
		files := []struct {
			name    string
			content string
		}{
			{"file1.txt", "content 1"},
			{"file2.txt", "content 2"},
		}

		for _, file := range files {
			filePath := filepath.Join(subDir, file.name)
			err = WriteFile(filePath, []byte(file.content))
			assert.NoError(t, err)
		}

		// Move the directory
		dstDirPath := filepath.Join(dstDir, "movedir")
		err = MoveDirectory(srcDir, dstDirPath)
		assert.NoError(t, err)

		// Assert the source directory no longer exists
		assert.False(t, PathExists(srcDir))

		// Assert the destination directory exists
		assert.True(t, PathExists(dstDirPath))

		// Verify that all files have been moved and their contents are correct
		for _, file := range files {
			movedFilePath := filepath.Join(dstDirPath, "subdir", file.name)
			data, err := ReadFile(movedFilePath)
			assert.NoError(t, err)
			assert.Equal(t, file.content, string(data))
		}
	})
}

func TestSanitizeArchivePath(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}

	baseDir := "/safe/directory"

	tests := []struct {
		name      string
		inputPath string
		expected  string
		expectErr bool
	}{
		{"Valid path", "file.txt", "/safe/directory/file.txt", false},
		{"Valid path in subdirectory", "subdir/file.txt", "/safe/directory/subdir/file.txt", false},
		{"Path with parent directory traversal", "../outside/file.txt", "", true},
		{"Absolute path outside base directory", "/etc/passwd", "/safe/directory/etc/passwd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeArchivePath(baseDir, tt.inputPath)
			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
				assert.Empty(t, result, "Expected empty result due to error")
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
				assert.Equal(t, tt.expected, result, "Sanitized path did not match expected")
			}
		})
	}
}

func TestListFilesInDir(t *testing.T) {
	tmpDir := TempDirPath()

	file1Path := filepath.Join(tmpDir, "public_file")
	file1, err := os.Create(file1Path)
	require.NoError(t, err)
	require.NoError(t, file1.Close())

	file2Path := filepath.Join(tmpDir, ".hidden_file")
	file2, err := os.Create(file2Path)
	require.NoError(t, err)
	require.NoError(t, file2.Close())

	err = os.Mkdir(filepath.Join(tmpDir, "directory"), 0o750)
	require.NoError(t, err)

	files, err := ListFilesInDir(tmpDir)
	require.NoError(t, err)

	assert.Len(t, files, 2)
	assert.Contains(t, files, file1Path)
	assert.Contains(t, files, file2Path)
}
