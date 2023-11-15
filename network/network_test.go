package network

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func alwaysPropagate(_ *GossipMessage) bool {
	return true
}

// Original code from:
// https://github.com/libp2p/go-libp2p/blob/master/p2p/host/autorelay/autorelay_test.go
func makeTestRelay(t *testing.T) lp2phost.Host {
	t.Helper()

	h, err := lp2p.New(
		lp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		lp2p.DisableRelay(),
		lp2p.EnableRelayService(),
		lp2p.ForceReachabilityPublic(),
		lp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr {
			return addrs
		}),
	)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		for _, p := range h.Mux().Protocols() {
			if p == "/libp2p/circuit/relay/0.2.0/hop" {
				return true
			}
		}
		return false
	}, time.Second, 10*time.Millisecond)
	return h
}

func makeTestNetwork(t *testing.T, conf *Config, opts []lp2p.Option) *network {
	t.Helper()

	log := logger.NewSubLogger("_network", nil)
	net, err := newNetwork(conf, log, opts)
	require.NoError(t, err)

	log.SetObj(testsuite.NewOverrideStringer(
		fmt.Sprintf("%s - %s: ", net.SelfID().ShortString(), t.Name()), net))

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic(alwaysPropagate))

	return net
}

func testConfig() *Config {
	return &Config{
		ListenAddrStrings:    []string{},
		NetworkKey:           util.TempFilePath(),
		BootstrapAddrStrings: []string{},
		MinConns:             4,
		MaxConns:             8,
		EnableNATService:     false,
		EnableUPnP:           false,
		EnableRelay:          false,
		EnableMdns:           false,
		ForcePrivateNetwork:  true,
		NetworkName:          "test",
		DefaultPort:          12345,
	}
}

func shouldReceiveEvent(t *testing.T, net *network, eventType EventType) Event {
	t.Helper()

	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldReceiveEvent Timeout, test: %v id:%s", t.Name(), net.SelfID().String()))
			return nil
		case e := <-net.EventChannel():
			if e.Type() == eventType {
				return e
			}
		}
	}
}

func shouldNotReceiveEvent(t *testing.T, net *network) {
	t.Helper()

	timeout := time.NewTimer(100 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case <-net.EventChannel():
			require.NoError(t, fmt.Errorf("shouldNotReceiveEvent, test: %v id:%s", t.Name(), net.SelfID().String()))
			return
		}
	}
}

func readData(t *testing.T, r io.ReadCloser, length int) []byte {
	t.Helper()

	buf := make([]byte, length)
	_, err := r.Read(buf)
	if !errors.Is(err, io.EOF) {
		assert.NoError(t, err)
		assert.NoError(t, r.Close())
	}

	return buf
}

func TestStoppingNetwork(t *testing.T) {
	net, err := NewNetwork(testConfig())
	assert.NoError(t, err)

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic(alwaysPropagate))

	// Should stop peacefully
	net.Stop()
}

