package store

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	if runtime.GOOS != "windows" {
		c.Path = util.TempDirPath()
		assert.NoError(t, c.SanityCheck())
		assert.Equal(t, c.StorePath(), c.Path+"/store.db")
	}
}
