package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	t.Run("Invalid mintbase address", func(t *testing.T) {
		c.MintbaseAddress = "invalid"
		assert.Error(t, c.SanityCheck())
	})
}
