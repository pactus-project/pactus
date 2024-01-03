package store

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	conf := DefaultConfig()

	err := conf.BasicCheck()
	assert.ErrorIs(t, InvalidConfigError{"cache size set to zero"}, err)

	conf.TxCacheSize = 1
	conf.SortitionCacheSize = 1
	err = conf.BasicCheck()
	assert.NoError(t, err)

	if runtime.GOOS != "windows" {
		conf.Path = util.TempDirPath()
		assert.NoError(t, conf.BasicCheck())
		assert.Equal(t, conf.StorePath(), conf.Path+"/store.db")
	}
}
