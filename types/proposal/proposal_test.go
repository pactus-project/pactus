package proposal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestProposalMarshaling(t *testing.T) {
	p1, _ := GenerateTestProposal(10, 10)
	bz1, err := p1.MarshalCBOR()
	assert.NoError(t, err)
	var p2 Proposal
	err = p2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
}

func TestProposalSignBytes(t *testing.T) {
	p, _ := GenerateTestProposal(util.RandInt32(100000), util.RandInt16(10))
	sb := p.Block().Hash().Bytes()
	sb = append(sb, util.Int32ToSlice(p.Height())...)
	sb = append(sb, util.Int16ToSlice(p.Round())...)

	assert.Equal(t, sb, p.SignBytes())
	assert.Equal(t, hash.CalcHash(sb), p.Hash())
}

func TestProposalSignature(t *testing.T) {
	signer := bls.GenerateTestSigner()

	p, prv := GenerateTestProposal(util.RandInt32(100000), util.RandInt16(10))
	pub := prv.PublicKey()
	assert.NoError(t, p.Verify(pub))
	assert.False(t, p.IsForBlock(hash.GenerateTestHash()))
	assert.True(t, p.IsForBlock(p.Block().Hash()))

	err := p.Verify(signer.PublicKey())
	assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)

	signer.SignMsg(p)
	err = p.Verify(pub)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)

	p.data.Signature = nil // No signature
	assert.Error(t, p.Verify(pub))
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
