package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	conf := DefaultConfig()

	assert.True(t, conf.LockMode)
	require.NoError(t, conf.BasicCheck())
}
