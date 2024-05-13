package message

import (
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

func (m *VoteMessage) String() string {
	return m.Vote.String()
}
