package event

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/logger"
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

// CreateBlockEvent creates an event when the new block is committed.
// The block event structure is like :
// <topic_id><tx_hash><height><sequence_number>
func CreateNewTransactionEvent(txHash tx.ID, height uint32) Event {
	buf := make([]byte, 0, 42)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElements(w, TopicNewTransaction, txHash, height)
	if err != nil {
		logger.Error("error on encoding event in transaction event", "err", err)
	}
	return w.Bytes()
}
