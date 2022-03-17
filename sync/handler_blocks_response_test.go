package sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util"
)

func TestOneBlockShorter(t *testing.T) {
	setup(t)

	lastBlockHash := tState.LastBlockHash()
	lastBlockheight := tState.LastBlockHeight()
	b1 := block.GenerateTestBlock(nil, &lastBlockHash)
	c1 := block.GenerateTestCertificate(b1.Hash())
	pid := util.RandomPeerID()

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid)

	t.Run("Peer is busy. Session should be closed", func(t *testing.T) {
		sid := tSync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeBusy, sid, 0, nil, nil)
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		assert.Nil(t, tSync.peerSet.FindSession(sid))
	})

	t.Run("Request is rejected. Session should be closed", func(t *testing.T) {
		sid := tSync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeRejected, sid, 0, nil, nil)
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		assert.Nil(t, tSync.peerSet.FindSession(sid))
	})

	t.Run("Commit one block", func(t *testing.T) {
		sid := tSync.peerSet.OpenSession(pid).SessionID()
		msg := message.NewBlocksResponseMessage(message.ResponseCodeSynced, sid, lastBlockheight+1, []*block.Block{b1}, c1)
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		assert.Nil(t, tSync.peerSet.FindSession(sid))
		assert.Equal(t, tState.LastBlockHeight(), lastBlockheight+1)
	})
}

// TestSyncing is an important test and try to test syncing process between two test nodes (Alice and Bob).
// In the real situation, more nodes are involved, but the procedure is almost the same.
func TestSyncing(t *testing.T) {
	configAlice := TestConfig()
	configBob := TestConfig()
	signerAlice := bls.GenerateTestSigner()
	signerBob := bls.GenerateTestSigner()
	stateAlice := state.MockingState()
	stateBob := state.MockingState()
	consensusAlice := consensus.MockingConsensus(stateAlice)
	consensusBob := consensus.MockingConsensus(stateBob)
	broadcastChAlice := make(chan message.Message, 1000)
	broadcastChBob := make(chan message.Message, 1000)
	networkAlice := network.MockingNetwork(util.RandomPeerID())
	networkBob := network.MockingNetwork(util.RandomPeerID())

	LatestBlockInterval = 30
	configBob.NodeNetwork = true
	networkAlice.AddAnotherNetwork(networkBob)
	networkBob.AddAnotherNetwork(networkAlice)
	stateBob.TestGenHash = stateAlice.TestGenHash
	testAddBlocks(t, stateBob, 100)

	sync1, err := NewSynchronizer(configAlice,
		signerAlice,
		stateAlice,
		consensusAlice,
		networkAlice,
		broadcastChAlice,
	)
	assert.NoError(t, err)
	syncAlice := sync1.(*synchronizer)

	sync2, err := NewSynchronizer(configBob,
		signerBob,
		stateBob,
		consensusBob,
		networkBob,
		broadcastChBob,
	)
	assert.NoError(t, err)
	syncBob := sync2.(*synchronizer)

	// -------------------------------
	// For better logging when testing
	overrideLogger := func(sync *synchronizer, name string) {
		sync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: fmt.Sprintf("%s - %s: ", name, t.Name()), sync: sync})
	}

	overrideLogger(syncAlice, "Alice")
	overrideLogger(syncBob, "Bob")
	// -------------------------------

	assert.NoError(t, syncAlice.Start())
	assert.NoError(t, syncBob.Start())

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeHello)

	// Hello-ack
	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeHello)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeHello)

	assert.Len(t, syncAlice.Peers(), 1)
	assert.Len(t, syncBob.Peers(), 1)

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 1-10
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 11-20
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 21-30
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 31-40
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 41-50
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 51-60
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 61-70
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 71-80
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 81-90
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // NoMoreBlock

	// Latest block requests
	shouldPublishMessageWithThisType(t, networkAlice, message.MessageTypeBlocksRequest)
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // 91-100
	shouldPublishMessageWithThisType(t, networkBob, message.MessageTypeBlocksResponse) // Synced
}
