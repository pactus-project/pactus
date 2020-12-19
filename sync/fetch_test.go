package sync

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
)

func TestRequestForBlocksInvalidLastBlocHash(t *testing.T) {
	setup(t)

	invHash := crypto.GenerateTestHash()

	// Send block request, but block hash is invalid, ignore it
	msg := message.NewBlocksReqMessage(7, 12, invHash)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)

	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, message.PayloadTypeBlocks)
}

func TestRequestForBlocks(t *testing.T) {
	setup(t)

	h := tState.Store.Blocks[7].Header().LastBlockHash()
	msg := message.NewBlocksReqMessage(7, 11, h)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)

	blocks := make([]*block.Block, 0)
	for i := 7; i <= 11; i++ {
		blocks = append(blocks, tState.Store.Blocks[i])
	}

	expectedMsg := message.NewBlocksMessage(7, blocks, nil)
	tNetAPI.waitingForMessage(t, expectedMsg)
}

func TestRequestForBlocksWithLastCommist(t *testing.T) {
	setup(t)

	h := tState.Store.Blocks[7].Header().LastBlockHash()
	msg := message.NewBlocksReqMessage(7, tState.LastBlockHeight(), h)
	data, _ := cbor.Marshal(msg)
	tSync.ParsMessage(data, tOurID)

	blocks := make([]*block.Block, 0)
	for i := 7; i <= tState.LastBlockHeight(); i++ {
		blocks = append(blocks, tState.Store.Blocks[i])
	}

	assert.NotNil(t, tState.LastBlockCommit)
	expectedMsg := message.NewBlocksMessage(7, blocks, tState.LastBlockCommit)
	tNetAPI.waitingForMessage(t, expectedMsg)
}
