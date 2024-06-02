package firewall

import (
	"bytes"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	firewall      *Firewall
	badPeerID     peer.ID
	goodPeerID    peer.ID
	unknownPeerID peer.ID
	network       *network.MockNetwork
	state         *state.MockState
}

func setup(t *testing.T, conf *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	if conf == nil {
		conf = DefaultConfig()
	}

	subLogger := logger.NewSubLogger("firewall", nil)
	peerSet := peerset.NewPeerSet(1 * time.Minute)
	st := state.MockingState(ts)
	net := network.MockingNetwork(ts, ts.RandPeerID())

	// TODO: It should be like: "/ip6/2a01:4f9:4a:1d85::2"
	conf.BlackListAddresses = []string{"/ip6/2a01:4f9:4a:1d85::2/tcp/21888"}
	require.NoError(t, conf.BasicCheck())
	firewall := NewFirewall(conf, net, peerSet, st, subLogger)
	assert.NotNil(t, firewall)
	badPeerID := ts.RandPeerID()
	goodPeerID := ts.RandPeerID()
	unknownPeerID := ts.RandPeerID()

	net.AddAnotherNetwork(network.MockingNetwork(ts, goodPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, unknownPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, badPeerID))

	firewall.peerSet.UpdateStatus(goodPeerID, peerset.StatusCodeKnown)
	firewall.peerSet.UpdateStatus(badPeerID, peerset.StatusCodeBanned)

	return &testData{
		TestSuite:     ts,
		firewall:      firewall,
		network:       net,
		state:         st,
		badPeerID:     badPeerID,
		goodPeerID:    goodPeerID,
		unknownPeerID: unknownPeerID,
	}
}

func TestInvalidBundlesCounter(t *testing.T) {
	td := setup(t, nil)

	assert.Nil(t, td.firewall.OpenGossipBundle([]byte("bad"), td.unknownPeerID))
	assert.Nil(t, td.firewall.OpenGossipBundle(nil, td.unknownPeerID))

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), -1, td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	d, _ := bdl.Encode()
	assert.Nil(t, td.firewall.OpenGossipBundle(d, td.unknownPeerID))

	p := td.firewall.peerSet.GetPeer(td.unknownPeerID)
	assert.Equal(t, p.InvalidBundles, 3)
}

func TestGossipMessage(t *testing.T) {
	t.Run("Message  from: unknown => should NOT close the connection", func(t *testing.T) {
		td := setup(t, nil)

		bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.unknownPeerID))
		assert.NotNil(t, td.firewall.OpenGossipBundle(d, td.unknownPeerID))
		assert.False(t, td.network.IsClosed(td.unknownPeerID))
	})

	t.Run("Message  from: bad => should close the connection", func(t *testing.T) {
		td := setup(t, nil)

		bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.badPeerID))
		assert.Nil(t, td.firewall.OpenGossipBundle(d, td.badPeerID))
		assert.True(t, td.network.IsClosed(td.badPeerID))
	})

	t.Run("Message is nil => should close the connection", func(t *testing.T) {
		td := setup(t, nil)

		assert.Nil(t, td.firewall.OpenGossipBundle(nil, td.unknownPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		td := setup(t, nil)

		bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.goodPeerID))
		assert.NotNil(t, td.firewall.OpenGossipBundle(d, td.goodPeerID))
		assert.False(t, td.network.IsClosed(td.goodPeerID))
	})
}

func TestStreamMessage(t *testing.T) {
	t.Run("Message is nil => should close the connection", func(t *testing.T) {
		td := setup(t, nil)

		assert.False(t, td.network.IsClosed(td.badPeerID))
		assert.Nil(t, td.firewall.OpenStreamBundle(bytes.NewReader(nil), td.badPeerID))
		assert.True(t, td.network.IsClosed(td.badPeerID))
	})

	t.Run("Message from: bad => should close the connection", func(t *testing.T) {
		td := setup(t, nil)

		bdl := bundle.NewBundle(message.NewBlocksRequestMessage(td.RandInt(100), 1, 100))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.badPeerID))
		assert.Nil(t, td.firewall.OpenStreamBundle(bytes.NewReader(d), td.badPeerID))
		assert.True(t, td.network.IsClosed(td.badPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		td := setup(t, nil)

		bdl := bundle.NewBundle(message.NewBlocksRequestMessage(td.RandInt(100), 1, 100))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.goodPeerID))
		assert.NotNil(t, td.firewall.OpenStreamBundle(bytes.NewReader(d), td.goodPeerID))
		assert.False(t, td.network.IsClosed(td.goodPeerID))
	})
}

func TestUpdateLastReceived(t *testing.T) {
	td := setup(t, nil)

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	d, _ := bdl.Encode()
	now := time.Now().UnixNano()
	assert.NotNil(t, td.firewall.OpenGossipBundle(d, td.goodPeerID))

	peerGood := td.firewall.peerSet.GetPeer(td.goodPeerID)
	assert.GreaterOrEqual(t, peerGood.LastReceived.UnixNano(), now)
}

func TestNetworkFlags(t *testing.T) {
	td := setup(t, nil)

	// TODO: add tests for Mainnet and Testnet flags
	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	assert.NoError(t, td.firewall.checkBundle(bdl))

	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	assert.Error(t, td.firewall.checkBundle(bdl))

	bdl.Flags = 0
	assert.Error(t, td.firewall.checkBundle(bdl))

	td.state.TestParams.BlockVersion = 0x3f // changing genesis hash
	bdl.Flags = 1
	assert.Error(t, td.firewall.checkBundle(bdl))
}

func TestBlackListAddress(t *testing.T) {
	td := setup(t, nil)

	testCases := []struct {
		addr        string
		blacklisted bool
	}{
		{
			addr:        "/ip4/115.193.157.138/tcp/21888",
			blacklisted: true,
		},
		// TOD: uncomment me
		// {
		// 	addr:        "/ip4/115.193.157.138",
		// 	blacklisted: true,
		// },
		{
			addr:        "/ip4/10.10.10.10",
			blacklisted: false,
		},
		{
			addr:        "/ip4/10.10.10.10/udp/21888",
			blacklisted: false,
		},
		{
			addr:        "/ip6/2a01:4f9:4a:1d85::2",
			blacklisted: false,
		},
	}

	for i, tc := range testCases {
		blacklisted := td.firewall.IsBlackListAddress(tc.addr)

		if tc.blacklisted {
			assert.True(t, blacklisted,
				"test %v failed, addr %v should blacklisted", i, tc.addr)
		} else {
			assert.False(t, blacklisted,
				"test %v failed, addr %v should not blacklisted", i, tc.addr)
		}
	}
}

func TestAllowBlockRequest(t *testing.T) {
	conf := DefaultConfig()
	conf.RateLimit.BlockTopic = 1

	td := setup(t, conf)

	assert.True(t, td.firewall.AllowBlockRequest())
	assert.False(t, td.firewall.AllowBlockRequest())
}

func TestAllowTransactionRequest(t *testing.T) {
	conf := DefaultConfig()
	conf.RateLimit.TransactionTopic = 1

	td := setup(t, conf)

	assert.True(t, td.firewall.AllowTransactionRequest())
	assert.False(t, td.firewall.AllowTransactionRequest())
}

func TestAllowConsensusRequest(t *testing.T) {
	conf := DefaultConfig()
	conf.RateLimit.ConsensusTopic = 1

	td := setup(t, conf)

	assert.True(t, td.firewall.AllowConsensusRequest())
	assert.False(t, td.firewall.AllowConsensusRequest())
}
