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
	signer       crypto.Signer
	genDoc       *genesis.Genesis
	ctx          context.Context
	txPool       *txpool.MockTxPool
	validTxs     []*tx.Tx
	validBlocks  []block.Block
	validCommits []block.Commit
	totalBlocks  int
	syncConf     *Config
	peerID       peer.ID
)

func setup(t *testing.T) {
	syncConf = TestConfig()
	val, key := validator.GenerateTestValidator(0)
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)
	signer = crypto.NewSigner(key)
	genDoc = genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val}, 1)
	ctx = context.Background()
	txPool = txpool.NewMockTxPool()

	_, _, st := newTestSynchronizer(&signer)

	blockCount := 12
	validBlocks = make([]block.Block, 0, blockCount)
	validTxs = make([]*tx.Tx, 0, blockCount)
	for i := 0; i < blockCount; i++ {
		b := st.ProposeBlock()
		id := b.TxIDs().IDs()[0]
		trx := txPool.PendingTx(id)
		validTxs = append(validTxs, trx)

		v := vote.NewPrecommit(i+1, 0, b.Hash(), signer.Address())
		signer.SignMsg(v)
		sig := v.Signature()
		c := block.NewCommit(0, []block.Committer{
			{Status: 1, Address: signer.Address()}},
			*sig)

		assert.NoError(t, st.ApplyBlock(i+1, b, *c))

		validBlocks = append(validBlocks, b)
		validCommits = append(validCommits, *c)
	}
	if st.LastBlockHeight() != blockCount {
		panic("Unable to commit blocks for test")
	}
	totalBlocks = len(validBlocks)

	peerID, _ = peer.IDFromString("12D3KooWDEWpKkZVxpc8hbLKQL1jvFfyBQDit9AR3ToU4k951Jyi")
}

func newTestSynchronizer(signer *crypto.Signer) (*Synchronizer, *mockNetworkAPI, state.State) {
	consConf := consensus.TestConfig()
	stateConf := state.TestConfig()
	loggerConfig := logger.TestConfig()

	logger.InitLogger(loggerConfig)

	ch := make(chan *message.Message, 10)
	go func() {
		for {
			<-ch
		}
	}()

	if signer == nil {
		_, _, key := crypto.GenerateTestKeyPair()
		s := crypto.NewSigner(key)
		signer = &s
	}

	st, _ := state.LoadOrNewState(stateConf, genDoc, *signer, txPool)
	cons, _ := consensus.NewConsensus(consConf, st, *signer, ch)
	api := mockingNetworkAPI()

	syncer := &Synchronizer{
		ctx:         context.Background(),
		config:      syncConf,
		store:       st.StoreReader(),
		state:       st,
		consensus:   cons,
		txPool:      txPool,
		broadcastCh: ch,
		networkAPI:  api,
	}

	logger := logger.NewLogger("_sync", syncer)

	syncer.logger = logger
	syncer.blockPool = NewBlockPool()
	syncer.stats = stats.NewStats(genDoc.Hash())

	return syncer, api, st
}

func startTestSynchronizer(t *testing.T, signer *crypto.Signer) (*Synchronizer, *mockNetworkAPI, state.State) {
	sync, api, st := newTestSynchronizer(signer)

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

func TestValidBlocks(t *testing.T) {
	sync, _, st := startTestSynchronizer(t, nil)

	blocks := make([]block.Block, totalBlocks)
	copy(blocks, validBlocks)

	// Send transactions
	msg := message.NewTxsMessage(validTxs)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	// Send blocks, all are valid
	msg = message.NewBlocksMessage(1, blocks[:totalBlocks], &validCommits[totalBlocks-1])
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

	// Send blocks, all are valid
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
	setup(t)
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
