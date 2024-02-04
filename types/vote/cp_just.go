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
		return &JustInitZero{}
	case JustTypeInitOne:
		return &JustInitOne{}
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

type JustInitZero struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}
type JustInitOne struct {
	//
}

type JustPreVoteSoft struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

type JustPreVoteHard struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}
type JustMainVoteConflict struct {
	Just0 Just
	Just1 Just
}
type JustMainVoteNoConflict struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

type JustDecided struct {
	QCert *certificate.Certificate `cbor:"1,keyasint"`
}

func (j *JustInitZero) Type() JustType {
	return JustTypeInitZero
}

func (j *JustInitOne) Type() JustType {
	return JustTypeInitOne
}

func (j *JustPreVoteSoft) Type() JustType {
	return JustTypePreVoteSoft
}

func (j *JustPreVoteHard) Type() JustType {
	return JustTypePreVoteHard
}

func (j *JustMainVoteConflict) Type() JustType {
	return JustTypeMainVoteConflict
}

func (j *JustMainVoteNoConflict) Type() JustType {
	return JustTypeMainVoteNoConflict
}

func (j *JustDecided) Type() JustType {
	return JustTypeDecided
}

func (j *JustInitZero) BasicCheck() error {
	return j.QCert.BasicCheck()
}

func (j *JustInitOne) BasicCheck() error {
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
