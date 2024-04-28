package vote

import (
	"github.com/pactus-project/pactus/types/certificate"
)

type JustType uint8

const (
	JustTypeInitZero           = JustType(1)
	JustTypeInitOne            = JustType(2)
	JustTypePreVoteSoft        = JustType(3)
	JustTypePreVoteHard        = JustType(4)
	JustTypeMainVoteConflict   = JustType(5)
	JustTypeMainVoteNoConflict = JustType(6)
	JustTypeDecided            = JustType(7)
)

func (t JustType) String() string {
	switch t {
	case JustTypeInitZero:
		return "JustInitZero"
	case JustTypeInitOne:
		return "JustInitOne"
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

func makeJust(t JustType) Just {
	switch t {
	case JustTypeInitZero:
		return &JustInitNo{}
	case JustTypeInitOne:
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
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
}
type JustInitYes struct {
	//
}

type JustPreVoteSoft struct {
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
}

type JustPreVoteHard struct {
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
}
type JustMainVoteConflict struct {
	Just0 Just
	Just1 Just
}
type JustMainVoteNoConflict struct {
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
}

type JustDecided struct {
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
}

func (j *JustInitNo) Type() JustType {
	return JustTypeInitZero
}

func (j *JustInitYes) Type() JustType {
	return JustTypeInitOne
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
	return nil
}

func (j *JustInitYes) BasicCheck() error {
	return nil
}

func (j *JustPreVoteSoft) BasicCheck() error {
	return j.QCert.BasicCheck()
}

func (j *JustPreVoteHard) BasicCheck() error {
	return j.QCert.BasicCheck()
}

func (j *JustMainVoteConflict) BasicCheck() error {
	if err := j.Just0.BasicCheck(); err != nil {
		return err
	}

	return j.Just1.BasicCheck()
}

func (j *JustMainVoteNoConflict) BasicCheck() error {
	return j.QCert.BasicCheck()
}

func (j *JustDecided) BasicCheck() error {
	return j.QCert.BasicCheck()
}
