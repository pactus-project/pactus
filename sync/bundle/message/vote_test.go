package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestVoteType(t *testing.T) {
	m := &VoteMessage{}
	assert.Equal(t, m.Type(), MessageTypeVote)
}

func TestVoteMessage(t *testing.T) {
	t.Run("Invalid vote", func(t *testing.T) {
		v, _ := vote.GenerateTestPrepareVote(100, -1)
		m := NewVoteMessage(v)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		v, _ := vote.GenerateTestPrepareVote(100, 0)
		m := NewVoteMessage(v)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), v.Fingerprint())
	})
}