// In this test, we are setting up a simulated network environment that consists of six nodes:
//   - R is a Relay node
//   - B is a Bootstrap node
//   - P is a Public node
//   - M, N, and X are Private Nodes behind a Network Address Translation (NAT)
//   - Both M and N are connected to the relay node R
//   - X is not connected to the relay node and does not join the consensus topic
//
// The test will evaluate the following scenarios:
//   - Connection establishment to the bootstrap node
//   - General and consensus topics and gossip message
//   - Direct and relayed stream communication between nodes
func TestNetwork(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// Relay
	nodeR := makeTestRelay(t)

	relayAddrs := []string{}
	for _, addr := range nodeR.Addrs() {
		addr2 := fmt.Sprintf("%s/p2p/%s", addr, nodeR.ID().String())
		fmt.Printf("relay address: %s\n", addr2)
		relayAddrs = append(relayAddrs, addr2)
	}

	bootstrapPort := ts.RandInt32(9999) + 10000
	publicPort := ts.RandInt32(9999) + 10000

	// Bootstrap node
	confB := testConfig()
	confB.ListenAddrStrings = []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", bootstrapPort),
	}
	fmt.Println("Starting Bootstrap node")
	networkB := makeTestNetwork(t, confB, []lp2p.Option{
		lp2p.ForceReachabilityPublic(),
	})
	bootstrapAddresses := []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%v", bootstrapPort, networkB.SelfID().String()),
	}

	// Public node
	confP := testConfig()
	confP.BootstrapAddrStrings = bootstrapAddresses
	confP.ListenAddrStrings = []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", publicPort),
	}
	fmt.Println("Starting Public node")
	networkP := makeTestNetwork(t, confP, []lp2p.Option{
		lp2p.ForceReachabilityPublic(),
	})
	publicAddrInfo, _ := lp2ppeer.AddrInfoFromString(
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", publicPort, networkP.SelfID()))

	// Private node M
	confM := testConfig()
	confM.EnableRelay = true
	confM.RelayAddrStrings = relayAddrs
	confM.BootstrapAddrStrings = bootstrapAddresses
	confM.ListenAddrStrings = []string{
		"/ip4/127.0.0.1/tcp/9987",
	}
	fmt.Println("Starting Private node M")
	networkM := makeTestNetwork(t, confM, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})

	// Private node N
	confN := testConfig()
	confN.EnableRelay = true
	confN.RelayAddrStrings = relayAddrs
	confN.BootstrapAddrStrings = bootstrapAddresses
	confN.ListenAddrStrings = []string{
		"/ip4/127.0.0.1/tcp/5678",
	}
	fmt.Println("Starting Private node N")
	networkN := makeTestNetwork(t, confN, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})

	// Private node X, doesn't join consensus topic and without relay address
	confX := testConfig()
	confX.EnableRelay = false
	confX.BootstrapAddrStrings = bootstrapAddresses
	confX.ListenAddrStrings = []string{
		"/ip4/127.0.0.1/tcp/0",
	}
	fmt.Println("Starting Private node X")
	networkX := makeTestNetwork(t, confX, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})

	assert.NoError(t, networkB.JoinConsensusTopic(alwaysPropagate))
	assert.NoError(t, networkP.JoinConsensusTopic(alwaysPropagate))
	assert.NoError(t, networkM.JoinConsensusTopic(alwaysPropagate))
	assert.NoError(t, networkN.JoinConsensusTopic(alwaysPropagate))
	// Network X doesn't join the consensus topic

	time.Sleep(2 * time.Second)

	t.Run("all nodes have at least one connection to the bootstrap node B", func(t *testing.T) {
		assert.GreaterOrEqual(t, networkP.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkB.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkM.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkN.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkX.NumConnectedPeers(), 1)
	})

	t.Run("Gossip: all nodes receive general gossip messages", func(t *testing.T) {
		msg := []byte("test-general-topic")

		require.NoError(t, networkP.Broadcast(msg, TopicIDGeneral))

		eB := shouldReceiveEvent(t, networkB, EventTypeGossip).(*GossipMessage)
		eM := shouldReceiveEvent(t, networkM, EventTypeGossip).(*GossipMessage)
		eN := shouldReceiveEvent(t, networkN, EventTypeGossip).(*GossipMessage)
		eX := shouldReceiveEvent(t, networkX, EventTypeGossip).(*GossipMessage)

		assert.Equal(t, eB.Data, msg)
		assert.Equal(t, eM.Data, msg)
		assert.Equal(t, eN.Data, msg)
		assert.Equal(t, eX.Data, msg)
	})

	t.Run("only nodes subscribed to the consensus topic receive consensus gossip messages", func(t *testing.T) {
		msg := []byte("test-consensus-topic")

		require.NoError(t, networkP.Broadcast(msg, TopicIDConsensus))

		eB := shouldReceiveEvent(t, networkB, EventTypeGossip).(*GossipMessage)
		eM := shouldReceiveEvent(t, networkM, EventTypeGossip).(*GossipMessage)
		eN := shouldReceiveEvent(t, networkN, EventTypeGossip).(*GossipMessage)
		shouldNotReceiveEvent(t, networkX) // Not joined the consensus topic

		assert.Equal(t, eB.Data, msg)
		assert.Equal(t, eM.Data, msg)
		assert.Equal(t, eN.Data, msg)
	})

	t.Run("node P (public) is directly accessible by nodes M and N (private behind NAT)", func(t *testing.T) {
		require.NoError(t, networkM.host.Connect(networkM.ctx, *publicAddrInfo))

		msgM := []byte("test-stream-from-m")
		require.NoError(t, networkM.SendTo(msgM, networkP.SelfID()))
		eP := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eP.From, networkM.SelfID())
		assert.Equal(t, readData(t, eP.Reader, len(msgM)), msgM)
	})

	t.Run("node P (public) is directly accessible by node X (private behind NAT, without relay)", func(t *testing.T) {
		require.NoError(t, networkX.host.Connect(networkX.ctx, *publicAddrInfo))

		msgX := []byte("test-stream-from-x")
		require.NoError(t, networkX.SendTo(msgX, networkP.SelfID()))
		eP := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eP.From, networkX.SelfID())
		assert.Equal(t, readData(t, eP.Reader, len(msgX)), msgX)
	})

	t.Run("node P (public) is directly accessible by node B (bootstrap)", func(t *testing.T) {
		msgB := []byte("test-stream-from-b")

		require.NoError(t, networkB.SendTo(msgB, networkP.SelfID()))
		eB := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eB.From, networkB.SelfID())
		assert.Equal(t, readData(t, eB.Reader, len(msgB)), msgB)
	})

	circuitAddrInfoN, _ := lp2ppeer.AddrInfoFromString(
		fmt.Sprintf("%s/p2p-circuit/p2p/%s", relayAddrs[0], networkN.SelfID()))

	t.Run("node X (private, not connected via relay) is not accessible by node M", func(t *testing.T) {
		require.Error(t, networkX.host.Connect(networkX.ctx, *circuitAddrInfoN))

		msgM := []byte("test-stream-from-m")
		require.Error(t, networkM.SendTo(msgM, networkX.SelfID()))
	})

	t.Run("nodes M and N (private, connected via relay) can communicate using the relay node R", func(t *testing.T) {
		require.NoError(t, networkM.host.Connect(networkM.ctx, *circuitAddrInfoN))

		// TODO: How to test this?
		// msgM := []byte("test-stream-from-m")
		// require.NoError(t, networkM.SendTo(msgM, networkN.SelfID()))
		// eM := shouldReceiveEvent(t, networkN, EventTypeStream).(*StreamMessage)
		// assert.Equal(t, eM.Source, networkM.SelfID())
		// assert.Equal(t, readData(t, eM.Reader, len(msgM)), msgM)
	})

	t.Run("closing connection", func(t *testing.T) {
		msgB := []byte("test-stream-from-b")

		networkP.Stop()
		networkB.CloseConnection(networkP.SelfID())
		e := shouldReceiveEvent(t, networkB, EventTypeDisconnect).(*DisconnectEvent)
		assert.Equal(t, e.PeerID, networkP.SelfID())
		require.Error(t, networkB.SendTo(msgB, networkP.SelfID()))
	})
}

