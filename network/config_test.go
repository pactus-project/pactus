package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	conf := DefaultConfig()
	conf.EnableRelay = true
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.ListenAddrStrings = []string{""}
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.ListenAddrStrings = []string{"127.0.0.1"}
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.ListenAddrStrings = []string{"/ip4"}
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.PublicAddrString = "/ip4"
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.RelayAddrStrings = []string{"/ip4/127.0.0.1/"}
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
	assert.Error(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.PublicAddrString = "/ip4/127.0.0.1/"
	assert.NoError(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.RelayAddrStrings = []string{"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"}
	assert.NoError(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.ListenAddrStrings = []string{"/ip4/127.0.0.1"}
	assert.NoError(t, conf.BasicCheck())

	conf = DefaultConfig()
	conf.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"}
	assert.NoError(t, conf.BasicCheck())
}
