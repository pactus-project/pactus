package certificate

import (
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
)

// BlockCertificate represents a certificate used for block validation,
// verifying if a block is signed by a majority of validators.
type BlockCertificate struct {
	baseCertificate
}

// NewBlockCertificate creates a new BlockCertificate.
func NewBlockCertificate(height uint32, round int16) *BlockCertificate {
	return &BlockCertificate{
		baseCertificate: baseCertificate{
			height: height,
			round:  round,
		},
	}
}

func (cert *BlockCertificate) SignBytes(blockHash hash.Hash) []byte {
	signBytes := blockHash.Bytes()
	signBytes = append(signBytes, util.Uint32ToSlice(cert.height)...)
	signBytes = append(signBytes, util.Int16ToSlice(cert.round)...)

	return signBytes
}

func (cert *BlockCertificate) Validate(validators []*validator.Validator, blockHash hash.Hash) error {
	signBytes := cert.SignBytes(blockHash)

	return cert.baseCertificate.validate(validators, signBytes, require2FPower)
}

func (cert *BlockCertificate) Clone() *BlockCertificate {
	cloned := &BlockCertificate{
		baseCertificate: baseCertificate{
			height:     cert.height,
			round:      cert.round,
			committers: make([]int32, len(cert.committers)),
			absentees:  make([]int32, len(cert.absentees)),
			signature:  new(bls.Signature),
		},
	}

	copy(cloned.committers, cert.committers)
	copy(cloned.absentees, cert.absentees)
	*cloned.signature = *cert.signature

	return cloned
}
