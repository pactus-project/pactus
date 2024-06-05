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

func (*changeProposer) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (cp *changeProposer) onTimeout(t *ticker) {
	if t.Target == tickerTargetQueryVotes {
		cp.queryVotes()
		cp.scheduleTimeout(t.Duration*2, cp.height, cp.round, tickerTargetQueryVotes)
	}
}

func (*changeProposer) cpCheckCPValue(value vote.CPValue, allowedValues ...vote.CPValue) error {
	for _, v := range allowedValues {
		if value == v {
			return nil
		}
	}

	return invalidJustificationError{
		Reason: fmt.Sprintf("invalid value: %v", value),
	}
}

func (cp *changeProposer) cpCheckJustInitNo(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustInitNo)
	if !ok {
		return invalidJustificationError{
			Reason: "invalid just data",
		}
	}

	if cpRound != 0 {
		return invalidJustificationError{
			Reason: fmt.Sprintf("invalid round: %v", cpRound),
		}
	}

	err := cp.cpCheckCPValue(cpValue, vote.CPValueNo)
	if err != nil {
		return err
	}

	err = j.QCert.ValidatePrepare(cp.validators, blockHash)
	if err != nil {
		return err
	}

	return nil
}

func (cp *changeProposer) cpCheckJustInitYes(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	_, ok := just.(*vote.JustInitYes)
	if !ok {
		return invalidJustificationError{
			Reason: "invalid just data",
		}
	}

	if cpRound != 0 {
		return invalidJustificationError{
			Reason: fmt.Sprintf("invalid round: %v", cpRound),
		}
	}

	err := cp.cpCheckCPValue(cpValue, vote.CPValueYes)
	if err != nil {
		return err
	}

	if !blockHash.IsUndef() {
		return invalidJustificationError{
			Reason: fmt.Sprintf("invalid block hash: %s", blockHash),
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
			Reason: "invalid just data",
		}
	}

	if cpRound == 0 {
		return invalidJustificationError{
			Reason: "invalid round: 0",
		}
	}

	err := cp.cpCheckCPValue(cpValue, vote.CPValueNo, vote.CPValueYes)
	if err != nil {
		return err
	}

	err = j.QCert.ValidateCPPreVote(cp.validators,
		blockHash, cpRound-1, byte(cpValue))
	if err != nil {
		return invalidJustificationError{
			Reason: err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustPreVoteSoft(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustPreVoteSoft)
	if !ok {
		return invalidJustificationError{
			Reason: "invalid just data",
		}
	}

	if cpRound == 0 {
		return invalidJustificationError{
			Reason: "invalid round: 0",
		}
	}

	err := cp.cpCheckCPValue(cpValue, vote.CPValueNo, vote.CPValueYes)
	if err != nil {
		return err
	}

	err = j.QCert.ValidateCPMainVote(cp.validators,
		blockHash, cpRound-1, byte(vote.CPValueAbstain))
	if err != nil {
		return invalidJustificationError{
			Reason: err.Error(),
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
			Reason: "invalid just data",
		}
	}
	err := cp.cpCheckCPValue(cpValue, vote.CPValueNo, vote.CPValueYes)
	if err != nil {
		return err
	}

	err = j.QCert.ValidateCPPreVote(cp.validators,
		blockHash, cpRound, byte(cpValue))
	if err != nil {
		return invalidJustificationError{
			Reason: err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustMainVoteConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustMainVoteConflict)
	if !ok {
		return invalidJustificationError{
			Reason: "invalid just data",
		}
	}

	err := cp.cpCheckCPValue(cpValue, vote.CPValueAbstain)
	if err != nil {
		return err
	}

	switch j.JustNo.Type() {
	case vote.JustTypeInitNo:
		err := cp.cpCheckJustInitNo(j.JustNo, blockHash, cpRound, vote.CPValueNo)
		if err != nil {
			return err
		}
	case vote.JustTypePreVoteHard:
		err := cp.cpCheckJustPreVoteHard(j.JustNo, blockHash, cpRound, vote.CPValueNo)
		if err != nil {
			return err
		}
	case vote.JustTypePreVoteSoft:
		err := cp.cpCheckJustPreVoteSoft(j.JustNo, blockHash, cpRound, vote.CPValueNo)
		if err != nil {
			return err
		}

	case vote.JustTypeInitYes,
		vote.JustTypeMainVoteConflict,
		vote.JustTypeMainVoteNoConflict,
		vote.JustTypeDecided:
		return invalidJustificationError{
			Reason: fmt.Sprintf("unexpected justification: %s", j.JustNo.Type()),
		}
	}

	switch j.JustYes.Type() {
	case vote.JustTypeInitYes:
		err := cp.cpCheckJustInitYes(j.JustYes, hash.UndefHash, cpRound, vote.CPValueYes)
		if err != nil {
			return err
		}

	case vote.JustTypePreVoteHard:
		err := cp.cpCheckJustPreVoteHard(j.JustYes, hash.UndefHash, cpRound, vote.CPValueYes)
		if err != nil {
			return err
		}

	case vote.JustTypeInitNo,
		vote.JustTypePreVoteSoft,
		vote.JustTypeMainVoteConflict,
		vote.JustTypeMainVoteNoConflict,
		vote.JustTypeDecided:
		return invalidJustificationError{
			Reason: fmt.Sprintf("unexpected justification: %s", j.JustNo.Type()),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJustDecide(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustDecided)
	if !ok {
		return invalidJustificationError{
			Reason: "invalid just data",
		}
	}

	err := cp.cpCheckCPValue(cpValue, vote.CPValueNo, vote.CPValueYes)
	if err != nil {
		return err
	}

	err = j.QCert.ValidateCPMainVote(cp.validators,
		blockHash, cpRound, byte(cpValue))
	if err != nil {
		return invalidJustificationError{
			Reason: err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) cpCheckJust(v *vote.Vote) error {
	switch v.CPJust().Type() {
	case vote.JustTypeInitYes:
		return cp.cpCheckJustInitYes(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeInitNo:
		return cp.cpCheckJustInitNo(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypePreVoteSoft:
		return cp.cpCheckJustPreVoteSoft(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypePreVoteHard:
		return cp.cpCheckJustPreVoteHard(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeMainVoteNoConflict:
		return cp.cpCheckJustMainVoteNoConflict(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeMainVoteConflict:
		return cp.cpCheckJustMainVoteConflict(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeDecided:
		return cp.cpCheckJustDecide(v.CPJust(),
			v.BlockHash(), v.CPRound(), v.CPValue())

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
