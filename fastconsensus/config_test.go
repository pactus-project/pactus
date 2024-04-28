package fastconsensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c1 := DefaultConfig()
	c2 := DefaultConfig()
	c3 := DefaultConfig()
	c4 := DefaultConfig()
	c5 := DefaultConfig()
	assert.NoError(t, c1.BasicCheck())

	c2.ChangeProposerDelta = 0 * time.Second
	assert.ErrorIs(t, c2.BasicCheck(), ConfigError{Reason: "change proposer delta can't be negative"})

	c3.ChangeProposerTimeout = 0 * time.Second
	assert.ErrorIs(t, c3.BasicCheck(), ConfigError{Reason: "timeout for change proposer can't be negative"})

	c4.ChangeProposerTimeout = -1 * time.Second
	assert.ErrorIs(t, c4.BasicCheck(), ConfigError{Reason: "timeout for change proposer can't be negative"})

	c5.MinimumAvailabilityScore = 1.5
	assert.ErrorIs(t, c5.BasicCheck(), ConfigError{Reason: "minimum availability score can't be negative or more than 1"})

	c5.MinimumAvailabilityScore = -0.8
	assert.ErrorIs(t, c5.BasicCheck(), ConfigError{Reason: "minimum availability score can't be negative or more than 1"})
}

func TestCalculateChangeProposerTimeout(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, c.CalculateChangeProposerTimeout(0), c.ChangeProposerTimeout)
	assert.Equal(t, c.CalculateChangeProposerTimeout(1), c.ChangeProposerTimeout+c.ChangeProposerDelta)
	assert.Equal(t, c.CalculateChangeProposerTimeout(4), c.ChangeProposerTimeout+(4*c.ChangeProposerDelta))
}
