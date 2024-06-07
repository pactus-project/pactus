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
func NewBlockCertificate(height uint32, round int16, fastPath bool) *BlockCertificate {
	return &BlockCertificate{
		baseCertificate: baseCertificate{
			height:   height,
			round:    round,
			fastPath: fastPath,
		},
	}
}

func (cert *BlockCertificate) SignBytes(blockHash hash.Hash) []byte {
	sb := blockHash.Bytes()
	sb = append(sb, util.Uint32ToSlice(cert.height)...)
	sb = append(sb, util.Int16ToSlice(cert.round)...)

	if cert.fastPath {
		sb = append(sb, util.StringToBytes("PREPARE")...)
	}

	return sb
}

func (cert *BlockCertificate) Validate(validators []*validator.Validator, blockHash hash.Hash) error {
	calcRequiredPowerFn := func(committeePower int64) int64 {
		t := (committeePower - 1) / 5
		p := (3 * t) + 1
		if cert.fastPath {
			p = (4 * t) + 1
		}

		return p
	}

	signBytes := cert.SignBytes(blockHash)

	return cert.baseCertificate.validate(validators, signBytes, calcRequiredPowerFn)
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
