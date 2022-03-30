package message

import (
	"github.com/zarbchain/zarb-go/consensus/vote"
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
	if err := m.Vote.SanityCheck(); err != nil {
		return err
	}
	return nil
}

func (m *VoteMessage) Type() Type {
	return MessageTypeVote
}

func (m *VoteMessage) Fingerprint() string {
	return m.Vote.Fingerprint()
}
