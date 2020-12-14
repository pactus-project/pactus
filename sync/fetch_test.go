package sync

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/message"
)

func TestSendBlocks(t *testing.T) {
	setup(t)
	sync, api, st := newTestSynchronizer(nil)

	len := 12 //len(validBlocks)
	// Update state
	for i := 0; i < len; i++ {
		err := st.ApplyBlock(i+1, validBlocks[i], validCommits[i])
		assert.NoError(t, err)
	}

	assert.NoError(t, sync.Start())

	// Stopping HeartBeat ticker
	sync.heartBeatTicker.Stop()

	expectedMsg := message.NewSalamMessage(genDoc.Hash(), len)
	api.waitingForMessage(t, expectedMsg)

	msg := message.NewSalamMessage(genDoc.Hash(), 0)
	data, _ := cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	expectedMsg = message.NewBlocksMessage(1, validBlocks[0:11], nil)
	api.waitingForMessage(t, expectedMsg)

	// Send block request, but block hash is invalid, ignore it
	msg = message.NewBlocksReqMessage(7, len, validBlocks[5].Hash())
	data, _ = cbor.Marshal(msg)
	sync.ParsMessage(data, peerID)

	expectedMsg = message.NewBlocksMessage(7, validBlocks[6:12], &validCommits[11])
	api.waitingForMessage(t, expectedMsg)
}
