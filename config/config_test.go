package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestSanityCheck(t *testing.T) {

	conf := DefaultConfig()
	assert.NoError(t, conf.SanityCheck())
}

func TestTOML(t *testing.T) {
	f := util.TempFilePath()
	f += ".toml"
	conf1 := DefaultConfig()
	conf1.Store.Path = "abc"
	conf1.Sync.Moniker = "Test1"
	conf1.Consensus.ChangeProposerTimeout = 22
	assert.NoError(t, conf1.SaveToFile(f))
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
