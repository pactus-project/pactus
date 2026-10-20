package vote

import (
	"fmt"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/certificate"
)

type JustType uint8

const (
	JustTypeInitNo             = JustType(1)
	JustTypeInitYes            = JustType(2)
	JustTypePreVoteSoft        = JustType(3)
	JustTypePreVoteHard        = JustType(4)
	JustTypeMainVoteConflict   = JustType(5)
	JustTypeMainVoteNoConflict = JustType(6)
	JustTypeDecided            = JustType(7)
)

func (t JustType) String() string {
	switch t {
	case JustTypeInitNo:
		return "JustInitNo"
	case JustTypeInitYes:
		return "JustInitYes"
	case JustTypePreVoteSoft:
		return "JustPreVoteSoft"
	case JustTypePreVoteHard:
		return "JustPreVoteHard"
	case JustTypeMainVoteConflict:
		return "JustMainVoteConflict"
	case JustTypeMainVoteNoConflict:
		return "JustMainVoteNoConflict"
	case JustTypeDecided:
		return "JustDecided"

	default:
		return "Unknown"
	}
}

type Just interface {
	Type() JustType
	BasicCheck() error
}

// checkJustHeightRound checks that the quorum certificates attached to the
// justification were issued at the same height and round as the vote. Without
// this binding, a valid certificate from another height or round could be
// replayed as the justification of a newly signed vote.
func checkJustHeightRound(just Just, height types.Height, round types.Round) error {
	switch j := just.(type) {
	case *JustInitNo:
		return checkQCertHeightRound(j.QCert, height, round)
	case *JustInitYes:
		return nil
	case *JustPreVoteSoft:
		return checkQCertHeightRound(j.QCert, height, round)
	case *JustPreVoteHard:
		return checkQCertHeightRound(j.QCert, height, round)
	case *JustMainVoteNoConflict:
		return checkQCertHeightRound(j.QCert, height, round)
	case *JustDecided:
		return checkQCertHeightRound(j.QCert, height, round)
	case *JustMainVoteConflict:
		if j.JustNo == nil || j.JustYes == nil {
			return BasicCheckError{Reason: "quorum certificate is not set"}
		}
		if err := checkJustHeightRound(j.JustNo, height, round); err != nil {
			return err
		}

		return checkJustHeightRound(j.JustYes, height, round)
	default:
		return BasicCheckError{Reason: "unknown justification"}
	}
}

func checkQCertHeightRound(qCert *certificate.Certificate,
	height types.Height, round types.Round,
) error {
	if qCert == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}

	if qCert.Height() != height {
		return BasicCheckError{
			Reason: fmt.Sprintf("certificate height %d does not match vote height %d",
				qCert.Height(), height),
		}
	}

	if qCert.Round() != round {
		return BasicCheckError{
			Reason: fmt.Sprintf("certificate round %d does not match vote round %d",
				qCert.Round(), round),
		}
	}

	return nil
}

func makeJust(t JustType) Just {
	switch t {
	case JustTypeInitNo:
		return &JustInitNo{}
	case JustTypeInitYes:
		return &JustInitYes{}
	case JustTypePreVoteSoft:
		return &JustPreVoteSoft{}
	case JustTypePreVoteHard:
		return &JustPreVoteHard{}
	case JustTypeMainVoteConflict:
		return &JustMainVoteConflict{}
	case JustTypeMainVoteNoConflict:
		return &JustMainVoteNoConflict{}
	case JustTypeDecided:
		return &JustDecided{}
	default:
		return nil
	}
}

type JustInitNo struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

type JustInitYes struct {
	//
}

type JustPreVoteSoft struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

type JustPreVoteHard struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}
type JustMainVoteConflict struct {
	JustNo  Just
	JustYes Just
}

type JustMainVoteConflictV2 struct {
	Just1 Just
	Just2 Just
}

type JustMainVoteNoConflict struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

type JustDecided struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

func (*JustInitNo) Type() JustType {
	return JustTypeInitNo
}

func (*JustInitYes) Type() JustType {
	return JustTypeInitYes
}

func (*JustPreVoteSoft) Type() JustType {
	return JustTypePreVoteSoft
}

func (*JustPreVoteHard) Type() JustType {
	return JustTypePreVoteHard
}

func (*JustMainVoteConflict) Type() JustType {
	return JustTypeMainVoteConflict
}

func (*JustMainVoteNoConflict) Type() JustType {
	return JustTypeMainVoteNoConflict
}

func (*JustDecided) Type() JustType {
	return JustTypeDecided
}

func (j *JustInitNo) BasicCheck() error {
	if j.QCert == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}
	err := j.QCert.BasicCheck()
	if err != nil {
		return err
	}

	return nil
}

func (*JustInitYes) BasicCheck() error {
	return nil
}

func (j *JustPreVoteSoft) BasicCheck() error {
	if j.QCert == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}

	return j.QCert.BasicCheck()
}

func (j *JustPreVoteHard) BasicCheck() error {
	if j.QCert == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}

	return j.QCert.BasicCheck()
}

func (j *JustMainVoteConflict) BasicCheck() error {
	if j.JustNo == nil || j.JustYes == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}
	if err := j.JustNo.BasicCheck(); err != nil {
		return err
	}

	return j.JustYes.BasicCheck()
}

func (j *JustMainVoteNoConflict) BasicCheck() error {
	if j.QCert == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}

	return j.QCert.BasicCheck()
}

func (j *JustDecided) BasicCheck() error {
	if j.QCert == nil {
		return BasicCheckError{Reason: "quorum certificate is not set"}
	}

	return j.QCert.BasicCheck()
}
