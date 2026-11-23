//go:build gtk

package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeedOnce(t *testing.T) {
	seed := &seedOnce{}

	first := seed.get()
	assert.Len(t, strings.Fields(first), 12)
	assert.Equal(t, first, seed.get(), "going back to the seed page shows the same seed")
	assert.NotEqual(t, first, (&seedOnce{}).get(), "another assistant has another seed")
}
