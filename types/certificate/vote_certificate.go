package certificate

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
)

// VoteCertificate represents a certificate used for consensus voting,
// checking if a majority of validators have voted in a consensus step.
type VoteCertificate struct {
	baseCertificate
}

// NewVoteCertificate creates a new VoteCertificate instance.
func NewVoteCertificate(height uint32, round int16) *VoteCertificate {
	return &VoteCertificate{
		baseCertificate: baseCertificate{
			height: height,
			round:  round,
		},
	}
}

// SignBytes returns the sign bytes for the vote certificate.
// This method provides the same data as the `SignBytes` function in vote struct.
func (cert *VoteCertificate) SignBytes(blockHash hash.Hash, extraData ...[]byte) []byte {
	sb := blockHash.Bytes()
	sb = append(sb, util.Uint32ToSlice(cert.height)...)
	sb = append(sb, util.Int16ToSlice(cert.round)...)
	for _, data := range extraData {
		sb = append(sb, data...)
	}

	return sb
}

func (cert *VoteCertificate) ValidatePrepare(validators []*validator.Validator,
	blockHash hash.Hash,
) error {
	signBytes := cert.SignBytes(blockHash,
		util.StringToBytes("PREPARE"))

	return cert.validate(validators, signBytes)
}

func (cert *VoteCertificate) ValidateCPPreVote(validators []*validator.Validator,
	blockHash hash.Hash, cpRound int16, cpValue byte,
) error {
	signBytes := cert.SignBytes(blockHash,
		util.StringToBytes("PRE-VOTE"),
		util.Int16ToSlice(cpRound),
		[]byte{cpValue})

	return cert.validate(validators, signBytes)
}

func (cert *VoteCertificate) ValidateCPMainVote(validators []*validator.Validator,
	blockHash hash.Hash, cpRound int16, cpValue byte,
) error {
	signBytes := cert.SignBytes(blockHash,
		util.StringToBytes("MAIN-VOTE"),
		util.Int16ToSlice(cpRound),
		[]byte{cpValue})

	return cert.validate(validators, signBytes)
}

func (cert *VoteCertificate) validate(validators []*validator.Validator, signBytes []byte) error {
	return cert.baseCertificate.validate(validators, signBytes, require2Fp1Power)
}
