package zmq

import (
	"encoding/binary"
	"fmt"
)

type Topic int16

const (
	TopicBlockInfo       Topic = 0x0001
	TopicTransactionInfo Topic = 0x0002
	TopicRawBlock        Topic = 0x0003
	TopicRawTransaction  Topic = 0x0004
	TopicAnchorInfo      Topic = 0x0005
)

func (t Topic) String() string {
	switch t {
	case TopicBlockInfo:
		return "block_info"

	case TopicTransactionInfo:
		return "transaction_info"

	case TopicRawBlock:
		return "raw_block"

	case TopicRawTransaction:
		return "raw_transaction"

	case TopicAnchorInfo:
		return "anchor_info"

	default:
		return fmt.Sprintf("topic-%d", t)
	}
}

func (t Topic) Bytes() []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(t))

	return b
}

func TopicFromBytes(b []byte) Topic {
	if len(b) < 2 {
		return 0
	}

	return Topic(binary.BigEndian.Uint16(b))
}
