package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	conf := DefaultConfig()

	conf.EnableRelay = true
	assert.Error(t, conf.SanityCheck())

	conf.EnableRelay = false
	assert.NoError(t, conf.SanityCheck())

	inValidAddresses := []string{
		"/ip4/1wg.1/tcp/21777/p2p8VccQ8888888Wyc",
		"",
		"/ip4",
		"127.0.0.1",
	}
	assert.Error(t, validateAddresses(inValidAddresses))

	validAddresses := []string{
		"/ip4/127.0.0.1",
	}
	assert.Equal(t, nil, validateAddresses(validAddresses))
}
