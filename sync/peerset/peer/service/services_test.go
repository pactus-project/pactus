package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesString(t *testing.T) {
	assert.Equal(t, "", New(None).String())
	assert.Equal(t, "FULL", New(FullNode).String())
	assert.Equal(t, "PRUNED", New(PrunedNode).String())
	assert.Equal(t, "FULL | PRUNED", New(FullNode, PrunedNode).String())
	assert.Equal(t, "FULL | 4", New(5).String())
	assert.Equal(t, "PRUNED | 4", New(6).String())
}

func TestAppend(t *testing.T) {
	s := New(FullNode)
	assert.True(t, s.IsFullNode())
	assert.False(t, s.IsPrunedNode())

	s.Append(PrunedNode)
	assert.True(t, s.IsFullNode())
	assert.True(t, s.IsPrunedNode())
}
