package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.Path = "/tmp/zarb"
	assert.NoError(t, c.SanityCheck())
	assert.Equal(t, c.StorePath(), "/tmp/zarb/data/store.db")

}
