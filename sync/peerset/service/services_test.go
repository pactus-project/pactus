package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesString(t *testing.T) {
	assert.Equal(t, New(None).String(), "")
	assert.Equal(t, New(Network).String(), "NETWORK")
	assert.Equal(t, New(Gossip).String(), "GOSSIP")
	assert.Equal(t, New(Network, Gossip).String(), "NETWORK | GOSSIP")
	assert.Equal(t, New(5).String(), "NETWORK | 4")
}

func TestAppend(t *testing.T) {
	s := New(Network)
	assert.True(t, s.IsNetwork())
	assert.False(t, s.IsGossip())

	s.Append(Gossip)
	assert.True(t, s.IsNetwork())
	assert.True(t, s.IsGossip())
}

func TestIsNetwork(t *testing.T) {
	assert.False(t, New(None).IsNetwork())
	assert.True(t, New(Network).IsNetwork())
}

func TestIsGossip(t *testing.T) {
	assert.False(t, New(None).IsGossip())
	assert.False(t, New(Network).IsGossip())
}
