package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestVoteMarshaling(t *testing.T) {
	v1, _ := GenerateTestPrecommitVote(10, 10)

	bz1, err := v1.MarshalCBOR()
	assert.NoError(t, err)
	var v2 Vote
	err = v2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	bz2, _ := v2.MarshalCBOR()

	assert.Equal(t, bz1, bz2)
	assert.Equal(t, v1.Hash(), v2.Hash())
}

func TestVoteSignature(t *testing.T) {
	h1 := crypto.GenerateTestHash()
	addr1, pb1, pv1 := crypto.GenerateTestKeyPair()
	addr2, pb2, pv2 := crypto.GenerateTestKeyPair()

	v1 := NewVote(VoteTypePrepare, 101, 5, h1, addr1)
	v2 := NewVote(VoteTypePrepare, 101, 5, h1, addr2)

	sig1 := pv1.Sign(v1.SignBytes())
	assert.Error(t, v1.Verify(pb1)) // No signature

	sig2 := pv2.Sign(v2.SignBytes())
	v1.SetSignature(sig1)
	assert.NoError(t, v1.Verify(pb1))

	v2.SetSignature(sig2)
	assert.Error(t, v2.Verify(pb1)) // invalid public key

	sig3 := pv1.Sign(v2.SignBytes())
	v2.SetSignature(sig3)
	assert.Error(t, v2.Verify(pb2)) // invalid signature
}

func TestVoteFingerprint(t *testing.T) {
	v, _ := GenerateTestPrecommitVote(1, 1)
	assert.Contains(t, v.Fingerprint(), v.Signer().Fingerprint())
}
