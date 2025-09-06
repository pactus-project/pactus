package consensusv2

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestCPChangeProposer(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestCPQueryVote(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)
	td.queryVoteTimeout(td.consP)

	td.shouldPublishQueryVote(t, td.consP, 2, 0)
}

func TestCPSetProposalAfterChangeProposer(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)

	prop := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(prop)
	assert.NotNil(t, td.consP.Proposal())
}

func TestCPChangeProposerAgreementYes(t *testing.T) {
	td := setup(t)

	height := uint32(1)
	round := int16(0)
	td.enterNewHeight(td.consP)
	td.checkHeightRound(t, td.consP, height, round)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
	td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, preVote0.CPJust(), tIndexY)

	mainVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, hash.UndefHash)
	td.addCPMainVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, mainVote0.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, mainVote0.CPJust(), tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPDecided, hash.UndefHash)
	td.checkHeightRound(t, td.consP, height, round+1)
}

func TestCPChangeProposerAgreementNo(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	height := uint32(2)
	round := int16(1)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.checkHeightRound(t, td.consP, height, round)

	prop := td.makeProposal(t, height, round)
	blockHash := prop.Block().Hash()

	td.consP.SetProposal(prop)
	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexX)
	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexY)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, blockHash)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, preVote0.CPJust(), tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, blockHash)
	td.shouldPublishBlockAnnounce(t, td.consP, blockHash)
	td.checkHeightRound(t, td.consP, height+1, 0)
}

// ConsP receives all `PRE-VOTE:no` votes before receiving a proposal or precommit votes.
// It should vote `PRE-VOTES:yes` but goes to the pre-commit phase.
func TestCPCrashOnTestnet(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	height := uint32(2)
	round := int16(0)
	td.consP.MoveToNewHeight()

	blockHash := td.RandHash()
	just0, _, _ := td.makeChangeProposerJusts(t, blockHash, height, round)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, just0, tIndexX)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, just0, tIndexY)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, just0, tIndexB)

	td.newHeightTimeout(td.consP)
	td.changeProposerTimeout(td.consP)

	preVote := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
	assert.Equal(t, vote.CPValueYes, preVote.CPValue())
	assert.Equal(t, "precommit", td.consP.currentState.name())
}

func TestCPMoveToNextRoundOnDecidedVoteYes(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)
	h := uint32(1)
	r := int16(0)

	td.checkHeightRound(t, td.consP, h, r)

	_, _, decideJust := td.makeChangeProposerJusts(t, hash.UndefHash, h, r)
	td.addCPDecidedVote(td.consP, hash.UndefHash, h, r, vote.CPValueYes, decideJust, tIndexX)

	td.checkHeightRound(t, td.consP, h, r+1)
}

func TestCPInvalidJustInitYes(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustInitYes{}

	t.Run("invalid value: no", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: no",
		})
	})

	t.Run("cp-round should be 0", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r, 1, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid round: 1",
		})
	})

	t.Run("invalid block hash", func(t *testing.T) {
		blockHash := td.RandHash()
		v := vote.NewCPPreVote(blockHash, h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid block hash: " + blockHash.String(),
		})
	})
}

func TestCPInvalidJustInitNo(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustInitNo{
		QCert: td.GenerateTestVoteCertificate(h),
	}

	t.Run("invalid value: yes", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: yes",
		})
	})

	t.Run("cp-round should be 0", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid round: 1",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.Error(t, err)
	})
}

func TestCPInvalidJustPreVoteHard(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustPreVoteHard{
		QCert: td.GenerateTestVoteCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be 0", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid round: 0",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestCPInvalidJustPreVoteSoft(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustPreVoteSoft{
		QCert: td.GenerateTestVoteCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be 0", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid round: 0",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestCPInvalidJustMainVoteNoConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustMainVoteNoConflict{
		QCert: td.GenerateTestVoteCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestCPInvalidJustMainVoteConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)

	blockHash := td.RandHash()
	properJustInitNo, _, _ := td.makeChangeProposerJusts(t, blockHash, h, r)
	properJustInitYes := &vote.JustInitYes{}

	t.Run("invalid value: no", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo:  properJustInitNo,
			JustYes: properJustInitYes,
		}
		v := vote.NewCPMainVote(blockHash, h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: no",
		})
	})

	t.Run("invalid value: yes", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo:  properJustInitNo,
			JustYes: properJustInitYes,
		}
		v := vote.NewCPMainVote(blockHash, h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: yes",
		})
	})

	t.Run("invalid value: unexpected justification (JustNo)", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo:  &vote.JustInitYes{},
			JustYes: properJustInitYes,
		}
		v := vote.NewCPMainVote(blockHash, h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "unexpected justification: JustInitYes",
		})
	})

	t.Run("invalid value: unexpected justification (JustYes)", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: properJustInitNo,
			JustYes: &vote.JustInitNo{
				QCert: td.GenerateTestVoteCertificate(h),
			},
		}
		v := vote.NewCPMainVote(blockHash, h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "unexpected justification: JustInitNo",
		})
	})

	t.Run("invalid certificate - No", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestVoteCertificate(h),
			},
			JustYes: properJustInitYes,
		}
		v := vote.NewCPMainVote(blockHash, h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid round: 0",
		})
	})

	t.Run("invalid certificate - Yes", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: properJustInitNo,
			JustYes: &vote.JustPreVoteHard{
				QCert: td.GenerateTestVoteCertificate(h),
			},
		}
		v := vote.NewCPMainVote(blockHash, h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid round: 0",
		})
	})
}

func TestCPInvalidJustDecided(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustDecided{
		QCert: td.GenerateTestVoteCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPDecidedVote(td.RandHash(), h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPDecidedVote(hash.UndefHash, h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consP.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			Reason: fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}
