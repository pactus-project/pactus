package sync

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/network_api"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/txpool"
)

var (
	tAliceConfig      *Config
	tBobConfig        *Config
	tTxPool           *txpool.MockTxPool
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
}

func setup(t *testing.T) {
	LatestBlockInterval = 10

	tAliceConfig = TestConfig()
	tBobConfig = TestConfig()
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

	// Alice has 12 and Bob has 6 blocks
	lastBlockHash := crypto.Hash{}
	for i := 0; i < 12; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCommit(b.Hash())
		lastBlockHash = b.Hash()
		tAliceState.AddBlock(i+1, b, trxs)
		tAliceState.LastBlockCommit = c

		if i < 6 {
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
	tAliceSync.peerSet = peerset.NewPeerSet()
	tAliceSync.firewall = firewall.NewFirewall(tAliceSync.peerSet, tAliceState)
	tAliceSync.consensusTopic = NewConsensusTopic(tAliceConfig, tAliceConsensus, tAliceSync.logger, tAliceSync.PublishMessage)
	tAliceSync.generalTopic = NewGeneralTopic(tAliceConfig, tAliceNetAPI.SelfID(), tAliceSync.signer.PublicKey(), tAliceSync.peerSet, tAliceState, tAliceSync.logger, tAliceSync.PublishMessage)
	tAliceSync.dataTopic = NewDataTopic(tAliceConfig, tAliceSync.cache, tAliceState, tAliceSync.logger, tAliceSync.PublishMessage)

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
	tBobSync.peerSet = peerset.NewPeerSet()
	tBobSync.firewall = firewall.NewFirewall(tBobSync.peerSet, tBobState)
	tBobSync.consensusTopic = NewConsensusTopic(tBobConfig, tBobConsensus, tBobSync.logger, tBobSync.PublishMessage)
	tBobSync.generalTopic = NewGeneralTopic(tBobConfig, tBobNetAPI.SelfID(), tBobSync.signer.PublicKey(), tBobSync.peerSet, tBobState, tBobSync.logger, tBobSync.PublishMessage)
	tBobSync.dataTopic = NewDataTopic(tBobConfig, tBobSync.cache, tBobState, tBobSync.logger, tBobSync.PublishMessage)

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
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocks)

	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
}
