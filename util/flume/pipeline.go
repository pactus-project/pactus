package flume

import (
	"context"
	"sync"

	"github.com/pactus-project/pactus/util/logger"
)

var _ Pipeline[int] = &pipeline[int]{}

type PipelineBase interface {
	Name() string
	Close()
	IsClosed() bool
}

// Pipeline defines the core interface for all pipelines
type Pipeline[T any] interface {
	PipelineBase

	Send(T)
	RegisterReceiver(func(T))
}

// pipeline implements Pipeline
type pipeline[T any] struct {
	sync.RWMutex

	ctx      context.Context
	cancel   context.CancelFunc
	name     string
	closed   bool
	ch       chan T
	receiver func(T)
}

// NewPipeline creates a new pipeline instance
func NewPipeline[T any](parentCtx context.Context,
	name string, bufferSize int,
) *pipeline[T] {
	ctx, cancel := context.WithCancel(parentCtx)

	return &pipeline[T]{
		ctx:    ctx,
		cancel: cancel,
		name:   name,
		closed: false,
		ch:     make(chan T, bufferSize),
	}
}

func (p *pipeline[T]) Name() string {
	return p.name
}

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
	}
}

func (p *pipeline[T]) RegisterReceiver(receiver func(T)) {
	p.receiver = receiver

	go p.receiveLoop()
}

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
			p.receiver(data)
		}
	}
}

func (p *pipeline[T]) Close() {
	p.Lock()
	defer p.Unlock()

	if !p.closed {
		p.cancel()

		close(p.ch)
		p.closed = true
	}
}

func (p *pipeline[T]) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()

	return p.closed
}
