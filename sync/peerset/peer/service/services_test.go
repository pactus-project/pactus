package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesString(t *testing.T) {
	assert.Equal(t, New(None).String(), "")
	assert.Equal(t, New(FullNode).String(), "FULL")
	assert.Equal(t, New(PrunedNode).String(), "PRUNED")
	assert.Equal(t, New(FullNode, PrunedNode).String(), "FULL | PRUNED")
	assert.Equal(t, New(5).String(), "FULL | 4")
	assert.Equal(t, New(6).String(), "PRUNED | 4")
}

func TestAppend(t *testing.T) {
	s := New(FullNode)
	assert.True(t, s.IsFullNode())
	assert.False(t, s.IsPrunedNode())

	s.Append(PrunedNode)
	assert.True(t, s.IsFullNode())
	assert.True(t, s.IsPrunedNode())
}
