package message

import (
	"fmt"

	"github.com/pactus-project/pactus/util/errors"
)

type BlocksRequestMessage struct {
	SessionID int    `cbor:"1,keyasint"`
	From      uint32 `cbor:"2,keyasint"`
	Count     uint32 `cbor:"3,keyasint"`
}

func NewBlocksRequestMessage(sid int, from, count uint32) *BlocksRequestMessage {
	return &BlocksRequestMessage{
		SessionID: sid,
		From:      from,
		Count:     count,
	}
}

func (m *BlocksRequestMessage) To() uint32 {
	return m.From + m.Count - 1
}

func (m *BlocksRequestMessage) BasicCheck() error {
	if m.From == 0 {
		return errors.Errorf(errors.ErrInvalidHeight, "height is zero")
	}

	if m.Count == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "count is zero")
	}

	return nil
}

func (m *BlocksRequestMessage) Type() Type {
	return TypeBlocksRequest
}

func (m *BlocksRequestMessage) String() string {
	return fmt.Sprintf("{⚓ %d %v:%v}", m.SessionID, m.From, m.To())
}
