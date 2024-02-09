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
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
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
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, message.ResponseCodeMoreBlocks.String(), sid,
			lastHeight+1, [][]byte{blkData}, cert)
		err := td.receivingNewMessage(td.sync, msg, pid)

		assert.ErrorIs(t, err, test.err)
	}
}

func shouldPublishBlockRequest(t *testing.T, net *network.MockNetwork, from uint32) {
	t.Helper()

	bdl := shouldPublishMessageWithThisType(t, net, message.TypeBlocksRequest)
	msg := bdl.Message.(*message.BlocksRequestMessage)
	require.Equal(t, from, msg.From)
}

func shouldPublishBlockResponse(t *testing.T, net *network.MockNetwork,
	from, count uint32, code message.ResponseCode,
) {
	t.Helper()

	bdl := shouldPublishMessageWithThisType(t, net, message.TypeBlocksResponse)
	msg := bdl.Message.(*message.BlocksResponseMessage)
	require.Equal(t, from, msg.From)
	require.Equal(t, count, msg.Count())
	require.Equal(t, code, msg.ResponseCode)
}

type networkAliceBob struct {
	*testsuite.TestSuite

	stateAlice   *state.MockState
	stateBob     *state.MockState
	networkAlice *network.MockNetwork
	networkBob   *network.MockNetwork
	syncAlice    *synchronizer
	syncBob      *synchronizer
}

func makeAliceAndBobNetworks(t *testing.T) *networkAliceBob {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	configAlice := testConfig()
	configBob := testConfig()

	valKeyAlice := []*bls.ValidatorKey{ts.RandValKey()}
	valKeyBob := []*bls.ValidatorKey{ts.RandValKey()}
	stateAlice := state.MockingState(ts)
	stateBob := state.MockingState(ts)
	consMgrAlice, _ := consensus.MockingManager(ts, valKeyAlice)
	consMgrBob, _ := consensus.MockingManager(ts, valKeyBob)
	internalMessageCh := make(chan message.Message, 1000)
	networkAlice := network.MockingNetwork(ts, ts.RandPeerID())
	networkBob := network.MockingNetwork(ts, ts.RandPeerID())

	configBob.NodeNetwork = true
	networkAlice.AddAnotherNetwork(networkBob)
	networkBob.AddAnotherNetwork(networkAlice)

	sync1, err := NewSynchronizer(configAlice,
		valKeyAlice,
		stateAlice,
		consMgrAlice,
		networkAlice,
		internalMessageCh,
	)
	assert.NoError(t, err)
	syncAlice := sync1.(*synchronizer)

	sync2, err := NewSynchronizer(configBob,
		valKeyBob,
		stateBob,
		consMgrBob,
		networkBob,
		internalMessageCh,
	)
	assert.NoError(t, err)
	syncBob := sync2.(*synchronizer)

	// -------------------------------
	// Better logging during testing
	overrideLogger := func(sync *synchronizer, name string) {
		sync.logger = logger.NewSubLogger("_sync",
			testsuite.NewOverrideStringer(fmt.Sprintf("%s - %s: ", name, t.Name()), sync))
	}

	overrideLogger(syncAlice, "Alice")
	overrideLogger(syncBob, "Bob")
	// -------------------------------

	assert.NoError(t, syncAlice.Start())
	assert.NoError(t, syncBob.Start())

	// Verify that Hello messages are exchanged between Alice and Bob
	syncAlice.sayHello(syncBob.SelfID())
	syncBob.sayHello(syncAlice.SelfID())

	shouldPublishMessageWithThisType(t, networkAlice, message.TypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeHello)

	shouldPublishMessageWithThisType(t, networkBob, message.TypeHelloAck)
	shouldPublishMessageWithThisType(t, networkAlice, message.TypeHelloAck)

	// Ensure peers are connected and block heights are correct
	require.Eventually(t, func() bool {
		return syncAlice.PeerSet().Len() == 1 &&
			syncBob.PeerSet().Len() == 1
	}, time.Second, 100*time.Millisecond)

	require.Equal(t, peerset.StatusCodeKnown, syncAlice.PeerSet().GetPeer(syncBob.SelfID()).Status)
	require.Equal(t, peerset.StatusCodeKnown, syncBob.PeerSet().GetPeer(syncAlice.SelfID()).Status)

	return &networkAliceBob{
		TestSuite:    ts,
		syncAlice:    syncAlice,
		stateAlice:   stateAlice,
		networkAlice: networkAlice,
		syncBob:      syncBob,
		stateBob:     stateBob,
		networkBob:   networkBob,
	}
}

// TestIdenticalBundles tests if two different peers publish the same message,
// whether the bundle data is also the same.
func TestIdenticalBundles(t *testing.T) {
	td := makeAliceAndBobNetworks(t)

	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)

	bdlAlice := td.syncAlice.prepareBundle(msg)
	bdlBob := td.syncBob.prepareBundle(msg)

	assert.Equal(t, bdlAlice, bdlBob)
}

