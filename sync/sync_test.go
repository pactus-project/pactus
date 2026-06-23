package sync

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/ezex-io/gopkg/pipeline"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	config    *Config
	state     *state.FakeState
	consV1Mgr *manager.FakeConsensusManager
	consV2Mgr *manager.FakeConsensusManager
	network   *network.MockNetwork
	sync      *synchronizer
}

func testConfig() *Config {
	return &Config{
		Moniker:             "test",
		SessionTimeoutStr:   "1s",
		MaxSessions:         4,
		BlockPerMessage:     11,
		BlockPerSession:     27,
		PruneWindow:         13,
		Firewall:            firewall.DefaultConfig(),
		LatestSupportingVer: DefaultConfig().LatestSupportingVer,
		Services:            service.New(service.FullNode, service.PrunedNode),
	}
}

func setup(t *testing.T, config *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	if config == nil {
		config = testConfig()
		config.Moniker = "Alice"
	}
	valKeys := []*bls.ValidatorKey{ts.RandValKey(), ts.RandValKey()}

	consV1Mgr := manager.NewFakeConsensusManager(ts)
	consV2Mgr := manager.NewFakeConsensusManager(ts)

	state := state.NewFakeState(ts, nil)
	state.StateParams.BlockIntervalInSecond = 1
	state.CommitTestBlocks(100)

	consV1Mgr.EXPECT().IsDeprecated().Return(false).AnyTimes()
	consV1Mgr.EXPECT().MoveToNewHeight().Return().AnyTimes()
	consV1Mgr.EXPECT().HeightRound().DoAndReturn(func() (types.Height, types.Round) {
		return state.LastHeight + 1, ts.RandRound()
	}).AnyTimes()

	mockNetwork := network.MockingNetwork(ts, ts.RandPeerID())
	broadcastPipe := pipeline.New[message.Message](t.Context())

	syncInst, err := NewSynchronizer(t.Context(), config, valKeys,
		state, consV1Mgr, consV2Mgr, mockNetwork, broadcastPipe, mockNetwork.EventPipe)
	require.NoError(t, err)
	sync := syncInst.(*synchronizer)

	td := &testData{
		TestSuite: ts,
		config:    config,
		state:     state,
		consV1Mgr: consV1Mgr,
		consV2Mgr: consV2Mgr,
		network:   mockNetwork,
		sync:      sync,
	}

	require.NoError(t, td.sync.Start())
	assert.Equal(t, config.Moniker, td.sync.Moniker())
	assert.Equal(t, config.Services, td.sync.Services())

	logger.Info("setup finished, running the tests", "name", t.Name())

	return td
}

func shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) *bundle.Bundle {
	t.Helper()

	timer := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-timer.C:
			require.NoError(t, fmt.Errorf("shouldPublishMessageWithThisType %v: Timeout, test: %v", msgType, t.Name()))

			return nil

		case data := <-net.PublishCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(data.Data))
			require.NoError(t, err)

			// -----------
			// Check flags
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagCarrierLibP2P), "invalid flag: %v", bdl)
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)

			if data.Target == nil {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted), "invalid flag: %v", bdl)
			}

			if bdl.Message.Type() == message.TypeHello ||
				bdl.Message.Type() == message.TypeHelloAck {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHandshaking), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHandshaking), "invalid flag: %v", bdl)
			}
			// -----------

			require.Equal(t, msgType, bdl.Message.Type(), "not expected message: %s", msgType)

			return bdl
		}
	}
}

func (td *testData) shouldPublishMessageWithThisType(t *testing.T, msgType message.Type,
) *bundle.Bundle {
	t.Helper()

	return shouldPublishMessageWithThisType(t, td.network, msgType)
}

func shouldNotPublishAnyMessage(t *testing.T, net *network.MockNetwork) {
	t.Helper()

	timer := time.NewTimer(3 * time.Millisecond)

	for {
		select {
		case <-timer.C:
			return

		case data := <-net.PublishCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(data.Data))
			require.NoError(t, err)
			require.Fail(t, "published unexpected message: "+bdl.Message.Type().String())
		}
	}
}

