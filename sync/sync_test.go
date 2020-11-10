package sync

import (
	"context"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	privVal     *validator.PrivValidator
	genDoc      *genesis.Genesis
	ctx         context.Context
	txPool      *txpool.TxPool
	validTxs    []tx.Tx
	validBlocks []block.Block
	totalBlocks int
	lastCommit  *block.Commit
	syncConf    *Config
	peerID      peer.ID
)

func init() {
	syncConf = TestConfig()
	val, key := validator.GenerateTestValidator()
	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)
	privVal = validator.NewPrivValidator(key)
	genDoc = genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
	ctx = context.Background()

	_, _, st := newTestSynchronizer(privVal)

	blockCount := 12
	validBlocks = make([]block.Block, 0, blockCount)
	validTxs = make([]tx.Tx, 0, blockCount)
	for i := 0; i < blockCount; i++ {
		b := st.ProposeBlock()
		txHash := b.TxHashes().Hashes()[0]
		trx := txPool.PendingTx(txHash)
		validBlocks = append(validBlocks, b)
		validTxs = append(validTxs, *trx)
		v := vote.NewPrecommit(i+1, 0, b.Hash(), privVal.Address())
		privVal.SignMsg(v)
		sig := v.Signature()
		lastCommit = block.NewCommit(0, []crypto.Address{privVal.Address()}, []crypto.Signature{*sig})

		st.ApplyBlock(i+1, b, *lastCommit)
	}
	if st.LastBlockHeight() != blockCount {
		panic("Unable to commit blocks for test")
	}
	totalBlocks = len(validBlocks)

	peerID, _ = peer.IDFromString("12D3KooWDEWpKkZVxpc8hbLKQL1jvFfyBQDit9AR3ToU4k951Jyi")
}

func newTestSynchronizer(pVal *validator.PrivValidator) (*Synchronizer, *mockNetworkApi, state.State) {
	consConf := consensus.TestConfig()
	stateConf := state.TestConfig()
	txPoolConf := txpool.TestConfig()
	loggerConfig := logger.TestConfig()

	logger.InitLogger(loggerConfig)

	ch := make(chan *message.Message, 10)
	go func() {
		for {
			select {
			case <-ch:
			default:
			}
		}
	}()

	if pVal == nil {
		_, _, key := crypto.GenerateTestKeyPair()
		pVal = validator.NewPrivValidator(key)
	}

	txPoolConf.WaitingTimeout = 100 * time.Millisecond
	txPool, _ = txpool.NewTxPool(txPoolConf, ch)
	st, _ := state.LoadOrNewState(stateConf, genDoc, pVal.Address(), txPool)
	cons, _ := consensus.NewConsensus(consConf, st, pVal, ch)
	api := mockingNetworkApi()

	syncer := &Synchronizer{
		ctx:         context.Background(),
		config:      syncConf,
		store:       st.StoreReader(),
		state:       st,
		consensus:   cons,
		txPool:      txPool,
		txkPool:     make(map[crypto.Hash]*tx.Tx),
		broadcastCh: ch,
		networkApi:  api,
	}

	logger := logger.NewLogger("_sync", syncer)

	syncer.logger = logger
	syncer.blockPool = NewBlockPool()
	syncer.stats = stats.NewStats(genDoc.Hash())

	return syncer, api, st
}

func startTestSynchronizer(t *testing.T, pVal *validator.PrivValidator) (*Synchronizer, *mockNetworkApi, state.State) {
	sync, api, st := newTestSynchronizer(pVal)

	assert.NoError(t, sync.Start())

	// Stopping HeartBeat ticker
	sync.heartBeatTicker.Stop()

	api.waitingForMessage(t, message.NewSalamMessage(genDoc.Hash(), 0))

	msg := message.NewSalamMessage(genDoc.Hash(), 12)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)
	api.waitingForMessage(t, message.NewBlocksReqMessage(1, 12, st.LastBlockHash()))

	return sync, api, st
}

func TestValidBlock(t *testing.T) {
	sync, _, st := startTestSynchronizer(t, nil)

	blocks := make([]block.Block, totalBlocks)
	copy(blocks, validBlocks)

	// Send transactions
	msg := message.NewTxsMessage(validTxs)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	// Send blocks, but one block is invalid
	msg = message.NewBlocksMessage(1, blocks[:totalBlocks], lastCommit)
	data, _ = cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	assert.Equal(t, st.LastBlockHeight(), totalBlocks)
}

func TestValidBlockNoCommit(t *testing.T) {
	sync, _, st := startTestSynchronizer(t, nil)

	blocks := make([]block.Block, totalBlocks)
	copy(blocks, validBlocks)

	// Send transactions
	msg := message.NewTxsMessage(validTxs)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	// Send blocks, but one block is invalid
	msg = message.NewBlocksMessage(1, blocks[:totalBlocks], nil)
	data, _ = cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	assert.Equal(t, st.LastBlockHeight(), totalBlocks-1)
}

func TestDuplicateBlock(t *testing.T) {
	sync, _, _ := startTestSynchronizer(t, nil)

	// Set invalid block 1
	// Because the first block is invalid, peer should ask for block 1 and 2 to resend
	invalidBlock, _ := block.GenerateTestBlock(nil)
	sync.blockPool.AppendBlock(1, invalidBlock)
	assert.Equal(t, sync.blockPool.blocks[1].Hash(), invalidBlock.Hash())

	// Send valid blocks 1
	// Block 1 should be replaced with the new one
	msg := message.NewBlocksMessage(1, validBlocks[:1], nil)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)
	assert.Equal(t, sync.blockPool.blocks[1].Hash(), validBlocks[0].Hash())
}

func TestInvalidBlock(t *testing.T) {
	sync, api, st := startTestSynchronizer(t, nil)

	blocks := make([]block.Block, totalBlocks)
	copy(blocks, validBlocks)

	invalidBlock, _ := block.GenerateTestBlock(nil)
	blocks[5] = invalidBlock

	// Send transactions
	msg := message.NewTxsMessage(validTxs)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	// Send blocks, but one block is invalid
	msg = message.NewBlocksMessage(1, blocks[:totalBlocks], nil)
	data, _ = cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	// Should request for block number 5 and 6 again
	api.waitingForMessage(t, message.NewBlocksReqMessage(5, 6, st.LastBlockHash()))

	assert.Equal(t, st.LastBlockHeight(), 4)

	// Now, Let's send valid block for height 5 and 6
	msg = message.NewBlocksMessage(5, validBlocks[4:6], nil)
	data, _ = cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	assert.Equal(t, st.LastBlockHeight(), totalBlocks-1)
}
