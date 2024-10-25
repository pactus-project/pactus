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
	services := New(FullNode)
	assert.True(t, services.IsFullNode())
	assert.False(t, services.IsPrunedNode())

	services.Append(PrunedNode)
	assert.True(t, services.IsFullNode())
	assert.True(t, services.IsPrunedNode())
}
