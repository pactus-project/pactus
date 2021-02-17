package sync

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/network_api"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/version"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	tTxPool           *txpool.MockTxPool
	tAliceConfig      *Config
	tBobConfig        *Config
	tAliceState       *state.MockState
	tBobState         *state.MockState
	tAliceConsensus   *consensus.MockConsensus
	tBobConsensus     *consensus.MockConsensus
	tAliceNetAPI      *network_api.MockNetworkAPI
	tBobNetAPI        *network_api.MockNetworkAPI
	tAliceSync        *synchronizer
	tBobSync          *synchronizer
	tAliceBroadcastCh chan *message.Message
	tBobBroadcastCh   chan *message.Message
	tAlicePeerID      peer.ID
	tBobPeerID        peer.ID
	tAnotherPeerID    peer.ID
)

type OverrideFingerprint struct {
	sync Synchronizer
	name string
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.sync.Fingerprint()
}

func init() {
	logger.InitLogger(logger.TestConfig())
	tAliceConfig = TestConfig()
	tBobConfig = TestConfig()

	tAliceConfig.Moniker = "alice"
	tBobConfig.Moniker = "bob"

	LatestBlockInterval = 20
	DownloadBlockInterval = 30
}

func setup(t *testing.T) {
	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	aliceSigner := crypto.NewSigner(priv1)
	bobSigner := crypto.NewSigner(priv2)

	tTxPool = txpool.MockingTxPool()

	tAlicePeerID = util.RandomPeerID()
	tBobPeerID = util.RandomPeerID()
	tAnotherPeerID = util.RandomPeerID()
	tAliceState = state.MockingState()
	tBobState = state.MockingState()
	tAliceConsensus = consensus.MockingConsensus()
	tBobConsensus = consensus.MockingConsensus()
	tAliceBroadcastCh = make(chan *message.Message, 100)
	tBobBroadcastCh = make(chan *message.Message, 100)
	tAliceNetAPI = network_api.MockingNetworkAPI(tAlicePeerID)
	tBobNetAPI = network_api.MockingNetworkAPI(tBobPeerID)
	aliceCache, _ := cache.NewCache(tAliceConfig.CacheSize, tAliceState.StoreReader())
	bobCache, _ := cache.NewCache(tBobConfig.CacheSize, tBobState.StoreReader())

	tBobState.GenHash = tAliceState.GenHash

	// Apply 20 blocks for both Alice and Bob
	lastBlockHash := crypto.Hash{}
	for i := 0; i < 21; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCommit(b.Hash())
		lastBlockHash = b.Hash()

		tAliceState.AddBlock(i+1, b, trxs)
		tAliceState.LastBlockCommit = c

		tBobState.AddBlock(i+1, b, trxs)
		tBobState.LastBlockCommit = c
	}

	tAliceSync = &synchronizer{
		ctx:         context.Background(),
		config:      tAliceConfig,
		signer:      aliceSigner,
		state:       tAliceState,
		consensus:   tAliceConsensus,
		cache:       aliceCache,
		txPool:      tTxPool,
		broadcastCh: tAliceBroadcastCh,
		networkAPI:  tAliceNetAPI,
	}
	tAliceSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "alice: ", sync: tAliceSync})
	tAliceSync.peerSet = peerset.NewPeerSet(tAliceConfig.SessionTimeout)
	tAliceSync.firewall = firewall.NewFirewall(tAliceSync.peerSet, tAliceState)
	tAliceSync.consensusSync = NewConsensusSync(tAliceConfig, tAlicePeerID, tAliceConsensus, tAliceSync.logger, tAliceSync.publishMessage)
	tAliceSync.stateSync = NewStateSync(tAliceConfig, tAlicePeerID, tAliceSync.cache, tAliceState, tTxPool, tAliceSync.peerSet, tAliceSync.logger, tAliceSync.publishMessage, tAliceSync.synced)

	tBobSync = &synchronizer{
		ctx:         context.Background(),
		config:      tBobConfig,
		signer:      bobSigner,
		state:       tBobState,
		consensus:   tBobConsensus,
		cache:       bobCache,
		txPool:      tTxPool,
		broadcastCh: tBobBroadcastCh,
		networkAPI:  tBobNetAPI,
	}
	tBobSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "bob: ", sync: tBobSync})
	tBobSync.peerSet = peerset.NewPeerSet(tBobConfig.SessionTimeout)
	tBobSync.firewall = firewall.NewFirewall(tBobSync.peerSet, tBobState)
	tBobSync.consensusSync = NewConsensusSync(tBobConfig, tBobPeerID, tBobConsensus, tBobSync.logger, tBobSync.publishMessage)
	tBobSync.stateSync = NewStateSync(tBobConfig, tBobPeerID, tBobSync.cache, tBobState, tTxPool, tBobSync.peerSet, tBobSync.logger, tBobSync.publishMessage, tBobSync.synced)

	tAliceNetAPI.ParsFn = tAliceSync.ParsMessage
	tAliceNetAPI.Firewall = tAliceSync.firewall
	tAliceNetAPI.OtherAPI = tBobNetAPI

	tBobNetAPI.ParsFn = tBobSync.ParsMessage
	tBobNetAPI.Firewall = tBobSync.firewall
	tBobNetAPI.OtherAPI = tAliceNetAPI

	assert.NoError(t, tAliceSync.Start())
	assert.NoError(t, tBobSync.Start())

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeSalam)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeSalam)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)

	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
}

