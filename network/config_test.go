package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	conf := DefaultConfig()

	conf.EnableRelay = true
	assert.Error(t, conf.BasicCheck())

	conf.Listens = []string{""}
	assert.Error(t, conf.BasicCheck())

	conf.Listens = []string{"127.0.0.1"}
	assert.Error(t, conf.BasicCheck())

	conf.Listens = []string{"/ip4"}
	assert.Error(t, conf.BasicCheck())

	conf.RelayAddrs = []string{"/ip4"}
	assert.Error(t, conf.BasicCheck())

	conf.RelayAddrs = []string{}
	conf.Listens = []string{}

	conf.RelayAddrs = []string{"/ip4/127.0.0.1"}
	assert.NoError(t, conf.BasicCheck())

	conf.Listens = []string{"/ip4/127.0.0.1"}
	assert.NoError(t, conf.BasicCheck())
}
