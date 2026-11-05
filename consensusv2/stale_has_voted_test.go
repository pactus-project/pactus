package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

// (v2 dormant): precommitState.hasVoted is only reset to false
// in precommitState.enter(), which proposeState.decide() skips entirely when
// the round's proposer has an availability score below
// config.MinimumAvailabilityScore (PIP-19): it calls startChangingProposer()
// and returns before ever entering precommitState for this round.
//
// If the validator had hasVoted == true from an earlier round (it precommitted
// there), that stale value survives into the new round untouched. When the CP
// phase's cpPreVoteState.decideFirstRound() then reads it, it takes the "we
// already precommitted something this round" branch and waits for 2f+1 power
// of (cp pre-votes + precommits) for this round. Once three OTHER, non-stale
// validators correctly pre-vote "Yes" (their own hasVoted is false, so they
// take the honest branch), that quorum is reached purely from CP pre-vote
// power, with zero real precommits. decideFirstRound then dereferences
// s.log.RoundProposal(s.round).Block() — but no proposal exists for this
// round, since proposeState returned before creating or waiting for one. The
// result is a nil-pointer panic reachable by an entirely ordinary event (a
// low-availability-score proposer, plus validators that did not all
// precommit the same way last round), with no attacker needed.
//
// This test expresses the desired behavior (no panic; a Yes pre-vote instead,
// since nothing has actually been precommitted this round) and currently
// fails with a panic on the third injected vote, proving the bug.
func TestStalePrecommitHasVotedPanicsOnLowScoreProposer(t *testing.T) {
	td := setup(t)

	// Simulate: this validator precommitted something in an earlier round,
	// while Y, B and P (below) did not — e.g. they joined the committee too
	// late to vote on the previous block.
	td.consX.precommitState.hasVoted = true

	// testConfig() leaves MinimumAvailabilityScore at its zero value, so set
	// it to the production default (consensusv2/config.go DefaultConfig) for
	// this test, and make this round's proposer score fall below it: as in
	// production, proposeState then skips precommitState.enter() (and
	// therefore the hasVoted reset) and goes straight into the
	// change-proposer phase.
	td.consX.config.MinimumAvailabilityScore = 0.666667
	td.stateX.FakeScore = 0.1

	td.enterNewHeight(td.consX)
	td.checkHeightRound(t, td.consX, 1, 0)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decideFirstRound panicked on a stale hasVoted with no round proposal: %v", r)
		}
	}()

	// Y, B and P are fresh (hasVoted == false): each honestly casts a Yes
	// pre-vote at cpRound 0. With n=4 equal-power validators, 3 shares meet
	// the 2f+1 threshold (f=1) purely from CP pre-vote power.
	td.addCPPreVote(t, td.consX, hash.UndefHash, 1, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)
	td.addCPPreVote(t, td.consX, hash.UndefHash, 1, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexB)
	td.addCPPreVote(t, td.consX, hash.UndefHash, 1, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexP)
}
