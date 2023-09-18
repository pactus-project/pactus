package network

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Original code from:
// https://github.com/libp2p/go-libp2p/blob/master/p2p/host/autorelay/autorelay_test.go
func makeTestRelay(t *testing.T) host.Host {
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

	net, err := newNetwork("test", conf, opts)
	require.NoError(t, err)

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic())

	return net
}

func testConfig() *Config {
	return &Config{
		Listens:     []string{},
		NetworkKey:  util.TempFilePath(),
		EnableNAT:   false,
		EnableRelay: false,
		EnableMdns:  false,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{},
			MinThreshold: 4,
			MaxThreshold: 8,
			Period:       2 * time.Second,
		},
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

func readData(t *testing.T, r io.ReadCloser, len int) []byte {
	t.Helper()

	buf := make([]byte, len)
	_, err := r.Read(buf)
	if !errors.Is(err, io.EOF) {
		assert.NoError(t, err)
		assert.NoError(t, r.Close())
	}

	return buf
}

func TestStoppingNetwork(t *testing.T) {
	net, err := NewNetwork("test", testConfig())
	assert.NoError(t, err)

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic())

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

	// Bootstrap node
	confB := testConfig()
	bootstrapPort := ts.RandInt32(9999) + 10000
	confB.Listens = []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", bootstrapPort),
		fmt.Sprintf("/ip6/::1/tcp/%v", bootstrapPort),
	}
	fmt.Println("Starting Bootstrap node")
	networkB := makeTestNetwork(t, confB, []lp2p.Option{})
	bootstrapAddresses := []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%v", bootstrapPort, networkB.SelfID().String()),
		fmt.Sprintf("/ip6/::1/tcp/%v/p2p/%v", bootstrapPort, networkB.SelfID().String()),
	}
	assert.NoError(t, networkB.JoinConsensusTopic())

	// Public node
	confP := testConfig()
	confP.EnableNAT = true
	confP.Bootstrap.Addresses = bootstrapAddresses
	confP.Listens = []string{
		"/ip4/127.0.0.1/tcp/0",
		"/ip6/::1/tcp/0",
	}
	fmt.Println("Starting Public node")
	networkP := makeTestNetwork(t, confP, []lp2p.Option{
		lp2p.ForceReachabilityPublic(),
	})
	assert.NoError(t, networkP.JoinConsensusTopic())

	// Private node M
	confM := testConfig()
	confM.EnableRelay = true
	confM.RelayAddrs = relayAddrs
	confM.Bootstrap.Addresses = bootstrapAddresses
	confM.Listens = []string{
		"/ip4/127.0.0.1/tcp/0",
		"/ip6/::1/tcp/0",
	}
	fmt.Println("Starting Private node M")
	networkM := makeTestNetwork(t, confM, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})
	assert.NoError(t, networkM.JoinConsensusTopic())

	// Private node N
	confN := testConfig()
	confN.EnableRelay = true
	confN.RelayAddrs = relayAddrs
	confN.Bootstrap.Addresses = bootstrapAddresses
	confN.Listens = []string{
		"/ip4/127.0.0.1/tcp/0",
		"/ip6/::1/tcp/0",
	}
	fmt.Println("Starting Private node N")
	networkN := makeTestNetwork(t, confN, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})
	assert.NoError(t, networkN.JoinConsensusTopic())

	// Private node X, doesn't join consensus topic
	confX := testConfig()
	confX.EnableRelay = false
	confX.Bootstrap.Addresses = bootstrapAddresses
	confX.Listens = []string{
		"/ip4/127.0.0.1/tcp/0",
		"/ip6/::1/tcp/0",
	}
	fmt.Println("Starting Private node X")
	networkX := makeTestNetwork(t, confX, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})
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

		assert.Equal(t, eB.Source, networkP.SelfID())
		assert.Equal(t, eM.Source, networkP.SelfID())
		assert.Equal(t, eN.Source, networkP.SelfID())
		assert.Equal(t, eX.Source, networkP.SelfID())

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
		shouldNotReceiveEvent(t, networkX)

		assert.Equal(t, eB.Source, networkP.SelfID())
		assert.Equal(t, eM.Source, networkP.SelfID())
		assert.Equal(t, eN.Source, networkP.SelfID())

		assert.Equal(t, eB.Data, msg)
		assert.Equal(t, eM.Data, msg)
		assert.Equal(t, eN.Data, msg)
	})

	t.Run("node P (public) is not directly accessible by nodes M and N (private behind NAT)", func(t *testing.T) {
		msgM := []byte("test-stream-from-m")

		require.Error(t, networkM.SendTo(msgM, networkP.SelfID()))
	})

	t.Run("node P (public) is directly accessible by node B (bootstrap)", func(t *testing.T) {
		msgB := []byte("test-stream-from-b")

		require.NoError(t, networkB.SendTo(msgB, networkP.SelfID()))
		eB := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eB.Source, networkB.SelfID())
		assert.Equal(t, readData(t, eB.Reader, len(msgB)), msgB)
	})

	t.Run("nodes M and N (private, connected via relay) can communicate using the relay node R", func(t *testing.T) {
		msgM := []byte("test-stream-from-m")

		require.NoError(t, networkM.SendTo(msgM, networkN.SelfID()))
		eM := shouldReceiveEvent(t, networkN, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eM.Source, networkM.SelfID())
		assert.Equal(t, readData(t, eM.Reader, len(msgM)), msgM)
	})

	t.Run("node X (private, not connected via relay) is not accessible by node M", func(t *testing.T) {
		msgM := []byte("test-stream-from-m")

		require.Error(t, networkM.SendTo(msgM, networkX.SelfID()))
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

func TestInvalidRelayAddress(t *testing.T) {
	conf := testConfig()
	conf.EnableRelay = true

	conf.RelayAddrs = []string{"127.0.0.1:4001"}
	_, err := NewNetwork("test", conf)
	assert.Error(t, err)

	conf.RelayAddrs = []string{"/ip4/127.0.0.1/tcp/4001"}
	_, err = NewNetwork("test", conf)
	assert.Error(t, err)
}
