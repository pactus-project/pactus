package consensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.DeltaDuration = -1 * time.Second
	assert.Error(t, c.SanityCheck())
	c.TimeoutPrecommit = -1 * time.Second
	assert.Error(t, c.SanityCheck())
	c.TimeoutPrepare = -1 * time.Second
	assert.Error(t, c.SanityCheck())
	c.TimeoutPropose = -1 * time.Second
	assert.Error(t, c.SanityCheck())
}
