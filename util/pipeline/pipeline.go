// Package pipeline provides a high-level abstraction for managing Go channels with
// built-in lifecycle management, error handling, and receiver callbacks.
//
// The pipeline pattern implemented here offers several advantages over raw channels:
// - Encapsulated channel management with controlled access
// - Context-aware cancellation and graceful shutdown
// - Thread-safe operations with proper synchronization
// - Simplified receiver registration pattern
// - Built-in logging for debugging and monitoring
package pipeline

import (
	"context"
	"sync"

	"github.com/pactus-project/pactus/util/logger"
)

var _ Pipeline[int] = &pipeline[int]{}

// Pipeline defines the contract for a managed channel pipeline.
// It provides type-safe channel operations with lifecycle management.
type Pipeline[T any] interface {
	// Name returns the identifier for this pipeline instance.
	Name() string

	// Close initiates a graceful shutdown of the pipeline.
	Close()

	// IsClosed reports whether the pipeline has been closed.
	IsClosed() bool

	// Send publishes a message to the pipeline (non-blocking).
	Send(T)

	// RegisterReceiver sets the handler function for incoming messages.
	RegisterReceiver(func(T))

	// UnsafeGetChannel provides direct read access to the underlying channel
	// WARNING: This bypasses pipeline management and should be used with caution.
	UnsafeGetChannel() <-chan T
}

// pipeline implements the Pipeline interface with proper synchronization
// and lifecycle management.
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
func New[T any](parentCtx context.Context, name string, bufferSize int) Pipeline[T] {
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
	defer p.RUnlock()

	if p.closed {
		logger.Debug("send on closed channel", "name", p.name)

		return
	}

	select {
	case <-p.ctx.Done():
		err := p.ctx.Err()
		switch err {
		case context.Canceled:
			logger.Debug("pipeline draining", "name", p.name)
		case context.DeadlineExceeded:
			logger.Warn("pipeline timeout", "name", p.name)
		default:
			logger.Error("pipeline error", "name", p.name, "error", err)
		}
	case p.ch <- data:
		// Successful send
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
		case data := <-p.ch:
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

// UnsafeGetChannel provides direct read access to the underlying channel.
// WARNING: Bypasses all pipeline safeguards.
//
// Returns:
//   - The underlying receive channel
func (p *pipeline[T]) UnsafeGetChannel() <-chan T {
	return p.ch
}
