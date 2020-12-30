package txpool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.MaxSize = 0
	assert.Error(t, c.SanityCheck())
	c.WaitingTimeout = -1 * time.Second
	assert.Error(t, c.SanityCheck())
}
