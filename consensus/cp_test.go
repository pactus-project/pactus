package consensus

import (
	"fmt"
	"testing"

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
	assert.NotNil(t, td.consP.RoundProposal(0))
}

func TestChangeProposerAgreement1(t *testing.T) {
	td := setup(t)

	h := uint32(1)
	r := int16(0)
	td.enterNewHeight(td.consP)
	td.checkHeightRound(t, td.consP, h, r)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueOne, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueOne, preVote0.CPJust(), tIndexY)

	mainVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, hash.UndefHash)
	td.addCPMainVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueOne, mainVote0.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueOne, mainVote0.CPJust(), tIndexY)

	preVote1 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 1, vote.CPValueOne, preVote1.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 1, vote.CPValueOne, preVote1.CPJust(), tIndexY)

	mainVote1 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, hash.UndefHash)
	td.addCPMainVote(td.consP, hash.UndefHash, h, r, 1, vote.CPValueOne, mainVote1.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, hash.UndefHash, h, r, 1, vote.CPValueOne, mainVote1.CPJust(), tIndexY)

	checkHeightRound(t, td.consP, 1, 1)
}

func TestChangeProposerAgreement0(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	h := uint32(2)
	r := int16(1)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.checkHeightRound(t, td.consP, h, r)

	p := td.makeProposal(t, h, r)

	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexB)

	td.changeProposerTimeout(td.consP)

	preVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, p.Block().Hash())
	td.addCPPreVote(td.consP, p.Block().Hash(), h, r, 0, vote.CPValueZero, preVote0.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, p.Block().Hash(), h, r, 0, vote.CPValueZero, preVote0.CPJust(), tIndexY)

	mainVote0 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, p.Block().Hash())
	td.addCPMainVote(td.consP, p.Block().Hash(), h, r, 0, vote.CPValueZero, mainVote0.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, p.Block().Hash(), h, r, 0, vote.CPValueZero, mainVote0.CPJust(), tIndexY)

	preVote1 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, p.Block().Hash())
	td.addCPPreVote(td.consP, p.Block().Hash(), h, r, 1, vote.CPValueZero, preVote1.CPJust(), tIndexX)
	td.addCPPreVote(td.consP, p.Block().Hash(), h, r, 1, vote.CPValueZero, preVote1.CPJust(), tIndexY)

	mainVote1 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPMainVote, p.Block().Hash())
	td.addCPMainVote(td.consP, p.Block().Hash(), h, r, 1, vote.CPValueZero, mainVote1.CPJust(), tIndexX)
	td.addCPMainVote(td.consP, p.Block().Hash(), h, r, 1, vote.CPValueZero, mainVote1.CPJust(), tIndexY)

	td.shouldPublishQueryProposal(t, td.consP, h, r)
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	checkHeightRound(t, td.consP, h, r)
}

func TestInvalidJustInitOne(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustInitOne{}

	t.Run("invalid value: zero", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r, 0, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: zero",
		})
	})

	t.Run("invalid block hash", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r, 1, vote.CPValueOne, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("cp-round should be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueOne, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid block hash",
		})
	})

	t.Run("with main-vote justification", func(t *testing.T) {
		invJust := &vote.JustMainVoteNoConflict{}
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueOne, invJust, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
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
	just := &vote.JustInitZero{
		QCert: td.GenerateTestCertificate(),
	}

	t.Run("invalid value: one", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueOne, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: one",
		})
	})

	t.Run("cp-round should be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate height is invalid (expected 1 got %v)", just.QCert.Height()),
		})
	})
}

func TestInvalidJustPreVoteHard(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustPreVoteHard{
		QCert: td.GenerateTestCertificate(),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate height is invalid (expected 1 got %v)", just.QCert.Height()),
		})
	})
}

func TestInvalidJustPreVoteSoft(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustPreVoteSoft{
		QCert: td.GenerateTestCertificate(),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("cp-round should not be zero", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 0, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid pre-vote justification",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPPreVote(td.RandHash(), h, r, 1, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate height is invalid (expected 1 got %v)", just.QCert.Height()),
		})
	})
}

func TestInvalidJustMainVoteNoConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)
	just := &vote.JustMainVoteNoConflict{
		QCert: td.GenerateTestCertificate(),
	}

	t.Run("invalid value: abstain", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: abstain",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   fmt.Sprintf("certificate height is invalid (expected 1 got %v)", just.QCert.Height()),
		})
	})
}

func TestInvalidJustMainVoteConflict(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	h := uint32(1)
	r := int16(0)

	t.Run("invalid value: zero", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustInitZero{
				QCert: td.GenerateTestCertificate(),
			},
			Just1: &vote.JustInitOne{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueZero, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: zero",
		})
	})

	t.Run("invalid value: one", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustInitZero{
				QCert: td.GenerateTestCertificate(),
			},
			Just1: &vote.JustInitOne{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueOne, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid value: one",
		})
	})

	t.Run("invalid value: unexpected justification (just0)", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestCertificate(),
			},
			Just1: &vote.JustInitOne{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "unexpected justification: JustPreVoteSoft",
		})
	})

	t.Run("invalid value: unexpected justification", func(t *testing.T) {
		just := &vote.JustMainVoteConflict{
			Just0: &vote.JustInitZero{
				QCert: td.GenerateTestCertificate(),
			},
			Just1: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestCertificate(),
			},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just.Type(),
			Reason:   "unexpected justification: JustInitZero",
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		just0 := &vote.JustInitZero{
			QCert: td.GenerateTestCertificate(),
		}
		just := &vote.JustMainVoteConflict{
			Just0: just0,
			Just1: &vote.JustInitOne{},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 0, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just0.Type(),
			Reason:   fmt.Sprintf("certificate height is invalid (expected 1 got %v)", just0.QCert.Height()),
		})
	})

	t.Run("invalid certificate", func(t *testing.T) {
		just0 := &vote.JustPreVoteSoft{
			QCert: td.GenerateTestCertificate(),
		}
		just := &vote.JustMainVoteConflict{
			Just0: just0,
			Just1: &vote.JustPreVoteSoft{
				QCert: td.GenerateTestCertificate(),
			},
		}
		v := vote.NewCPMainVote(td.RandHash(), h, r, 1, vote.CPValueAbstain, just, td.signers[tIndexB].Address())

		err := td.consX.checkJust(v)
		assert.ErrorIs(t, err, invalidJustificationError{
			JustType: just0.Type(),
			Reason:   fmt.Sprintf("certificate height is invalid (expected 1 got %v)", just0.QCert.Height()),
		})
	})
}