func TestConnections(t *testing.T) {
	t.Parallel() // run the tests in parallel

	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		bootstrapAddr string
		peerAddr      string
	}{
		{"/ip4/127.0.0.1/tcp/%d", "/ip4/127.0.0.1/tcp/0"},
		{"/ip4/127.0.0.1/udp/%d/quic-v1", "/ip4/127.0.0.1/udp/0/quic-v1"},
		{"/ip6/::1/tcp/%d", "/ip6/::1/tcp/0"},
		{"/ip6/::1/udp/%d/quic-v1", "/ip6/::1/udp/0/quic-v1"},
	}

	for i, test := range tests {
		// Bootstrap node
		confB := testConfig()
		bootstrapPort := ts.RandInt32(9999) + 10000
		bootstrapAddr := fmt.Sprintf(test.bootstrapAddr, bootstrapPort)
		confB.ListenAddrStrings = []string{bootstrapAddr}
		fmt.Println("Starting Bootstrap node")
		networkB := makeTestNetwork(t, confB, []lp2p.Option{
			lp2p.ForceReachabilityPublic(),
		})

		// Public node
		confP := testConfig()
		confP.BootstrapAddrStrings = []string{
			fmt.Sprintf("%s/p2p/%v", bootstrapAddr, networkB.SelfID().String()),
		}
		confP.ListenAddrStrings = []string{test.peerAddr}
		fmt.Println("Starting Public node")
		networkP := makeTestNetwork(t, confP, []lp2p.Option{
			lp2p.ForceReachabilityPublic(),
		})

		t.Run(fmt.Sprintf("Running test %d: %s <-> %s ... ",
			i, test.bootstrapAddr, test.peerAddr), func(t *testing.T) {
			t.Parallel() // run the tests in parallel

			testConnection(t, networkP, networkB)
		})
	}
}

func testConnection(t *testing.T, networkP, networkB *network) {
	t.Helper()

	// Ensure that peers are connected to each other
	for i := 0; i < 20; i++ {
		if networkP.NumConnectedPeers() >= 1 &&
			networkB.NumConnectedPeers() >= 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	assert.Equal(t, networkB.NumConnectedPeers(), 1)
	assert.Equal(t, networkP.NumConnectedPeers(), 1)

	msg := []byte("test-msg")

	require.NoError(t, networkP.SendTo(msg, networkB.SelfID()))
	e := shouldReceiveEvent(t, networkB, EventTypeStream).(*StreamMessage)
	assert.Equal(t, e.From, networkP.SelfID())
	assert.Equal(t, readData(t, e.Reader, len(msg)), msg)

	networkB.Stop()
	networkP.Stop()
}
