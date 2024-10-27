package event

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/logger"
)

const (
	TopicBlock         = uint16(0x0101)
	TopicTransaction   = uint16(0x0201)
	TopicAccountChange = uint16(0x0301)
)

type Event []byte

// CreateBlockEvent creates an event when the new block is committed.
// The block event structure is like :
// <topic_id><block_hash><height><sequence_number>.
func CreateBlockEvent(blockHash hash.Hash, height uint32) Event {
	buf := bytes.NewBuffer(make([]byte, 0, 42))
	err := encoding.WriteElements(buf, TopicBlock, blockHash, height)
	if err != nil {
		logger.Error("error on encoding event in new block", "error", err)
	}

	return buf.Bytes()
}

// CreateTransactionEvent creates an event when a new transaction sent.
// The new transaction event structure is like :
// <topic_id><tx_hash><height><sequence_number>.
func CreateTransactionEvent(txHash tx.ID, height uint32) Event {
	buf := bytes.NewBuffer(make([]byte, 0, 42))
	err := encoding.WriteElements(buf, TopicTransaction, txHash, height)
	if err != nil {
		logger.Error("error on encoding event in new transaction", "error", err)
	}

	return buf.Bytes()
}

// CreateAccountChangeEvent creates an event when the new account is created.
// The account event structure is like :
// <topic_id><account_address><height><sequence_number>.
func CreateAccountChangeEvent(accountAddr crypto.Address, height uint32) Event {
	buf := bytes.NewBuffer(make([]byte, 0, 42))
	err := encoding.WriteElements(buf, TopicAccountChange, accountAddr, height)
	if err != nil {
		logger.Error("error on encoding event in new account", "error", err)
	}

	return buf.Bytes()
}
