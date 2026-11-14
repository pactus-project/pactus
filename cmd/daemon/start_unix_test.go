//go:build !windows

package main

import (
	"context"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A signal that arrives while the node starts cancels the start, and the node is stopped
// once it is started. A second signal exits at once.
func TestShutdownWhileStarting(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	released := atomic.Bool{}
	exits := make(chan int, 2)
	shut := newShutdown(cancel, func() { released.Store(true) }, func(code int) { exits <- code })

	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGTERM))
	select {
	case <-ctx.Done():
	case <-time.After(5 * time.Second):
		require.FailNow(t, "the start is not cancelled")
	}

	// The node was starting: it is not stopped, and the process doesn't exit yet.
	select {
	case code := <-exits:
		require.FailNow(t, "exit before the node is stopped", "code %d", code)
	case <-time.After(100 * time.Millisecond):
	}

	node := &fakeNode{}
	shut.setNode(node)
	assert.True(t, node.stopped.Load())
	assert.True(t, released.Load())
	select {
	case code := <-exits:
		assert.Equal(t, 143, code)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "the node is not stopped once it is started")
	}

	// A second signal forces the exit, even if stopping the node blocks.
	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGTERM))
	select {
	case code := <-exits:
		assert.Equal(t, 143, code)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "a second signal doesn't exit")
	}
}

// A signal that arrives once the node is started stops it, and the process exits.
func TestShutdownOnceStarted(t *testing.T) {
	_, cancel := context.WithCancel(t.Context())
	released := atomic.Bool{}
	exits := make(chan int, 2)
	shut := newShutdown(cancel, func() { released.Store(true) }, func(code int) { exits <- code })

	node := &fakeNode{}
	shut.setNode(node)
	assert.False(t, node.stopped.Load())

	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGINT))
	select {
	case code := <-exits:
		assert.Equal(t, 130, code)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "the node is not stopped")
	}
	assert.True(t, node.stopped.Load())
	assert.True(t, released.Load())
}
