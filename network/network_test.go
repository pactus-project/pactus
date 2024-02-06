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

func alwaysPropagate(_ *GossipMessage) bool {
	return true
}

func makeTestNetwork(t *testing.T, conf *Config, opts []lp2p.Option) *network {
	t.Helper()

	log := logger.NewSubLogger("_network", nil)
	net, err := newNetwork(conf, log, opts)
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
		DefaultPort:          12345,
	}
}

func shouldReceiveEvent(t *testing.T, net *network, eventType EventType) Event {
	t.Helper()

	timeout := time.NewTimer(4 * time.Second)

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

// In this test, we are setting up a simulated network environment that consists of these nodes:
//   - B is a Bootstrap node
//   - P is a Public and relay node
//   - M, N, and X are Private Nodes behind a Network Address Translation (NAT)
//   - M and N have relay enabled, while X does not.
//
// The test will evaluate the following scenarios:
//   - Connection establishment to the bootstrap node
//   - General and consensus topics and gossip message
//   - Direct and relayed stream communication between nodes
func TestNetwork(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

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

	// Public and relay node
	confP := testConfig()
	confP.BootstrapAddrStrings = bootstrapAddresses
	confP.EnableRelay = false
	confP.EnableRelayService = true
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

	assert.NoError(t, networkB.JoinGeneralTopic(alwaysPropagate))
	assert.NoError(t, networkP.JoinGeneralTopic(alwaysPropagate))
	assert.NoError(t, networkM.JoinGeneralTopic(alwaysPropagate))
	assert.NoError(t, networkN.JoinGeneralTopic(alwaysPropagate))
	assert.NoError(t, networkX.JoinGeneralTopic(alwaysPropagate))

	assert.NoError(t, networkB.JoinConsensusTopic(alwaysPropagate))
	assert.NoError(t, networkP.JoinConsensusTopic(alwaysPropagate))
	assert.NoError(t, networkM.JoinConsensusTopic(alwaysPropagate))
	assert.NoError(t, networkN.JoinConsensusTopic(alwaysPropagate))
	// Network X doesn't join the consensus topic

	time.Sleep(2 * time.Second)

	t.Run("Supported Protocols", func(t *testing.T) {
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkM.Protocols()
			assert.Contains(t, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(t, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkN.Protocols()
			assert.Contains(t, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(t, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkP.Protocols()
			assert.NotContains(t, protos, lp2pproto.ProtoIDv2Stop)
			assert.Contains(t, protos, lp2pproto.ProtoIDv2Hop)
		}, time.Second, 100*time.Millisecond)
	})

	t.Run("all nodes have at least one connection to the bootstrap node B", func(t *testing.T) {
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

	t.Run("Gossip: all nodes receive general gossip messages", func(t *testing.T) {
		msg := ts.RandBytes(64)

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
		msg := ts.RandBytes(64)

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

		msgM := ts.RandBytes(64)
		require.NoError(t, networkM.SendTo(msgM, networkP.SelfID()))
		eP := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eP.From, networkM.SelfID())
		assert.Equal(t, readData(t, eP.Reader, len(msgM)), msgM)
	})

	t.Run("node P (public) is directly accessible by node X (private behind NAT, without relay)", func(t *testing.T) {
		require.NoError(t, networkX.host.Connect(networkX.ctx, *publicAddrInfo))

		msgX := ts.RandBytes(64)
		require.NoError(t, networkX.SendTo(msgX, networkP.SelfID()))
		eP := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eP.From, networkX.SelfID())
		assert.Equal(t, readData(t, eP.Reader, len(msgX)), msgX)
	})

	t.Run("node P (public) is directly accessible by node B (bootstrap)", func(t *testing.T) {
		msgB := ts.RandBytes(64)

		require.NoError(t, networkB.SendTo(msgB, networkP.SelfID()))
		eB := shouldReceiveEvent(t, networkP, EventTypeStream).(*StreamMessage)
		assert.Equal(t, eB.From, networkB.SelfID())
		assert.Equal(t, readData(t, eB.Reader, len(msgB)), msgB)
	})

	t.Run("Ignore broadcasting identical messages", func(t *testing.T) {
		msg := ts.RandBytes(64)

		require.NoError(t, networkM.Broadcast(msg, TopicIDGeneral))
		require.NoError(t, networkN.Broadcast(msg, TopicIDGeneral))

		eX := shouldReceiveEvent(t, networkX, EventTypeGossip).(*GossipMessage)

		assert.Equal(t, eX.Data, msg)
		assert.NotEqual(t, eX.From, networkM.SelfID(), "network X has no direct connection with M")
		assert.NotEqual(t, eX.From, networkN.SelfID(), "network X has no direct connection with N")

		shouldNotReceiveEvent(t, networkX)
	})

	t.Run("node X (private, not connected via relay) is not accessible by node M", func(t *testing.T) {
		msgM := ts.RandBytes(64)
		require.Error(t, networkM.SendTo(msgM, networkX.SelfID()))
	})

	// TODO: How to test this?
	// t.Run("nodes M and N (private, connected via relay) can communicate using the relay node R", func(t *testing.T) {
	// 	msgM := ts.RandBytes(64)
	// 	require.NoError(t, networkM.SendTo(msgM, networkN.SelfID()))
	// 	eM := shouldReceiveEvent(t, networkN, EventTypeStream).(*StreamMessage)
	// 	assert.Equal(t, readData(t, eM.Reader, len(msgM)), msgM)
	// })

	t.Run("closing connection", func(t *testing.T) {
		msgB := ts.RandBytes(64)

		networkP.Stop()
		networkB.CloseConnection(networkP.SelfID())
		e := shouldReceiveEvent(t, networkB, EventTypeDisconnect).(*DisconnectEvent)
		assert.Equal(t, e.PeerID, networkP.SelfID())
		require.Error(t, networkB.SendTo(msgB, networkP.SelfID()))
	})

	t.Run("Reachability Status", func(t *testing.T) {
		assert.Equal(t, networkP.ReachabilityStatus(), "Public")
		assert.Equal(t, networkB.ReachabilityStatus(), "Public")
		assert.Equal(t, networkM.ReachabilityStatus(), "Private")
		assert.Equal(t, networkN.ReachabilityStatus(), "Private")
		assert.Equal(t, networkX.ReachabilityStatus(), "Private")
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
