package sync

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidBlockData(t *testing.T) {
	td := setup(t, nil)

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
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, message.ResponseCodeMoreBlocks.String(), sid,
			lastHeight+1, [][]byte{test.data}, cert)

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
	td.addPeer(t, pub, pid, services.New(services.None))

	sid := td.sync.peerSet.OpenSession(pid).SessionID()
	msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, t.Name(), sid,
		lastHeight+1, [][]byte{d1}, cert1)
	assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

	assert.Nil(t, td.sync.peerSet.FindSession(sid))
	assert.Equal(t, td.state.LastBlockHeight(), lastHeight+1)
}

func TestStrippedPublicKey(t *testing.T) {
	td := setup(t, nil)

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
	td.addPeer(t, peerPubKey, pid, services.New(services.None))

	for _, test := range tests {
		blkData, _ := test.blk.Bytes()
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
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
	ts := testsuite.NewTestSuite(t)

	configAlice := testConfig()
	configBob := testConfig()

	valKeyAlice := []*bls.ValidatorKey{ts.RandValKey()}
	valKeyBob := []*bls.ValidatorKey{ts.RandValKey()}
	stateAlice := state.MockingState(ts)
	stateBob := state.MockingState(ts)
	consMgrAlice, _ := consensus.MockingManager(ts, valKeyAlice)
	consMgrBob, _ := consensus.MockingManager(ts, valKeyBob)
	broadcastChAlice := make(chan message.Message, 1000)
	broadcastChBob := make(chan message.Message, 1000)
	networkAlice := network.MockingNetwork(ts, ts.RandPeerID())
	networkBob := network.MockingNetwork(ts, ts.RandPeerID())

	configBob.NodeNetwork = true
	networkAlice.AddAnotherNetwork(networkBob)
	networkBob.AddAnotherNetwork(networkAlice)

	// Adding 100 blocks for Bob
	blockInterval := stateBob.Genesis().Params().BlockInterval()
	blockTime := util.RoundNow(int(blockInterval.Seconds()))
	for i := uint32(0); i < 100; i++ {
		blk, cert := ts.GenerateTestBlockWithTime(i+1, blockTime)
		assert.NoError(t, stateBob.CommitBlock(blk, cert))

		blockTime = blockTime.Add(blockInterval)
	}

	sync1, err := NewSynchronizer(configAlice,
		valKeyAlice,
		stateAlice,
		consMgrAlice,
		networkAlice,
		broadcastChAlice,
	)
	assert.NoError(t, err)
	syncAlice := sync1.(*synchronizer)

	sync2, err := NewSynchronizer(configBob,
		valKeyBob,
		stateBob,
		consMgrBob,
		networkBob,
		broadcastChBob,
	)
	assert.NoError(t, err)
	syncBob := sync2.(*synchronizer)

	// -------------------------------
	// Better logging during testing
	overrideLogger := func(sync *synchronizer, name string) {
		sync.logger = logger.NewSubLogger("_sync", &OverrideStringer{
			name: fmt.Sprintf("%s - %s: ", name, t.Name()), sync: sync,
		})
	}

	overrideLogger(syncAlice, "Alice")
	overrideLogger(syncBob, "Bob")
	// -------------------------------

	assert.NoError(t, syncAlice.Start())
	assert.NoError(t, syncBob.Start())

	// Verify that Hello messages are exchanged between Alice and Bob
	assert.NoError(t, syncAlice.sayHello(syncBob.SelfID()))
	assert.NoError(t, syncBob.sayHello(syncAlice.SelfID()))

	shouldPublishMessageWithThisType(t, networkAlice, message.TypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeHello)

	shouldPublishMessageWithThisType(t, networkBob, message.TypeHelloAck)
	shouldPublishMessageWithThisType(t, networkAlice, message.TypeHelloAck)

	time.Sleep(1 * time.Second)

	// Ensure peers are connected and block heights are correct
	assert.Equal(t, 1, syncAlice.PeerSet().Len())
	assert.Equal(t, 1, syncBob.PeerSet().Len())
	require.Equal(t, peerset.StatusCodeKnown, syncAlice.PeerSet().GetPeer(syncBob.SelfID()).Status)
	require.Equal(t, peerset.StatusCodeKnown, syncBob.PeerSet().GetPeer(syncAlice.SelfID()).Status)
	assert.Equal(t, uint32(0), syncAlice.state.LastBlockHeight())
	assert.Equal(t, uint32(100), syncBob.state.LastBlockHeight())

	// Alice receives a BlockAnnounce message and starts updating its blockchain
	syncAlice.updateBlockchain()

	// Perform block syncing
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
