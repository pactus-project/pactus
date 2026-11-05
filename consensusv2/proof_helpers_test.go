package consensusv2

import (
	"slices"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/vote"
)

// qCert builds a certificate signed only by `signers`.
// Other committee members are marked as absentees (the committee numbers are the indexes).
// `signBytes` returns the bytes that each signer must sign.
func (td *testData) qCert(t *testing.T, height types.Height, round types.Round,
	signers []int, signBytes func(cert *certificate.Certificate) []byte,
) *certificate.Certificate {
	t.Helper()

	cert := certificate.NewCertificate(height, round)
	signed := signBytes(cert)
	committers := make([]int32, 0, len(td.valKeys))
	absentees := make([]int32, 0)
	sigs := make([]*bls.Signature, 0)
	for i := range td.valKeys {
		committers = append(committers, int32(i))
		if slices.Contains(signers, i) {
			sigs = append(sigs, td.valKeys[i].Sign(signed))
		} else {
			absentees = append(absentees, int32(i))
		}
	}
	aggSig, _ := bls.SignatureAggregate(sigs...)
	cert.SetSignature(committers, absentees, aggSig)

	return cert
}

// precommitCert returns a certificate of precommit votes for `blockHash`.
func (td *testData) precommitCert(t *testing.T, blockHash hash.Hash,
	height types.Height, round types.Round, signers ...int,
) *certificate.Certificate {
	t.Helper()

	return td.qCert(t, height, round, signers, func(c *certificate.Certificate) []byte {
		return c.SignBytesPrecommit(blockHash)
	})
}

// preVoteCert returns a certificate of cp:pre-vote votes.
func (td *testData) preVoteCert(t *testing.T, blockHash hash.Hash,
	height types.Height, round types.Round, cpRound int16, value vote.CPValue, signers ...int,
) *certificate.Certificate {
	t.Helper()

	return td.qCert(t, height, round, signers, func(c *certificate.Certificate) []byte {
		return c.SignBytesCPPreVote(blockHash, cpRound, byte(value))
	})
}

// send signs `v` with the key of validator `valID` and delivers it to `cons`.
func (td *testData) send(t *testing.T, cons *consensusV2, v *vote.Vote, valID int) {
	t.Helper()

	td.addVote(t, cons, v, valID)
}
