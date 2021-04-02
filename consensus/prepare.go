package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrepare(round int) {
	if cs.isPrepared || round > cs.hrs.Round() {
		cs.logger.Trace("Prepare: Precommitted, prepared or invalid round/step", "round", round)
		return
	}
	cs.hrs.UpdateStep(hrs.StepTypePrepare)
	cs.scheduleTimeout(cs.config.PrecommitTimeout(round), cs.hrs.Height(), round, hrs.StepTypePrecommit)

	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Warn("Prepare: No proposal")
		cs.signAddVote(vote.VoteTypePrepare, round, crypto.UndefHash)
		return
	}

	// Everything is good
	cs.isPrepared = true
	cs.logger.Info("Prepare: Proposal approved", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrepare, round, roundProposal.Block().Hash())
}
