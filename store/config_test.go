package store

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
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
