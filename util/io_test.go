package util

import (
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	p := TempDirPath()
	d := []byte("some-data")
	assert.NoError(t, WriteFile(p+"/d.dat", d))
	assert.NoError(t, WriteFile(p+"/another-folder/d.dat", d))
}

func TestEmptyPath(t *testing.T) {
	p := TempDirPath()
	assert.True(t, IsDirEmpty(p))

	f := TempFilePath()
	d := []byte("zarb")
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
	t.Run("Should panic because it doesn't exists", func(t *testing.T) {
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
	if !isRoot() {
		assert.False(t, IsValidDirPath("/root"))
		assert.False(t, IsValidDirPath("/test"))
	}
	assert.False(t, IsValidDirPath("./io_test.go"))
	assert.True(t, IsValidDirPath("/tmp"))
	assert.True(t, IsValidDirPath("/tmp/zarb"))
}
