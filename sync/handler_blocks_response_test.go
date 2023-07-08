package sync

import (
	"context"
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestInvalidBlockData(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandomPeerID()
	sid := td.sync.peerSet.OpenSession(pid).SessionID()
	msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, sid,
		0, [][]byte{{1, 2, 3}}, nil)

	assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))
}

func TestOneBlockShorter(t *testing.T) {
	td := setup(t, nil)

	lastBlockHash := td.state.LastBlockHash()
	lastBlockHeight := td.state.LastBlockHeight()
	b1 := td.GenerateTestBlock(nil, &lastBlockHash)
	c1 := td.GenerateTestCertificate(b1.Hash())
	d1, _ := b1.Bytes()
	pid := td.RandomPeerID()

	pub, _ := td.RandomBLSKeyPair()
	td.addPeer(t, pub, pid, false)

	t.Run("Peer is busy. Session should be closed", func(t *testing.T) {
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeBusy, sid,
			0, nil, nil)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		assert.Nil(t, td.sync.peerSet.FindSession(sid))
	})

	t.Run("Request is rejected. Session should be closed", func(t *testing.T) {
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeRejected, sid,
			0, nil, nil)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		assert.Nil(t, td.sync.peerSet.FindSession(sid))
	})

	t.Run("Commit one block", func(t *testing.T) {
		sid := td.sync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, sid,
			lastBlockHeight+1, [][]byte{d1}, c1)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		assert.Nil(t, td.sync.peerSet.FindSession(sid))
		assert.Equal(t, td.state.LastBlockHeight(), lastBlockHeight+1)
	})
}

// TestSyncing is an important test to verify the syncing process between two
// test nodes, Alice and Bob. In real-world scenarios, multiple nodes are typically
// involved, but the procedure remains similar.
func TestSyncing(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	configAlice := testConfig()
	configBob := testConfig()
	signersAlice := []crypto.Signer{ts.RandomSigner()}
	signersBob := []crypto.Signer{ts.RandomSigner()}
	stateAlice := state.MockingState(ts)
	stateBob := state.MockingState(ts)
	consMgrAlice, _ := consensus.MockingManager(ts, signersAlice)
	consMgrBob, _ := consensus.MockingManager(ts, signersBob)
	broadcastChAlice := make(chan message.Message, 1000)
	broadcastChBob := make(chan message.Message, 1000)
	networkAlice := network.MockingNetwork(ts, ts.RandomPeerID())
	networkBob := network.MockingNetwork(ts, ts.RandomPeerID())

	ctx := context.Background()

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
		ctx,
	)
	assert.NoError(t, err)
	syncAlice := sync1.(*synchronizer)

	sync2, err := NewSynchronizer(configBob,
		signersBob,
		stateBob,
		consMgrBob,
		networkBob,
		broadcastChBob,
		ctx,
	)
	assert.NoError(t, err)
	syncBob := sync2.(*synchronizer)

	// -------------------------------
	// For better logging when testing
	overrideLogger := func(sync *synchronizer, name string) {
		sync.logger = logger.NewLogger("_sync", &OverrideFingerprint{
			name: fmt.Sprintf("%s - %s: ", name, t.Name()), sync: sync})
	}

	overrideLogger(syncAlice, "Alice")
	overrideLogger(syncBob, "Bob")
	// -------------------------------

	assert.NoError(t, syncAlice.Start())
	assert.NoError(t, syncBob.Start())

	// Verify that Hello messages are exchanged between Alice and Bob
	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeHello)

	// Verify that Hello-ack messages are exchanged between Alice and Bob
	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeHello)

	// Ensure peers are connected and block heights are correct
	assert.Equal(t, syncAlice.PeerSet().Len(), 1)
	assert.Equal(t, syncBob.PeerSet().Len(), 1)
	assert.Equal(t, syncAlice.state.LastBlockHeight(), uint32(0))
	assert.Equal(t, syncBob.state.LastBlockHeight(), uint32(100))

	// Perform block syncing
	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 1-11
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 12-22
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 23-23
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 24-34
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 35-45
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 46-46
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 47-57
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 58-68
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 69-69
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 70-80
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 81-91
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 92-92
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	// Last block requests
	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse)        // 93-100
	bdl := shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // Synced

	assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
	// Alice needs more time to process all the bundles,
	// but the block height should be greater than zero
	assert.Greater(t, syncAlice.state.LastBlockHeight(), uint32(0))
	assert.Equal(t, syncBob.state.LastBlockHeight(), uint32(100))
}
