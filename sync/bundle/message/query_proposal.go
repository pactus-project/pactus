package message

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
)

type QueryProposalMessage struct {
	Height  uint32         `cbor:"1,keyasint"`
	Querier crypto.Address `cbor:"2,keyasint"`
}

func NewQueryProposalMessage(height uint32, querier crypto.Address) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height:  height,
		Querier: querier,
	}
}

func (m *QueryProposalMessage) BasicCheck() error {
	return nil
}

func (m *QueryProposalMessage) Type() Type {
	return TypeQueryProposal
}

func (m *QueryProposalMessage) String() string {
	return fmt.Sprintf("{%v %s}", m.Height, m.Querier.ShortString())
}