func (td *testData) shouldNotPublishAnyMessage(t *testing.T) {
	t.Helper()

	shouldNotPublishAnyMessage(t, td.network)
}

func (*testData) receivingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) {
	bdl := bundle.NewBundle(msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)

	sync.processIncomingBundle(bdl, from)
}

func (td *testData) addPeer(t *testing.T, status status.Status, services service.Services) peer.ID {
	t.Helper()

	pid := td.RandPeerID()
	pub, _ := td.RandBLSKeyPair()

	td.sync.peerSet.UpdateInfo(pid, t.Name(),
		version.NodeAgent.String(), []*bls.PublicKey{pub}, services)
	td.sync.peerSet.UpdateStatus(pid, status)

	return pid
}

func (td *testData) checkPeerStatus(t *testing.T, pid peer.ID, code status.Status) {
	t.Helper()

	require.Equal(t, code, td.sync.peerSet.GetPeerStatus(pid))
}

func TestStop(t *testing.T) {
	td := setup(t, nil)

	// Should stop gracefully.
	td.sync.Stop()
}

func TestConnectEvent(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	remoteAddr := "/ip4/2.2.2.2/tcp/21888"
	ce := &network.ConnectEvent{
		PeerID:        pid,
		RemoteAddress: remoteAddr,
		Direction:     lp2pnetwork.DirInbound,
	}
	td.network.EventPipe.Send(ce)

	assert.Eventually(t, func() bool {
		return td.sync.peerSet.HasPeer(pid)
	}, time.Second, 100*time.Millisecond)

	peer := td.sync.peerSet.GetPeer(pid)
	assert.Equal(t, status.StatusConnected, peer.Status)
	assert.Equal(t, remoteAddr, peer.Address)
	assert.Equal(t, lp2pnetwork.DirInbound, peer.Direction)
}

func TestDisconnectEvent(t *testing.T) {
	td := setup(t, nil)
	pid := td.RandPeerID()
	de := &network.DisconnectEvent{
		PeerID: pid,
	}
	td.network.EventPipe.Send(de)

	assert.Eventually(t, func() bool {
		return td.sync.peerSet.HasPeer(pid)
	}, time.Second, 100*time.Millisecond)

	td.checkPeerStatus(t, pid, status.StatusDisconnected)
}

func TestProtocolsEvent(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	pe := &network.ProtocolsEvents{
		PeerID:    pid,
		Protocols: []string{"protocol-1", "protocol-2"},
	}
	td.network.EventPipe.Send(pe)

	assert.Eventually(t, func() bool {
		return td.sync.peerSet.HasPeer(pid)
	}, time.Second, 100*time.Millisecond)

	peer := td.sync.peerSet.GetPeer(pid)
	assert.Equal(t, []string{"protocol-1", "protocol-2"}, peer.Protocols)
}

