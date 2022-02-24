package message

import (
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/errors"
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
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}
	return nil
}

func (m *VoteMessage) Type() Type {
	return MessageTypeVote
}

func (m *VoteMessage) Fingerprint() string {
	return m.Vote.Fingerprint()
}
