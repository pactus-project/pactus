package network

import (
	"context"
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
	"github.com/pactus-project/pactus/util/pipeline"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func alwaysPropagate(_ *GossipMessage) PropagationPolicy {
	return Propagate
}

func makeTestNetwork(t *testing.T, conf *Config, opts []lp2p.Option) *network {
	t.Helper()

	pipe := pipeline.MockingPipeline[Event]()
	log := logger.NewSubLogger("_network", nil)
	net, err := makeNetwork(context.Background(), conf, log, pipe, opts)
	require.NoError(t, err)

	log.SetObj(testsuite.NewOverrideStringer(
		fmt.Sprintf("%s - %s: ", net.SelfID().ShortString(), t.Name()), net))

	assert.NoError(t, net.Start())

	return net
}

func testConfig() *Config {
	return &Config{
		ListenAddrStrings:         []string{},
		NetworkKey:                util.TempFilePath(),
		BootstrapAddrStrings:      []string{},
		MaxConns:                  16,
		EnableUDP:                 true,
		EnableNATService:          false,
		EnableUPnP:                false,
		EnableRelay:               false,
		EnableRelayService:        false,
		EnableMdns:                false,
		ForcePrivateNetwork:       true,
		NetworkName:               "test",
		DefaultPort:               testsuite.FindFreePort(),
		PeerStorePath:             util.TempFilePath(),
		StreamTimeout:             10 * time.Second,
		CheckConnectivityInterval: 1 * time.Second,
		MaxGossipMessageSize:      1 * 1024 * 1024, // 1 MB
		MaxStreamMessageSize:      8 * 1024 * 1024, // 8 MB
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

		case e := <-net.networkPipe.UnsafeGetChannel():
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

		case <-net.networkPipe.UnsafeGetChannel():
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
	net := makeTestNetwork(t, testConfig(), nil)

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
	fmt.Println("Starting Bootstrap node")
	networkB := makeTestNetwork(t, confB, []lp2p.Option{
		lp2p.ForceReachabilityPublic(),
	})
	bootstrapAddresses := []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%v", confB.DefaultPort, networkB.SelfID().String()),
		fmt.Sprintf("/ip4/127.0.0.1/udp/%v/quic-v1/p2p/%v", confB.DefaultPort, networkB.SelfID().String()),
		fmt.Sprintf("/ip6/::1/tcp/%v/p2p/%v", confB.DefaultPort, networkB.SelfID().String()),
		fmt.Sprintf("/ip6/::1/udp/%v/quic-v1/p2p/%v", confB.DefaultPort, networkB.SelfID().String()),
	}

	// Public and relay node
	confP := testConfig()
	confP.BootstrapAddrStrings = bootstrapAddresses
	confP.EnableRelay = false
	confP.EnableRelayService = true
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
	fmt.Println("Starting Private node M")
	networkM := makeTestNetwork(t, confM, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})

	// Private node N
	confN := testConfig()
	confN.EnableRelay = true
	confN.BootstrapAddrStrings = bootstrapAddresses
	fmt.Println("Starting Private node N")
	networkN := makeTestNetwork(t, confN, []lp2p.Option{
		lp2p.ForceReachabilityPrivate(),
	})

	// Private node X, doesn't join consensus topic and disabled relay
	confX := testConfig()
	confX.EnableRelay = false
	confX.BootstrapAddrStrings = bootstrapAddresses
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
		}, 2*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkN.Protocols()
			assert.Contains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, 2*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkP.Protocols()
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.Contains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, 2*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			protos := networkX.Protocols()
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Stop)
			assert.NotContains(c, protos, lp2pproto.ProtoIDv2Hop)
		}, 2*time.Second, 100*time.Millisecond)
	})

	t.Run("Reachability", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkB.ReachabilityStatus()
			require.Equal(c, "Public", reachability)
		}, time.Second*2, 500*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkM.ReachabilityStatus()
			require.Equal(c, "Private", reachability)
		}, time.Second*2, 500*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkN.ReachabilityStatus()
			require.Equal(c, "Private", reachability)
		}, time.Second*2, 500*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkP.ReachabilityStatus()
			require.Equal(c, "Public", reachability)
		}, time.Second*2, 500*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			reachability := networkP.ReachabilityStatus()
			require.Equal(c, "Public", reachability)
		}, time.Second*2, 100*time.Millisecond)
	})

	t.Run("all nodes have at least one connection to the bootstrap node B", func(t *testing.T) {
		fmt.Printf("Running %s\n", t.Name())

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			require.GreaterOrEqual(c, networkP.NumConnectedPeers(), 1) // Connected to B, M, N, X
		}, 5*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			require.GreaterOrEqual(c, networkB.NumConnectedPeers(), 1) // Connected to P, M, N, X
		}, 5*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			require.GreaterOrEqual(c, networkM.NumConnectedPeers(), 1) // Connected to B, P, N?
		}, 5*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			require.GreaterOrEqual(c, networkN.NumConnectedPeers(), 1) // Connected to B, P, M?
		}, 5*time.Second, 100*time.Millisecond)

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			require.GreaterOrEqual(c, networkX.NumConnectedPeers(), 1) // Connected to B, P
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

		require.Equal(t, "Public", networkP.ReachabilityStatus())
		require.Equal(t, "Public", networkB.ReachabilityStatus())
		require.Equal(t, "Private", networkM.ReachabilityStatus())
		require.Equal(t, "Private", networkN.ReachabilityStatus())
		require.Equal(t, "Private", networkX.ReachabilityStatus())
	})
}

