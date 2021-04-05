package consensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.Timeout = -1 * time.Second
	assert.Error(t, c.SanityCheck())
}
