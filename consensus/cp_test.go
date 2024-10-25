package consensus

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestChangeProposer(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestSetProposalAfterChangeProposer(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)

	p := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(p)
	assert.NotNil(t, td.consP.Proposal())
}

func TestChangeProposerAgreement1(t *testing.T) {
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

func TestChangeProposerAgreement0(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	height := uint32(2)
	round := int16(1)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.checkHeightRound(t, td.consP, height, round)

	prop := td.makeProposal(t, height, round)

	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, prop.Block().Hash())
	td.addCPPreVote(td.consP, prop.Block().Hash(), height, round, vote.CPValueNo, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, prop.Block().Hash(), height, round, vote.CPValueNo, preVote0.CPJust(), tIndexY)

	mainVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, prop.Block().Hash())
	td.addCPMainVote(td.consP, prop.Block().Hash(), height, round, vote.CPValueNo, mainVote0.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, prop.Block().Hash(), height, round, vote.CPValueNo, mainVote0.CPJust(), tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPDecided, prop.Block().Hash())
	td.shouldPublishQueryProposal(t, td.consP, height)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.checkHeightRound(t, td.consP, height, round)
}

// ConsP receives all PRE-VOTE:0 votes before receiving a proposal or prepare votes.
// It should vote PRE-VOTES:1 and MAIN-VOTE:0.
func TestCrashOnTestnet(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	height := uint32(2)
	round := int16(0)
	td.consP.MoveToNewHeight()

	blockHash := td.RandHash()
	vote1 := vote.NewPrepareVote(blockHash, height, round, td.consX.valKey.Address())
	vote2 := vote.NewPrepareVote(blockHash, height, round, td.consY.valKey.Address())
	vote3 := vote.NewPrepareVote(blockHash, height, round, td.consB.valKey.Address())

	td.HelperSignVote(td.consX.valKey, vote1)
	td.HelperSignVote(td.consY.valKey, vote2)
	td.HelperSignVote(td.consB.valKey, vote3)

	votes := map[crypto.Address]*vote.Vote{}
	votes[vote1.Signer()] = vote1
	votes[vote2.Signer()] = vote2
	votes[vote3.Signer()] = vote3

	qCert := td.consP.makeVoteCertificate(votes)
	just0 := &vote.JustInitNo{QCert: qCert}
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, just0, tIndexX)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, just0, tIndexY)
	td.addCPPreVote(td.consP, blockHash, height, round, vote.CPValueNo, just0, tIndexB)

	td.newHeightTimeout(td.consP)
	td.changeProposerTimeout(td.consP)

	preVote := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
	mainVote := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, blockHash)
	assert.Equal(t, vote.CPValueYes, preVote.CPValue())
	assert.Equal(t, vote.CPValueNo, mainVote.CPValue())
}

func TestInvalidJustInitOne(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)
	just := &vote.JustInitYes{}

	t.Run("invalid value: no", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(hash.UndefHash, height, round, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: no",
		})
	})

	t.Run("invalid block hash", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(hash.UndefHash, height, round, 1, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("with main-vote justification", func(t *testing.T) {
		invJust := &vote.JustMainVoteNoConflict{}
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 0, vote.CPValueYes, invJust, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: invJust.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})
}

func TestInvalidJustInitZero(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)
	just := &vote.JustInitNo{
		QCert: td.GenerateTestPrepareCertificate(height),
	}

	t.Run("invalid value: yes", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), height, round, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: yes",
		})
	})

	t.Run("cp-round should be zero", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustPreVoteHard(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)
	just := &vote.JustPreVoteHard{
		QCert: td.GenerateTestPrepareCertificate(height),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be zero", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustPreVoteSoft(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)
	just := &vote.JustPreVoteSoft{
		QCert: td.GenerateTestPrepareCertificate(height),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be zero", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(td.RandHash(), height, round, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustMainVoteNoConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)
	just := &vote.JustMainVoteNoConflict{
		QCert: td.GenerateTestPrepareCertificate(height),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), height, round, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		cpVote := vote.NewCPMainVote(td.RandHash(), height, round, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustMainVoteConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)

	t.Run("invalid value: no", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: &vote.JustInitNo{
				QCert: td.GenerateTestPrepareCertificate(height),
			},
			JustYes: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), height, round, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: no",
		})
	})

	t.Run("invalid value: yes", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: &vote.JustInitNo{
				QCert: td.GenerateTestPrepareCertificate(height),
			},
			JustYes: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), height, round, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: yes",
		})
	})

	t.Run("invalid value: unexpected justification (just0)", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestPrepareCertificate(height),
			},
			JustYes: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), height, round, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: vote.JustTypePreVoteSoft,
			Reason:   "invalid just data",
		})
	})

	t.Run("invalid value: unexpected justification", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			JustNo: &vote.JustInitNo{
				QCert: td.GenerateTestPrepareCertificate(height),
			},
			JustYes: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestPrepareCertificate(height),
			},
		}
		v := vote.NewCPMainVote(td.RandHash(), height, round, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "unexpected justification: JustInitNo",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		just0 := &vote.JustInitNo{
			QCert: td.GenerateTestPrepareCertificate(height),
		}
		just := &vote.JustMainVoteConflict{
			JustNo:  just0,
			JustYes: &vote.JustInitYes{},
		}
		cpVote := vote.NewCPMainVote(td.RandHash(), height, round, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just0.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just0.QCert.Committers()),
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		just0 := &vote.JustPreVoteSoft{
			QCert: td.GenerateTestPrepareCertificate(height),
		}
		just := &vote.JustMainVoteConflict{
			JustNo: just0,
			JustYes: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestPrepareCertificate(height),
			},
		}
		cpVote := vote.NewCPMainVote(td.RandHash(), height, round, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(cpVote)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just0.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just0.QCert.Committers()),
		})
	})
}

func TestInvalidJustDecided(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	height := uint32(1)
	round := int16(0)
	just := &vote.JustDecided{
		QCert: td.GenerateTestPrepareCertificate(height),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPDecidedVote(td.RandHash(), height, round, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPDecidedVote(td.RandHash(), height, round, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.checkJust(v)
		assert.ErrorIs(t, err, InvalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}
