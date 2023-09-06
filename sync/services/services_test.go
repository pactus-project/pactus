package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesString(t *testing.T) {
	assert.Equal(t, New(None).String(), "")
	assert.Equal(t, New(Network).String(), "NETWORK")
	assert.Equal(t, New(2).String(), "2")
	assert.Equal(t, New(3).String(), "NETWORK | 2")
}

func TestIsNetwork(t *testing.T) {
	assert.False(t, New(None).IsNetwork())
	assert.True(t, New(Network).IsNetwork())
}