func addMoreBlocksForBobAndSendBlockAnnounceMessage(t *testing.T, count int) {

	lastBlockHash := tBobState.LastBlockHash()
	for i := 0; i < count; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCommit(b.Hash())
		lastBlockHash = b.Hash()

		tBobState.AddBlock(tBobState.LastBlockHeight()+1, b, trxs)
		tBobState.LastBlockCommit = c
	}

	msg := message.NewBlockAnnounceMessage(
		tBobPeerID,
		tBobState.LastBlockHeight(),
		tBobState.Store.Blocks[len(tBobState.Store.Blocks)-1],
		tBobState.LastBlockCommit)

	tBobBroadcastCh <- msg
}
func disableHeartbeat(t *testing.T) {
	tAliceSync.heartBeatTicker.Stop()
	tBobSync.heartBeatTicker.Stop()
}

func joinAliceToTheSet(t *testing.T) {
	val := validator.NewValidator(tAliceSync.signer.PublicKey(), 4, tAliceState.LastBlockHeight())
	assert.NoError(t, tAliceState.ValSet.UpdateTheSet(0, []*validator.Validator{val}))
	assert.NoError(t, tBobState.ValSet.UpdateTheSet(0, []*validator.Validator{val}))
}

func joinBobToTheSet(t *testing.T) {
	val := validator.NewValidator(tBobSync.signer.PublicKey(), 5, tBobState.LastBlockHeight())
	assert.NoError(t, tAliceState.ValSet.UpdateTheSet(0, []*validator.Validator{val}))
	assert.NoError(t, tBobState.ValSet.UpdateTheSet(0, []*validator.Validator{val}))
}

func TestAccessors(t *testing.T) {
	setup(t)

	assert.Equal(t, tAliceSync.PeerID(), tAlicePeerID)
	assert.Equal(t, len(tAliceSync.Peers()), 1)
}

func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("bad-genesis", pub, tAnotherPeerID, invGenHash, 0, 0)
	tAliceNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)

	msg2 := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.ResponseCode, payload.ResponseCodeRejected)
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 3, 0x1)
	tAliceNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)

	msg2 := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.ResponseCode, payload.ResponseCodeOK)
	assert.Equal(t, tBobSync.peerSet.MaxClaimedHeight(), tAliceState.LastBlockHeight())

	p := tAliceSync.peerSet.GetPeer(tAnotherPeerID)
	assert.Equal(t, p.NodeVersion(), version.NodeVersion)
	assert.Equal(t, p.Moniker(), "kitty")
	assert.True(t, pub.EqualsTo(p.PublicKey()))
	assert.Equal(t, p.PeerID(), tAnotherPeerID)
	assert.Equal(t, p.Height(), 3)
	assert.Equal(t, p.InitialBlockDownload(), true)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop normally
	tAliceSync.Stop()
	tBobSync.Stop()
}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	claimedHeight := tAliceState.LastBlockHeight() + 5
	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, claimedHeight, 0)
	tAliceNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	msg := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty", pub, tAnotherPeerID, 1, 0)
	tAliceNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)

	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
}

func TestSendAleykPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	claimedHeight := tAliceState.LastBlockHeight() + 5
	msg := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty", pub, tAnotherPeerID, claimedHeight, 0)
	tAliceNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
}

func TestSendAleykPeerSameHeight(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	claimedHeight := tAliceState.LastBlockHeight()
	msg := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty", pub, tAnotherPeerID, claimedHeight, 0)
	tAliceNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)

	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
}

func TestIncreaseHeight(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	msg1 := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenesisHash(), 103, 0)
	tAliceSync.ParsMessage(msg1, tAnotherPeerID)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 103)

	msg2 := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty-2", pub, tAnotherPeerID, 104, 0)
	tAliceSync.ParsMessage(msg2, tAnotherPeerID)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 104)

	msg3 := message.NewHeartBeatMessage(tAnotherPeerID, crypto.GenerateTestHash(), hrs.NewHRS(106, 0, 1))
	tAliceSync.ParsMessage(msg3, tAnotherPeerID)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 105)
}

