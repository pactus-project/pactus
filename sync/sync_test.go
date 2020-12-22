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
	tState       *state.MockState
	tConsensus   *consensus.MockConsensus
	tTxPool      *txpool.MockTxPool
	tNetAPI      *mockNetworkAPI
	tSync        *Synchronizer
	tCache       *cache.Cache
	tBroadcastCh chan *message.Message
	tPeerID      peer.ID
	tSelfID      peer.ID
)

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	syncConf := TestConfig()
	_, _, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)

	tPeerID, _ = peer.IDB58Decode("12D3KooWLQ8GKaLdKU8Ms6AkMYjDWCr5UTPvdewag3tcarxh7saC")
	tSelfID, _ = peer.IDB58Decode("12D3KooWHyepEGGdeSk3nPZrEamxLNba7tFZJKWbyEdZ654fHJdk")

	tTxPool = txpool.NewMockTxPool()
	tState = state.NewMockStore()
	tConsensus = consensus.NewMockConsensus()
	tNetAPI = mockingNetworkAPI()
	tCache, _ = cache.NewCache(syncConf.CacheSize, tState.StoreReader())
	tBroadcastCh = make(chan *message.Message, 100)

	// State has some block
	for i := 0; i < 12; i++ {
		b, trxs := block.GenerateTestBlock(nil)
		tState.AddBlock(i+1, b, trxs)
		tState.LastBlockCommit = block.GenerateTestCommit(b.Hash())
	}

	tSync = &Synchronizer{
		ctx:         context.Background(),
		config:      syncConf,
		signer:      signer,
		state:       tState,
		consensus:   tConsensus,
		cache:       tCache,
		txPool:      tTxPool,
		broadcastCh: tBroadcastCh,
		networkAPI:  tNetAPI,
	}

	logger := logger.NewLogger("_sync", tSync)

	tSync.logger = logger
	tSync.stats = stats.NewStats(tState.GenHash)

	assert.NoError(t, tSync.Start())

	tNetAPI.waitingForMessage(t, message.NewSalamMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, tState.LastBlockHeight()))
	tNetAPI.waitingForMessage(t, message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, tState.LastBlockHeight(), 0, "Welcome!"))
}

func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	msg := message.NewSalamMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), invGenHash, 0)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tPeerID)
	tNetAPI.waitingForMessage(t, message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, tState.LastBlockHeight(), 1, "Invalid genesis hash"))
}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	msg := message.NewSalamMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, 1)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tPeerID)
	tNetAPI.waitingForMessage(t, message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, tState.LastBlockHeight(), 0, "Welcome!"))
	assert.Equal(t, tSync.stats.MaxHeight(), tState.LastBlockHeight())
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)

	msg := message.NewSalamMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, 111)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tPeerID)
	tNetAPI.waitingForMessage(t, message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, tState.LastBlockHeight(), 0, "Welcome!"))
	tNetAPI.waitingForMessage(t, message.NewBlocksReqMessage(tState.LastBlockHeight()+1, 111, tState.LastBlockHash()))
	assert.Equal(t, tSync.stats.MaxHeight(), 111)
}
func TestSendAleykPeerAhead(t *testing.T) {
	setup(t)

	msg := message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, 1, 0, "Welcome!")
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tPeerID)
	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, payload.PayloadTypeBlocksReq)
	assert.Equal(t, tSync.stats.MaxHeight(), tState.LastBlockHeight())
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)

	msg := message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, 111, 0, "Welcome!")
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tPeerID)
	tNetAPI.waitingForMessage(t, message.NewBlocksReqMessage(tState.LastBlockHeight()+1, 111, tState.LastBlockHash()))
	assert.Equal(t, tSync.stats.MaxHeight(), 111)
}

func TestSendAleykPeerSameHeight(t *testing.T) {
	setup(t)

	msg := message.NewAleykMessage(tSync.config.Moniker, tSync.signer.PublicKey(), tNetAPI.SelfID(), tState.GenHash, tState.LastBlockHeight(), 0, "Welcome!")
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tPeerID)
	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, payload.PayloadTypeBlocksReq)
	assert.Equal(t, tSync.stats.MaxHeight(), tState.LastBlockHeight())

}

func TestCacheBlocksAndTransactions(t *testing.T) {
	setup(t)

	b, trxs := block.GenerateTestBlock(nil)

	// Send transactions
	tSync.broadcastTxs(trxs)

	// Send blocks
	tSync.broadcastBlocks(1001, []*block.Block{b}, nil)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeBlocks)

	assert.NotNil(t, tCache.GetBlock(1001))
	assert.NotNil(t, tCache.GetCommit(b.Header().LastBlockHash()))
	assert.NotNil(t, tCache.GetTransaction(trxs[0].ID()))
}

func TestDuplicatingBlock(t *testing.T) {
	setup(t)

	b1, _ := block.GenerateTestBlock(nil)
	b2, _ := block.GenerateTestBlock(nil)

	// Send block 1
	tSync.broadcastBlocks(1001, []*block.Block{b1}, nil)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeBlocks)
	assert.Equal(t, tCache.GetBlock(1001).Hash(), b1.Hash())

	// Send block 1 again, should overwrite the first one in cache
	tSync.broadcastBlocks(1001, []*block.Block{b2}, nil)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeBlocks)
	assert.Equal(t, tCache.GetBlock(1001).Hash(), b2.Hash())
}

func TestCheckTxsInCache(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSortitionTx()

	tCache.AddTransaction(trx1)
	msg := message.NewTxsReqMessage([]crypto.Hash{trx1.ID(), trx2.ID()})
	tBroadcastCh <- msg
	tNetAPI.waitingForMessage(t, message.NewTxsReqMessage([]crypto.Hash{trx2.ID()}))

	tCache.AddTransaction(trx3)
	msg = message.NewTxsReqMessage([]crypto.Hash{trx3.ID()})
	tBroadcastCh <- msg
	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, payload.PayloadTypeTxsReq)
}

// TODO: test these situation in consensus
func TestRequestForProposal(t *testing.T) {
	setup(t)

	hrs := hrs.NewHRS(100, 1, 6)
	p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
	tConsensus.SetProposal(p)
	tConsensus.HRS_ = hrs
	t.Run("Request proposal at same HRS", func(t *testing.T) {
		msg := message.NewProposalReqMessage(hrs.Height(), hrs.Round())
		tBroadcastCh <- msg
		tNetAPI.waitingForMessage(t, message.NewProposalMessage(p))
	})

	t.Run("Request proposal at same height but lower HRS", func(t *testing.T) {
		msg := message.NewProposalReqMessage(hrs.Height(), hrs.Round()-1)
		tBroadcastCh <- msg
		tNetAPI.waitingForMessage(t, message.NewProposalMessage(p))
	})

	t.Run("Request proposal at same height but higher HRS", func(t *testing.T) {
		msg := message.NewProposalReqMessage(hrs.Height(), hrs.Round()+1)
		tBroadcastCh <- msg
		tNetAPI.waitingForMessage(t, message.NewProposalReqMessage(hrs.Height(), hrs.Round()+1))
	})

}
