package firewall

import (
	"bytes"
	"testing"
	"time"

	"github.com/pactus-project/pactus/genesis"
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
	bannedPeerID  peer.ID
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
	bannedPeerID := ts.RandPeerID()
	goodPeerID := ts.RandPeerID()
	unknownPeerID := ts.RandPeerID()

	net.AddAnotherNetwork(network.MockingNetwork(ts, goodPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, unknownPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, bannedPeerID))

	firewall.peerSet.UpdateStatus(goodPeerID, status.StatusKnown)
	firewall.peerSet.UpdateStatus(bannedPeerID, status.StatusBanned)

	return &testData{
		TestSuite:     ts,
		firewall:      firewall,
		network:       net,
		state:         st,
		bannedPeerID:  bannedPeerID,
		goodPeerID:    goodPeerID,
		unknownPeerID: unknownPeerID,
	}
}

func (td *testData) testGossipBundle() []byte {
	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	d, _ := bdl.Encode()

	return d
}

func (td *testData) testStreamBundle() []byte {
	bdl := bundle.NewBundle(message.NewBlocksRequestMessage(td.RandInt(100), 1, 100))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	d, _ := bdl.Encode()

	return d
}

func TestDecodeBundles(t *testing.T) {
	td := setup(t, nil)

	testCases := []struct {
		name    string
		data    string
		peerID  string
		wantErr bool
	}{
		{
			name:    "invalid data",
			data:    "bad0",
			wantErr: true,
		},
		{
			name:    "nil data",
			data:    "",
			wantErr: true,
		},
		{
			name: "invalid bundle (round is -1)",
			data: "a4" + // Map with 4 key-value pairs
				"01" + "01" + // Key 1 (Flags), Value: 1 (Mainnet)
				"02" + "06" + // Key 2 (Message Type), Value: 6 (QueryVote)
				"03" + "581d" + // Key 2 (Message), Value: 30 Bytes
				"" + "a3" + // Map with 3 key-value pairs
				"" + "01" + "1864" + // Key 1 (Height), Value: 100
				"" + "02" + "20" + // Key 2 (Round), Value: -1
				"" + "03" + "5501aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + // Key 3 (Querier), Value: 21 Bytes
				"04" + "05", // Key 4 (Sequence number), Value: 5
			wantErr: true,
		},

		{
			name: "valid bundle (invalid network, Testnet)",
			data: "a4" + // Map with 4 key-value pairs
				"01" + "02" + // Key 1 (Flags), Value: 1 (Testnet)
				"02" + "06" + // Key 2 (Message Type), Value: 6 (QueryVote)
				"03" + "581d" + // Key 2 (Message), Value: 30 Bytes
				"" + "a3" + // Map with 3 key-value pairs
				"" + "01" + "1864" + // Key 1 (Height), Value: 100
				"" + "02" + "00" + // Key 2 (Round), Value: 0
				"" + "03" + "5501aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + // Key 3 (Querier), Value: 21 Bytes
				"04" + "05", // Key 4 (Sequence number), Value: 5
			wantErr: true,
		},
		{
			name: "valid bundle",
			data: "a4" + // Map with 4 key-value pairs
				"01" + "01" + // Key 1 (Flags), Value: 1 (Mainnet)
				"02" + "06" + // Key 2 (Message Type), Value: 6 (QueryVote)
				"03" + "581d" + // Key 2 (Message), Value: 30 Bytes
				"" + "a3" + // Map with 3 key-value pairs
				"" + "01" + "1864" + // Key 1 (Height), Value: 100
				"" + "02" + "00" + // Key 2 (Round), Value: 0
				"" + "03" + "5501aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + // Key 3 (Querier), Value: 21 Bytes
				"04" + "05", // Key 4 (Sequence number), Value: 5
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs := td.DecodingHex(tc.data)
			_, err := td.firewall.OpenGossipBundle(bs, td.unknownPeerID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	p := td.firewall.peerSet.GetPeer(td.unknownPeerID)
	assert.Equal(t, 5, p.ReceivedBundles)
	assert.Equal(t, 4, p.InvalidBundles)
}

func TestGossipMessage(t *testing.T) {
	t.Run("Message is nil", func(t *testing.T) {
		td := setup(t, nil)

		_, err := td.firewall.OpenGossipBundle(nil, td.unknownPeerID)
		require.Error(t, err)
		assert.False(t, td.network.IsClosed(td.unknownPeerID))
	})

	t.Run("Message from banned peer", func(t *testing.T) {
		td := setup(t, nil)

		data := td.testGossipBundle()

		assert.False(t, td.network.IsClosed(td.bannedPeerID))
		_, err := td.firewall.OpenGossipBundle(data, td.bannedPeerID)
		require.ErrorIs(t, err, PeerBannedError{
			PeerID:  td.bannedPeerID,
			Address: "",
		})
		assert.True(t, td.network.IsClosed(td.bannedPeerID))
	})

	t.Run("Stream message as gossip message", func(t *testing.T) {
		td := setup(t, nil)

		data := td.testStreamBundle()

		assert.False(t, td.network.IsClosed(td.unknownPeerID))
		_, err := td.firewall.OpenGossipBundle(data, td.unknownPeerID)
		require.ErrorIs(t, err, ErrGossipMessage)
		assert.True(t, td.network.IsClosed(td.unknownPeerID))
	})

	t.Run("Ok", func(t *testing.T) {
		td := setup(t, nil)

		data := td.testGossipBundle()

		assert.False(t, td.network.IsClosed(td.goodPeerID))
		_, err := td.firewall.OpenGossipBundle(data, td.goodPeerID)
		require.NoError(t, err)
		assert.False(t, td.network.IsClosed(td.goodPeerID))
	})
}

func TestStreamMessage(t *testing.T) {
	t.Run("Message is nil", func(t *testing.T) {
		td := setup(t, nil)

		assert.False(t, td.network.IsClosed(td.unknownPeerID))
		_, err := td.firewall.OpenStreamBundle(bytes.NewReader(nil), td.unknownPeerID)
		assert.Error(t, err)
	})

	t.Run("Message from banned peer", func(t *testing.T) {
		td := setup(t, nil)

		data := td.testStreamBundle()

		assert.False(t, td.network.IsClosed(td.bannedPeerID))
		_, err := td.firewall.OpenStreamBundle(bytes.NewReader(data), td.bannedPeerID)
		assert.ErrorIs(t, err, PeerBannedError{
			PeerID:  td.bannedPeerID,
			Address: "",
		})

		assert.True(t, td.network.IsClosed(td.bannedPeerID))
	})

	t.Run("Gossip message as direct message", func(t *testing.T) {
		td := setup(t, nil)

		data := td.testGossipBundle()

		assert.False(t, td.network.IsClosed(td.unknownPeerID))
		_, err := td.firewall.OpenStreamBundle(bytes.NewReader(data), td.unknownPeerID)
		require.ErrorIs(t, err, ErrStreamMessage)
		assert.True(t, td.network.IsClosed(td.unknownPeerID))
	})

	t.Run("Ok", func(t *testing.T) {
		td := setup(t, nil)

		data := td.testStreamBundle()

		assert.False(t, td.network.IsClosed(td.goodPeerID))
		_, err := td.firewall.OpenStreamBundle(bytes.NewReader(data), td.goodPeerID)
		require.NoError(t, err)
		assert.False(t, td.network.IsClosed(td.goodPeerID))
	})
}

func TestUpdateLastReceived(t *testing.T) {
	td := setup(t, nil)

	data := td.testGossipBundle()
	now := time.Now().UnixNano()
	_, err := td.firewall.OpenGossipBundle(data, td.goodPeerID)
	require.NoError(t, err)

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
		peerID := td.RandPeerID()
		td.firewall.peerSet.UpdateAddress(peerID, tc.addr, "inbound")
		data := td.testGossipBundle()
		_, err := td.firewall.OpenGossipBundle(data, peerID)

		if tc.banned {
			expectedErr := PeerBannedError{
				PeerID:  peerID,
				Address: tc.addr,
			}
			assert.ErrorIs(t, err, expectedErr,
				"test %v failed, addr %v should be banned", i, tc.addr)
		} else {
			assert.NoError(t, err,
				"test %v failed, addr %v should not be banned", i, tc.addr)
		}
	}
}

func TestNetworkFlagsMainnet(t *testing.T) {
	td := setup(t, nil)

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	assert.NoError(t, td.firewall.checkBundle(bdl))

	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	assert.Error(t, td.firewall.checkBundle(bdl))

	bdl.Flags = 0
	assert.Error(t, td.firewall.checkBundle(bdl))
}

func TestNetworkFlagsTestnet(t *testing.T) {
	td := setup(t, nil)
	td.state.TestGenesis = genesis.TestnetGenesis()

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	assert.NoError(t, td.firewall.checkBundle(bdl))

	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	assert.Error(t, td.firewall.checkBundle(bdl))

	bdl.Flags = 0
	assert.Error(t, td.firewall.checkBundle(bdl))
}

func TestNetworkFlagsLocalnet(t *testing.T) {
	td := setup(t, nil)
	td.state.TestParams.BlockVersion = 0x3f // changing genesis hash

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	assert.Error(t, td.firewall.checkBundle(bdl))

	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	assert.Error(t, td.firewall.checkBundle(bdl))

	bdl.Flags = 0
	assert.NoError(t, td.firewall.checkBundle(bdl))
}

func TestParseP2PAddr(t *testing.T) {
	td := setup(t, nil)

	tests := []struct {
		name        string
		address     string
		expectedIP  string
		expectError bool
	}{
		{
			name:       "Valid IPv4 with p2p",
			address:    "/ip4/84.247.165.249/tcp/21888/p2p/12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			expectedIP: "84.247.165.249",
		},
		{
			name:       "Valid IPv4 without p2p",
			address:    "/ip4/115.193.157.138/tcp/21888",
			expectedIP: "115.193.157.138",
		},
		{
			name: "Valid IPv6 with p2p",
			address: "/ip6/240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf/tcp/21888/p2p/" +
				"12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			expectedIP: "240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf",
		},
		{
			name:        "Invalid address",
			address:     "/invalid/address",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := td.firewall.getIPFromMultiAddress(tt.address)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIP, ip)
			}
		})
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
