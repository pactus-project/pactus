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

func (m *VoteMessage) SanityCheck() error {
	return m.Vote.SanityCheck()
}

func (m *VoteMessage) Type() Type {
	return MessageTypeVote
}

func (m *VoteMessage) Fingerprint() string {
	return m.Vote.Fingerprint()
}
