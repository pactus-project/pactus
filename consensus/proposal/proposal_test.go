package proposal

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
	bz2, err := p2.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, bz2)
}

func TestProposalSignature(t *testing.T) {
	signer := crypto.GenerateTestSigner()

	p, pv := GenerateTestProposal(5, 5)
	pb := pv.PublicKey()
	assert.NoError(t, p.Verify(pb))
	assert.False(t, p.IsForBlock(crypto.GenerateTestHash()))
	assert.True(t, p.IsForBlock(p.Block().Hash()))

	assert.Error(t, p.Verify(signer.PublicKey())) // invalid public key

	signer.SignMsg(p)
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
