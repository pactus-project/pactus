package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zarbchain/zarb-go/block"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/txpool"
)

var (
	tState       *state.MockState
	tConsensus   *consensus.MockConsensus
	tTxPool      *txpool.MockTxPool
	tNetAPI      *mockNetworkAPI
	tSync        *Synchronizer
	tCache       *cache.Cache
	tBroadcastCh <-chan *message.Message
	tOurID       peer.ID
)

func setup(t *testing.T) {
	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)

	tOurID, _ = peer.IDFromString("12D3KooWDEWpKkZVxpc8hbLKQL1jvFfyBQDit9AR3ToU4k951Jyi")

	tTxPool = txpool.NewMockTxPool()
	tState = state.NewMockStore()
	tConsensus := consensus.NewMockConsensus()
	tNetAPI = mockingNetworkAPI()
	tCache, _ = cache.NewCache(10, tState.StoreReader())
	tBroadcastCh = make(chan *message.Message, 100)

	// State has some block
	for i := 0; i < 12; i++ {
		b, trxs := block.GenerateTestBlock(nil)
		tState.AddBlock(i+1, b, trxs)
		tState.LastBlockCommit = block.GenerateTestCommit(b.Hash())
	}

	tSync = &Synchronizer{
		ctx:         context.Background(),
		config:      TestConfig(),
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

	tNetAPI.waitingForMessage(t, message.NewSalamMessage(tState.GenHash, tState.LastBlockHeight()))

}

func (mock *mockNetworkAPI) waitingForBroadcaseMessage(t *testing.T, msg *message.Message) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			assert.NoError(t, fmt.Errorf("Timeout"))
			return
		case bMsg := <-tBroadcastCh:
			logger.Info("comparing messages", "broadcastMsg", bMsg, "msg", msg)

		}
	}
}
func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	msg := message.NewSalamMessage(invGenHash, 0)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)
	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, message.PayloadTypeAleyk)

}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	msg := message.NewSalamMessage(tState.GenHash, 0)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)
	tNetAPI.waitingForMessage(t, message.NewAleykMessage(tState.GenHash, tState.LastBlockHeight()))
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)

	msg := message.NewSalamMessage(tState.GenHash, 111)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)
	tNetAPI.waitingForMessage(t, message.NewAleykMessage(tState.GenHash, tState.LastBlockHeight()))
	tNetAPI.waitingForMessage(t, message.NewBlocksReqMessage(tState.LastBlockHeight()+1, 111, tState.LastBlockHash()))
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)

	msg := message.NewAleykMessage(tState.GenHash, 111)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)
	tNetAPI.waitingForMessage(t, message.NewBlocksReqMessage(tState.LastBlockHeight()+1, 111, tState.LastBlockHash()))
}

func TestCacheBlocksAndTransactions(t *testing.T) {
	setup(t)

	b, trxs := block.GenerateTestBlock(nil)

	// Send transactions
	msg := message.NewTxsMessage(trxs)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)

	// Send blocks, all are valid
	msg = message.NewBlocksMessage(10001, []*block.Block{b}, nil)
	data, _ = cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)

	assert.NotNil(t, tCache.GetBlock(10001))
	assert.NotNil(t, tCache.GetTransaction(trxs[0].ID()))
}

// func TestValidBlockNoCommit(t *testing.T) {
// 	setup(t)

// 	blocks := make([]block.Block, totalBlocks)
// 	copy(blocks, validBlocks)

// 	// Send transactions
// 	msg := message.NewTxsMessage(validTxs)
// 	data, _ := cbor.Marshal(msg)
// 	tSync.ParsMessage(data, tOurID)

// 	// Send blocks, all are valid
// 	msg = message.NewBlocksMessage(1, blocks[:totalBlocks], nil)
// 	data, _ = cbor.Marshal(msg)
// 	tSync.ParsMessage(data, tOurID)

// 	assert.Equal(t, st.LastBlockHeight(), totalBlocks-1)
// }

// func TestDuplicateBlock(t *testing.T) {
// 	setup(t)

// 	// Set invalid block 1
// 	// Because the first block is invalid, peer should ask for block 1 and 2 to resend
// 	invalidBlock, _ := block.GenerateTestBlock(nil)
// 	tSync.blockPool.AppendBlock(1, invalidBlock)
// 	assert.Equal(t, tSync.blockPool.blocks[1].Hash(), invalidBlock.Hash())

// 	// Send valid blocks 1
// 	// Block 1 should be replaced with the new one
// 	msg := message.NewBlocksMessage(1, validBlocks[:1], nil)
// 	data, _ := cbor.Marshal(msg)
// 	tSync.ParsMessage(data, tOurID)
// 	assert.Equal(t, tSync.blockPool.blocks[1].Hash(), validBlocks[0].Hash())
// }

// func TestInvalidBlock(t *testing.T) {
// 	setup(t)
// 	setup(t)

// 	blocks := make([]block.Block, totalBlocks)
// 	copy(blocks, validBlocks)

// 	invalidBlock, _ := block.GenerateTestBlock(nil)
// 	blocks[5] = invalidBlock

// 	// Send transactions
// 	msg := message.NewTxsMessage(validTxs)
// 	data, _ := cbor.Marshal(msg)
// 	tSync.ParsMessage(data, tOurID)

// 	// Send blocks, but one block is invalid
// 	msg = message.NewBlocksMessage(1, blocks[:totalBlocks], nil)
// 	data, _ = cbor.Marshal(msg)
// 	tSync.ParsMessage(data, tOurID)

// 	// Should request for block number 5 and 6 again
// 	api.waitingForMessage(t, message.NewBlocksReqMessage(5, 6, st.LastBlockHash()))

// 	assert.Equal(t, st.LastBlockHeight(), 4)

// 	// Now, Let's send valid block for height 5 and 6
// 	msg = message.NewBlocksMessage(5, validBlocks[4:6], nil)
// 	data, _ = cbor.Marshal(msg)
// 	tSync.ParsMessage(data, tOurID)

// 	assert.Equal(t, st.LastBlockHeight(), totalBlocks-1)
// }
