// Package pipeline provides managed Go channels.
//
// Pipeline completely encapsulates Go channels internally,
// exposing only clean Pipeline interfaces to consumers.
// This design keeps all channel management internal to the system.
//
// Key Features:
// - Designed for long-lived channels shared between separate modules
// - Automatic lifecycle management via context
// - Receiving messages through callback functions
package pipeline

import (
	"context"
	"sync"

	"github.com/pactus-project/pactus/util/logger"
)

var _ Pipeline[int] = &pipeline[int]{}

// Pipeline defines the core interface for all pipelines.
// It provides type-safe channel operations with lifecycle management.
type Pipeline[T any] interface {
	// Returns the name of the pipeline
	Name() string
	// Closes the pipeline gracefully
	Close()
	// Checks if pipeline is closed
	IsClosed() bool
	// Sends data through the pipeline
	Send(T)
	// Registers a receiver callback
	RegisterReceiver(func(T))
}

// pipeline is the internal implementation of the Pipeline interface.
// It manages the channel lifecycle and handles message passing.
type pipeline[T any] struct {
	sync.RWMutex

	ctx      context.Context
	cancel   context.CancelFunc
	name     string
	closed   bool
	ch       chan T
	receiver func(T)
}

// New creates and initializes a new pipeline instance.
//
// Parameters:
//   - parentCtx: The parent context for lifecycle management
//   - name: Identifier for the pipeline (used in logging)
//   - bufferSize: Size of the channel buffer (0 for unbuffered)
//
// Returns:
//   - A new pipeline instance ready for use
func New[T any](parentCtx context.Context, name string, bufferSize int) *pipeline[T] {
	ctx, cancel := context.WithCancel(parentCtx)

	return &pipeline[T]{
		ctx:    ctx,
		cancel: cancel,
		name:   name,
		closed: false,
		ch:     make(chan T, bufferSize),
	}
}

// Name returns the identifier name of the pipeline.
func (p *pipeline[T]) Name() string {
	return p.name
}

// Send writes data to the pipeline channel in a thread-safe manner.
// It handles various context cancellation scenarios and logs appropriate messages.
//
// Parameters:
//   - data: The data to send through the pipeline
func (p *pipeline[T]) Send(data T) {
	p.RLock()
	closed := p.closed
	p.RUnlock()

	if closed {
		logger.Warn("send on closed channel", "name", p.name)
		return
	}

	select {
	case <-p.ctx.Done():
		err := p.ctx.Err()
		switch err {
		case context.Canceled:
			logger.Debug("pipeline draining, normal shutdown", "name", p.name)
		case context.DeadlineExceeded:
			logger.Warn("pipeline timeout, possible stall", "name", p.name)
		default:
			logger.Error("pipeline aborted, unexpected error", "name", p.name, "error", err)
		}
	case p.ch <- data:
		// Data successfully sent
	}
}

// RegisterReceiver sets the callback function for processing received data
// and starts the receive loop in a separate goroutine.
//
// Parameters:
//   - receiver: The callback function that will process received data
func (p *pipeline[T]) RegisterReceiver(receiver func(T)) {
	p.receiver = receiver

	go p.receiveLoop()
}

// receiveLoop continuously listens for incoming data and processes it
// using the registered receiver callback until the pipeline is closed.
func (p *pipeline[T]) receiveLoop() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case data, ok := <-p.ch:
			if !ok {
				logger.Warn("channel is closed", "name", p.name)
				return
			}
			// Process received data with the callback
			p.receiver(data)
		}
	}
}

// Close shuts down the pipeline gracefully.
// It cancels the context, closes the channel, and marks the pipeline as closed.
// This method is idempotent - subsequent calls have no effect.
func (p *pipeline[T]) Close() {
	p.Lock()
	defer p.Unlock()

	if !p.closed {
		p.cancel()

		// Close the channel and mark pipeline as closed
		close(p.ch)
		p.closed = true
	}
}

// IsClosed checks if the pipeline has been closed.
//
// Returns:
//   - true if the pipeline is closed, false otherwise
func (p *pipeline[T]) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()

	return p.closed
}
