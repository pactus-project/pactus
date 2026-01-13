//go:build gtk

package gtkutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsContextDone(t *testing.T) {
	// Test with active context
	ctx := context.Background()
	assert.False(t, IsContextDone(ctx))

	// Test with cancelled context
	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	assert.True(t, IsContextDone(cancelledCtx))
}