func TestSendHello(t *testing.T) {
	td := setup(t, nil)

	t.Run("Peer with unknown Direction", func(t *testing.T) {
		pid := td.RandPeerID()
		pe := &network.ProtocolsEvents{
			PeerID:    pid,
			Protocols: []string{"protocol-1"},
		}
		td.network.EventPipe.Send(pe)

		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("Peer with inbound Direction", func(t *testing.T) {
		pid := td.RandPeerID()
		td.sync.peerSet.UpdateAddress(pid, "test-address", lp2pnetwork.DirInbound)

		pe := &network.ProtocolsEvents{
			PeerID:    pid,
			Protocols: []string{"protocol-1"},
		}
		td.network.EventPipe.Send(pe)

		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("Peer with outbound Direction", func(t *testing.T) {
		pid := td.RandPeerID()
		td.sync.peerSet.UpdateAddress(pid, "test-address", lp2pnetwork.DirOutbound)

		pe := &network.ProtocolsEvents{
			PeerID:    pid,
			Protocols: []string{"protocol-1"},
		}
		td.network.EventPipe.Send(pe)

		td.shouldPublishMessageWithThisType(t, message.TypeHello)
	})
}

func TestTestNetFlags(t *testing.T) {
	td := setup(t, nil)

	td.state.GenDoc = genesis.TestnetGenesis()

	bdl := td.sync.prepareBundle(message.NewQueryProposalMessage(
		td.RandHeight(), td.RandRound(), td.RandValAddress(),
	))

	require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)
	require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)
}

func TestDownload(t *testing.T) {
	conf := testConfig()
	// Let's not allow `GetRandomPeer` to disappoint us!
	conf.MaxSessions = 32

	td := setup(t, conf)

	td.consV1Mgr.EXPECT().MoveToNewHeight().Return().AnyTimes()
	td.state.EXPECT().Genesis().Return(genesis.MainnetGenesis()).AnyTimes()
	td.state.EXPECT().Params().Return(param.FromGenesis(genesis.MainnetGenesis())).AnyTimes()

	t.Run("try to download blocks, but the peer is not known", func(t *testing.T) {
		pid := td.addPeer(t, status.StatusConnected, service.New(service.None))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert, nil)
		td.receivingNewMessage(td.sync, baMsg, pid)

		td.shouldNotPublishAnyMessage(t)
		td.network.IsClosed(pid)
	})

	t.Run("try to download blocks, but the peer is not a network node", func(t *testing.T) {
		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert, nil)
		td.receivingNewMessage(td.sync, baMsg, pid)

		td.shouldNotPublishAnyMessage(t)
		td.network.IsClosed(pid)
	})

	t.Run("try to download blocks and the peer is a network node", func(t *testing.T) {
		pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert, nil)
		td.receivingNewMessage(td.sync, baMsg, pid)

		td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
	})

	t.Run("download request is rejected", func(t *testing.T) {
		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))
		from := td.sync.stateHeight() + 1
		count := uint32(123)
		sid := td.sync.peerSet.OpenSession(pid, from, count)
		msg := message.NewBlocksResponseMessage(message.ResponseCodeRejected, t.Name(),
			sid, 1, nil, nil)
		td.receivingNewMessage(td.sync, msg, pid)

		assert.False(t, td.sync.peerSet.HasOpenSession(pid))
	})
}

func TestBroadcastBlockAnnounce(t *testing.T) {
	td := setup(t, nil)

	t.Run("Should announce the block", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		msg := message.NewBlockAnnounceMessage(blk, cert, nil)

		td.sync.broadcast(msg)

		td.shouldPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	})

	t.Run("Should NOT announce the block", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		msg := message.NewBlockAnnounceMessage(blk, cert, nil)

		td.sync.cache.AddBlock(blk)
		td.sync.broadcast(msg)

		td.shouldNotPublishAnyMessage(t)
	})
}

func TestAllBlocksInCache(t *testing.T) {
	td := setup(t, nil)

	blk100, _ := td.GenerateTestBlock(100)
	blk101, _ := td.GenerateTestBlock(101)
	blk102, _ := td.GenerateTestBlock(102)

	td.sync.cache.AddBlock(blk100)
	td.sync.cache.AddBlock(blk101)
	td.sync.cache.AddBlock(blk102)

	res := td.sync.sendBlockRequestToRandomPeer(100, 3, true)
	assert.True(t, res)
}

func TestCommitMissedBlock(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	lastHeight := td.state.LastBlockHeight()

	blk1, cert1 := td.GenerateTestBlock(lastHeight + 1)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1, nil)

	blk2, cert2 := td.GenerateTestBlock(lastHeight + 2)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2, nil)

	t.Run("Receiving block announce message, without committing previous block", func(t *testing.T) {
		td.receivingNewMessage(td.sync, msg2, pid)

		consHeight, _ := td.sync.getConsMgr().HeightRound()
		assert.Equal(t, lastHeight+1, consHeight)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		td.receivingNewMessage(td.sync, msg1, pid)

		newHeight := td.state.LastBlockHeight()
		assert.Equal(t, lastHeight+2, newHeight)
	})
}

