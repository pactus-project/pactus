package network

import (
	"context"
	"testing"

	lp2p "github.com/libp2p/go-libp2p"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An in-memory host: it does not listen on any address.
func notifeeFactory(t *testing.T) func() *NotifeeService {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	host, err := lp2p.New(lp2p.NoListenAddrs, lp2p.NoTransports, lp2p.DisableRelay())
	require.NoError(t, err)
	t.Cleanup(func() { _ = host.Close() })

	return func() *NotifeeService {
		log := logger.NewSubLogger(ctx, "_notifee", nil)
		pipe := pipeline.New[Event](ctx)

		mgr := newPeerMgr(ctx, host, testConfig(), log) //nolint:contextcheck // the config only needs a free port.
		notifee, err := newNotifeeService(ctx, host, pipe, mgr, "/test/stream/v1", log)
		require.NoError(t, err)

		return notifee
	}
}

type conn struct {
	lp2pnetwork.Conn // nil: only the methods below are used by the notifee

	pid lp2ppeer.ID
	ma  multiaddr.Multiaddr
}

func (c *conn) RemotePeer() lp2ppeer.ID              { return c.pid }
func (c *conn) RemoteMultiaddr() multiaddr.Multiaddr { return c.ma }
func (*conn) Stat() lp2pnetwork.ConnStats {
	return lp2pnetwork.ConnStats{Stats: lp2pnetwork.Stats{Direction: lp2pnetwork.DirInbound}}
}

// libp2p calls Connected / Disconnected once per *connection*, and two peers can
// hold several connections at the same time (simultaneous dial, TCP + QUIC).
// The first Disconnected must not make the node believe the peer is gone while a
// second connection is still open.
func TestNotifeeTwoConnectionsSamePeer(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	notifee := notifeeFactory(t)()

	pid := ts.RandPeerID()
	ma, err := multiaddr.NewMultiaddr("/ip4/2.2.2.2/tcp/21888")
	require.NoError(t, err)
	conn1 := &conn{pid: pid, ma: ma}
	conn2 := &conn{pid: pid, ma: ma}

	notifee.Connected(nil, conn1)
	notifee.Connected(nil, conn2)
	notifee.Disconnected(nil, conn1) // conn2 is still open

	assert.Equal(t, 1, notifee.peerMgr.NumInbound(),
		"one inbound connection is still open")

	disconnects := 0
	ch := notifee.networkPipe.UnsafeGetChannel()
drain:
	for {
		select {
		case evt := <-ch:
			if evt.Type() == EventTypeDisconnect {
				disconnects++
			}
		default:
			break drain
		}
	}
	assert.Zero(t, disconnects,
		"a DisconnectEvent was sent to the sync module while the peer still has an open connection")
}