func TestQueryTransaction(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSendTx()

	// Alice has trx1 in his cache
	tAliceSync.cache.AddTransaction(trx1)
	tAliceSync.cache.AddTransaction(trx3)
	tBobSync.cache.AddTransaction(trx2)
	msg := message.NewQueryTransactionsMessage(tAlicePeerID, []crypto.Hash{trx2.ID(), trx3.ID()})

	t.Run("Alice should not send query transaction message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- msg
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)
		assert.True(t, tTxPool.HasTx(trx3.ID()))
	})

	t.Run("Bob should not process alice message because he is not an active validator", func(t *testing.T) {
		msg := msg
		tBobNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)
		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeTransactions)
	})

	joinBobToTheSet(t)

	t.Run("Bob should not process alice message because she is not an active validator", func(t *testing.T) {
		msg := msg
		tBobNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)
		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeTransactions)
	})

	joinAliceToTheSet(t)

	t.Run("Alice sends query transaction message", func(t *testing.T) {
		tAliceBroadcastCh <- msg
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)
	})

	t.Run("Alice sends query transaction message, but she has it in the cache", func(t *testing.T) {
		tAliceBroadcastCh <- message.NewQueryTransactionsMessage(tAlicePeerID, []crypto.Hash{trx1.ID()})
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)
	})

	t.Run("Bob processes alice message", func(t *testing.T) {
		msg := msg
		tBobNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeTransactions)
	})
}

func TestQueryProposal(t *testing.T) {
	setup(t)

	p1, _ := vote.GenerateTestProposal(106, 0)
	p2, _ := vote.GenerateTestProposal(106, 1)

	tAliceConsensus.HRS_ = hrs.NewHRS(106, 0, 1)
	tBobConsensus.HRS_ = hrs.NewHRS(106, 1, 1)

	tAliceSync.cache.AddProposal(p1)
	tBobConsensus.SetProposal(p2)
	msg := message.NewQueryProposalMessage(tAlicePeerID, 106, 1)

	t.Run("Alice should not send query proposal message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- msg
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob should not process alice message because he is not an active validator", func(t *testing.T) {
		tBobNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)
		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeTransactions)
	})

	joinBobToTheSet(t)

	t.Run("Bob should not process alice message because she is not an active validator", func(t *testing.T) {
		tBobNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)
		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeTransactions)
	})

	joinAliceToTheSet(t)

	t.Run("Alice sends query transaction message", func(t *testing.T) {
		tAliceBroadcastCh <- msg
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice sends query transaction message, but she has it in her cache", func(t *testing.T) {
		tAliceBroadcastCh <- message.NewQueryProposalMessage(tAlicePeerID, 106, 0)
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob processes alice message", func(t *testing.T) {
		tBobNetAPI.CheckAndParsMessage(msg, tAnotherPeerID)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})
}
func TestHeartbeatNotInSet(t *testing.T) {
	setup(t)

	// Alice is not in validator set
	tAliceSync.broadcastHeartBeat()
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)

	msg := message.NewHeartBeatMessage(tAnotherPeerID, crypto.GenerateTestHash(), hrs.NewHRS(106, 1, 1))
	tAliceSync.ParsMessage(msg, tAnotherPeerID)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 105)

	tAliceConsensus.HRS_ = hrs.NewHRS(106, 0, 1)
	tAliceSync.ParsMessage(msg, tAnotherPeerID)
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryVotes)
}

func TestBlockAnnounceMessage(t *testing.T) {
	setup(t)

	tAliceConsensus.Started = false

	t.Run("Bob should not broadcast block announce message because he is not an active validator", func(t *testing.T) {
		addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 1)

		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeBlockAnnounce)
	})

	joinBobToTheSet(t)

	t.Run("Bob should broadcast block announce message because he is an active validator", func(t *testing.T) {
		addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 1)

		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeBlockAnnounce)
		assert.True(t, tAliceConsensus.Started)
	})
}

func TestRequestForBlock(t *testing.T) {
	setup(t)

	joinBobToTheSet(t)

	t.Run("Bob claims that he has one more block", func(t *testing.T) {
		addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 1)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeBlockAnnounce)
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	})

	t.Run("Bob claims that he has two more blocks", func(t *testing.T) {
		addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 1)

		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeBlockAnnounce)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	})
}

func TestNotActiveValidator(t *testing.T) {
	setup(t)

	t.Run("Alice is not an active validator, She should not send query proposal message", func(t *testing.T) {
		tAliceSync.queryProposal(1, 1)
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice is not an active validator, She should not send query votes message", func(t *testing.T) {
		tAliceSync.queryVotes(1, 1)
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryVotes)
	})

	joinAliceToTheSet(t)

	t.Run("Alice is an active validator, She can send query proposal message", func(t *testing.T) {
		tAliceSync.queryProposal(1, 1)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice is an active validator, She can send query votes message", func(t *testing.T) {
		tAliceSync.queryVotes(1, 1)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryVotes)
	})
}
