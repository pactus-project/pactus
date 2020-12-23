package sync

import (
	"context"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	tTxPool           *txpool.MockTxPool
	tAliceState       *state.MockState
	tBobState         *state.MockState
	tAliceConsensus   *consensus.MockConsensus
	tBobConsensus     *consensus.MockConsensus
	tAliceNetAPI      *mockNetworkAPI
	tBobNetAPI        *mockNetworkAPI
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
	syncConf := TestConfig()
	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	aliceSigner := crypto.NewSigner(priv1)
	bobSigner := crypto.NewSigner(priv2)

	tTxPool = txpool.NewMockTxPool()

	tAlicePeerID, _ = peer.IDB58Decode("12D3KooWLQ8GKaLdKU8Ms6AkMYjDWCr5UTPvdewag3tcarxh7saC")
	tBobPeerID, _ = peer.IDB58Decode("12D3KooWHyepEGGdeSk3nPZrEamxLNba7tFZJKWbyEdZ654fHJdk")
	tAnotherPeerID, _ = peer.IDB58Decode("12D3KooWM4dZKiZ8y21biCZXuAJYD5db8vSr1hfMpSgBSqpekY4Q")
	tAliceState = state.NewMockState()
	tBobState = state.NewMockState()
	tAliceConsensus = consensus.NewMockConsensus()
	tBobConsensus = consensus.NewMockConsensus()
	tAliceBroadcastCh = make(chan *message.Message, 100)
	tBobBroadcastCh = make(chan *message.Message, 100)
	tAliceNetAPI = mockingNetworkAPI(tAlicePeerID)
	tBobNetAPI = mockingNetworkAPI(tBobPeerID)
	aliceCache, _ := cache.NewCache(syncConf.CacheSize, tAliceState.StoreReader())
	bobCache, _ := cache.NewCache(syncConf.CacheSize, tBobState.StoreReader())

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
		config:      syncConf,
		signer:      aliceSigner,
		state:       tAliceState,
		consensus:   tAliceConsensus,
		cache:       aliceCache,
		txPool:      tTxPool,
		broadcastCh: tAliceBroadcastCh,
		networkAPI:  tAliceNetAPI,
		stats:       stats.NewStats(tAliceState.GenHash),
	}
	tAliceSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "alice: ", sync: tAliceSync})

	tBobSync = &Synchronizer{
		ctx:         context.Background(),
		config:      syncConf,
		signer:      bobSigner,
		state:       tBobState,
		consensus:   tBobConsensus,
		cache:       bobCache,
		txPool:      tTxPool,
		broadcastCh: tBobBroadcastCh,
		networkAPI:  tBobNetAPI,
		stats:       stats.NewStats(tBobState.GenHash),
	}
	tBobSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "bob: ", sync: tBobSync})

	tAliceNetAPI.peerSync = tBobSync
	tBobNetAPI.peerSync = tAliceSync

	assert.NoError(t, tAliceSync.Start())
	assert.NoError(t, tBobSync.Start())

	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeSalam)
	tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeSalam)

	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)

	tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocksReq)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocks)
}

func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("bad-genesis", pub, tAnotherPeerID, invGenHash, 0)
	data, _ := cbor.Marshal(msg)
	tAliceSync.ParsMessage(data, tAnotherPeerID)
	msg2 := tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.Response.Status, payload.SalamResponseCodeRejected)
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 0)
	data, _ := cbor.Marshal(msg)
	tAliceSync.ParsMessage(data, tAnotherPeerID)
	msg2 := tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.Response.Status, payload.SalamResponseCodeOK)
	assert.Equal(t, tBobSync.stats.MaxClaimedHeight(), tAliceState.LastBlockHeight())
}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 111)
	data, _ := cbor.Marshal(msg)
	tAliceSync.ParsMessage(data, tAnotherPeerID)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	tAliceNetAPI.shouldPublishThisMessage(t, message.NewBlocksReqMessage(tAliceState.LastBlockHeight()+1, 111, tAliceState.LastBlockHash()))

	assert.Equal(t, tAliceSync.stats.MaxClaimedHeight(), 111)
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 1, 0, "Welcome!")
	data, _ := cbor.Marshal(msg)
	tAliceSync.ParsMessage(data, tAnotherPeerID)
	tAliceNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeBlocksReq)
}

func TestSendAleykPeerAhead(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 111, 0, "Welcome!")
	data, _ := cbor.Marshal(msg)
	tAliceSync.ParsMessage(data, tAnotherPeerID)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocksReq)
	assert.Equal(t, tAliceSync.stats.MaxClaimedHeight(), 111)
}

func TestSendAleykPeerSameHeight(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, tAliceState.LastBlockHeight(), 0, "Welcome!")
	data, _ := cbor.Marshal(msg)
	tAliceSync.ParsMessage(data, tAnotherPeerID)
	tAliceNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeBlocksReq)
}

func TestAddBlockToCache(t *testing.T) {
	setup(t)

	b1, _ := block.GenerateTestBlock(nil, nil)

	// Alice send block to bob, bob should cache it
	tAliceSync.broadcastBlocks(1001, []*block.Block{b1}, nil)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocks)
	assert.Equal(t, tBobSync.cache.GetBlock(1001).Hash(), b1.Hash())
}

func TestAddTxToCache(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()

	// Alice send transaction to bob, bob should cache it
	tAliceSync.broadcastTxs([]*tx.Tx{trx1})
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeTxs)
	assert.NotNil(t, tBobSync.cache.GetTransaction(trx1.ID()))
}

func TestRequestForProposal(t *testing.T) {

	t.Run("Alice and bob are in same height. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(100, 1, 6)
		p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
		tAliceConsensus.SetProposal(p)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewProposalReqMessage(hrs.Height(), hrs.Round())
		tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposalReq)
		tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), tBobConsensus.Proposal.Hash())
	})

	t.Run("Alice and bob are in same height. Alice doesn't have have proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(101, 1, 6)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewProposalReqMessage(hrs.Height(), hrs.Round())
		tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposalReq)

		// Alice doesn't respond
		tAliceNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

	t.Run("Alice and bob are in same height. Alice is in next round. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(102, 1, 6)
		p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
		tAliceConsensus.SetProposal(p)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewProposalReqMessage(hrs.Height(), hrs.Round()-1)
		tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposalReq)
		tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), tBobConsensus.Proposal.Hash())
	})

	t.Run("Alice and bob are in same height. Alice is in previous round. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(103, 1, 6)
		p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
		tAliceConsensus.SetProposal(p)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewProposalReqMessage(hrs.Height(), hrs.Round()+1)
		tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposalReq)

		// Alice doesn't respond
		tAliceNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

}