// TestSyncing is an important test to verify the syncing process between two
// test nodes, Alice and Bob. In real-world scenarios, multiple nodes are typically
// involved, but the procedure remains similar.
func TestSyncing(t *testing.T) {
	td := makeAliceAndBobNetworks(t)

	// Adding 100 blocks for Bob
	blockInterval := td.syncBob.state.Genesis().Params().BlockInterval()
	blockTime := td.syncBob.state.Genesis().GenesisTime()
	for i := uint32(0); i < 100; i++ {
		blk, cert := td.GenerateTestBlockWithTime(i+1, blockTime)
		assert.NoError(t, td.syncBob.state.CommitBlock(blk, cert))

		blockTime = blockTime.Add(blockInterval)
	}

	assert.Equal(t, uint32(0), td.syncAlice.state.LastBlockHeight())
	assert.Equal(t, uint32(100), td.syncBob.state.LastBlockHeight())

	// Announcing a block
	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)
	td.syncBob.broadcast(msg)
	shouldPublishMessageWithThisType(t, td.networkBob, message.TypeBlockAnnounce)

	// Perform block syncing
	assert.Equal(t, uint32(11), td.syncAlice.config.BlockPerMessage)
	assert.Equal(t, uint32(23), td.syncAlice.config.LatestBlockInterval)

	shouldNotPublishMessageWithThisType(t, td.networkBob, message.TypeBlocksRequest)
	shouldPublishBlockRequest(t, td.networkAlice, 1)
	shouldPublishBlockResponse(t, td.networkBob, 1, 11, message.ResponseCodeMoreBlocks)  // 1-11
	shouldPublishBlockResponse(t, td.networkBob, 12, 11, message.ResponseCodeMoreBlocks) // 12-22
	shouldPublishBlockResponse(t, td.networkBob, 23, 1, message.ResponseCodeMoreBlocks)  // 23-23
	shouldPublishBlockResponse(t, td.networkBob, 0, 0, message.ResponseCodeNoMoreBlocks) // NoMoreBlock

	shouldPublishBlockRequest(t, td.networkAlice, 24)
	shouldPublishBlockResponse(t, td.networkBob, 24, 11, message.ResponseCodeMoreBlocks) // 24-34
	shouldPublishBlockResponse(t, td.networkBob, 35, 11, message.ResponseCodeMoreBlocks) // 35-45
	shouldPublishBlockResponse(t, td.networkBob, 46, 1, message.ResponseCodeMoreBlocks)  // 46-46
	shouldPublishBlockResponse(t, td.networkBob, 0, 0, message.ResponseCodeNoMoreBlocks) // NoMoreBlock

	shouldPublishBlockRequest(t, td.networkAlice, 47)
	shouldPublishBlockResponse(t, td.networkBob, 47, 11, message.ResponseCodeMoreBlocks) // 47-57
	shouldPublishBlockResponse(t, td.networkBob, 58, 11, message.ResponseCodeMoreBlocks) // 58-68
	shouldPublishBlockResponse(t, td.networkBob, 69, 1, message.ResponseCodeMoreBlocks)  // 69-69
	shouldPublishBlockResponse(t, td.networkBob, 0, 0, message.ResponseCodeNoMoreBlocks) // NoMoreBlock

	shouldPublishBlockRequest(t, td.networkAlice, 70)
	shouldPublishBlockResponse(t, td.networkBob, 70, 11, message.ResponseCodeMoreBlocks) // 70-80
	shouldPublishBlockResponse(t, td.networkBob, 81, 11, message.ResponseCodeMoreBlocks) // 81-91
	shouldPublishBlockResponse(t, td.networkBob, 92, 1, message.ResponseCodeMoreBlocks)  // 92-92
	shouldPublishBlockResponse(t, td.networkBob, 0, 0, message.ResponseCodeNoMoreBlocks) // NoMoreBlock

	// Last block requests
	shouldPublishBlockRequest(t, td.networkAlice, 93)                                   // 93-116
	shouldPublishBlockResponse(t, td.networkBob, 93, 8, message.ResponseCodeMoreBlocks) // 93-100
	shouldPublishBlockResponse(t, td.networkBob, 100, 0, message.ResponseCodeSynced)    // Synced

	assert.Eventually(t, func() bool {
		return td.syncAlice.state.LastBlockHeight() == uint32(100)
	}, 10*time.Second, 1*time.Second)
}

func TestSyncingHasBlockInCache(t *testing.T) {
	td := makeAliceAndBobNetworks(t)

	// Adding 100 blocks for Bob
	blockInterval := td.syncBob.state.Genesis().Params().BlockInterval()
	blockTime := td.syncBob.state.Genesis().GenesisTime()
	for i := uint32(0); i < 23; i++ {
		blk, cert := td.GenerateTestBlockWithTime(i+1, blockTime)
		assert.NoError(t, td.syncBob.state.CommitBlock(blk, cert))

		blockTime = blockTime.Add(blockInterval)
	}

	assert.Equal(t, uint32(0), td.syncAlice.state.LastBlockHeight())
	assert.Equal(t, uint32(23), td.syncBob.state.LastBlockHeight())

	// Adding some blocs to the cache
	blk1 := td.stateBob.TestStore.Blocks[1]
	blk2 := td.stateBob.TestStore.Blocks[2]
	blk3 := td.stateBob.TestStore.Blocks[3]
	td.syncAlice.cache.AddBlock(blk1)
	td.syncAlice.cache.AddBlock(blk2)
	td.syncAlice.cache.AddBlock(blk3)

	// Announcing a block
	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)
	td.syncBob.broadcast(msg)
	shouldPublishMessageWithThisType(t, td.networkBob, message.TypeBlockAnnounce)

	shouldNotPublishMessageWithThisType(t, td.networkBob, message.TypeBlocksRequest)
	// blocks 1-2 are inside the cache
	shouldPublishBlockRequest(t, td.networkAlice, 4)
	shouldPublishBlockResponse(t, td.networkBob, 4, 11, message.ResponseCodeMoreBlocks) // 4-14
	shouldPublishBlockResponse(t, td.networkBob, 15, 9, message.ResponseCodeMoreBlocks) // 15-23
	shouldPublishBlockResponse(t, td.networkBob, 23, 0, message.ResponseCodeSynced)     // Synced
}
