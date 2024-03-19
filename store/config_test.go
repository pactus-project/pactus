package store

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	conf := DefaultConfig()

	conf.TxCacheSize = 0
	err := conf.BasicCheck()
	assert.ErrorIs(t, ConfigError{"cache size set to zero"}, err)

	conf.TxCacheSize = 1
	err = conf.BasicCheck()
	assert.NoError(t, err)

	conf.Path = util.TempDirPath()
	assert.NoError(t, conf.BasicCheck())

	if runtime.GOOS != "windows" {
		assert.Equal(t, conf.StorePath(), conf.Path+"/store.db")
	} else {
		assert.Equal(t, conf.StorePath(), conf.Path+"\\store.db")
	}
}
