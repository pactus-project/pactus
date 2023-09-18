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
	TopicNewBlock       = uint16(0x0101)
	TopicNewTransaction = uint16(0x0201)
	TopicNewValidator   = uint16(0x0301)
	TopicNewAccount     = uint16(0x0401)
)

type Event []byte

// CreateBlockEvent creates an event when the new block is committed.
// The new block event structure is like :
// <topic_id><block_hash><height><sequence_number>.
func CreateBlockEvent(blockHash hash.Hash, height uint32) Event {
	buf := make([]byte, 0, 42)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElements(w, TopicNewBlock, blockHash, height)
	if err != nil {
		logger.Error("error on encoding event in new block", "error", err)
	}
	return w.Bytes()
}

// CreateNewTransactionEvent creates an event when a new transaction sent.
// The new transaction event structure is like :
// <topic_id><tx_hash><height><sequence_number>.
func CreateNewTransactionEvent(txHash tx.ID, height uint32) Event {
	buf := make([]byte, 0, 42)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElements(w, TopicNewTransaction, txHash, height)
	if err != nil {
		logger.Error("error on encoding event in new transaction", "error", err)
	}
	return w.Bytes()
}

// CreateValidatorEvent creates an event when the new Validator is created.
// The validator event structure is like :
// <topic_id><validator_address><height><sequence_number>.
func CreateValidatorEvent(validatorAddr crypto.Address, height uint32) Event {
	buf := make([]byte, 0, 42)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElements(w, TopicNewValidator, validatorAddr, height)
	if err != nil {
		logger.Error("error on encoding event in new validator", "error", err)
	}
	return w.Bytes()
}

// CreateAccountEvent creates an event when the new account is created.
// The account event structure is like :
// <topic_id><account_address><height><sequence_number>.
func CreateAccountEvent(accountAddr crypto.Address, height uint32) Event {
	buf := make([]byte, 0, 42)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElements(w, TopicNewAccount, accountAddr, height)
	if err != nil {
		logger.Error("error on encoding event in new account", "error", err)
	}
	return w.Bytes()
}
