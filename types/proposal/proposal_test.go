package proposal_test

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestProposalMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	p1, _ := ts.GenerateTestProposal(10, 10)
	bz1, err := p1.MarshalCBOR()
	assert.NoError(t, err)
	var p2 proposal.Proposal
	err = p2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
}

func TestProposalSignBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	p, _ := ts.GenerateTestProposal(ts.RandUint32(100000), ts.RandInt16(10))
	sb := p.Block().Hash().Bytes()
	sb = append(sb, util.Uint32ToSlice(p.Height())...)
	sb = append(sb, util.Int16ToSlice(p.Round())...)

	assert.Equal(t, sb, p.SignBytes())
	assert.Equal(t, hash.CalcHash(sb), p.Hash())
}

func TestProposalSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signer := ts.RandomSigner()
	p, prv := ts.GenerateTestProposal(ts.RandUint32(100000), ts.RandInt16(10))
	pub := prv.PublicKey()
	assert.NoError(t, p.Verify(pub))
	assert.False(t, p.IsForBlock(ts.RandomHash()))
	assert.True(t, p.IsForBlock(p.Block().Hash()))

	err := p.Verify(signer.PublicKey())
	assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)

	signer.SignMsg(p)
	err = p.Verify(pub)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
}

func TestSanityCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No block", func(t *testing.T) {
		p := &proposal.Proposal{}
		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid height", func(t *testing.T) {
		p, _ := ts.GenerateTestProposal(0, 0)
		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p, _ := ts.GenerateTestProposal(1, -1)
		assert.Error(t, p.SanityCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		pub, _ := ts.RandomBLSKeyPair()
		d, _ := hex.DecodeString(
			"a401186402000358c30140da9b641551048b59a859946ca7f9ab95c9cf84da488a1a5c49ba643b29b653dc223bc20a4e9ff03158165f3d42" +
				"4e2a74677bfe24a7295d1ce2e55ca3644cbe9a5a5e7d913b8e1ba6a020afbd5a25024a12b37cf8e1ed0b9498f91d75b294db0f95123d8593" +
				"05aa5deea3d4216777e74310b6a601bb4d4d6b13c9b295781ab1533aea032978d4f8930504060f1b23010fab4f72234cc7c12048bbbc616c" +
				"005573d8ad4d5c6997996d6f488946cdd78410f0a400c4a7f9bdb41506bdf717a892fa0004f6")
		p := &proposal.Proposal{}
		err := cbor.Unmarshal(d, &p)

		assert.NoError(t, err)
		assert.Error(t, p.SanityCheck())
		assert.Error(t, p.Verify(pub))
	})

	t.Run("Ok", func(t *testing.T) {
		p, _ := ts.GenerateTestProposal(100, 0)
		assert.NoError(t, p.SanityCheck())
	})
}
