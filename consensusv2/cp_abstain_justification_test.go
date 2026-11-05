package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

// TestPreVoteAbstainWithConflictJustification (variant of the known issue "cpCheckJust
// checks the justification type, not the vote type"): a `cp:PRE-VOTE` carrying the
// justification of a main-vote (JustMainVoteConflict, value `abstain`) is accepted and stored in the
// `abstain` box of the pre-vote set. When the 2f+1 pre-votes of the node contain neither
// a `no` quorum nor a `yes` quorum, cpMainVoteState.decide() takes the "no quorum" branch and calls
// vote0.CPJust() on the result of GetRandomVote(cpRound, no), which is nil when no `no`
// pre-vote exists: nil pointer dereference.
func TestPreVoteAbstainWithConflictJustification(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1
	height := types.Height(2)
	td.enterNewHeight(td.consX)

	blockHash := td.RandHash()
	justNo := &vote.JustInitNo{QCert: td.precommitCert(t, blockHash, height, 0, tIndexY, tIndexB, tIndexP)}
	conflict := &vote.JustMainVoteConflict{JustNo: justNo, JustYes: &vote.JustInitYes{}}

	// Y pre-votes yes, B (Byzantine) pre-votes `abstain` with a valid main-vote justification.
	td.addCPPreVote(t, td.consX, hash.UndefHash, height, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)
	byz := vote.NewCPPreVote(blockHash, height, 0, 0, vote.CPValueAbstain, conflict, td.valKeys[tIndexB].Address())

	assert.NotPanics(t, func() {
		td.send(t, td.consX, byz, tIndexB)
		td.changeProposerTimeout(td.consX) // X pre-votes yes: 3 pre-votes, no quorum
	})
}
