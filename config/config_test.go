package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestTOML(t *testing.T) {
	f := util.TempFilePath()
	f += ".toml"
	conf1 := DefaultConfig()
	assert.NoError(t, conf1.SaveToFile(f))
	assert.NoError(t, conf1.Check())
	conf2, err := LoadFromFile(f)
	assert.NoError(t, err)
	assert.Equal(t, conf1, conf2)
}

func TestJSON(t *testing.T) {
	f := util.TempFilePath()
	f += ".json"
	conf1 := DefaultConfig()
	assert.NoError(t, conf1.SaveToFile(f))
	conf2, err := LoadFromFile(f)
	assert.NoError(t, err)
	assert.Equal(t, conf1, conf2)
}

func TestInvalidFile(t *testing.T) {
	f := util.TempFilePath()
	f += ".inv"
	conf1 := DefaultConfig()
	assert.Error(t, conf1.SaveToFile(f))
	_, err := LoadFromFile(f)
	assert.Error(t, err)
}
