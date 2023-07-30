package sync

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	config      *Config
	state       *state.MockState
	consMgr     consensus.Manager
	consMocks   []*consensus.MockConsensus
	network     *network.MockNetwork
	sync        *synchronizer
	broadcastCh chan message.Message
}

type OverrideStringer struct {
	sync *synchronizer
	name string
}

func init() {
	LatestBlockInterval = 23
}

func (o *OverrideStringer) String() string {
	return o.name + o.sync.String()
}

func testConfig() *Config {
	return &Config{
		Moniker:         "test",
		HeartBeatTimer:  0, // Disabling heartbeat
		SessionTimeout:  time.Second * 1,
		NodeNetwork:     true,
		BlockPerMessage: 11,
		MaxOpenSessions: 4,
		CacheSize:       1000,
		Firewall:        firewall.DefaultConfig(),
	}
}

func setup(t *testing.T, config *Config) *testData {
	ts := testsuite.NewTestSuite(t)

	if config == nil {
		config = testConfig()
		config.Moniker = "Alice"
	}
	signers := []crypto.Signer{ts.RandomSigner(), ts.RandomSigner()}
	state := state.MockingState(ts)
	consMgr, consMocks := consensus.MockingManager(ts, signers)
	broadcastCh := make(chan message.Message, 1000)
	network := network.MockingNetwork(ts, ts.RandomPeerID())

	Sync, err := NewSynchronizer(config,
		signers,
		state,
		consMgr,
		network,
		broadcastCh,
	)
	assert.NoError(t, err)
	sync := Sync.(*synchronizer)

	td := &testData{
		TestSuite:   ts,
		config:      config,
		state:       state,
		consMgr:     consMgr,
		consMocks:   consMocks,
		network:     network,
		sync:        sync,
		broadcastCh: broadcastCh,
	}

	td.addBlocks(t, state, 21)

	assert.NoError(t, td.sync.Start())
	assert.Equal(t, td.sync.Moniker(), config.Moniker)

	td.shouldPublishMessageWithThisType(t, network, message.TypeHello) // Alice key 1
	td.shouldPublishMessageWithThisType(t, network, message.TypeHello) // Alice key 2

	logger.Info("setup finished, running the tests", "name", t.Name())

	return td
}

func shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) *bundle.Bundle {
	timeout := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldPublishMessageWithThisType %v: Timeout, test: %v", msgType, t.Name()))
			return nil
		case b := <-net.BroadcastCh:
			net.SendToOthers(b.Data, b.Target)
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.Equal(t, bdl.Initiator, net.ID)

			// -----------
			// Check flags
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagCarrierLibP2P), "invalid flag: %v", bdl)
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)

			if b.Target == nil {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted), "invalid flag: %v", bdl)
			}

			if bdl.Message.Type() == message.TypeHello {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHelloMessage), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHelloMessage), "invalid flag: %v", bdl)
			}
			// -----------

			if bdl.Message.Type() == msgType {
				logger.Info("shouldPublishMessageWithThisType", "bundle", bdl, "type", msgType.String())
				return bdl
			}
		}
	}
}

func (td *testData) shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork,
	msgType message.Type) *bundle.Bundle {
	return shouldPublishMessageWithThisType(t, net, msgType)
}

func (td *testData) shouldNotPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case b := <-net.BroadcastCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.NotEqual(t, bdl.Message.Type(), msgType)
		}
	}
}

func (td *testData) receivingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) error {
	bdl := bundle.NewBundle(from, msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)
	return sync.processIncomingBundle(bdl)
}

func addBlocks(t *testing.T, state *state.MockState, count int) {
	h := state.LastBlockHeight()
	state.CommitTestBlocks(count)
	assert.Equal(t, h+uint32(count), state.LastBlockHeight())
}

func (td *testData) addBlocks(t *testing.T, state *state.MockState, count int) {
	addBlocks(t, state, count)
}

