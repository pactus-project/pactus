package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrepare(round int) {
	if cs.isPreCommitted || cs.isPrepared || round != cs.hrs.Round() {
		cs.logger.Debug("Prepare: Precommitted, prepared or invalid round/step", "round", round)
		return
	}
	cs.updateStep(hrs.StepTypePrepare)

	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Warn("Prepare: No proposal")
		cs.signAddVote(vote.VoteTypePrepare, crypto.UndefHash)
		return
	}

	// Note:after receiving proposal we can update the vote

	// Everything is good
	cs.isPrepared = true
	cs.logger.Info("Prepare: Proposal is validated", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrepare, roundProposal.Block().Hash())
}

func (cs *consensus) enterPrepareWait(round int) {
	cs.logger.Info("PrepareWait: Wait for some more prepares") //then enter precommit

	cs.scheduleTimeout(cs.config.PrepareTimeout(round), cs.hrs.Height(), round, hrs.StepTypePrecommit)
}
