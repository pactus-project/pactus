package sync

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
	"time"

	beevikntp "github.com/beevik/ntp"
	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/ntp"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errTest = errors.New("test error")

// topicNetwork records the topics that the synchronizer joins, and it can refuse one of them.
type topicNetwork struct {
	*network.MockNetwork

	failTopic  network.TopicID
	evaluators map[network.TopicID]network.PropagationEvaluator
}

func newTopicNetwork(mock *network.MockNetwork, failTopic network.TopicID) *topicNetwork {
	return &topicNetwork{
		MockNetwork: mock,
		failTopic:   failTopic,
		evaluators:  make(map[network.TopicID]network.PropagationEvaluator),
	}
}

func (n *topicNetwork) JoinTopic(topicID network.TopicID, eval network.PropagationEvaluator) error {
	if topicID == n.failTopic {
		return errTest
	}

	n.evaluators[topicID] = eval

	return nil
}

// fakeQuerier returns the same NTP response for any pool.
type fakeQuerier struct {
	res *beevikntp.Response
	err error
}

func (q *fakeQuerier) Query(string) (*beevikntp.Response, error) {
	return q.res, q.err
}

// closeTracker reports if a reader is closed.
type closeTracker struct {
	io.Reader

	closed bool
}

func (c *closeTracker) Close() error {
	c.closed = true

	return nil
}

func gossipWithConsensusHeight(from peer.ID, height uint32) *network.GossipMessage {
	// The consensus height is the last 4 bytes of the bundle.
	data := make([]byte, 8)
	binary.BigEndian.PutUint32(data[4:], height)

	return &network.GossipMessage{From: from, Data: data}
}

func encodeBundle(t *testing.T, msg message.Message) []byte {
	t.Helper()

	bdl := bundle.NewBundle(msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)
	data, err := bdl.Encode()
	require.NoError(t, err)

	return data
}

// registerConnection registers the peer in the mocked network,
// so that it is possible to know if the connection is closed.
func (td *testData) registerConnection(pid peer.ID) {
	td.network.OtherNets[pid] = network.MockingNetwork(td.TestSuite, pid)
}

func TestStartFailures(t *testing.T) {
	topics := []network.TopicID{
		network.TopicIDBlock,
		network.TopicIDTransaction,
		network.TopicIDConsensus,
	}

	for _, topic := range topics {
		t.Run(topic.String(), func(t *testing.T) {
			td := setup(t, nil)
			td.sync.network = newTopicNetwork(td.network, topic)

			require.ErrorIs(t, td.sync.Start(), errTest)
		})
	}
}

func TestTopicEvaluators(t *testing.T) {
	td := setup(t, nil)
	topicNet := newTopicNetwork(td.network, network.TopicIDUnspecified)
	td.sync.network = topicNet
	require.NoError(t, td.sync.Start())
	require.Len(t, topicNet.evaluators, 3)

	from := td.RandPeerID()
	height := uint32(td.state.LastBlockHeight())

	t.Run("block topic", func(t *testing.T) {
		eval := topicNet.evaluators[network.TopicIDBlock]

		assert.Equal(t, network.Drop, eval(gossipWithConsensusHeight(from, height-10)), "expired message")
		assert.Equal(t, network.Propagate, eval(gossipWithConsensusHeight(from, height)))
		assert.Equal(t, network.DropButConsume, eval(gossipWithConsensusHeight(from, height)), "rate limit")
	})

	t.Run("transaction topic", func(t *testing.T) {
		eval := topicNet.evaluators[network.TopicIDTransaction]

		for i := 0; i < td.config.Firewall.RateLimit.TransactionTopic; i++ {
			assert.Equal(t, network.Propagate, eval(gossipWithConsensusHeight(from, 0)))
		}
		assert.Equal(t, network.DropButConsume, eval(gossipWithConsensusHeight(from, 0)), "rate limit")
	})

	t.Run("consensus topic", func(t *testing.T) {
		eval := topicNet.evaluators[network.TopicIDConsensus]

		assert.Equal(t, network.Drop, eval(gossipWithConsensusHeight(from, height-10)), "expired message")
		assert.Equal(t, network.Propagate, eval(gossipWithConsensusHeight(from, height)))
		assert.Equal(t, network.Propagate, eval(gossipWithConsensusHeight(from, height+1)))
	})
}

