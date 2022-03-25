package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestVoteMarshaling(t *testing.T) {
	v1, _ := GenerateTestPrepareVote(10, 10)

	bz1, err := v1.MarshalCBOR()
	assert.NoError(t, err)
	var v2 Vote
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
	h1 := hash.GenerateTestHash()
	pb1, pv1 := bls.GenerateTestKeyPair()
	pb2, pv2 := bls.GenerateTestKeyPair()

	v1 := NewVote(VoteTypePrepare, 101, 5, h1, pb1.Address())
	v2 := NewVote(VoteTypePrepare, 101, 5, h1, pb2.Address())

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

func TestVoteSanityCheck(t *testing.T) {
	v, _ := GenerateTestChangeProposerVote(5, 5)
	assert.NoError(t, v.SanityCheck())
	v.data.Type = 4
	assert.Error(t, v.SanityCheck())
	v.data.Type = VoteTypePrepare
	v.data.Round = -1
	assert.Error(t, v.SanityCheck())
	v.data.Round = 0
	v.data.Height = 0
	assert.Error(t, v.SanityCheck())
	v.data.Height = 1
	v.data.Signature = nil
	assert.Error(t, v.SanityCheck())
}

func TestSignBytes(t *testing.T) {
	signer := crypto.GenerateTestAddress()
	blockHash := hash.GenerateTestHash()
	height := util.RandInt32(100000)
	round := util.RandInt16(10)

	v1 := NewVote(VoteTypePrepare, height, round, blockHash, signer)
	v2 := NewVote(VoteTypeChangeProposer, height, round, blockHash, signer)
	v3 := NewVote(VoteTypePrecommit, height, round, blockHash, signer)

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
