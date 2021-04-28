package consensus

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
	assert.NoError(t, c1.SanityCheck())

	c2.ChangeProposerDelta = 0 * time.Second
	assert.Error(t, c2.SanityCheck())

	c3.ChangeProposerTimeout = 0 * time.Second
	assert.Error(t, c3.SanityCheck())

	c4.ChangeProposerTimeout = -1 * time.Second
	assert.Error(t, c4.SanityCheck())

	c5.QueryProposalTimeout = -1 * time.Second
	assert.Error(t, c5.SanityCheck())
}

func TestCalculateChangeProposerTimeout(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, c.CalculateChangeProposerTimeout(0), c.ChangeProposerTimeout)
	assert.Equal(t, c.CalculateChangeProposerTimeout(1), c.ChangeProposerTimeout+c.ChangeProposerDelta)
	assert.Equal(t, c.CalculateChangeProposerTimeout(4), c.ChangeProposerTimeout+(4*c.ChangeProposerDelta))
}
