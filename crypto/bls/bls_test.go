package bls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestSigning(t *testing.T) {
	msg := []byte("zarb")
	prv, _ := PrivateKeyFromString("PRV1Q95DE0UXSYEA8KA56Y4QC2G8EXCF8HL0AAKNS4DVKESZAHNQ5HRDQ3KQEWF")
	pub, _ := PublicKeyFromString("af0f74917f5065af94727ae9541b0ddcfb5b828a9e016b02498f477ed37fb44d5d882495afb6fd4f9773e4ea9deee436030c4d61c6e3a1151585e1d838cae1444a438d089ce77e10c492a55f6908125c5be9b236a246e4082d08de564e111e65")
	sig, _ := SignatureFromString("ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792")
	addr, _ := crypto.AddressFromString("zc1qxset4l7ewjvvd5r9uwmpngw59xmdkuvwy7va0tl")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig1.Bytes(), sig.Bytes())
	assert.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, pub.Address(), addr)
}

func TestAggregation(t *testing.T) {
	msg := []byte("zarb")
	prv1, _ := PrivateKeyFromString("PRV1Q95DE0UXSYEA8KA56Y4QC2G8EXCF8HL0AAKNS4DVKESZAHNQ5HRDQ3KQEWF")
	prv2, _ := PrivateKeyFromString("PRV1Q9H3SH6NFLPE4PEF5WKJYS83V750XD74PA5F85XGR7343LZEHK06S57MECW")
	agg, _ := SignatureFromString("a390ffec7061827b7e89193a26841dd9e3537b5db0af55661b624e8b93b855e9f65278850002ea72fb3098e674220eca")
	sig1 := prv1.Sign(msg).(*Signature)
	sig2 := prv2.Sign(msg).(*Signature)
	assert.True(t, Aggregate([]*Signature{sig1, sig2}).EqualsTo(agg))
}

func TestAggregationFailed(t *testing.T) {
	pub1, prv1 := GenerateTestKeyPair()
	pub2, prv2 := GenerateTestKeyPair()
	pub3, prv3 := GenerateTestKeyPair()
	pub4, prv4 := GenerateTestKeyPair()
	msg1 := []byte("zarb")
	msg2 := []byte("zarb0")

	sig1 := prv1.Sign(msg1).(*Signature)
	sig11 := prv1.Sign(msg2).(*Signature)
	sig2 := prv2.Sign(msg1).(*Signature)
	sig3 := prv3.Sign(msg1).(*Signature)
	sig4 := prv4.Sign(msg1).(*Signature)

	agg1 := Aggregate([]*Signature{sig1, sig2, sig3})
	agg2 := Aggregate([]*Signature{sig1, sig2, sig4})
	agg3 := Aggregate([]*Signature{sig11, sig2, sig3})
	agg4 := Aggregate([]*Signature{sig1, sig2})
	agg5 := Aggregate([]*Signature{sig3, sig2, sig1})

	pubs1 := []*PublicKey{pub1, pub2, pub3}
	pubs2 := []*PublicKey{pub1, pub2, pub4}
	pubs3 := []*PublicKey{pub1, pub2}
	pubs4 := []*PublicKey{pub3, pub2, pub1}

	assert.NoError(t, pub1.Verify(msg1, sig1))
	assert.NoError(t, pub2.Verify(msg1, sig2))
	assert.NoError(t, pub3.Verify(msg1, sig3))
	assert.Error(t, pub2.Verify(msg1, sig1))
	assert.Error(t, pub3.Verify(msg1, sig1))
	assert.Error(t, pub1.Verify(msg1, agg1))
	assert.Error(t, pub2.Verify(msg1, agg1))
	assert.Error(t, pub3.Verify(msg1, agg1))

	assert.True(t, VerifyAggregated(agg1, pubs1, msg1))
	assert.False(t, VerifyAggregated(agg1, pubs1, msg2))
	assert.False(t, VerifyAggregated(agg2, pubs1, msg1))
	assert.False(t, VerifyAggregated(agg1, pubs2, msg1))
	assert.True(t, VerifyAggregated(agg2, pubs2, msg1))
	assert.False(t, VerifyAggregated(agg2, pubs2, msg2))
	assert.False(t, VerifyAggregated(agg3, pubs1, msg1))
	assert.False(t, VerifyAggregated(agg3, pubs1, msg2))
	assert.False(t, VerifyAggregated(agg4, pubs1, msg1))
	assert.False(t, VerifyAggregated(agg1, pubs3, msg1))
	assert.True(t, VerifyAggregated(agg5, pubs1, msg1))
	assert.True(t, VerifyAggregated(agg1, pubs4, msg1))
}

func TestAggregationOnlyOneSignature(t *testing.T) {
	_, prv1 := GenerateTestKeyPair()
	msg1 := []byte("zarb")
	sig1 := prv1.Sign(msg1).(*Signature)
	agg1 := Aggregate([]*Signature{sig1})

	assert.True(t, agg1.EqualsTo(sig1))
}

// TODO: should we check for duplication here?
func TestDuplicatedAggregation(t *testing.T) {
	pub1, prv1 := GenerateTestKeyPair()
	pub2, prv2 := GenerateTestKeyPair()

	msg1 := []byte("zarb")

	sig1 := prv1.Sign(msg1).(*Signature)
	sig2 := prv2.Sign(msg1).(*Signature)

	agg1 := Aggregate([]*Signature{sig1, sig2, sig1})
	agg2 := Aggregate([]*Signature{sig1, sig2})
	assert.False(t, agg1.EqualsTo(agg2))

	pubs1 := []*PublicKey{pub1, pub2}
	assert.False(t, VerifyAggregated(agg1, pubs1, msg1))

	pubs2 := []*PublicKey{pub1, pub2, pub1}
	assert.True(t, VerifyAggregated(agg1, pubs2, msg1))
}