func (td *testData) addPeer(t *testing.T, pub crypto.PublicKey, pid peer.ID, nodeNetwork bool) {
	td.sync.peerSet.UpdatePeerInfo(pid, peerset.StatusCodeKnown, t.Name(),
		version.Agent(), pub.(*bls.PublicKey), nodeNetwork)
}

func (td *testData) addPeerToCommittee(t *testing.T, pid peer.ID, pub crypto.PublicKey) {
	if pub == nil {
		pub, _ = td.RandomBLSKeyPair()
	}
	td.addPeer(t, pub, pid, true)
	val := validator.NewValidator(pub.(*bls.PublicKey), td.RandInt32(1000))
	// Note: This may not be completely accurate, but it poses no harm for testing purposes.
	val.UpdateLastJoinedHeight(td.state.TestCommittee.Proposer(0).LastJoinedHeight() + 1)
	td.state.TestStore.UpdateValidator(val)
	td.state.TestCommittee.Update(0, []*validator.Validator{val})
	require.True(t, td.state.TestCommittee.Contains(pub.Address()))

	for _, cons := range td.consMocks {
		cons.SetActive(cons.Signer.PublicKey().EqualsTo(pub))
	}
}

func (td *testData) checkPeerStatus(t *testing.T, pid peer.ID, code peerset.StatusCode) {
	require.Equal(t, td.sync.peerSet.GetPeer(pid).Status, code)
}

func TestStop(t *testing.T) {
	td := setup(t, nil)

	// Should stop gracefully.
	td.sync.Stop()
}

func TestTestNetFlags(t *testing.T) {
	td := setup(t, nil)

	td.state.TestParams.BlockVersion = 0x3f
	bdl := td.sync.prepareBundle(message.NewHeartBeatMessage(1, 0, td.RandomHash()))
	require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)
	require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)
}

func TestDownload(t *testing.T) {
	td := setup(t, nil)

	ourBlockHeight := td.state.LastBlockHeight()
	b := td.GenerateTestBlock(nil, nil)
	c := td.GenerateTestCertificate(b.Hash())
	pid := td.RandomPeerID()
	msg := message.NewBlockAnnounceMessage(ourBlockHeight+LatestBlockInterval+1, b, c)

	t.Run("try to query latest blocks, but the peer is not known", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
	})

	t.Run("try to download blocks, but the peer is not known", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
	})

	t.Run("try to download blocks, but the peer is not a network node", func(t *testing.T) {
		pub, _ := td.RandomBLSKeyPair()
		td.addPeer(t, pub, pid, false)

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendSuccess, 0)
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendFailed, 0)
	})

	t.Run("try to download blocks and the peer is a network node", func(t *testing.T) {
		pub, _ := td.RandomBLSKeyPair()
		td.addPeer(t, pub, pid, true)

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendSuccess, 1)
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendFailed, 0)

		// Let's close the opened session
		td.sync.peerSet.CloseSession(bdl.Message.(*message.BlocksRequestMessage).SessionID)
	})

	t.Run("download request is rejected", func(t *testing.T) {
		session := td.sync.peerSet.OpenSession(pid)
		msg2 := message.NewBlocksResponseMessage(message.ResponseCodeRejected, session.SessionID(), 1, nil, nil)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg2, pid))
		bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)

		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendSuccess, 2)
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendFailed, 1)

		// Let's close the opened session
		td.sync.peerSet.CloseSession(bdl.Message.(*message.BlocksRequestMessage).SessionID)
	})

	t.Run("testing send failure", func(t *testing.T) {
		td.network.SendError = fmt.Errorf("send error")

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendSuccess, 2)
		// Since the test pid is the only peer in the peerSet list, it always tries to connect to it.
		// After the second attempt, the requested height is higher than that of the test peer.
		// So we have two sends, and therefore two failures.
		assert.Equal(t, td.sync.peerSet.GetPeer(pid).SendFailed, 3)
	})
}
