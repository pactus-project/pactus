package message

import (
	"fmt"
)

type QueryProposalMessage struct {
	Height uint32 `cbor:"1,keyasint"`
}

func NewQueryProposalMessage(h uint32) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height: h,
	}
}

func (m *QueryProposalMessage) BasicCheck() error {
	return nil
}

func (m *QueryProposalMessage) Type() Type {
	return TypeQueryProposal
}

func (m *QueryProposalMessage) String() string {
	return fmt.Sprintf("{%v}", m.Height)
}
