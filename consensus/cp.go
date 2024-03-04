package consensus

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
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

func (cp *changeProposer) checkCPValue(vte *vote.Vote, allowedValues ...vote.CPValue) error {
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

func (cp *changeProposer) checkJustInitZero(just vote.Just, blockHash hash.Hash) error {
	j, ok := just.(*vote.JustInitZero)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypePrepare.String())...)

	err := j.QCert.Validate(cp.height, cp.validators, sb)
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

func (cp *changeProposer) checkJustInitOne(just vote.Just) error {
	_, ok := just.(*vote.JustInitOne)
	if !ok {
		return invalidJustificationError{
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
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPPreVote.String())...)
	sb = append(sb, util.Int16ToSlice(cpRound-1)...)
	sb = append(sb, byte(cpValue))

	err := j.QCert.Validate(cp.height, cp.validators, sb)
	if err != nil {
		return invalidJustificationError{
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
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPMainVote.String())...)
	sb = append(sb, util.Int16ToSlice(cpRound-1)...)
	sb = append(sb, byte(vote.CPValueAbstain))

	err := j.QCert.Validate(cp.height, cp.validators, sb)
	if err != nil {
		return invalidJustificationError{
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
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPPreVote.String())...)
	sb = append(sb, util.Int16ToSlice(cpRound)...)
	sb = append(sb, byte(cpValue))

	err := j.QCert.Validate(cp.height, cp.validators, sb)
	if err != nil {
		return invalidJustificationError{
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
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	if cpRound == 0 {
		err := cp.checkJustInitZero(j.Just0, blockHash)
		if err != nil {
			return err
		}

		err = cp.checkJustInitOne(j.Just1)
		if err != nil {
			return err
		}

		return nil
	}

	// Just0 can be for Zero or Abstain values.
	switch j.Just0.Type() {
	case vote.JustTypePreVoteSoft:
		err := cp.checkJustPreVoteSoft(j.Just0, blockHash, cpRound)
		if err != nil {
			return err
		}
	case vote.JustTypePreVoteHard:
		err := cp.checkJustPreVoteHard(j.Just0, blockHash, cpRound, vote.CPValueZero)
		if err != nil {
			return err
		}
	default:
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("unexpected justification: %s", j.Just0.Type()),
		}
	}

	err := cp.checkJustPreVoteHard(j.Just1, hash.UndefHash, cpRound, vote.CPValueOne)
	if err != nil {
		return err
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJustPreVote(v *vote.Vote) error {
	just := v.CPJust()
	if v.CPRound() == 0 {
		switch just.Type() {
		case vote.JustTypeInitZero:
			err := cp.checkCPValue(v, vote.CPValueZero)
			if err != nil {
				return err
			}

			return cp.checkJustInitZero(just, v.BlockHash())

		case vote.JustTypeInitOne:
			err := cp.checkCPValue(v, vote.CPValueOne)
			if err != nil {
				return err
			}

			return cp.checkJustInitOne(just)
		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	} else {
		switch just.Type() {
		case vote.JustTypePreVoteSoft:
			err := cp.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
			if err != nil {
				return err
			}

			return cp.checkJustPreVoteSoft(just, v.BlockHash(), v.CPRound())

		case vote.JustTypePreVoteHard:
			err := cp.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
			if err != nil {
				return err
			}

			return cp.checkJustPreVoteHard(just, v.BlockHash(), v.CPRound(), v.CPValue())

		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	}
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJustMainVote(v *vote.Vote) error {
	just := v.CPJust()
	switch just.Type() {
	case vote.JustTypeMainVoteNoConflict:
		err := cp.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
		if err != nil {
			return err
		}

		return cp.checkJustMainVoteNoConflict(just, v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeMainVoteConflict:
		err := cp.checkCPValue(v, vote.CPValueAbstain)
		if err != nil {
			return err
		}

		return cp.checkJustMainVoteConflict(just, v.BlockHash(), v.CPRound())

	default:
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid main-vote justification",
		}
	}
}

func (cp *changeProposer) checkJustDecide(v *vote.Vote) error {
	err := cp.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
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

	sb := certificate.BlockCertificateSignBytes(v.BlockHash(),
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPMainVote.String())...)
	sb = append(sb, util.Int16ToSlice(v.CPRound())...)
	sb = append(sb, byte(v.CPValue()))

	err = j.QCert.Validate(cp.height, cp.validators, sb)
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}

	return nil
}

//nolint:exhaustive // refactor me; check just by just_type, not vote_type
func (cp *changeProposer) checkJust(v *vote.Vote) error {
	switch v.Type() {
	case vote.VoteTypeCPPreVote:
		return cp.checkJustPreVote(v)
	case vote.VoteTypeCPMainVote:
		return cp.checkJustMainVote(v)
	case vote.VoteTypeCPDecided:
		return cp.checkJustDecide(v)
	default:
		panic("unreachable")
	}
}

func (cp *changeProposer) strongTermination() {
	cpDecided := cp.log.CPDecidedVoteVoteSet(cp.round)
	if cpDecided.HasAnyVoteFor(cp.cpRound, vote.CPValueZero) {
		cp.cpDecide(vote.CPValueZero)
	} else if cpDecided.HasAnyVoteFor(cp.cpRound, vote.CPValueOne) {
		cp.cpDecide(vote.CPValueOne)
	}
}

func (cp *changeProposer) cpDecide(cpValue vote.CPValue) {
	if cpValue == vote.CPValueOne {
		cp.round++
		cp.cpDecided = 1
		cp.enterNewState(cp.proposeState)
	} else if cpValue == vote.CPValueZero {
		roundProposal := cp.log.RoundProposal(cp.round)
		if roundProposal == nil {
			cp.queryProposal()
		}
		cp.cpDecided = 0
		cp.enterNewState(cp.prepareState)
	}
}
