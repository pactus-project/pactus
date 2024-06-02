package message

import (
	"fmt"
)

type HelloAckMessage struct {
	ResponseCode ResponseCode `cbor:"1,keyasint"`
	Reason       string       `cbor:"2,keyasint"`
	Height       uint32       `cbor:"3,keyasint"`
}

func NewHelloAckMessage(code ResponseCode, reason string, height uint32) *HelloAckMessage {
	return &HelloAckMessage{
		ResponseCode: code,
		Reason:       reason,
		Height:       height,
	}
}

func (m *HelloAckMessage) BasicCheck() error {
	return nil
}

func (m *HelloAckMessage) Type() Type {
	return TypeHelloAck
}

func (m *HelloAckMessage) String() string {
	return fmt.Sprintf("{%s: %s %v}", m.ResponseCode, m.Reason, m.Height)
}
