package sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

func TestOneBlockShorter(t *testing.T) {
	setup(t)

	lastBlockHash := tState.LastBlockHash()
	lastBlockheight := tState.LastBlockHeight()
	b1, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
	c1 := block.GenerateTestCertificate(b1.Hash())
	pid := util.RandomPeerID()

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid, peerset.StatusCodeKnown)
	sid := tSync.peerSet.OpenSession(pid).SessionID()

	t.Run("Request is rejected. Session should be closed", func(t *testing.T) {
		pld := payload.NewBlocksResponsePayload(payload.ResponseCodeRejected, sid, 0, nil, nil, nil)
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		assert.Nil(t, tSync.peerSet.FindSession(sid))
	})

	t.Run("Commit one block", func(t *testing.T) {
		pld := payload.NewBlocksResponsePayload(payload.ResponseCodeSynced, sid, lastBlockheight+1, []*block.Block{b1}, trxs, c1)
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		assert.Nil(t, tSync.peerSet.FindSession(sid))
		assert.Equal(t, tState.LastBlockHeight(), lastBlockheight+1)
	})
}

// TestSyncing try to test syncing process
func TestSyncing(t *testing.T) {
	configAlice := TestConfig()
	configBob := TestConfig()
	signerAlice := bls.GenerateTestSigner()
	signerBob := bls.GenerateTestSigner()
	committeeAlice, _ := committee.GenerateTestCommittee()
	committeeBob, _ := committee.GenerateTestCommittee()
	stateAlice := state.MockingState(committeeAlice)
	stateBob := state.MockingState(committeeBob)
	consensusAlice := consensus.MockingConsensus(stateAlice)
	consensusBob := consensus.MockingConsensus(stateBob)
	broadcastChAlice := make(chan payload.Payload, 1000)
	broadcastChBob := make(chan payload.Payload, 1000)
	networkAlice := network.MockingNetwork(util.RandomPeerID())
	networkBob := network.MockingNetwork(util.RandomPeerID())

	LatestBlockInterval = 30
	configBob.InitialBlockDownload = true
	networkAlice.AddAnotherNetwork(networkBob)
	networkBob.AddAnotherNetwork(networkAlice)
	stateBob.GenHash = stateAlice.GenHash
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

	shouldPublishPayloadWithThisType(t, networkAlice, payload.PayloadTypeHello)
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeHello)

	// Hello-ack
	shouldPublishPayloadWithThisType(t, networkAlice, payload.PayloadTypeHello)
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeHello)

	assert.Len(t, syncAlice.Peers(), 1)
	assert.Len(t, syncBob.Peers(), 1)

	shouldPublishPayloadWithThisType(t, networkAlice, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 1-10
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 11-20
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 21-30
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, networkAlice, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 31-40
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 41-50
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 51-60
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, networkAlice, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 61-70
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 71-80
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 81-90
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // NoMoreBlock

	// Latest block requests
	shouldPublishPayloadWithThisType(t, networkAlice, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // 91-100
	shouldPublishPayloadWithThisType(t, networkBob, payload.PayloadTypeBlocksResponse) // Synced
}
