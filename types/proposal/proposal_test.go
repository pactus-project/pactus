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

	prop := ts.GenerateTestProposal(10, 10)
	data, err := prop.MarshalCBOR()
	assert.NoError(t, err)
	var p2 proposal.Proposal
	err = p2.UnmarshalCBOR(data)
	assert.NoError(t, err)
}

func TestProposalSignBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	prop := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound())
	sb := prop.Block().Hash().Bytes()
	sb = append(sb, util.Uint32ToSlice(prop.Height())...)
	sb = append(sb, util.Int16ToSlice(prop.Round())...)

	assert.Equal(t, sb, prop.SignBytes())
	assert.Equal(t, hash.CalcHash(sb), prop.Hash())
}

func TestIsForBlock(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	prop := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound())
	assert.False(t, prop.IsForBlock(ts.RandHash()))
	assert.True(t, prop.IsForBlock(prop.Block().Hash()))
}

func TestProposalSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	proposerKey := ts.RandValKey()
	prop := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound(), testsuite.ProposalWithKey(proposerKey))

	err := prop.Verify(proposerKey.PublicKey())
	assert.NoError(t, err)

	rndValKey := ts.RandValKey()
	err = prop.Verify(rndValKey.PublicKey())
	assert.ErrorIs(t, err, crypto.AddressMismatchError{
		Expected: rndValKey.Address(),
		Got:      prop.Block().Header().ProposerAddress(),
	})

	ts.HelperSignProposal(rndValKey, prop)
	err = prop.Verify(proposerKey.PublicKey())
	assert.ErrorIs(t, crypto.ErrInvalidSignature, err)
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No block", func(t *testing.T) {
		p := &proposal.Proposal{}
		err := p.BasicCheck()
		assert.ErrorIs(t, err, proposal.BasicCheckError{
			Reason: "no block",
		})
	})

	t.Run("Invalid height", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		p := proposal.NewProposal(0, 0, blk)
		p.SetSignature(ts.RandBLSSignature())
		err := p.BasicCheck()
		assert.ErrorIs(t, err, proposal.BasicCheckError{
			Reason: "invalid height",
		})
	})

	t.Run("Invalid round", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		p := proposal.NewProposal(ts.RandHeight(), -1, blk)
		p.SetSignature(ts.RandBLSSignature())
		err := p.BasicCheck()
		assert.ErrorIs(t, err, proposal.BasicCheckError{
			Reason: "invalid round",
		})
	})

	t.Run("No signature", func(t *testing.T) {
		pub, _ := ts.RandBLSKeyPair()
		data := ts.DecodingHex(
			"a401186402000358c80140da9b641551048b59a859946ca7f9ab95c9cf84da488a1a5c49ba643b29b653dc223bc20a4e9ff03158165f3d42" +
				"4e2a74677bfe24a7295d1ce2e55ca3644cbe9a5a5e7d913b8e1ba6a020afbd5a25024a12b37cf8e1ed0b9498f91d75b294db0f95123d8593" +
				"05aa5deea3d4216777e74310b6a601bb4d4d6b13c9b295781ab1533aea032978d4f89305000000010004060f1b23010fab4f72234cc7c120" +
				"48bbbc616c005573d8ad4d5c6997996d6f488946cdd78410f0a400c4a7f9bdb41506bdf717a892fa0004f6")
		prop := &proposal.Proposal{}
		err := cbor.Unmarshal(data, &prop)
		assert.NoError(t, err)

		err = prop.BasicCheck()
		assert.ErrorIs(t, err, proposal.BasicCheckError{
			Reason: "no signature",
		})

		err = prop.Verify(pub)
		assert.Error(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Ok", func(t *testing.T) {
		prop := ts.GenerateTestProposal(100, 0)
		assert.NoError(t, prop.BasicCheck())
	})
}
