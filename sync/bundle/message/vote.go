package message

import (
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/types/vote"
)

type VoteMessage struct {
	Vote *vote.Vote `cbor:"1,keyasint"`
}

func NewVoteMessage(v *vote.Vote) *VoteMessage {
	return &VoteMessage{
		Vote: v,
	}
}

func (m *VoteMessage) BasicCheck() error {
	return m.Vote.BasicCheck()
}

func (*VoteMessage) Type() Type {
	return TypeVote
}

func (*VoteMessage) TopicID() network.TopicID {
	return network.TopicIDConsensus
}

func (*VoteMessage) ShouldBroadcast() bool {
	return true
}

func (m *VoteMessage) ConsensusHeight() uint32 {
	return m.Vote.Height()
}

// LogString returns a concise string representation intended for use in logs.
func (m *VoteMessage) LogString() string {
	return m.Vote.LogString()
}
