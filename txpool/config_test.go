package txpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.MaxSize = 0
	assert.Error(t, c.SanityCheck())
}