func TestNewSynchronizerFailures(t *testing.T) {
	t.Run("invalid banned network in the firewall", func(t *testing.T) {
		conf := testConfig()
		conf.Firewall.BannedNets = []string{"not-a-network"}

		ts := testsuite.NewTestSuite(t)
		_, err := NewSynchronizer(t.Context(), conf, nil, nil, nil, nil,
			network.MockingNetwork(ts, ts.RandPeerID()), nil, nil)
		require.Error(t, err)
	})
}

func TestBroadcastMessageByPipeline(t *testing.T) {
	td := setup(t, nil)

	msg := message.NewQueryProposalMessage(td.RandHeight(), td.RandRound(), td.RandValAddress())
	td.sync.broadcastPipe.Send(msg)

	bdl := td.shouldPublishMessageWithThisType(t, message.TypeQueryProposal)
	assert.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted))
}

func TestProcessNetworkMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("garbage gossip message from an unknown peer", func(t *testing.T) {
		pid := td.RandPeerID()

		require.NotPanics(t, func() {
			td.sync.processGossipMessage(&network.GossipMessage{From: pid, Data: []byte{1, 2, 3}})
		})

		require.True(t, td.sync.peerSet.HasPeer(pid))
		assert.Equal(t, int64(1), td.sync.peerSet.GetPeer(pid).Metric.TotalInvalid.Bundles)
	})

	t.Run("garbage stream message from an unknown peer", func(t *testing.T) {
		pid := td.RandPeerID()
		reader := &closeTracker{Reader: bytes.NewReader([]byte{1, 2, 3})}

		require.NotPanics(t, func() {
			td.sync.processStreamMessage(&network.StreamMessage{From: pid, Reader: reader})
		})

		assert.True(t, reader.closed, "the stream should always be closed")
		assert.Equal(t, int64(1), td.sync.peerSet.GetPeer(pid).Metric.TotalInvalid.Bundles)
	})

	t.Run("stream message received as gossip, the connection is closed", func(t *testing.T) {
		pid := td.RandPeerID()
		td.registerConnection(pid)
		data := encodeBundle(t, message.NewBlocksRequestMessage(1, 5, 10))

		td.sync.processGossipMessage(&network.GossipMessage{From: pid, Data: data})

		assert.True(t, td.network.IsClosed(pid))
		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("gossip message received as stream, the connection is closed", func(t *testing.T) {
		pid := td.RandPeerID()
		td.registerConnection(pid)
		data := encodeBundle(t, message.NewTransactionsMessage([]*tx.Tx{td.GenerateTestTransferTx()}))
		reader := &closeTracker{Reader: bytes.NewReader(data)}

		td.sync.processStreamMessage(&network.StreamMessage{From: pid, Reader: reader})

		assert.True(t, td.network.IsClosed(pid))
		assert.True(t, reader.closed)
	})

	t.Run("valid gossip message is handled", func(t *testing.T) {
		pid := td.RandPeerID()
		td.registerConnection(pid)
		trx := td.GenerateTestTransferTx()
		data := encodeBundle(t, message.NewTransactionsMessage([]*tx.Tx{trx}))

		// The transaction is decoded from the bundle, so it is not the same instance.
		td.state.EXPECT().AddPendingTx(gomock.Any()).DoAndReturn(func(got *tx.Tx) error {
			assert.Equal(t, trx.ID(), got.ID())

			return nil
		}).Times(1)
		td.sync.processGossipMessage(&network.GossipMessage{From: pid, Data: data})

		assert.False(t, td.network.IsClosed(pid))
	})

	t.Run("valid stream message is handled", func(t *testing.T) {
		// A peer that has not done the handshake can't request blocks.
		pid := td.RandPeerID()
		data := encodeBundle(t, message.NewBlocksRequestMessage(1, 5, 10))
		reader := &closeTracker{Reader: bytes.NewReader(data)}

		td.sync.processStreamMessage(&network.StreamMessage{From: pid, Reader: reader})

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
		resp, ok := bdl.Message.(*message.BlocksResponseMessage)
		require.True(t, ok)
		assert.True(t, resp.IsRequestRejected())
		assert.True(t, reader.closed)
	})

	t.Run("bundle with an unknown handler is ignored", func(t *testing.T) {
		delete(td.sync.handlers, message.TypeQueryVote)
		msg := message.NewQueryVoteMessage(td.RandHeight(), td.RandRound(), td.RandValAddress())

		require.NotPanics(t, func() {
			td.receivingNewMessage(td.sync, msg, td.RandPeerID())
		})
	})
}

func TestTransactionsHandlerContinuesAfterInvalidTransaction(t *testing.T) {
	td := setup(t, nil)

	invalidTrx := td.GenerateTestTransferTx()
	validTrx := td.GenerateTestTransferTx()
	msg := message.NewTransactionsMessage([]*tx.Tx{invalidTrx, validTrx})

	td.state.EXPECT().AddPendingTx(invalidTrx).Return(errTest).Times(1)
	td.state.EXPECT().AddPendingTx(validTrx).Return(nil).Times(1)

	td.receivingNewMessage(td.sync, msg, td.RandPeerID())
}

func TestClockOffset(t *testing.T) {
	newChecker := func(td *testData, querier *fakeQuerier) {
		td.sync.ntp.Stop()
		td.sync.ntp = ntp.NewNtpChecker(t.Context(),
			ntp.WithQuerier(querier),
			ntp.WithInterval(10*time.Millisecond),
			ntp.WithThreshold(time.Second))
		td.sync.ntp.Start()
	}

	t.Run("clock is in sync", func(t *testing.T) {
		td := setup(t, nil)
		newChecker(td, &fakeQuerier{res: &beevikntp.Response{Stratum: 1, ClockOffset: 100 * time.Millisecond}})

		require.Eventually(t, func() bool {
			offset, _ := td.sync.ClockOffset()

			return offset == 100*time.Millisecond
		}, 5*time.Second, 20*time.Millisecond)
		assert.False(t, td.sync.IsClockOutOfSync())
	})

	t.Run("clock is out of sync", func(t *testing.T) {
		td := setup(t, nil)
		newChecker(td, &fakeQuerier{res: &beevikntp.Response{Stratum: 1, ClockOffset: 3 * time.Second}})

		require.Eventually(t, func() bool {
			offset, _ := td.sync.ClockOffset()

			return offset == 3*time.Second
		}, 5*time.Second, 20*time.Millisecond)
		assert.True(t, td.sync.IsClockOutOfSync())
	})

	t.Run("unable to query", func(t *testing.T) {
		td := setup(t, nil)
		newChecker(td, &fakeQuerier{err: errTest})

		require.Eventually(t, func() bool {
			_, err := td.sync.ClockOffset()

			return err != nil
		}, 5*time.Second, 20*time.Millisecond)
		assert.True(t, td.sync.IsClockOutOfSync())
	})
}

func TestConsensusManagerSelection(t *testing.T) {
	t.Run("version 1 is active", func(t *testing.T) {
		td := setup(t, nil)

		assert.Same(t, td.consV1Mgr, td.sync.getConsMgr())
	})

	t.Run("version 1 is deprecated", func(t *testing.T) {
		td := setup(t, nil)

		deprecated := manager.NewFakeConsensusManager(td.TestSuite)
		deprecated.EXPECT().IsDeprecated().Return(true).AnyTimes()
		td.sync.consV1Mgr = deprecated

		assert.Same(t, td.consV2Mgr, td.sync.getConsMgr())

		// The consensus is moved to a new height by the active manager.
		td.consV2Mgr.EXPECT().MoveToNewHeight().Times(1)
		td.sync.moveConsensusToNewHeight()
	})
}

func TestSayHello(t *testing.T) {
	td := setup(t, nil)

	t.Run("handshake is already done", func(t *testing.T) {
		pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

		td.sync.sayHello(pid)

		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("handshake is not done yet", func(t *testing.T) {
		pid := td.addPeer(t, status.StatusConnected, service.New(service.FullNode))

		td.sync.sayHello(pid)

		td.shouldPublishMessageWithThisType(t, message.TypeHello)
	})
}

func TestUpdateBlockchain(t *testing.T) {
	conf := testConfig()
	// Let's not allow `GetRandomPeer` to disappoint us!
	conf.MaxSessions = 32
	conf.SessionTimeoutStr = "10s"

	t.Run("uncompleted session is requested again", func(t *testing.T) {
		td := setup(t, conf)
		pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

		sid := td.sync.peerSet.OpenSession(pid, 101, 10)
		td.sync.peerSet.SetSessionUncompleted(sid)

		td.sync.updateBlockchain()

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
		msg := bdl.Message.(*message.BlocksRequestMessage)
		assert.Equal(t, uint32(101), uint32(msg.From))
		assert.Equal(t, uint32(10), msg.Count)
	})

	t.Run("open sessions are not duplicated", func(t *testing.T) {
		td := setup(t, conf)
		pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
		td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

		td.sync.peerSet.OpenSession(pid, 101, 10)

		td.sync.updateBlockchain()

		td.shouldNotPublishAnyMessage(t)
		assert.Equal(t, 1, td.sync.peerSet.NumberOfSessions())
	})

	t.Run("block of the next height is in the cache", func(t *testing.T) {
		td := setup(t, conf)
		td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

		nextHeight := td.state.LastBlockHeight() + 1
		blk, _ := td.GenerateTestBlock(nextHeight)
		td.sync.cache.AddBlock(blk)

		td.sync.updateBlockchain()

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
		msg := bdl.Message.(*message.BlocksRequestMessage)
		assert.Equal(t, uint32(nextHeight+1), uint32(msg.From), "the cached block should not be downloaded again")
	})

	t.Run("chain is far behind, only full nodes are asked", func(t *testing.T) {
		td := setup(t, conf)
		pruned := td.addPeer(t, status.StatusKnown, service.New(service.PrunedNode))
		td.registerConnection(pruned)

		td.sync.updateBlockchain()

		td.shouldNotPublishAnyMessage(t)
		assert.True(t, td.network.IsClosed(pruned), "a peer without all the blocks is disconnected")
	})
}

func TestTryCommitBlocksFailures(t *testing.T) {
	t.Run("block can't be committed, it is removed from the cache", func(t *testing.T) {
		td := setup(t, nil)
		td.state.ErrCommit = errTest

		height := td.state.LastBlockHeight() + 1
		blk, cert := td.GenerateTestBlock(height)
		td.sync.cache.AddBlock(blk)
		td.sync.cache.AddCertificate(cert)
		require.True(t, td.sync.cache.HasBlockInCache(height))

		td.sync.tryCommitBlocks()

		assert.False(t, td.sync.cache.HasBlockInCache(height))
		assert.Nil(t, td.sync.cache.GetCertificate(height))
	})

	t.Run("public key of the signer is unknown, block is removed from the cache", func(t *testing.T) {
		td := setup(t, nil)

		trx := td.GenerateTestTransferTx()
		trx.StripPublicKey()
		height := td.state.LastBlockHeight() + 1
		blk, cert := td.GenerateTestBlock(height, testsuite.BlockWithTransactions(block.Txs{trx}))
		td.sync.cache.AddBlock(blk)
		td.sync.cache.AddCertificate(cert)

		td.state.EXPECT().PublicKey(trx.Payload().Signer()).Return(nil, errTest).Times(1)
		td.sync.tryCommitBlocks()

		assert.False(t, td.sync.cache.HasBlockInCache(height))
		assert.Equal(t, height-1, td.state.LastBlockHeight(), "the block should not be committed")
	})

	t.Run("block without certificate stays in the cache", func(t *testing.T) {
		td := setup(t, nil)

		height := td.state.LastBlockHeight() + 1
		blk, _ := td.GenerateTestBlock(height)
		td.sync.cache.AddBlock(blk)

		td.sync.tryCommitBlocks()

		assert.True(t, td.sync.cache.HasBlockInCache(height))
		assert.Equal(t, height-1, td.state.LastBlockHeight())
	})
}

func TestBlocksResponseSessionStatus(t *testing.T) {
	td := setup(t, nil)

	pid1 := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	pid2 := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

	openSession := func(pid peer.ID) int {
		return td.sync.peerSet.OpenSession(pid, td.state.LastBlockHeight()+1, 10)
	}

	t.Run("response code OK keeps the session open and refreshes its activity", func(t *testing.T) {
		td.sync.peerSet.RemoveAllSessions()
		sid := openSession(pid1)
		for _, ssn := range td.sync.peerSet.Sessions() {
			ssn.LastActivity = time.Now().Add(-time.Hour)
		}

		msg := message.NewBlocksResponseMessage(message.ResponseCodeOK, "ok", sid, 0, nil, nil)
		td.receivingNewMessage(td.sync, msg, pid1)

		require.True(t, td.sync.peerSet.HasOpenSession(pid1))
		for _, ssn := range td.sync.peerSet.Sessions() {
			assert.WithinDuration(t, time.Now(), ssn.LastActivity, time.Minute)
		}
	})

	t.Run("no more blocks completes the session, the other sessions continue", func(t *testing.T) {
		td.sync.peerSet.RemoveAllSessions()
		sid1 := openSession(pid1)
		openSession(pid2)

		msg := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks, "done", sid1, 0, nil, nil)
		td.receivingNewMessage(td.sync, msg, pid1)

		stats := td.sync.peerSet.SessionStats()
		assert.Equal(t, 1, stats.Completed)
		assert.Equal(t, 1, stats.Open)
		assert.True(t, td.sync.peerSet.HasOpenSession(pid2))
	})

	t.Run("no more blocks with no other open session starts a new download", func(t *testing.T) {
		td.sync.peerSet.RemoveAllSessions()
		sid := openSession(pid1)

		msg := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks, "done", sid, 0, nil, nil)
		td.receivingNewMessage(td.sync, msg, pid1)

		// The completed session is removed, and since the node is still behind the network,
		// new sessions are opened.
		stats := td.sync.peerSet.SessionStats()
		assert.Zero(t, stats.Completed, "completed sessions are removed")
		assert.Positive(t, stats.Open)
		td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
	})

	t.Run("synced completes the session and moves the consensus to a new height", func(t *testing.T) {
		td.sync.peerSet.RemoveAllSessions()
		sid := openSession(pid1)

		deprecated := manager.NewFakeConsensusManager(td.TestSuite)
		deprecated.EXPECT().IsDeprecated().Return(true).AnyTimes()
		td.sync.consV1Mgr = deprecated
		td.consV2Mgr.EXPECT().MoveToNewHeight().Times(1)

		msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, "synced", sid, 0, nil, nil)
		td.receivingNewMessage(td.sync, msg, pid1)

		assert.Equal(t, 1, td.sync.peerSet.SessionStats().Completed)
	})
}

