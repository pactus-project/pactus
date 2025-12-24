package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
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
		return BasicCheckError{Reason: "invalid height"}
	}
	if m.Count == 0 {
		return BasicCheckError{Reason: "count is zero"}
	}

	return nil
}

func (*BlocksRequestMessage) Type() Type {
	return TypeBlocksRequest
}

func (*BlocksRequestMessage) TopicID() network.TopicID {
	return network.TopicIDUnspecified
}

func (*BlocksRequestMessage) ShouldBroadcast() bool {
	return false
}

func (*BlocksRequestMessage) ConsensusHeight() uint32 {
	return 0
}

// LogString returns a concise string representation intended for use in logs.
func (m *BlocksRequestMessage) LogString() string {
	return fmt.Sprintf("{⚓ %d %v:%v}", m.SessionID, m.From, m.To())
}
