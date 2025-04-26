// Package flume provides managed Go channels inspired by real-world flumes -
// human-made water channels that control flow.
//
// Flume completely encapsulates Go channels internally,
// exposing only clean Pipeline interfaces to consumers.
// This design keeps all channel management internal to the system.
//
// Key Features:
// - Designed for long-lived channels shared between separate modules
// - Automatic lifecycle management via context
// - Receiving messages through callback functions
// - Graceful shutdown handling
package flume

import "context"

// Flume manages multiple pipelines with context support
type Flume struct {
	ctx    context.Context
	cancel context.CancelFunc
	pipes  []PipelineBase
}

// New creates a Flume instance
func New(ctx context.Context) *Flume {
	ctx, cancel := context.WithCancel(ctx)
	return &Flume{
		ctx:    ctx,
		cancel: cancel,
		pipes:  make([]PipelineBase, 0),
	}
}

func CreatePipeline[T any](f *Flume, name string, bufferSize int) *pipeline[T] {
	pipe := NewPipeline[T](f.ctx, name, bufferSize)
	f.pipes = append(f.pipes, pipe)

	return pipe
}

// CloseAll gracefully shuts down all pipelines
func (f *Flume) CloseAll() {
	f.cancel()

	for _, pipe := range f.pipes {
		pipe.Close()
	}
}
