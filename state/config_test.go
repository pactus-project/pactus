package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	t.Run("Invalid reward address", func(t *testing.T) {
		c.RewardAddress = "invalid"
		assert.Error(t, c.SanityCheck())
	})
}
