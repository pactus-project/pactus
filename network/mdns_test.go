package network

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMDNS(t *testing.T) {
	if runtime.GOOS == "darwin" {
		// Disable this test on darwin (macOS)
		// Read more here: https://github.com/pactus-project/pactus/issues/1860
		return
	}

	conf1 := testConfig()
	conf1.ListenAddrStrings = []string{
		"/ip6/::1/tcp/0", "/ip6/::1/udp/0/quic-v1",
		"/ip4/127.0.0.1/tcp/0", "/ip4/127.0.0.1/udp/0/quic-v1",
	}
	conf1.EnableMdns = true
	net1 := makeTestNetwork(t, conf1, nil)

	conf2 := testConfig()
	conf2.ListenAddrStrings = []string{
		"/ip6/::1/tcp/0", "/ip6/::1/udp/0/quic-v1",
		"/ip4/127.0.0.1/tcp/0", "/ip4/127.0.0.1/udp/0/quic-v1",
	}
	conf2.EnableMdns = true
	net2 := makeTestNetwork(t, conf2, nil)

	require.NoError(t, net1.Start())
	time.Sleep(250 * time.Millisecond)

	require.NoError(t, net2.Start())
	time.Sleep(250 * time.Millisecond)

	msg := []byte("test-mdns")
	net1.SendTo(msg, net2.SelfID())

	se := shouldReceiveEvent(t, net2, EventTypeStream).(*StreamMessage)
	assert.Equal(t, net1.SelfID(), se.From)
	assert.Equal(t, msg, readData(t, se.Reader, len(msg)))

	net1.Stop()
	net2.Stop()
}

// fakeMdnsHost is a host whose Connect blocks until it is released.
type fakeMdnsHost struct {
	lp2phost.Host

	selfID     lp2ppeer.ID
	release    chan struct{}
	started    atomic.Int32
	running    atomic.Int32
	maxRunning atomic.Int32
}

func (h *fakeMdnsHost) ID() lp2ppeer.ID {
	return h.selfID
}

func (h *fakeMdnsHost) Connect(_ context.Context, _ lp2ppeer.AddrInfo) error {
	h.started.Add(1)
	running := h.running.Add(1)
	for {
		maxRunning := h.maxRunning.Load()
		if running <= maxRunning || h.maxRunning.CompareAndSwap(maxRunning, running) {
			break
		}
	}
	<-h.release
	h.running.Add(-1)

	return nil
}

func TestMDNSBoundsConcurrentConnects(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	host := &fakeMdnsHost{selfID: ts.RandPeerID(), release: make(chan struct{})}
	log := logger.NewSubLogger(t.Context(), "_network", nil)
	mdns := newMdnsService(t.Context(), host, log)

	const numPeers = 100
	var returned atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < numPeers; i++ {
		addrInfo := lp2ppeer.AddrInfo{ID: ts.RandPeerID()}
		wg.Go(func() {
			mdns.HandlePeerFound(addrInfo)
			returned.Add(1)
		})
	}

	// Announcements over the limit are skipped without waiting for the dials.
	require.Eventually(t, func() bool {
		return returned.Load() == numPeers-maxMdnsConnects
	}, 10*time.Second, 10*time.Millisecond)

	close(host.release)
	wg.Wait()

	assert.Equal(t, int32(maxMdnsConnects), host.started.Load())
	assert.LessOrEqual(t, host.maxRunning.Load(), int32(maxMdnsConnects))
}

func TestMDNSSkipsPeerBeingDialed(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	host := &fakeMdnsHost{selfID: ts.RandPeerID(), release: make(chan struct{})}
	log := logger.NewSubLogger(t.Context(), "_network", nil)
	mdns := newMdnsService(t.Context(), host, log)

	addrInfo := lp2ppeer.AddrInfo{ID: ts.RandPeerID()}

	var wg sync.WaitGroup
	wg.Go(func() {
		mdns.HandlePeerFound(addrInfo)
	})
	require.Eventually(t, func() bool {
		return host.started.Load() == 1
	}, 10*time.Second, 10*time.Millisecond)

	// The same peer announced again while the first dial is running.
	mdns.HandlePeerFound(addrInfo)
	assert.Equal(t, int32(1), host.started.Load())

	close(host.release)
	wg.Wait()

	// Once the dial has finished the peer can be dialed again.
	mdns.HandlePeerFound(addrInfo)
	assert.Equal(t, int32(2), host.started.Load())
}

func TestMDNSIgnoresSelf(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	host := &fakeMdnsHost{selfID: ts.RandPeerID(), release: make(chan struct{})}
	log := logger.NewSubLogger(t.Context(), "_network", nil)
	mdns := newMdnsService(t.Context(), host, log)

	mdns.HandlePeerFound(lp2ppeer.AddrInfo{ID: host.selfID})
	assert.Equal(t, int32(0), host.started.Load())
}

func TestMDNSWarning(t *testing.T) {
	tests := []struct {
		name                string
		enableMdns          bool
		forcePrivateNetwork bool
		hasWarning          bool
	}{
		{"mdns disabled", false, false, false},
		{"mdns disabled on a private network", false, true, false},
		{"mdns on a private network", true, true, false},
		{"mdns without a private network", true, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := testConfig()
			conf.EnableMdns = tt.enableMdns
			conf.ForcePrivateNetwork = tt.forcePrivateNetwork

			warning := mdnsWarning(conf)
			if tt.hasWarning {
				assert.Contains(t, warning, "force_private_network")
			} else {
				assert.Empty(t, warning)
			}

			// Whatever the warning says, the configuration stays valid.
			assert.NoError(t, conf.BasicCheck())
		})
	}
}
