package network

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pproto "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/proto"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func alwaysPropagate(_ *GossipMessage) PropagationPolicy {
	return Propagate
}

func makeTestNetwork(t *testing.T, conf *Config, opts []lp2p.Option) *network {
	t.Helper()

	log := logger.NewSubLogger("_network", nil)
	net, err := makeNetwork(conf, log, opts)
	require.NoError(t, err)

	log.SetObj(testsuite.NewOverrideStringer(
		fmt.Sprintf("%s - %s: ", net.SelfID().ShortString(), t.Name()), net))

	assert.NoError(t, net.Start())

	return net
}

func testConfig() *Config {
	return &Config{
		ListenAddrStrings:    []string{},
		NetworkKey:           util.TempFilePath(),
		BootstrapAddrStrings: []string{},
		MaxConns:             8,
		EnableUDP:            true,
		EnableNATService:     false,
		EnableUPnP:           false,
		EnableRelay:          false,
		EnableMdns:           false,
		ForcePrivateNetwork:  true,
		NetworkName:          "test",
		DefaultPort:          FindFreePort(),
		PeerStorePath:        util.TempFilePath(),
		StreamTimeout:        10 * time.Second,
	}
}

func shouldReceiveEvent(t *testing.T, net *network, eventType EventType) Event {
	t.Helper()

	timer := time.NewTimer(10 * time.Second)

	for {
		select {
		case <-timer.C:
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

	timer := time.NewTimer(100 * time.Millisecond)

	for {
		select {
		case <-timer.C:
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
	assert.NoError(t, net.JoinTopic(TopicIDBlock, alwaysPropagate))

	// Should stop peacefully
	net.Stop()
}

// In this test, we are setting up a simulated network environment that consists of these nodes:
//   - B is a Bootstrap node
//   - P is a Public and relay node
//   - M, N, and X are Private Nodes behind a Network Address Translation (NAT)
//   - M and N have relay enabled, while X does not.
//
// The test evaluates the following scenarios:
//   - Get supporting protocols
//   - Connection establishment to the bootstrap node
//   - Receiving gossip message
//   - Receiving direct message
//   - Receiving relayed message (Not covered yet!)
func TestNetwork(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// Bootstrap node
	confB := testConfig()
	confB.ListenAddrStrings = []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", confB.DefaultPort),
	}
	fmt.Println("Starting Bootstrap node")
	networkB := makeTestNetwork(t, confB, []lp2p.Option{
		lp2p.ForceReachabilityPublic(),
	})
	bootstrapAddresses := []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%v", confB.DefaultPort, networkB.SelfID().String()),
	}

	// Public and relay node
	confP := testConfig()
	confP.BootstrapAddrStrings = bootstrapAddresses
	confP.EnableRelay = false
	confP.EnableRelayService = true
	confP.ListenAddrStrings = []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", confP.DefaultPort),
	}
	fmt.Println("Starting Public node")
	networkP := makeTestNetwork(t, confP, []lp2p.Option{
		lp2p.ForceReachabilityPublic(),
	})
	publicAddrInfo, _ := lp2ppeer.AddrInfoFromString(
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", confP.DefaultPort, networkP.SelfID()))

	// Private node M
	confM := testConfig()
	confM.EnableRelay = true
	confM.BootstrapAddrStrings = bootstrapAddresses
	confM.ListenAddrStrings = []string{
		"/ip4/127.0.0.1/tcp/0",
	}
	fmt.Println("Starting Private node M")
	networkM := makeTestNetwork(t, confM, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})

	// Private node N
	confN := testConfig()
	confN.EnableRelay = true
	confN.BootstrapAddrStrings = bootstrapAddresses
	confN.ListenAddrStrings = []string{
		"/ip4/127.0.0.1/tcp/0",
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

	assert.NoError(t, networkB.JoinTopic(TopicIDBlock, alwaysPropagate))
	assert.NoError(t, networkP.JoinTopic(TopicIDBlock, alwaysPropagate))
	assert.NoError(t, networkM.JoinTopic(TopicIDBlock, alwaysPropagate))
	assert.NoError(t, networkN.JoinTopic(TopicIDBlock, alwaysPropagate))
	assert.NoError(t, networkX.JoinTopic(TopicIDBlock, alwaysPropagate))

	assert.NoError(t, networkB.JoinTopic(TopicIDConsensus, alwaysPropagate))
	assert.NoError(t, networkP.JoinTopic(TopicIDConsensus, alwaysPropagate))
	assert.NoError(t, networkM.JoinTopic(TopicIDConsensus, alwaysPropagate))
	assert.NoError(t, networkN.JoinTopic(TopicIDConsensus, alwaysPropagate))
	// Network X doesn't join the consensus topic

	time.Sleep(4 * time.Second)

	t.Run("Supported Protocols", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkM.Protocols()
			assert.Contains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkN.Protocols()
			assert.Contains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkP.Protocols()
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.Contains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkX.Protocols()
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)
	})

	t.Run("Reachability", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkB.ReachabilityStatus()
			assert.Equal(c, "Public", reachability)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkM.ReachabilityStatus()
			assert.Equal(c, "Private", reachability)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkN.ReachabilityStatus()
			assert.Equal(c, "Private", reachability)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkP.ReachabilityStatus()
			assert.Equal(c, "Public", reachability)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkP.ReachabilityStatus()
			assert.Equal(c, "Public", reachability)
		}, time.Second, 100*time.Millisecond)
	})

	t.Run("all nodes have at least one connection to the bootstrap node B", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.GreaterOrEqual(c, networkP.NumConnectedPeers(), 1) // Connected to B, M, N, X
		}, 5*time.Second, 100*time.Millisecond)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.GreaterOrEqual(c, networkB.NumConnectedPeers(), 1) // Connected to P, M, N, X
		}, 5*time.Second, 100*time.Millisecond)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.GreaterOrEqual(c, networkM.NumConnectedPeers(), 1) // Connected to B, P, N?
		}, 5*time.Second, 100*time.Millisecond)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.GreaterOrEqual(c, networkN.NumConnectedPeers(), 1) // Connected to B, P, M?
		}, 5*time.Second, 100*time.Millisecond)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.GreaterOrEqual(c, networkX.NumConnectedPeers(), 1) // Connected to B, P
		}, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("Gossip: all nodes receive gossip messages", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		msg := ts.RandBytes(64)

		networkP.Broadcast(msg, TopicIDBlock)

		eB := shouldReceiveEvent(t, networkB, EventTypeGossip).(*GossipMessage)
		eM := shouldReceiveEvent(t, networkM, EventTypeGossip).(*GossipMessage)
		eN := shouldReceiveEvent(t, networkN, EventTypeGossip).(*GossipMessage)
		eX := shouldReceiveEvent(t, networkX, EventTypeGossip).(*GossipMessage)

		assert.Equal(t, msg, eB.Data)
		assert.Equal(t, msg, eM.Data)
		assert.Equal(t, msg, eN.Data)
		assert.Equal(t, msg, eX.Data)
	})

	t.Run("only nodes subscribed to the consensus topic receive consensus gossip messages", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		msg := ts.RandBytes(64)

		networkP.Broadcast(msg, TopicIDConsensus)

		eB := shouldReceiveEvent(t, networkB, EventTypeGossip).(*GossipMessage)
		eM := shouldReceiveEvent(t, networkM, EventTypeGossip).(*GossipMessage)
		eN := shouldReceiveEvent(t, networkN, EventTypeGossip).(*GossipMessage)
		shouldNotReceiveEvent(t, networkX) // Not joined the consensus topic

		assert.Equal(t, msg, eB.Data)
		assert.Equal(t, msg, eM.Data)
		assert.Equal(t, msg, eN.Data)
	})

	t.Run("node P (public) is directly accessible by nodes M and N (private behind NAT)", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		require.NoError(t, networkM.host.Connect(networkM.ctx, *publicAddrInfo))

		msgM := ts.RandBytes(64)
		networkM.SendTo(msgM, networkP.SelfID())
		eP := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, networkM.SelfID(), eP.From)
		assert.Equal(t, msgM, readData(t, eP.Reader, len(msgM)))
	})

	t.Run("node P (public) is directly accessible by node X (private behind NAT, without relay)", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		require.NoError(t, networkX.host.Connect(networkX.ctx, *publicAddrInfo))

		msgX := ts.RandBytes(64)
		networkX.SendTo(msgX, networkP.SelfID())
		eP := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, networkX.SelfID(), eP.From)
		assert.Equal(t, msgX, readData(t, eP.Reader, len(msgX)))
	})

	t.Run("node P (public) is directly accessible by node B (bootstrap)", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		msgB := ts.RandBytes(64)

		networkB.SendTo(msgB, networkP.SelfID())
		eB := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, networkB.SelfID(), eB.From)
		assert.Equal(t, msgB, readData(t, eB.Reader, len(msgB)))
	})

	t.Run("Ignore broadcasting identical messages", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		msg := ts.RandBytes(64)

		networkM.Broadcast(msg, TopicIDBlock)
		networkN.Broadcast(msg, TopicIDBlock)

		eX := shouldReceiveEvent(t, networkX, EventTypeGossip).(*GossipMessage)

		assert.Equal(t, msg, eX.Data)
		assert.NotEqual(t, eX.From, networkM.SelfID(), "network X has no direct connection with M")
		assert.NotEqual(t, eX.From, networkN.SelfID(), "network X has no direct connection with N")

		shouldNotReceiveEvent(t, networkX)
	})

	t.Run("node X (private, not connected via relay) is not accessible by node M", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		msgM := ts.RandBytes(64)
		networkM.SendTo(msgM, networkX.SelfID())
	})

	// TODO: How to test this?
	// t.Run("nodes M and N (private, connected via relay) can communicate using the relay node R", func(t *testing.T) {
	// 	msgM := ts.RandBytes(64)
	// 	networkM.SendTo(msgM, networkN.SelfID())
	// 	eM := shouldReceiveEvent(t, networkN, EventTypeStream).(*StreamMessage)
	// 	assert.Equal(t, msgM, readData(t, eM.Reader, len(msgM)))
	// })

	t.Run("closing connection", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		msgB := ts.RandBytes(64)

		networkP.Stop()
		networkB.CloseConnection(networkP.SelfID())
		e := shouldReceiveEvent(t, networkB, EventTypeDisconnect).(*DisconnectEvent)
		assert.Equal(t, networkP.SelfID(), e.PeerID)
		networkB.SendTo(msgB, networkP.SelfID())
	})

	t.Run("Reachability Status", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		assert.Equal(t, "Public", networkP.ReachabilityStatus())
		assert.Equal(t, "Public", networkB.ReachabilityStatus())
		assert.Equal(t, "Private", networkM.ReachabilityStatus())
		assert.Equal(t, "Private", networkN.ReachabilityStatus())
		assert.Equal(t, "Private", networkX.ReachabilityStatus())
	})
}

