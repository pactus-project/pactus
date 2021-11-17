package vote

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
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
	addr1, pb1, pv1 := bls.GenerateTestKeyPair()
	addr2, pb2, pv2 := bls.GenerateTestKeyPair()

	v1 := NewVote(VoteTypePrepare, 101, 5, h1, addr1)
	v2 := NewVote(VoteTypePrepare, 101, 5, h1, addr2)

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
	v1, _ := GenerateTestPrepareVote(10, 1)
	v2, _ := GenerateTestPrecommitVote(10, 1)

	bz1 := v1.SignBytes()
	bz2 := v2.SignBytes()

	assert.Contains(t, string(bz1), "prepare")
	assert.NotContains(t, string(bz2), "prepare")
}

func TestSignBytesMatchWithCommit(t *testing.T) {
	// Find this data in commit tests
	d, _ := hex.DecodeString("a20158201c8f67440c5d2fcaec3176cde966e8b46ec744c836f643612bec96eb6a83c1fe0206")
	s := new(signVote)
	assert.NoError(t, cbor.Unmarshal(d, s))
	v := Vote{data: voteData{
		Type:      VoteTypePrecommit,
		Round:     s.Round,
		BlockHash: s.BlockHash},
	}

	fmt.Printf("%x", v.SignBytes())
	assert.Equal(t, v.SignBytes(), d)
}
