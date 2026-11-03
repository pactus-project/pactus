package network

import (
	"context"
	"sync"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/gopkg/logger"
)

// maxPendingMessages is the maximum number of messages that are waiting to be
// sent to a single peer. Messages over the limit are dropped.
const maxPendingMessages = 256

// orderedSender sends the messages of each destination one after another, in
// the order they were queued, without blocking the caller.
// Different destinations are independent: a slow or unreachable peer only delays
// the messages that are queued for it.
type orderedSender struct {
	ctx    context.Context
	send   func(msg []byte, pid lp2ppeer.ID) error
	logger *logger.SubLogger

	lk sync.Mutex
	// queues holds the pending messages of every peer that has a running worker.
	// The message that is being sent stays at the head of its queue until it is sent.
	queues map[lp2ppeer.ID][][]byte
}

func newOrderedSender(ctx context.Context, send func(msg []byte, pid lp2ppeer.ID) error,
	log *logger.SubLogger,
) *orderedSender {
	return &orderedSender{
		ctx:    ctx,
		send:   send,
		logger: log,
		queues: make(map[lp2ppeer.ID][][]byte),
	}
}

// enqueue queues the message for the peer and returns immediately.
func (s *orderedSender) enqueue(msg []byte, pid lp2ppeer.ID) {
	if s.ctx.Err() != nil {
		return
	}

	s.lk.Lock()
	defer s.lk.Unlock()

	queue, running := s.queues[pid]
	if len(queue) >= maxPendingMessages {
		s.logger.Warn("too many pending messages, dropping a message", "pid", pid)

		return
	}
	s.queues[pid] = append(queue, msg)

	if !running {
		go s.run(pid)
	}
}

// run sends the queued messages of the peer and returns once its queue is empty.
func (s *orderedSender) run(pid lp2ppeer.ID) {
	for {
		s.lk.Lock()
		queue := s.queues[pid]
		if len(queue) == 0 || s.ctx.Err() != nil {
			delete(s.queues, pid)
			s.lk.Unlock()

			return
		}
		msg := queue[0]
		s.lk.Unlock()

		if err := s.send(msg, pid); err != nil {
			s.logger.Warn("error on sending msg", "pid", pid, "error", err)
		}

		s.lk.Lock()
		queue = s.queues[pid]
		queue[0] = nil
		s.queues[pid] = queue[1:]
		s.lk.Unlock()
	}
}
