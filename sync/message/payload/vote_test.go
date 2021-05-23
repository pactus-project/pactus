package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

func TestVoteType(t *testing.T) {
	p := &VotePayload{}
	assert.Equal(t, p.Type(), PayloadTypeVote)
}

func TestVotePayload(t *testing.T) {
	t.Run("Invalid vote", func(t *testing.T) {
		v, _ := vote.GenerateTestPrepareVote(100, -1)
		p := NewVotePayload(v)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		v, _ := vote.GenerateTestPrepareVote(100, 0)
		p := NewVotePayload(v)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), v.Fingerprint())
	})
}
