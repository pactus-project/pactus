package sync

import (
	"io"
	"testing"
	"time"

	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidBlockData(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(10)

	lastHeight := td.state.LastBlockHeight()
	prevCert := td.GenerateTestCertificate(lastHeight)
	cert := td.GenerateTestCertificate(lastHeight + 1)
	blk := block.MakeBlock(1, time.Now(), nil, td.RandHash(), td.RandHash(),
		prevCert, td.RandSeed(), td.RandValAddress())
	data, _ := blk.Bytes()
	tests := []struct {
		data []byte
		err  error
	}{
		{
			td.RandBytes(16),
			io.ErrUnexpectedEOF,
		},
		{
			data,
			block.BasicCheckError{
				Reason: "no subsidy transaction",
			},
		},
	}

	for _, test := range tests {
		pid := td.RandPeerID()
		sid := td.RandInt(1000)
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
			message.ResponseCodeMoreBlocks.String(),
			sid, lastHeight+1, [][]byte{test.data}, cert)

		err := td.receivingNewMessage(td.sync, msg, pid)
		assert.ErrorIs(t, err, test.err)
	}
}

func TestOneBlockShorter(t *testing.T) {
	td := setup(t, nil)

	lastHeight := td.state.LastBlockHeight()
	blk1, cert1 := td.GenerateTestBlock(lastHeight + 1)
	d1, _ := blk1.Bytes()
	pid := td.RandPeerID()

	pub, _ := td.RandBLSKeyPair()
	td.addPeer(t, pub, pid, service.New(service.None))

	sid := td.RandInt(1000)
	msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, t.Name(), sid,
		lastHeight+1, [][]byte{d1}, cert1)
	assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

	assert.Equal(t, td.state.LastBlockHeight(), lastHeight+1)
}

func TestStrippedPublicKey(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(10)

	lastHeight := td.state.LastBlockHeight()

	// Add a new block and keep the signer key
	indexedPub, indexedPrv := td.RandBLSKeyPair()
	trx0 := tx.NewTransferTx(lastHeight, indexedPub.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	td.HelperSignTransaction(indexedPrv, trx0)
	trxs0 := []*tx.Tx{trx0}
	blk0 := block.MakeBlock(1, time.Now(), trxs0, td.RandHash(), td.RandHash(),
		td.state.LastCertificate(), td.RandSeed(), td.RandValAddress())
	cert0 := td.GenerateTestCertificate(lastHeight + 1)
	err := td.state.CommitBlock(blk0, cert0)
	require.NoError(t, err)
	lastHeight++
	// -----

	rndPub, rndPrv := td.RandBLSKeyPair()
	trx1 := tx.NewTransferTx(lastHeight, rndPub.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	td.HelperSignTransaction(rndPrv, trx1)
	trx1.StripPublicKey()
	trxs1 := []*tx.Tx{trx1}
	blk1 := block.MakeBlock(1, time.Now(), trxs1, td.RandHash(), td.RandHash(),
		cert0, td.RandSeed(), td.RandValAddress())

	trx2 := tx.NewTransferTx(lastHeight, indexedPub.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	td.HelperSignTransaction(indexedPrv, trx2)
	trx2.StripPublicKey()
	trxs2 := []*tx.Tx{trx2}
	blk2 := block.MakeBlock(1, time.Now(), trxs2, td.RandHash(), td.RandHash(),
		cert0, td.RandSeed(), td.RandValAddress())

	tests := []struct {
		blk *block.Block
		err error
	}{
		{
			blk1,
			store.ErrNotFound,
		},
		{
			blk2,
			nil,
		},
	}

	// Add a peer
	pid := td.RandPeerID()
	peerPubKey, _ := td.RandBLSKeyPair()
	td.addPeer(t, peerPubKey, pid, service.New(service.None))

	for _, test := range tests {
		blkData, _ := test.blk.Bytes()
		sid := td.RandInt(1000)
		cert := td.GenerateTestCertificate(lastHeight + 1)
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, message.ResponseCodeRejected.String(), sid,
			lastHeight+1, [][]byte{blkData}, cert)
		err := td.receivingNewMessage(td.sync, msg, pid)

		assert.ErrorIs(t, err, test.err)
	}
}

// TestSyncing is an important test to verify the syncing process between two
// test nodes, Alice and Bob. In real-world scenarios, multiple nodes are typically
// involved, but the procedure remains similar.
func TestSyncing(t *testing.T) {
	ts, syncAlice, networkAlice, syncBob, networkBob := makeAliceAndBobNetworks(t)

	// Adding 100 blocks for Bob
	blockInterval := syncBob.state.Genesis().Params().BlockInterval()
	blockTime := util.RoundNow(int(blockInterval.Seconds()))
	for i := uint32(0); i < 100; i++ {
		blk, cert := ts.GenerateTestBlockWithTime(i+1, blockTime)
		assert.NoError(t, syncBob.state.CommitBlock(blk, cert))

		blockTime = blockTime.Add(blockInterval)
	}

	assert.Equal(t, uint32(0), syncAlice.state.LastBlockHeight())
	assert.Equal(t, uint32(100), syncBob.state.LastBlockHeight())

	// Announcing a block
	blk, cert := ts.GenerateTestBlock(ts.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)
	syncBob.broadcast(msg)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlockAnnounce)

	// Perform block syncing
	shouldNotPublishMessageWithThisType(t, networkBob, message.TypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkAlice, message.TypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 1-11
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 12-22
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 23-23
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.TypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 24-34
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 35-45
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 46-46
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.TypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 47-57
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 58-68
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 69-69
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.TypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 70-80
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 81-91
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // 92-92
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // NoMoreBlock

	// Last block requests
	shouldPublishMessageWithThisType(t, networkAlice, message.TypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse)        // 93-100
	bdl := shouldPublishMessageWithThisType(t, networkBob, message.TypeBlocksResponse) // Synced
	assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)

	// Alice needs more time to process all the bundles,
	// but the block height should be greater than zero
	assert.Greater(t, syncAlice.state.LastBlockHeight(), uint32(20))
	assert.Equal(t, syncBob.state.LastBlockHeight(), uint32(100))
}
