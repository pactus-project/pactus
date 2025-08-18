package vote

import (
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

type ConflictReason uint8

const (
	ConflictVotes  = ConflictReason(1)
	ConflictBiases = ConflictReason(2)
)

func (r ConflictReason) String() string {
	switch r {
	case ConflictVotes:
		return "conflict-votes"
	case ConflictBiases:
		return "conflict-biases"
	default:
		return "Unknown"
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
	JustNo  Just
	JustYes Just
}

type JustMainVoteConflictV2 struct {
	Just1 Just
	Just2 Just
}

type JustMainVoteNoConflict struct {
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
}

type JustDecided struct {
	QCert *certificate.VoteCertificate `cbor:"1,keyasint"`
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
	return j.QCert.BasicCheck()
}

func (j *JustPreVoteHard) BasicCheck() error {
	return j.QCert.BasicCheck()
}

func (j *JustMainVoteConflict) BasicCheck() error {
	if err := j.JustNo.BasicCheck(); err != nil {
		return err
	}

	return j.JustYes.BasicCheck()
}

func (j *JustMainVoteNoConflict) BasicCheck() error {
	return j.QCert.BasicCheck()
}

func (j *JustDecided) BasicCheck() error {
	return j.QCert.BasicCheck()
}
