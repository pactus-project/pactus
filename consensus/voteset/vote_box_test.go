package voteset

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	signer := ts.RandValAddress()
	power := ts.RandInt64(1000)

	v := vote.NewPrepareVote(hash, height, round, signer)

	vb := newVoteBox()

	vb.addVote(v, power)
	vb.addVote(v, power)

	assert.Equal(t, vb.votedPower, power)
}
