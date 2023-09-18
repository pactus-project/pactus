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
	assert.NoError(t, c1.BasicCheck())

	c2.ChangeProposerDelta = 0 * time.Second
	assert.Error(t, c2.BasicCheck())

	c3.ChangeProposerTimeout = 0 * time.Second
	assert.Error(t, c3.BasicCheck())

	c4.ChangeProposerTimeout = -1 * time.Second
	assert.Error(t, c4.BasicCheck())
}

func TestCalculateChangeProposerTimeout(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, c.CalculateChangeProposerTimeout(0), c.ChangeProposerTimeout)
	assert.Equal(t, c.CalculateChangeProposerTimeout(1), c.ChangeProposerTimeout+c.ChangeProposerDelta)
	assert.Equal(t, c.CalculateChangeProposerTimeout(4), c.ChangeProposerTimeout+(4*c.ChangeProposerDelta))
}
