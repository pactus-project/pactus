package proposal_test

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/util"
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

	p, _ := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound())
	sb := p.Block().Hash().Bytes()
	sb = append(sb, util.Uint32ToSlice(p.Height())...)
	sb = append(sb, util.Int16ToSlice(p.Round())...)

	assert.Equal(t, sb, p.SignBytes())
	assert.Equal(t, hash.CalcHash(sb), p.Hash())
}

func TestProposalSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	rndValKey := ts.RandValKey()
	p, valKey := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound())
	pub := valKey.PublicKey()
	assert.NoError(t, p.Verify(pub))
	assert.False(t, p.IsForBlock(ts.RandHash()))
	assert.True(t, p.IsForBlock(p.Block().Hash()))

	err := p.Verify(rndValKey.PublicKey())
	assert.ErrorIs(t, err, crypto.AddressMismatchError{
		Expected: rndValKey.Address(),
		Got:      valKey.Address(),
	})

	ts.HelperSignProposal(rndValKey, p)
	err = p.Verify(pub)
	assert.ErrorIs(t, crypto.ErrInvalidSignature, err)
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No block", func(t *testing.T) {
		p := &proposal.Proposal{}
		assert.Error(t, p.BasicCheck())
	})

	t.Run("Invalid height", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		p := proposal.NewProposal(0, 0, blk)
		p.SetSignature(ts.RandBLSSignature())
		assert.Error(t, p.BasicCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		p := proposal.NewProposal(ts.RandHeight(), -1, blk)
		p.SetSignature(ts.RandBLSSignature())
		assert.Error(t, p.BasicCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		pub, _ := ts.RandBLSKeyPair()
		d := ts.DecodingHex(
			"a401186402000358c80140da9b641551048b59a859946ca7f9ab95c9cf84da488a1a5c49ba643b29b653dc223bc20a4e9ff03158165f3d42" +
				"4e2a74677bfe24a7295d1ce2e55ca3644cbe9a5a5e7d913b8e1ba6a020afbd5a25024a12b37cf8e1ed0b9498f91d75b294db0f95123d8593" +
				"05aa5deea3d4216777e74310b6a601bb4d4d6b13c9b295781ab1533aea032978d4f89305000000010004060f1b23010fab4f72234cc7c120" +
				"48bbbc616c005573d8ad4d5c6997996d6f488946cdd78410f0a400c4a7f9bdb41506bdf717a892fa0004f6")
		p := &proposal.Proposal{}
		err := cbor.Unmarshal(d, &p)

		assert.NoError(t, err)
		assert.Error(t, p.BasicCheck())
		assert.Error(t, p.Verify(pub))
	})

	t.Run("Ok", func(t *testing.T) {
		p, _ := ts.GenerateTestProposal(100, 0)
		assert.NoError(t, p.BasicCheck())
	})
}
