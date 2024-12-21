package message

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
)

type QueryVoteMessage struct {
	Height  uint32         `cbor:"1,keyasint"`
	Round   int16          `cbor:"2,keyasint"`
	Querier crypto.Address `cbor:"3,keyasint"`
}

func NewQueryVoteMessage(height uint32, round int16, querier crypto.Address) *QueryVoteMessage {
	return &QueryVoteMessage{
		Height:  height,
		Round:   round,
		Querier: querier,
	}
}

func (m *QueryVoteMessage) BasicCheck() error {
	if m.Round < 0 {
		return BasicCheckError{Reason: "invalid round"}
	}

	return nil
}

func (*QueryVoteMessage) Type() Type {
	return TypeQueryVote
}

func (*QueryVoteMessage) TopicID() network.TopicID {
	return network.TopicIDConsensus
}

func (*QueryVoteMessage) ShouldBroadcast() bool {
	return true
}

func (m *QueryVoteMessage) ConsensusHeight() uint32 {
	return m.Height
}

func (m *QueryVoteMessage) String() string {
	return fmt.Sprintf("{%d/%d %s}", m.Height, m.Round, m.Querier.ShortString())
}
