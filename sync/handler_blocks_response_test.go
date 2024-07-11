package sync

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
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
	blk, cert := td.GenerateTestBlock(lastHeight+1, testsuite.BlockWithPrevCert(nil))
	data, _ := blk.Bytes()
	tests := []struct {
		data []byte
	}{
		{data: td.RandBytes(16)},
		{data: data},
	}

	for _, test := range tests {
		pid := td.RandPeerID()
		sid := td.RandInt(1000)
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
			message.ResponseCodeMoreBlocks.String(),
			sid, lastHeight+1, [][]byte{test.data}, cert)

		td.receivingNewMessage(td.sync, msg, pid)
		assert.Nil(t, td.sync.cache.GetBlock(msg.From))
	}
}

func TestOneBlockShorter(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(10)

	lastHeight := td.state.LastBlockHeight()
	blk1, cert1 := td.GenerateTestBlock(lastHeight + 1)
	d1, _ := blk1.Bytes()
	pid := td.addPeer(t, status.StatusKnown, service.New(service.None))

	sid := td.RandInt(1000)
	msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, t.Name(), sid,
		lastHeight+1, [][]byte{d1}, cert1)
	td.receivingNewMessage(td.sync, msg, pid)

	assert.Equal(t, td.state.LastBlockHeight(), lastHeight+1)
}

func TestStrippedPublicKey(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(10)

	lastHeight := td.state.LastBlockHeight()

	// Add a new block and keep the signer key
	_, indexedPrv := td.RandBLSKeyPair()
	trx0 := td.GenerateTestTransferTx(testsuite.TransactionWithSigner(indexedPrv))
	trxs0 := []*tx.Tx{trx0}
	blk0, cert0 := td.GenerateTestBlock(lastHeight+1, testsuite.BlockWithTransactions(trxs0))
	err := td.state.CommitBlock(blk0, cert0)
	require.NoError(t, err)
	lastHeight++
	// -----

	_, rndPrv := td.RandBLSKeyPair()
	trx1 := td.GenerateTestTransferTx(testsuite.TransactionWithSigner(rndPrv))
	trx1.StripPublicKey()
	trxs1 := []*tx.Tx{trx1}
	blk1, _ := td.GenerateTestBlock(lastHeight+1, testsuite.BlockWithTransactions(trxs1))

	trx2 := td.GenerateTestTransferTx(testsuite.TransactionWithSigner(indexedPrv))
	trx2.StripPublicKey()
	trxs2 := []*tx.Tx{trx2}
	blk2, _ := td.GenerateTestBlock(lastHeight+1, testsuite.BlockWithTransactions(trxs2))

	tests := []struct {
		receivedBlock *block.Block
		shouldFail    bool
	}{
		{
			receivedBlock: blk1,
			shouldFail:    true,
		},
		{
			receivedBlock: blk2,
			shouldFail:    false,
		},
	}

	// Add a peer
	pid := td.addPeer(t, status.StatusKnown, service.New(service.None))

	for _, tc := range tests {
		blkData, _ := tc.receivedBlock.Bytes()
		sid := td.RandInt(1000)
		cert := td.GenerateTestBlockCertificate(lastHeight + 1)
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, message.ResponseCodeMoreBlocks.String(), sid,
			lastHeight+1, [][]byte{blkData}, cert)
		td.receivingNewMessage(td.sync, msg, pid)

		if tc.shouldFail {
			assert.Nil(t, td.sync.cache.GetBlock(msg.From))
		} else {
			assert.NotNil(t, td.sync.cache.GetBlock(msg.From))
		}
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
	consMgrAlice, _ := consensus.MockingManager(ts, stateAlice, valKeyAlice)
	consMgrBob, _ := consensus.MockingManager(ts, stateBob, valKeyBob)
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
	}, 2*time.Second, 100*time.Millisecond)

	require.Equal(t, status.StatusKnown, syncAlice.PeerSet().GetPeerStatus(syncBob.SelfID()))
	require.Equal(t, status.StatusKnown, syncBob.PeerSet().GetPeerStatus(syncAlice.SelfID()))

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
		blk, cert := td.GenerateTestBlock(i+1, testsuite.BlockWithTime(blockTime))
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
		blk, cert := td.GenerateTestBlock(i+1, testsuite.BlockWithTime(blockTime))
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
