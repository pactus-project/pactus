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
	"github.com/zarbchain/zarb-go/version"
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
	tAliceSync        *Synchronizer
	tBobSync          *Synchronizer
	tAliceBroadcastCh chan *message.Message
	tBobBroadcastCh   chan *message.Message
	tAlicePeerID      peer.ID
	tBobPeerID        peer.ID
	tAnotherPeerID    peer.ID
)

type OverrideFingerprint struct {
	sync *Synchronizer
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

	tAlicePeerID, _ = peer.IDB58Decode("12D3KooWLQ8GKaLdKU8Ms6AkMYjDWCr5UTPvdewag3tcarxh7saC")
	tBobPeerID, _ = peer.IDB58Decode("12D3KooWHyepEGGdeSk3nPZrEamxLNba7tFZJKWbyEdZ654fHJdk")
	tAnotherPeerID, _ = peer.IDB58Decode("12D3KooWM4dZKiZ8y21biCZXuAJYD5db8vSr1hfMpSgBSqpekY4Q")
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

	// Alice has 100 and Bob has 92 blocks
	lastBlockHash := crypto.Hash{}
	for i := 0; i < 100; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCommit(b.Hash())
		lastBlockHash = b.Hash()
		tAliceState.AddBlock(i+1, b, trxs)
		tAliceState.LastBlockCommit = c

		if i < 82 {
			tBobState.AddBlock(i+1, b, trxs)
			tBobState.LastBlockCommit = c
		}
	}

	tAliceSync = &Synchronizer{
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
	tAliceSync.consensusSync = NewConsensusSync(tAliceConfig, tAliceConsensus, tAliceSync.logger, tAliceSync.publishMessage)
	tAliceSync.stateSync = NewStateSync(tAliceConfig, tAlicePeerID, tAliceSync.cache, tAliceState, tTxPool, tAliceSync.peerSet, tAliceSync.logger, tAliceSync.publishMessage, tAliceSync.synced)

	tBobSync = &Synchronizer{
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
	tBobSync.consensusSync = NewConsensusSync(tBobConfig, tBobConsensus, tBobSync.logger, tBobSync.publishMessage)
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

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // blocks 83-92
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // blocks 93-100
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // last commit + sync response

	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
}

func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("bad-genesis", pub, tAnotherPeerID, invGenHash, 0, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	msg2 := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.ResponseCode, payload.ResponseCodeRejected)
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 3, 0x1)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	msg2 := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.ResponseCode, payload.ResponseCodeOK)
	assert.Equal(t, tBobSync.peerSet.MaxClaimedHeight(), tAliceState.LastBlockHeight())

	p := tAliceSync.peerSet.GetPeer(tAnotherPeerID)
	assert.Equal(t, p.NodeVersion(), version.NodeVersion)
	assert.Equal(t, p.Moniker(), "kitty")
	assert.True(t, pub.EqualsTo(p.PublicKey()))
	assert.Equal(t, p.PeerID(), tAnotherPeerID)
	assert.Equal(t, p.Address(), pub.Address())
	assert.Equal(t, p.Height(), 3)
	assert.Equal(t, p.InitialBlockDownload(), true)
}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 111, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)

	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 111)
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty", pub, tAnotherPeerID, 1, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
}

func TestSendAleykPeerAhead(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty", pub, tAnotherPeerID, 111, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 111)
}

func TestSendAleykPeerSameHeight(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage(payload.ResponseCodeOK, "Welcome!", "kitty", pub, tAnotherPeerID, tAliceState.LastBlockHeight(), 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
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

func TestQueryForTransaction(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()

	// Alice has trx1 in his cache
	tAliceSync.cache.AddTransaction(trx1)
	tBobSync.cache.AddTransaction(trx2)

	tAliceBroadcastCh <- message.NewQueryTransactionsMessage([]crypto.Hash{trx1.ID()})
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)

	tAliceBroadcastCh <- message.NewQueryTransactionsMessage([]crypto.Hash{trx1.ID(), trx2.ID()})
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeTransactions)

	assert.NotNil(t, tAliceSync.cache.GetTransaction(trx2.ID()))
}

func TestHeartbeatNotInSet(t *testing.T) {
	setup(t)

	tAliceConsensus.HRS_ = hrs.NewHRS(106, 0, 1)

	// Alice is not in validator set
	tAliceSync.broadcastHeartBeat()
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)

	msg := message.NewHeartBeatMessage(tAnotherPeerID, crypto.GenerateTestHash(), hrs.NewHRS(106, 1, 1)) // peer is in same height but further round
	tAliceSync.ParsMessage(msg, tAnotherPeerID)
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeVoteSet)
}
