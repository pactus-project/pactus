package main

import (
	"sync/atomic"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeNode struct {
	stopped atomic.Bool
}

func (n *fakeNode) Stop() { n.stopped.Store(true) }

// otherSignal is a signal that is not a system signal.
type otherSignal struct{}

func (otherSignal) String() string { return "other" }
func (otherSignal) Signal()        {}

func TestExitCode(t *testing.T) {
	assert.Equal(t, 130, exitCode(syscall.SIGINT))
	assert.Equal(t, 143, exitCode(syscall.SIGTERM))
	assert.Equal(t, 1, exitCode(otherSignal{}))
}