func TestHostAddrs(t *testing.T) {
	conf := testConfig()
	net, err := NewNetwork(conf)
	assert.NoError(t, err)

	addrs := net.HostAddrs()
	assert.Contains(t, addrs, fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", conf.DefaultPort))
	assert.Contains(t, addrs, fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic-v1", conf.DefaultPort))
}

func TestNetworkName(t *testing.T) {
	conf := testConfig()
	net, err := NewNetwork(conf)
	assert.NoError(t, err)

	assert.Equal(t, conf.NetworkName, net.Name())
}

func TestConnections(t *testing.T) {
	t.Parallel() // run the tests in parallel

	tests := []struct {
		bootstrapAddr string
		peerAddr      string
	}{
		{"/ip4/127.0.0.1/tcp/%d", "/ip4/127.0.0.1/tcp/0"},
		{"/ip6/::1/tcp/%d", "/ip6/::1/tcp/0"},
		{"/ip4/127.0.0.1/udp/%d/quic-v1", "/ip4/127.0.0.1/udp/0/quic-v1"},
		{"/ip6/::1/udp/%d/quic-v1", "/ip6/::1/udp/0/quic-v1"},
	}

	for no, tt := range tests {
		// Bootstrap node
		confB := testConfig()
		bootstrapAddr := fmt.Sprintf(tt.bootstrapAddr, confB.DefaultPort)
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
		confP.ListenAddrStrings = []string{tt.peerAddr}
		fmt.Println("Starting Public node")
		networkP := makeTestNetwork(t, confP, []lp2p.Option{
			lp2p.ForceReachabilityPublic(),
		})

		t.Run(fmt.Sprintf("Running test %d: %s <-> %s ... ",
			no, bootstrapAddr, tt.peerAddr), func(t *testing.T) {
			t.Parallel() // run the tests in parallel

			checkConnection(t, networkP, networkB)
		})
	}
}

func checkConnection(t *testing.T, networkP, networkB *network) {
	t.Helper()

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.GreaterOrEqual(c, networkP.NumConnectedPeers(), 1)
		assert.GreaterOrEqual(c, networkB.NumConnectedPeers(), 1)
	}, 5*time.Second, 100*time.Millisecond)

	msg := []byte("test-msg")
	networkP.SendTo(msg, networkB.SelfID())
	e := shouldReceiveEvent(t, networkB, EventTypeStream).(*StreamMessage)
	assert.Equal(t, networkP.SelfID(), e.From)
	assert.Equal(t, msg, readData(t, e.Reader, len(msg)))

	networkB.Stop()
	networkP.Stop()
}
