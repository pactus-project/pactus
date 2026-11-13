package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	conf := DefaultConfig()

	require.NoError(t, conf.BasicCheck())
	assert.False(t, conf.Enable)
	assert.Empty(t, conf.Listen)
	assert.Empty(t, conf.Origins)
}
