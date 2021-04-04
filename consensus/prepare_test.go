package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrepareTimedout(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsY)

	s := &prepareState{tConsY, false}

	// Invalid target
	s.timedout(&ticker{Height: 2, Target: tickerTargetPrecommit})
	assert.False(t, s.hasTimedout)

	s.timedout(&ticker{Height: 2, Target: tickerTargetPrepare})
	assert.True(t, s.hasTimedout)

	// Add votes calls execute
	v, _ := vote.GenerateTestPrepareVote(2, 0)
	s.voteAdded(v)
	shouldPublishVote(t, tConsY, vote.VoteTypePrepare, crypto.UndefHash)
}

func TestPrepareNullVote(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	shouldPublishQueryProposal(t, tConsP, 2, 0)

	s := &prepareState{tConsP, false}
	s.vote()
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)
}
