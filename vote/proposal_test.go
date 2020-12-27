package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/crypto"
)

func TestProposalMarshaling(t *testing.T) {
	p1, _ := GenerateTestProposal(10, 10)
	bz1, err := p1.MarshalCBOR()
	assert.NoError(t, err)
	var p2 Proposal
	err = p2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	var p3 Proposal
	bz2, err := p2.MarshalJSON()
	assert.NoError(t, err)
	assert.NoError(t, p3.UnmarshalJSON(bz2))
	assert.Equal(t, p1.Hash(), p3.Hash())
}

func TestProposalSignature(t *testing.T) {
	_, pb0, pv0 := crypto.GenerateTestKeyPair()

	p, pv := GenerateTestProposal(5, 5)
	pb := pv.PublicKey()
	assert.NoError(t, p.Verify(pb))

	assert.Error(t, p.Verify(pb0)) // invalid public key

	sig0 := pv0.Sign(p.SignBytes())
	p.SetSignature(sig0)
	assert.Error(t, p.Verify(pb)) // invalid signature

	p.data.Signature = nil // No signature
	assert.Error(t, p.Verify(pb))
}
func TestProposalSanityCheck(t *testing.T) {
	p, _ := GenerateTestProposal(5, 5)
	p.data.Round = -1
	assert.Error(t, p.SanityCheck())
	p.data.Round = 0
	p.data.Height = 0
	assert.Error(t, p.SanityCheck())
	p.data.Height = 1
	p.data.Signature = nil
	assert.Error(t, p.SanityCheck())
}

func TestProposalFingerprint(t *testing.T) {
	p, _ := GenerateTestProposal(1, 1)
	assert.Contains(t, p.Fingerprint(), p.Block().Fingerprint())
}
