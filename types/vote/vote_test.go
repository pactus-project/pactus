package vote_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestVoteMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	v1, _ := ts.GenerateTestPrepareVote(10, 10)

	bz1, err := v1.MarshalCBOR()
	assert.NoError(t, err)
	var v2 vote.Vote
	err = v2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	bz2, _ := v2.MarshalCBOR()

	assert.Equal(t, bz1, bz2)
	assert.Equal(t, v1.Hash(), v2.Hash())
	assert.Equal(t, v1.Height(), v2.Height())
	assert.Equal(t, v1.Round(), v2.Round())
	assert.Equal(t, v1.BlockHash(), v2.BlockHash())
	assert.Equal(t, v1.Signer(), v2.Signer())
}

func TestVoteSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h1 := ts.RandomHash()
	pb1, pv1 := ts.RandomBLSKeyPair()
	pb2, pv2 := ts.RandomBLSKeyPair()

	v1 := vote.NewVote(vote.VoteTypePrepare, 101, 5, h1, pb1.Address())
	v2 := vote.NewVote(vote.VoteTypePrepare, 101, 5, h1, pb2.Address())

	assert.Error(t, v1.Verify(pb1), "No signature")

	sig1 := pv1.Sign(v1.SignBytes())
	v1.SetSignature(sig1)
	assert.NoError(t, v1.Verify(pb1), "Ok")

	sig2 := pv2.Sign(v2.SignBytes())
	v2.SetSignature(sig2)
	assert.Error(t, v2.Verify(pb1), "invalid public key")

	sig3 := pv1.Sign(v2.SignBytes())
	v2.SetSignature(sig3)
	assert.Error(t, v2.Verify(pb2), "invalid signature")
}

func TestSanityCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid type", func(t *testing.T) {
		v := vote.NewVote(4, 100, 0, ts.RandomHash(), ts.RandomAddress())

		err := v.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidVote)
	})

	t.Run("Invalid height", func(t *testing.T) {
		v := vote.NewVote(vote.VoteTypePrepare, 0, 0, ts.RandomHash(), ts.RandomAddress())

		err := v.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	t.Run("Invalid round", func(t *testing.T) {
		v := vote.NewVote(vote.VoteTypePrepare, 100, -1, ts.RandomHash(), ts.RandomAddress())

		err := v.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidRound)
	})

	t.Run("No signature", func(t *testing.T) {
		v := vote.NewVote(vote.VoteTypePrepare, 100, 0, ts.RandomHash(), ts.RandomAddress())

		err := v.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidVote)
	})

	t.Run("Ok", func(t *testing.T) {
		v, _ := ts.GenerateTestChangeProposerVote(5, 5)
		assert.NoError(t, v.SanityCheck())
	})
}

func TestSignBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signer := ts.RandomAddress()
	blockHash := ts.RandomHash()
	height := ts.RandUint32(100000)
	round := ts.RandInt16(10)

	v1 := vote.NewVote(vote.VoteTypePrepare, height, round, blockHash, signer)
	v2 := vote.NewVote(vote.VoteTypeChangeProposer, height, round, blockHash, signer)
	v3 := vote.NewVote(vote.VoteTypePrecommit, height, round, blockHash, signer)

	sb1 := v1.SignBytes()
	sb2 := v2.SignBytes()
	sb3 := v3.SignBytes()
	sb4 := block.CertificateSignBytes(blockHash, round)

	assert.Contains(t, string(sb1), "prepare")
	assert.Contains(t, string(sb2), "change-proposer")
	assert.Equal(t, len(sb3), 34)
	assert.NotEqual(t, sb1, sb2)
	assert.NotEqual(t, sb1, sb3)
	assert.NotEqual(t, sb2, sb3)
	assert.Equal(t, sb3, sb4)
}
