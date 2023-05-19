package network

import (
	"fmt"
	"io"
	"testing"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// https://github.com/libp2p/go-libp2p/blob/5a0411b8eba4276b907c21fefc3adebde19096f1/p2p/host/autorelay/autorelay_test.go
func makeTestRelay(t *testing.T) host.Host {
	t.Helper()
	h, err := lp2p.New(
		lp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
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
	Net, err := newNetwork(conf, opts)
	assert.NoError(t, err)

	assert.NoError(t, Net.Start())
	assert.NoError(t, Net.JoinGeneralTopic())

	return Net.(*network)
}

func testConfig() *Config {
	return &Config{
		Name:        "test-network",
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

func shouldReceiveEvent(t *testing.T, net *network) Event {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldReceiveEvent Timeout, test: %v id:%s", t.Name(), net.SelfID().String()))
			return nil
		case e := <-net.EventChannel():
			net.logger.Debug("received an event", "event", e, "id", net.SelfID().String())
			return e
		}
	}
}

func shouldNotReceiveEvent(t *testing.T, net *network) {
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
	buf := make([]byte, len)
	_, err := r.Read(buf)
	assert.NoError(t, err)

	return buf
}

func TestStoppingNetwork(t *testing.T) {
	net, err := NewNetwork(testConfig())
	assert.NoError(t, err)

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic())

	// Should stop peacefully
	net.Stop()
}

// We create 6 nodes:
//   - R is relay node
//   - B is bootstrap node,
//   - P is a public node
//   - M is a private Node behind NAT that is connected to the relay R
//   - N is a private Node behind NAT that is connected to the relay R
//   - X is a private Node behind NAT
//
// Let's test different scenarios here
func TestNetwork(t *testing.T) {
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
	bootstrapPort := util.RandInt32(9999) + 10000
	confB.Listens = []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", bootstrapPort),
		fmt.Sprintf("/ip6/::/tcp/%v", bootstrapPort),
	}
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
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
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
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
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
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
	networkN := makeTestNetwork(t, confN, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})
	assert.NoError(t, networkN.JoinConsensusTopic())

	// Private node X, doesn't join consensus topic
	confX := testConfig()
	confX.EnableRelay = false
	confX.Bootstrap.Addresses = bootstrapAddresses
	confX.Listens = []string{
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
	networkX := makeTestNetwork(t, confX, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})
	time.Sleep(1 * time.Second)

	t.Run("All nodes should have at least one connection to bootstrap node", func(t *testing.T) {
		assert.GreaterOrEqual(t, networkP.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkB.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkM.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkN.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(t, networkX.NumConnectedPeers(), 1)
	})

	t.Run("Gossip: All nodes should recede general gossip messages", func(t *testing.T) {
		msg := []byte("test-general-topic")

		require.NoError(t, networkP.Broadcast(msg, TopicIDGeneral))

		eB := shouldReceiveEvent(t, networkB).(*GossipMessage)
		eM := shouldReceiveEvent(t, networkM).(*GossipMessage)
		eN := shouldReceiveEvent(t, networkN).(*GossipMessage)
		eX := shouldReceiveEvent(t, networkX).(*GossipMessage)

		assert.Equal(t, eB.Source, networkP.SelfID())
		assert.Equal(t, eM.Source, networkP.SelfID())
		assert.Equal(t, eN.Source, networkP.SelfID())
		assert.Equal(t, eX.Source, networkP.SelfID())

		assert.Equal(t, eB.Data, msg)
		assert.Equal(t, eM.Data, msg)
		assert.Equal(t, eN.Data, msg)
		assert.Equal(t, eX.Data, msg)
	})

	t.Run("Gossip: Only nodes that subscribed to consensus topic should receive consensus gossip messages", func(t *testing.T) {
		msg := []byte("test-consensus-topic")

		require.NoError(t, networkP.Broadcast(msg, TopicIDConsensus))

		eB := shouldReceiveEvent(t, networkB).(*GossipMessage)
		eM := shouldReceiveEvent(t, networkM).(*GossipMessage)
		eN := shouldReceiveEvent(t, networkN).(*GossipMessage)
		shouldNotReceiveEvent(t, networkX)

		assert.Equal(t, eB.Source, networkP.SelfID())
		assert.Equal(t, eM.Source, networkP.SelfID())
		assert.Equal(t, eN.Source, networkP.SelfID())

		assert.Equal(t, eB.Data, msg)
		assert.Equal(t, eM.Data, msg)
		assert.Equal(t, eN.Data, msg)
	})

	t.Run("Stream: Node P should NOT be accessible by node M or N", func(t *testing.T) {
		msgM := []byte("test-stream-from-m")

		require.Error(t, networkM.SendTo(msgM, networkP.SelfID()))
	})

	t.Run("Stream: Node P should be accessible by node B", func(t *testing.T) {
		msgB := []byte("test-stream-from-b")

		require.NoError(t, networkB.SendTo(msgB, networkP.SelfID()))
		eB := shouldReceiveEvent(t, networkP).(*StreamMessage)
		assert.Equal(t, eB.Source, networkB.SelfID())
		assert.Equal(t, readData(t, eB.Reader, len(msgB)), msgB)
	})

	t.Run("Stream: Node M should be accessible by node N using relay node", func(t *testing.T) {
		msgM := []byte("test-stream-from-m")

		require.NoError(t, networkM.SendTo(msgM, networkN.SelfID()))
		eM := shouldReceiveEvent(t, networkN).(*StreamMessage)
		assert.Equal(t, eM.Source, networkM.SelfID())
		assert.Equal(t, readData(t, eM.Reader, len(msgM)), msgM)
	})

	t.Run("Stream: Node X should be NOT accessible by node M", func(t *testing.T) {
		msgM := []byte("test-stream-from-m")

		require.Error(t, networkM.SendTo(msgM, networkX.SelfID()))
	})

	// WHY?
	// t.Run("Stream: Node P closes connection", func(t *testing.T) {
	// 	msgB := []byte("test-stream-from-b")

	// 	networkB.CloseConnection(networkP.SelfID())

	// 	require.Error(t, networkB.SendTo(msgB, networkP.SelfID()))
	// })
}

func TestInvalidTopic(t *testing.T) {
	net, err := NewNetwork(testConfig())
	assert.NoError(t, err)

	msg := []byte("test-invalid-topic")

	require.Error(t, net.Broadcast(msg, -1))
}
