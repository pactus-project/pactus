package fastconsensus

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

	prop := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(prop)
	assert.NotNil(t, td.consP.Proposal())
}

func TestChangeProposerAgreement1(t *testing.T) {
	td := setup(t)

	h := uint32(1)
	r := int16(0)
	td.enterNewHeight(td.consP)
	td.checkHeightRound(t, td.consP, h, r)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueYes, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueYes, preVote0.CPJust(), tIndexY)

	mainVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, hash.UndefHash)
	td.addCPMainVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueYes, mainVote0.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueYes, mainVote0.CPJust(), tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPDecided, hash.UndefHash)
	checkHeightRound(t, td.consP, h, r+1)
}

func TestChangeProposerAgreement0(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	h := uint32(2)
	r := int16(1)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.checkHeightRound(t, td.consP, h, r)

	prop := td.makeProposal(t, h, r)
	blockHash := prop.Block().Hash()

	td.consP.SetProposal(prop)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexY)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexM)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, blockHash)
	td.addCPPreVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, preVote0.CPJust(), tIndexY)

	mainVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, blockHash)
	td.addCPMainVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, mainVote0.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, mainVote0.CPJust(), tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPDecided, blockHash)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexY)
	checkHeightRound(t, td.consP, h, r)
}

// ConsP receives all PRE-VOTE:0 votes before receiving a proposal or prepare votes.
// It should vote PRE-VOTES:1 and MAIN-VOTE:0.
func TestCrashOnTestnet(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	h := uint32(2)
	r := int16(0)
	td.consP.MoveToNewHeight()

	blockHash := td.RandHash()
	v1 := vote.NewPrepareVote(blockHash, h, r, td.consX.valKey.Address())
	v2 := vote.NewPrepareVote(blockHash, h, r, td.consY.valKey.Address())
	v3 := vote.NewPrepareVote(blockHash, h, r, td.consB.valKey.Address())

	td.HelperSignVote(td.consX.valKey, v1)
	td.HelperSignVote(td.consY.valKey, v2)
	td.HelperSignVote(td.consB.valKey, v3)

	votes := map[crypto.Address]*vote.Vote{}
	votes[v1.Signer()] = v1
	votes[v2.Signer()] = v2
	votes[v3.Signer()] = v3

	cert := td.consP.makeVoteCertificate(votes)
	just0 := &vote.JustInitNo{QCert: cert}
	td.addCPPreVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, just0, tIndexX)
	td.addCPPreVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, just0, tIndexY)
	td.addCPPreVote(td.consP, blockHash, h, r, 0, vote.CPValueNo, just0, tIndexB)

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
	h := uint32(1)
	r := int16(0)
	just := &vote.JustInitYes{}

	t.Run("invalid value: no", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: no",
		})
	})

	t.Run("invalid block hash", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r, 1, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("with main-vote justification", func(t *testing.T) {
		invJust := &vote.JustMainVoteNoConflict{}
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueYes, invJust, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: invJust.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})
}

func TestInvalidJustInitZero(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustInitNo{
		QCert: td.GenerateTestPrepareCertificate(h),
	}

	t.Run("invalid value: yes", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: yes",
		})
	})

	t.Run("cp-round should be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.Error(t, err)
	})
}

func TestInvalidJustPreVoteHard(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustPreVoteHard{
		QCert: td.GenerateTestPrepareCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustPreVoteSoft(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustPreVoteSoft{
		QCert: td.GenerateTestPrepareCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustMainVoteNoConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustMainVoteNoConflict{
		QCert: td.GenerateTestPrepareCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}

func TestInvalidJustMainVoteConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)

	t.Run("invalid value: no", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustInitNo{
				QCert: td.GenerateTestPrepareCertificate(h),
			},
			Just1: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueNo, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: no",
		})
	})

	t.Run("invalid value: yes", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustInitNo{
				QCert: td.GenerateTestPrepareCertificate(h),
			},
			Just1: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: yes",
		})
	})

	t.Run("invalid value: unexpected justification (just0)", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestPrepareCertificate(h),
			},
			Just1: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: vote.JustTypePreVoteSoft,
			Reason:   "invalid just data",
		})
	})

	t.Run("invalid value: unexpected justification", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustInitNo{
				QCert: td.GenerateTestPrepareCertificate(h),
			},
			Just1: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestPrepareCertificate(h),
			},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "unexpected justification: JustInitZero",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		just0 := &vote.JustInitNo{
			QCert: td.GenerateTestPrepareCertificate(h),
		}
		just := &vote.JustMainVoteConflict{
			Just0: just0,
			Just1: &vote.JustInitYes{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.Error(t, err)
	})

	t.Run("invalid certificate", func(t *testing.T) {
		just0 := &vote.JustPreVoteSoft{
			QCert: td.GenerateTestPrepareCertificate(h),
		}
		just := &vote.JustMainVoteConflict{
			Just0: just0,
			Just1: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestPrepareCertificate(h),
			},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just0.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just0.QCert.Committers()),
		})
	})
}

func TestInvalidJustDecided(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustDecided{
		QCert: td.GenerateTestPrepareCertificate(h),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPDecidedVote(td.RandHash(), h, r, 0, vote.CPValueAbstain, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPDecidedVote(td.RandHash(), h, r, 0, vote.CPValueYes, just, td.consB.valKey.Address())

		err := td.consX.changeProposer.cpCheckJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate has an unexpected committers: %v", just.QCert.Committers()),
		})
	})
}
