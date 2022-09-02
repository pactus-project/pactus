package event

import (
	"bytes"

	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/logger"
)

const TopicNewBlock = uint16(0x0101)
const TopicNewTransaction = uint16(0x0201)

type Event []byte

// CreateBlockEvent creates an event when the new block is committed.
// The block event structure is like :
// <topic_id><block_hash><height><sequence_number>
func CreateBlockEvent(blockHash hash.Hash, height uint32) Event {
	buf := make([]byte, 0, 42)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElements(w, TopicNewBlock, blockHash, height)
	if err != nil {
		logger.Error("error on encoding event", "err", err)
	}
	return w.Bytes()
}
