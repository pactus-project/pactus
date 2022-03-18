package store

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	if runtime.GOOS != "windows" {
		c.Path = "/tmp/zarb/data"
		assert.NoError(t, c.SanityCheck())
		assert.Equal(t, c.StorePath(), "/tmp/zarb/data/store.db")
	}
}
