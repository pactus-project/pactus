package fastconsensus

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type changeProposer struct {
	*consensus
}

func (cp *changeProposer) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (cp *changeProposer) onTimeout(t *ticker) {
	if t.Target == tickerTargetQueryVotes {
		cp.queryVotes()
		cp.scheduleTimeout(t.Duration*2, cp.height, cp.round, tickerTargetQueryVotes)
	}
}

func (cp *changeProposer) cpCheckCPValue(vte *vote.Vote, allowedValues ...vote.CPValue) error {
	for _, v := range allowedValues {
		if vte.CPValue() == v {
			return nil
		}
	}

	return invalidJustificationError{
		JustType: vte.CPJust().Type(),
		Reason:   fmt.Sprintf("invalid value: %v", vte.CPValue()),
	}
}

func (cp *changeProposer) cpCheckJustInitZero(just vote.Just, blockHash hash.Hash) error {
	j, ok := just.(*vote.JustInitNo)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidatePrepare(cp.validators, blockHash)
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustInitOne(just vote.Just) error {
	_, ok := just.(*vote.JustInitYes)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustPreVoteHard(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustPreVoteHard)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidateCPPreVote(cp.validators,
		blockHash, cpRound-1, byte(cpValue))
	if err != nil {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustPreVoteSoft(just vote.Just,
	blockHash hash.Hash, cpRound int16,
) error {
	j, ok := just.(*vote.JustPreVoteSoft)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidateCPMainVote(cp.validators,
		blockHash, cpRound-1, byte(vote.CPValueAbstain))
	if err != nil {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustMainVoteNoConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustMainVoteNoConflict)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidateCPPreVote(cp.validators,
		blockHash, cpRound, byte(cpValue))
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) cpCheckJustMainVoteConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16,
) error {
	j, ok := just.(*vote.JustMainVoteConflict)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	if cpRound == 0 {
		err := cp.cpCheckJustInitZero(j.Just0, blockHash)
		if err != nil {
			return err
		}

		err = cp.cpCheckJustInitOne(j.Just1)
		if err != nil {
			return err
		}

		return nil
	}

	// Just0 can be for Zero or Abstain values.
	switch j.Just0.Type() {
	case vote.JustTypePreVoteSoft:
		err := cp.cpCheckJustPreVoteSoft(j.Just0, blockHash, cpRound)
		if err != nil {
			return err
		}
	case vote.JustTypePreVoteHard:
		err := cp.cpCheckJustPreVoteHard(j.Just0, blockHash, cpRound, vote.CPValueNo)
		if err != nil {
			return err
		}
	default:
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("unexpected justification: %s", j.Just0.Type()),
		}
	}

	err := cp.cpCheckJustPreVoteHard(j.Just1, hash.UndefHash, cpRound, vote.CPValueYes)
	if err != nil {
		return err
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) cpCheckJustPreVote(v *vote.Vote) error {
	just := v.CPJust()
	if v.CPRound() == 0 {
		switch just.Type() {
		case vote.JustTypeInitZero:
			err := cp.cpCheckCPValue(v, vote.CPValueNo)
			if err != nil {
				return err
			}

			return cp.cpCheckJustInitZero(just, v.BlockHash())

		case vote.JustTypeInitOne:
			err := cp.cpCheckCPValue(v, vote.CPValueYes)
			if err != nil {
				return err
			}

			return cp.cpCheckJustInitOne(just)
		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	} else {
		switch just.Type() {
		case vote.JustTypePreVoteSoft:
			err := cp.cpCheckCPValue(v, vote.CPValueNo, vote.CPValueYes)
			if err != nil {
				return err
			}

			return cp.cpCheckJustPreVoteSoft(just, v.BlockHash(), v.CPRound())

		case vote.JustTypePreVoteHard:
			err := cp.cpCheckCPValue(v, vote.CPValueNo, vote.CPValueYes)
			if err != nil {
				return err
			}

			return cp.cpCheckJustPreVoteHard(just, v.BlockHash(), v.CPRound(), v.CPValue())

		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	}
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) cpCheckJustMainVote(v *vote.Vote) error {
	just := v.CPJust()
	switch just.Type() {
	case vote.JustTypeMainVoteNoConflict:
		err := cp.cpCheckCPValue(v, vote.CPValueNo, vote.CPValueYes)
		if err != nil {
			return err
		}

		return cp.cpCheckJustMainVoteNoConflict(just, v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeMainVoteConflict:
		err := cp.cpCheckCPValue(v, vote.CPValueAbstain)
		if err != nil {
			return err
		}

		return cp.cpCheckJustMainVoteConflict(just, v.BlockHash(), v.CPRound())

	default:
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid main-vote justification",
		}
	}
}

func (cp *changeProposer) cpCheckJustDecide(v *vote.Vote) error {
	err := cp.cpCheckCPValue(v, vote.CPValueNo, vote.CPValueYes)
	if err != nil {
		return err
	}
	j, ok := v.CPJust().(*vote.JustDecided)
	if !ok {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   "invalid just data",
		}
	}

	err = j.QCert.ValidateCPMainVote(cp.validators,
		v.BlockHash(), int16(v.CPValue()), byte(v.CPRound()))
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) cpCheckJust(v *vote.Vote) error {
	switch v.Type() {
	case vote.VoteTypeCPPreVote:
		return cp.cpCheckJustPreVote(v)
	case vote.VoteTypeCPMainVote:
		return cp.cpCheckJustMainVote(v)
	case vote.VoteTypeCPDecided:
		return cp.cpCheckJustDecide(v)
	default:
		panic("unreachable")
	}
}

// cpStrongTermination decides if the Change Proposer phase should be terminated.
// If there is only one proper and justified `decided` vote, the validators can
// move on to the next phase.
// If the decided vote is for "No," then validators move to the precommit step and
// wait for committing the current proposal by gathering enough precommit votes.
// If the decided vote is for "Yes," then the validator moves to the propose step
// and starts a new round.
func (cp *changeProposer) cpStrongTermination() {
	cpDecided := cp.log.CPDecidedVoteVoteSet(cp.round)
	if cpDecided.HasAnyVoteFor(cp.cpRound, vote.CPValueNo) {
		cp.cpDecided = 0
		cp.enterNewState(cp.precommitState)
	} else if cpDecided.HasAnyVoteFor(cp.cpRound, vote.CPValueYes) {
		cp.round++
		cp.cpDecided = 1
		cp.enterNewState(cp.proposeState)
	}
}
