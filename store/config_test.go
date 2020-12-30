package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.Path = "/abc"
	assert.Error(t, c.SanityCheck())
	c.Path = "."
	assert.NoError(t, c.SanityCheck())
}