func TestNewSynchronizerWithInvalidCacheSize(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	conf := testConfig()
	conf.BlockPerSession = 0

	_, err := NewSynchronizer(t.Context(), conf, nil, state.NewFakeState(ts), nil, nil,
		network.MockingNetwork(ts, ts.RandPeerID()), pipeline.New[message.Message](t.Context()),
		pipeline.New[network.Event](t.Context()))

	require.Error(t, err)
}

func TestLocalnetBundleFlags(t *testing.T) {
	td := setup(t, nil)

	td.state.GenDoc = genesis.MakeGenesis(time.Now(), nil, nil, genesis.DefaultGenesisParams())
	require.Equal(t, genesis.Localnet, td.state.GenDoc.ChainType())

	bdl := td.sync.prepareBundle(message.NewQueryProposalMessage(
		td.RandHeight(), td.RandRound(), td.RandValAddress(),
	))

	assert.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet))
	assert.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet))
	assert.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagCarrierLibP2P))
}

func TestUpdateBlockchainWithoutPeers(t *testing.T) {
	td := setup(t, nil)

	sid := td.sync.peerSet.OpenSession(td.RandPeerID(), 101, 10)
	td.sync.peerSet.SetSessionUncompleted(sid)

	require.NotPanics(t, td.sync.updateBlockchain)

	td.shouldNotPublishAnyMessage(t)
}

func TestUpdateBlockchainLittleBehind(t *testing.T) {
	conf := testConfig()
	conf.MaxSessions = 32
	td := setup(t, conf)

	// The node is a few blocks behind the network, so it can download blocks from any node.
	prunedNode := td.addPeer(t, status.StatusKnown, service.New(service.PrunedNode))
	td.registerConnection(prunedNode)

	blockInterval := td.state.Params().BlockInterval()
	blkTime := util.RoundNow(int(blockInterval.Seconds())).Add(-5 * blockInterval)
	blk, cert := td.GenerateTestBlock(td.state.LastBlockHeight()+1, testsuite.BlockWithTime(blkTime))
	require.NoError(t, td.state.CommitBlock(blk, cert))
	require.Less(t, uint32(5), conf.PruneWindow)

	td.sync.updateBlockchain()

	bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
	msg := bdl.Message.(*message.BlocksRequestMessage)
	assert.Equal(t, uint32(td.state.LastBlockHeight()+1), uint32(msg.From))
	assert.False(t, td.network.IsClosed(prunedNode), "the connection should stay open")
}
