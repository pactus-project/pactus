package firewall

import (
	"bytes"
	"testing"
	"time"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
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

	subLogger := logger.NewSubLogger("firewall", nil)
	peerSet := peerset.NewPeerSet(1 * time.Minute)
	st := state.MockingState(ts)
	net := network.MockingNetwork(ts, ts.RandPeerID())

	if conf == nil {
		conf = DefaultConfig()
	}
	require.NoError(t, conf.BasicCheck())
	firewall, err := NewFirewall(conf, net, peerSet, st, subLogger)
	if err != nil {
		return nil
	}

	assert.NotNil(t, firewall)
	badPeerID := ts.RandPeerID()
	goodPeerID := ts.RandPeerID()
	unknownPeerID := ts.RandPeerID()
	duplicatePeerID := peer.ID("peerX")
	duplicatePeerID2 := ts.RandPeerID()

	net.AddAnotherNetwork(network.MockingNetwork(ts, goodPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, unknownPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, badPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, duplicatePeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, duplicatePeerID2))

	firewall.peerSet.UpdateStatus(goodPeerID, status.StatusKnown)
	firewall.peerSet.UpdateStatus(duplicatePeerID, status.StatusKnown)
	firewall.peerSet.UpdateStatus(duplicatePeerID2, status.StatusConnected)
	firewall.peerSet.UpdateStatus(badPeerID, status.StatusBanned)

	firewall.peerSet.UpdateAddress(duplicatePeerID2, "/ip4/1.1.1.1", "")
	firewall.peerSet.UpdateAddress(duplicatePeerID, "/ip4/2.1.1.1", "")

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
	t.Run("Message from: unknown => should NOT close the connection", func(t *testing.T) {
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

func TestBannedAddress(t *testing.T) {
	conf := &Config{
		BannedNets: []string{
			"115.193.0.0/16",
			"240e:390:8a1:ae80:0000:0000:0000:0000/64",
		},
	}
	td := setup(t, conf)

	testCases := []struct {
		addr   string
		banned bool
	}{
		{
			addr:   "/ip4/115.193.157.138/tcp/21888",
			banned: true,
		},
		{
			addr:   "/ip4/10.10.10.10",
			banned: false,
		},
		{
			addr:   "/ip6/240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf/udp/21888",
			banned: true,
		},
		{
			addr:   "/ip6/2a01:4f9:4a:1d85::2",
			banned: false,
		},
	}

	for i, tc := range testCases {
		banned := td.firewall.IsAddressBanned(tc.addr, "")

		if tc.banned {
			assert.True(t, banned,
				"test %v failed, addr %v should be banned", i, tc.addr)
		} else {
			assert.False(t, banned,
				"test %v failed, addr %v should not be banned", i, tc.addr)
		}
	}

	td.firewall.config.DisallowDuplicateAddress = true

	tests := []struct {
		peerID     peer.ID
		remoteAddr string
		banned     bool
	}{
		{
			peerID:     "peerX",
			remoteAddr: "/ip4/1.2.3.4",
			banned:     true,
		},
		{
			peerID:     "peerZ",
			remoteAddr: "/ip4/1.1.1.1",
			banned:     true,
		},
		{
			peerID:     "peerZZ",
			remoteAddr: "/ip4/1.2.3.7",
			banned:     false,
		},
	}

	for _, tt := range tests {
		banned := td.firewall.IsAddressBanned(tt.remoteAddr, tt.peerID)
		if tt.banned {
			assert.True(t, banned)
		} else {
			assert.False(t, banned)
		}
	}
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
