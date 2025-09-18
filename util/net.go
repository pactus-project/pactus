package util

import (
	"context"
	"net"
	"time"
)

// NetworkListen listens on a network address with a context.
func NetworkListen(ctx context.Context, network, address string) (net.Listener, error) {
	var lc net.ListenConfig

	return lc.Listen(ctx, network, address)
}

// NetworkDialTimeout dials a network address with a timeout and a context.
func NetworkDialTimeout(ctx context.Context, network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}

	return d.DialContext(ctx, network, address)
}
