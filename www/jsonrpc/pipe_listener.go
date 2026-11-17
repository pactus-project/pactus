package jsonrpc

import (
	"context"
	"net"
	"sync"
)

// pipeListener is a listener whose connections are made in memory, by dial.
// The JSON-RPC library serves on it, and the front server of the node passes it
// the requests that it accepts. Nothing outside the process can connect to it.
type pipeListener struct {
	conns     chan net.Conn
	done      chan struct{}
	closeOnce sync.Once
}

func newPipeListener() *pipeListener {
	return &pipeListener{
		conns: make(chan net.Conn),
		done:  make(chan struct{}),
	}
}

func (l *pipeListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.conns:
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *pipeListener) Close() error {
	l.closeOnce.Do(func() { close(l.done) })

	return nil
}

func (*pipeListener) Addr() net.Addr {
	return pipeAddr{}
}

// dial makes a connection to the listener.
func (l *pipeListener) dial(ctx context.Context, _, _ string) (net.Conn, error) {
	client, server := net.Pipe()
	select {
	case l.conns <- server:
		return client, nil
	case <-l.done:
		_ = client.Close()
		_ = server.Close()

		return nil, net.ErrClosed
	case <-ctx.Done():
		_ = client.Close()
		_ = server.Close()

		return nil, ctx.Err()
	}
}

type pipeAddr struct{}

func (pipeAddr) Network() string { return "pipe" }
func (pipeAddr) String() string  { return "jsonrpc-pipe" }