func TestHostAddrs(t *testing.T) {
	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	addrs := net.HostAddrs()
	assert.Contains(t, addrs, fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", conf.DefaultPort))
	assert.Contains(t, addrs, fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic-v1", conf.DefaultPort))
}

func TestNetworkName(t *testing.T) {
	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

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

func TestLoadOrCreateKey(t *testing.T) {
	t.Run("Should load same network key after being created", func(t *testing.T) {
		keyPath := util.TempFilePath()

		// Create new valid key
		validKey, err := loadOrCreateKey(keyPath)
		assert.NoError(t, err)

		// Retrieve previously created valid key, the file path exists
		previousValidKey, err := loadOrCreateKey(keyPath)
		assert.NoError(t, err)

		assert.Equal(t, validKey.GetPublic(), previousValidKey.GetPublic())
	})

	t.Run("Should return error when file contains invalid data", func(t *testing.T) {
		tempFilePath := util.TempFilePath()

		err := util.WriteFile(tempFilePath, []byte("invalid_data"))
		assert.NoError(t, err)

		key, err := loadOrCreateKey(tempFilePath)
		assert.Error(t, err)
		assert.Nil(t, key)
	})

	t.Run("Should return error when file contains invalid private key", func(t *testing.T) {
		tempFilePath := util.TempFilePath()

		err := util.WriteFile(tempFilePath, []byte("00"))
		assert.NoError(t, err)

		key, err := loadOrCreateKey(tempFilePath)
		assert.Error(t, err)
		assert.Nil(t, key)
	})

	t.Run("Should return error when input path is directory", func(t *testing.T) {
		tempDir := t.TempDir()
		key, err := loadOrCreateKey(tempDir)
		assert.Error(t, err)
		assert.Nil(t, key)
	})

	t.Run("Should return error when input path is invalid", func(t *testing.T) {
		invalidPath := string([]byte{0x00})
		key, err := loadOrCreateKey(invalidPath)
		assert.Error(t, err)
		assert.Nil(t, key)
	})

	t.Run("Should trim spaces", func(t *testing.T) {
		tempFilePath := util.TempFilePath()

		err := util.WriteFile(tempFilePath,
			[]byte(" 080112406da99c6b29ac8093fad3a92327aaf87acf22dbb60927786db25880f025c0"+
				"4cb6f80873898709981d9b75795a191eab2d29bc7983ebcb4826e2b44566c85ea194 \r\n"))
		assert.NoError(t, err)

		key, err := loadOrCreateKey(tempFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, key)
	})
}