func TestBlockAnnounceCacheBadAnnounce(t *testing.T) {
	td := setup(t, nil)

	lastHeight := td.state.LastBlockHeight()

	// Create a block announce where the certificate height does not match
	// the block height (simulating an invalid/bad announce).
	blk, _ := td.GenerateTestBlock(lastHeight + 1)
	badCert := td.GenerateTestCertificate(lastHeight + 2) // wrong height

	badMsg := message.NewBlockAnnounceMessage(blk, badCert, nil)
	pid := td.RandPeerID()
	td.receivingNewMessage(td.sync, badMsg, pid)

	// The block is added at height lastHeight+1 (via AddBlock),
	// and the cert is added at height lastHeight+2 (via AddCertificate).
	// tryCommitBlocks checks for block+cert at lastHeight+1:
	// block found but cert NOT found at that height, so commit can't proceed.
	// The block stays in cache (orphaned).
	assert.True(t, td.sync.cache.HasBlockInCache(lastHeight+1))
	assert.NotNil(t, td.sync.cache.GetBlock(lastHeight+1))

	// State should be unchanged since nothing was committed.
	assert.Equal(t, lastHeight, td.state.LastBlockHeight())

	// Removing the orphaned block from cache should work,
	// simulating the path tryCommitBlocks takes on error.
	td.sync.cache.RemoveBlock(lastHeight + 1)
	assert.False(t, td.sync.cache.HasBlockInCache(lastHeight+1))
}

func TestBlockAnnounceCacheCommittedBefore(t *testing.T) {
	td := setup(t, nil)

	// Commit one more block so the node is ahead.
	td.state.CommitTestBlocks(1)
	lastHeight := td.state.LastBlockHeight()

	// Now receive a block announce for an already-committed height.
	blk, cert := td.GenerateTestBlock(lastHeight)
	msg := message.NewBlockAnnounceMessage(blk, cert, nil)
	pid := td.RandPeerID()
	td.receivingNewMessage(td.sync, msg, pid)

	// The block gets added to cache, but tryCommitBlocks starts
	// at lastHeight+1, so the already-committed block is ignored.
	// State should remain unchanged.
	assert.Equal(t, lastHeight, td.state.LastBlockHeight())

	// Block is in cache but doesn't interfere with state.
	assert.True(t, td.sync.cache.HasBlockInCache(lastHeight))
}

func TestBlockAnnouncementEntropyDelay(t *testing.T) {
	td := setup(t, nil)

	td.state.StateParams.BlockIntervalInSecond = 10

	t.Run("recent block returns delay within block interval", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight(),
			testsuite.BlockWithTime(time.Now()))
		msg := message.NewBlockAnnounceMessage(blk, cert, nil)

		delay := td.sync.blockAnnouncementEntropyDelay(msg)
		assert.GreaterOrEqual(t, delay, 0)
		assert.LessOrEqual(t, delay, 10)
	})

	t.Run("future block returns zero delay", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight(),
			testsuite.BlockWithTime(time.Now().Add(time.Minute)))
		msg := message.NewBlockAnnounceMessage(blk, cert, nil)

		delay := td.sync.blockAnnouncementEntropyDelay(msg)
		assert.Zero(t, delay)
	})

	t.Run("stale block returns zero delay", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight(),
			testsuite.BlockWithTime(time.Now().Add(-time.Minute)))
		msg := message.NewBlockAnnounceMessage(blk, cert, nil)

		delay := td.sync.blockAnnouncementEntropyDelay(msg)
		assert.Zero(t, delay)
	})
}
