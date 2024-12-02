package consensus

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type changeProposer struct {
	*consensus
}

func (*changeProposer) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (cp *changeProposer) onTimeout(t *ticker) {
	if t.Target == tickerTargetQueryVote {
		cp.queryVote()
		cp.scheduleTimeout(t.Duration*2, cp.height, cp.round, tickerTargetQueryVote)
	}
}

func (*changeProposer) checkCPValue(vote *vote.Vote, allowedValues ...vote.CPValue) error {
	for _, v := range allowedValues {
		if vote.CPValue() == v {
			return nil
		}
	}

	return InvalidJustificationError{
		JustType: vote.CPJust().Type(),
		Reason:   fmt.Sprintf("invalid value: %v", vote.CPValue()),
	}
}

func (cp *changeProposer) checkJustInitNo(just vote.Just, blockHash hash.Hash) error {
	j, ok := just.(*vote.JustInitNo)
	if !ok {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidatePrepare(cp.validators, blockHash)
	if err != nil {
		return InvalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (*changeProposer) checkJustInitYes(just vote.Just) error {
	_, ok := just.(*vote.JustInitYes)
	if !ok {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	return nil
}

func (cp *changeProposer) checkJustPreVoteHard(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustPreVoteHard)
	if !ok {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidateCPPreVote(cp.validators,
		blockHash, cpRound-1, byte(cpValue))
	if err != nil {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) checkJustPreVoteSoft(just vote.Just,
	blockHash hash.Hash, cpRound int16,
) error {
	j, ok := just.(*vote.JustPreVoteSoft)
	if !ok {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidateCPMainVote(cp.validators,
		blockHash, cpRound-1, byte(vote.CPValueAbstain))
	if err != nil {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) checkJustMainVoteNoConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustMainVoteNoConflict)
	if !ok {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	err := j.QCert.ValidateCPPreVote(cp.validators,
		blockHash, cpRound, byte(cpValue))
	if err != nil {
		return InvalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJustMainVoteConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16,
) error {
	j, ok := just.(*vote.JustMainVoteConflict)
	if !ok {
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	if cpRound == 0 {
		err := cp.checkJustInitNo(j.JustNo, blockHash)
		if err != nil {
			return err
		}

		err = cp.checkJustInitYes(j.JustYes)
		if err != nil {
			return err
		}

		return nil
	}

	// Just0 can be for Zero or Abstain values.
	switch j.JustNo.Type() {
	case vote.JustTypePreVoteSoft:
		err := cp.checkJustPreVoteSoft(j.JustNo, blockHash, cpRound)
		if err != nil {
			return err
		}
	case vote.JustTypePreVoteHard:
		err := cp.checkJustPreVoteHard(j.JustNo, blockHash, cpRound, vote.CPValueNo)
		if err != nil {
			return err
		}
	default:
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("unexpected justification: %s", j.JustNo.Type()),
		}
	}

	err := cp.checkJustPreVoteHard(j.JustYes, hash.UndefHash, cpRound, vote.CPValueYes)
	if err != nil {
		return err
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJustPreVote(vte *vote.Vote) error {
	just := vte.CPJust()
	if vte.CPRound() == 0 {
		switch just.Type() {
		case vote.JustTypeInitNo:
			err := cp.checkCPValue(vte, vote.CPValueNo)
			if err != nil {
				return err
			}

			return cp.checkJustInitNo(just, vte.BlockHash())

		case vote.JustTypeInitYes:
			err := cp.checkCPValue(vte, vote.CPValueYes)
			if err != nil {
				return err
			}

			return cp.checkJustInitYes(just)
		default:
			return InvalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	} else {
		switch just.Type() {
		case vote.JustTypePreVoteSoft:
			err := cp.checkCPValue(vte, vote.CPValueNo, vote.CPValueYes)
			if err != nil {
				return err
			}

			return cp.checkJustPreVoteSoft(just, vte.BlockHash(), vte.CPRound())

		case vote.JustTypePreVoteHard:
			err := cp.checkCPValue(vte, vote.CPValueNo, vote.CPValueYes)
			if err != nil {
				return err
			}

			return cp.checkJustPreVoteHard(just, vte.BlockHash(), vte.CPRound(), vte.CPValue())

		default:
			return InvalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	}
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJustMainVote(vte *vote.Vote) error {
	just := vte.CPJust()
	switch just.Type() {
	case vote.JustTypeMainVoteNoConflict:
		err := cp.checkCPValue(vte, vote.CPValueNo, vote.CPValueYes)
		if err != nil {
			return err
		}

		return cp.checkJustMainVoteNoConflict(just, vte.BlockHash(), vte.CPRound(), vte.CPValue())

	case vote.JustTypeMainVoteConflict:
		err := cp.checkCPValue(vte, vote.CPValueAbstain)
		if err != nil {
			return err
		}

		return cp.checkJustMainVoteConflict(just, vte.BlockHash(), vte.CPRound())

	default:
		return InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid main-vote justification",
		}
	}
}

func (cp *changeProposer) checkJustDecide(vte *vote.Vote) error {
	err := cp.checkCPValue(vte, vote.CPValueNo, vote.CPValueYes)
	if err != nil {
		return err
	}
	j, ok := vte.CPJust().(*vote.JustDecided)
	if !ok {
		return InvalidJustificationError{
			JustType: j.Type(),
			Reason:   "invalid just data",
		}
	}

	err = j.QCert.ValidateCPMainVote(cp.validators,
		vte.BlockHash(), vte.CPRound(), byte(vte.CPValue()))
	if err != nil {
		return InvalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJust(vte *vote.Vote) error {
	switch vte.Type() {
	case vote.VoteTypeCPPreVote:
		return cp.checkJustPreVote(vte)
	case vote.VoteTypeCPMainVote:
		return cp.checkJustMainVote(vte)
	case vote.VoteTypeCPDecided:
		return cp.checkJustDecide(vte)
	default:
		panic("unreachable")
	}
}

func (cp *changeProposer) cpStrongTermination(round, cpRound int16) {
	cpDecided := cp.log.CPDecidedVoteSet(round)
	if cpDecided.HasAnyVoteFor(cpRound, vote.CPValueNo) {
		cp.round = round
		cp.cpDecided = 0

		roundProposal := cp.log.RoundProposal(round)
		if roundProposal == nil {
			cp.queryProposal()
		}
		cp.enterNewState(cp.prepareState)
	} else if cpDecided.HasAnyVoteFor(cpRound, vote.CPValueYes) {
		cp.round = round + 1
		cp.cpDecided = 1

		cp.enterNewState(cp.proposeState)
	}
}
