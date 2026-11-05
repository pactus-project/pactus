package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainVoteNoDoesNotCrashHonestNode checks that a justified `cp:MAIN-VOTE:no`
// vote sent by a single Byzantine validator cannot make an honest node panic.
//
// Honest v2 validators never emit a main-vote for `no` (they go to the precommit step
// on a quorum of `no` pre-votes), but cpCheckJustMainVoteNoConflict accepts it when it
// carries a valid certificate of `no` pre-votes. When the honest node X gets 2f+1 main-votes
// that are not a quorum of `yes`, it moves to cpRound 1 and cpPreVoteState.decideNextRounds()
// calls logger.Panic ("unreachable state") if ANY main-vote of the previous cpRound is `no`.
//
// Network: X (honest, no proposal seen), Y and B (honest, they pre-voted `no`), P (Byzantine).
func TestMainVoteNoDoesNotCrashHonestNode(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1
	height := types.Height(2)
	round := types.Round(0)

	td.enterNewHeight(td.consX)

	blockHash := td.RandHash()

	// Y, B and P signed a precommit for blockHash (2f+1 = 3 signatures, X is absent).
	precommitCert := td.precommitCert(t, blockHash, height, round, tIndexY, tIndexB, tIndexP)
	justNo := &vote.JustInitNo{QCert: precommitCert}

	// Y and B pre-vote `no`. The pre-vote of P is withheld from X.
	td.send(t, td.consX, vote.NewCPPreVote(blockHash, height, round, 0,
		vote.CPValueNo, justNo, td.valKeys[tIndexY].Address()), tIndexY)
	td.send(t, td.consX, vote.NewCPPreVote(blockHash, height, round, 0,
		vote.CPValueNo, justNo, td.valKeys[tIndexB].Address()), tIndexB)

	// X has seen no proposal: it pre-votes `yes`, sees {yes, no, no} and main-votes `abstain`.
	td.changeProposerTimeout(td.consX)
	td.shouldPublishVote(t, td.consX, vote.VoteTypeCPMainVote, blockHash)
	require.Equal(t, "cp:decide", td.consX.currentState.name())

	// P (Byzantine) reveals a valid certificate of 3 `no` pre-votes (Y, B, P) in a main-vote.
	preVoteNoCert := td.preVoteCert(t, blockHash, height, round, 0, vote.CPValueNo,
		tIndexY, tIndexB, tIndexP)
	mainVoteNo := vote.NewCPMainVote(blockHash, height, round, 0, vote.CPValueNo,
		&vote.JustMainVoteNoConflict{QCert: preVoteNoCert}, td.valKeys[tIndexP].Address())
	td.send(t, td.consX, mainVoteNo, tIndexP)

	// Y (honest) saw {no, no, yes} before receiving the pre-vote of P: it main-votes `abstain`.
	justYes := &vote.JustInitYes{}
	mainVoteAbstain := vote.NewCPMainVote(blockHash, height, round, 0, vote.CPValueAbstain,
		&vote.JustMainVoteConflict{JustNo: justNo, JustYes: justYes}, td.valKeys[tIndexY].Address())

	// The third main-vote gives X a 2f+1 set without quorum: X moves to cpRound 1.
	assert.NotPanics(t, func() {
		td.send(t, td.consX, mainVoteAbstain, tIndexY)
	}, "an honest node must not panic because of a justified main-vote for `no`")
}
