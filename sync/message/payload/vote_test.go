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
		p1 := NewVotePayload(v)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		v, _ := vote.GenerateTestPrepareVote(100, 0)
		p2 := NewVotePayload(v)
		assert.NoError(t, p2.SanityCheck())
	})
}
