package sync

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestInvalidBlockData(t *testing.T) {
	td := setup(t, nil)

	blk := block.MakeBlock(1, time.Now(), nil, td.RandHash(), td.RandHash(),
		td.GenerateTestCertificate(), td.RandSeed(), td.RandAddress())
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
			td.RandHeight(), [][]byte{test.data}, nil)

		err := td.receivingNewMessage(td.sync, msg, pid)
		assert.ErrorIs(t, err, test.err)
	}
}

func TestOneBlockShorter(t *testing.T) {
	td := setup(t, nil)

	lastBlockHeight := td.state.LastBlockHeight()
	b1 := td.GenerateTestBlock(nil)
	c1 := td.GenerateTestCertificate()
	d1, _ := b1.Bytes()
	pid := td.RandPeerID()

	pub, _ := td.RandBLSKeyPair()
	td.addPeer(t, pub, pid, false)

	t.Run("Peer is busy. Session should be closed", func(t *testing.T) {
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeRejected, message.ResponseCodeRejected.String(), sid,
			0, nil, nil)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		assert.Nil(t, td.sync.peerSet.FindSession(sid))
	})

	t.Run("Request is rejected. Session should be closed", func(t *testing.T) {
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeRejected, message.ResponseCodeRejected.String(), sid,
			0, nil, nil)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		assert.Nil(t, td.sync.peerSet.FindSession(sid))
	})

	t.Run("Commit one block", func(t *testing.T) {
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, message.ResponseCodeRejected.String(), sid,
			lastBlockHeight+1, [][]byte{d1}, c1)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		assert.Nil(t, td.sync.peerSet.FindSession(sid))
		assert.Equal(t, td.state.LastBlockHeight(), lastBlockHeight+1)
	})
}

func TestStrippedPublicKey(t *testing.T) {
	td := setup(t, nil)

	// Add a peer
	pid := td.RandPeerID()
	pub, _ := td.RandBLSKeyPair()
	td.addPeer(t, pub, pid, false)

	blk1 := td.GenerateTestBlock(nil)
	trx := *td.state.TestStore.Blocks[1].Transactions()[0]
	trxs := []*tx.Tx{&trx}
	blk2 := block.MakeBlock(1, time.Now(), trxs, td.RandHash(), td.RandHash(),
		td.GenerateTestCertificate(), td.RandSeed(), td.RandAddress())

	tests := []struct {
		blk *block.Block
		err error
	}{
		{
			blk1,
			store.PublicKeyNotFoundError{
				Address: blk1.Transactions()[0].PublicKey().Address(),
			},
		},
		{
			blk2,
			nil,
		},
	}

	for _, test := range tests {
		assert.NoError(t, test.blk.BasicCheck())
		trx1 := test.blk.Transactions()[0]
		trx1.StripPublicKey()
		d1, _ := test.blk.Bytes()
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, message.ResponseCodeRejected.String(), sid,
			td.RandHeight(), [][]byte{d1}, nil)
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
	signersAlice := []crypto.Signer{ts.RandSigner()}
	signersBob := []crypto.Signer{ts.RandSigner()}
	stateAlice := state.MockingState(ts)
	stateBob := state.MockingState(ts)
	consMgrAlice, _ := consensus.MockingManager(ts, signersAlice)
	consMgrBob, _ := consensus.MockingManager(ts, signersBob)
	broadcastChAlice := make(chan message.Message, 1000)
	broadcastChBob := make(chan message.Message, 1000)
	networkAlice := network.MockingNetwork(ts, ts.RandPeerID())
	networkBob := network.MockingNetwork(ts, ts.RandPeerID())

	configBob.NodeNetwork = true
	networkAlice.AddAnotherNetwork(networkBob)
	networkBob.AddAnotherNetwork(networkAlice)
	addBlocks(t, stateBob, 100)

	sync1, err := NewSynchronizer(configAlice,
		signersAlice,
		stateAlice,
		consMgrAlice,
		networkAlice,
		broadcastChAlice,
	)
	assert.NoError(t, err)
	syncAlice := sync1.(*synchronizer)

	sync2, err := NewSynchronizer(configBob,
		signersBob,
		stateBob,
		consMgrBob,
		networkBob,
		broadcastChBob,
	)
	assert.NoError(t, err)
	syncBob := sync2.(*synchronizer)

	// -------------------------------
	// For better logging when testing
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

	syncAlice.sayHello(false, syncBob.SelfID())

	// Verify that Hello messages are exchanged between Alice and Bob
	shouldPublishMessageWithThisType(t, networkAlice, message.TypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.TypeHello)

	// Ensure peers are connected and block heights are correct
	assert.Equal(t, syncAlice.PeerSet().Len(), 1)
	assert.Equal(t, syncBob.PeerSet().Len(), 1)
	assert.Equal(t, syncAlice.state.LastBlockHeight(), uint32(0))
	assert.Equal(t, syncBob.state.LastBlockHeight(), uint32(100))

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
	assert.Greater(t, syncAlice.state.LastBlockHeight(), uint32(0))
	assert.Equal(t, syncBob.state.LastBlockHeight(), uint32(100))
}
