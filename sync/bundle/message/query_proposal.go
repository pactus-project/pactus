package message

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
)

type QueryProposalMessage struct {
	Height  uint32         `cbor:"1,keyasint"`
	Round   int16          `cbor:"3,keyasint"`
	Querier crypto.Address `cbor:"2,keyasint"`
}

func NewQueryProposalMessage(height uint32, round int16, querier crypto.Address) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height:  height,
		Round:   round,
		Querier: querier,
	}
}

func (m *QueryProposalMessage) BasicCheck() error {
	if m.Round < 0 {
		return BasicCheckError{Reason: "invalid round"}
	}

	return nil
}

func (*QueryProposalMessage) Type() Type {
	return TypeQueryProposal
}

func (*QueryProposalMessage) TopicID() network.TopicID {
	return network.TopicIDConsensus
}

func (*QueryProposalMessage) ShouldBroadcast() bool {
	return true
}

func (m *QueryProposalMessage) ConsensusHeight() uint32 {
	return m.Height
}

// LogString returns a concise string representation intended for use in logs.
func (m *QueryProposalMessage) LogString() string {
	return fmt.Sprintf("{%v %s}", m.Height, m.Querier.LogString())
}
