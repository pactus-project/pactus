package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf(t *testing.T) {
	conf := DefaultConfig()

	conf.EnableRelay = true
	assert.Error(t, conf.SanityCheck())

	conf.EnableRelay = false
	assert.NoError(t, conf.SanityCheck())
}
