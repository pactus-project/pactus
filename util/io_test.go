package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyPaht(t *testing.T) {
	p := TempDirPath()
	assert.True(t, IsDirEmpty(p))

	f := TempFilePath()
	d := []byte("zarb")
	assert.NoError(t, WriteFile(f, d))
	o, err := ReadFile(f)
	assert.NoError(t, err)
	assert.Equal(t, d, o)
}
