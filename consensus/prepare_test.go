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
	tConsX.lk.Lock()
	s.timedout(&ticker{Height: 2, Target: tickerTargetPrecommit})
	tConsX.lk.Unlock()
	assert.False(t, s.hasTimedout)

	tConsX.lk.Lock()
	s.timedout(&ticker{Height: 2, Target: tickerTargetPrepare})
	tConsX.lk.Unlock()
	assert.True(t, s.hasTimedout)

	// Add votes calls execute
	v, _ := vote.GenerateTestPrepareVote(2, 0)
	tConsX.lk.Lock()
	s.voteAdded(v)
	tConsX.lk.Unlock()
	shouldPublishVote(t, tConsY, vote.VoteTypePrepare, crypto.UndefHash)
}

func TestPrepareNullVote(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	shouldPublishQueryProposal(t, tConsP, 2, 0)

	s := &prepareState{tConsP, false}
	tConsX.lk.Lock()
	s.vote()
	tConsX.lk.Unlock()
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)
}
