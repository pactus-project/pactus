package bls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregation(t *testing.T) {
	msg := []byte("zarb")
	prv1, _ := PrivateKeyFromString("68dcbf868133d3dbb4d12a0c2907c9b093dfefef6d3855acb6602ede60a5c6d0")
	prv2, _ := PrivateKeyFromString("6f185f534fc39a8729a3ad2240f167a8f337d50f6893d0c81fa358fc59bd9fa8")
	agg, _ := SignatureFromString("959ad0810024121fc3d537053b33963bde9d372b2c940d262c53aeec876170ed91a6a42b24f95fae1d526c788047efbe")
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

	assert.True(t, pub1.Verify(msg1, sig1))
	assert.True(t, pub2.Verify(msg1, sig2))
	assert.True(t, pub3.Verify(msg1, sig3))
	assert.False(t, pub2.Verify(msg1, sig1))
	assert.False(t, pub3.Verify(msg1, sig1))
	assert.False(t, pub1.Verify(msg1, agg1))
	assert.False(t, pub2.Verify(msg1, agg1))
	assert.False(t, pub3.Verify(msg1, agg1))

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
	assert.Equal(t, agg1.RawBytes(), sig1.RawBytes())
}

func TestAggregateTheAggregated(t *testing.T) {
	pub1, prv1 := GenerateTestKeyPair()
	pub2, prv2 := GenerateTestKeyPair()
	pub3, prv3 := GenerateTestKeyPair()

	msg1 := []byte("zarb")

	sig1 := prv1.Sign(msg1).(*Signature)
	sig2 := prv2.Sign(msg1).(*Signature)
	sig3 := prv3.Sign(msg1).(*Signature)

	agg1 := Aggregate([]*Signature{sig1, sig2, sig3})
	agg2 := Aggregate([]*Signature{sig1, sig2})
	agg3 := Aggregate([]*Signature{agg2, sig3})

	assert.Equal(t, agg1.RawBytes(), agg3.RawBytes())

	pubs2 := []*PublicKey{pub1, pub2}
	pubs3 := []*PublicKey{pub1, pub2, pub3}

	assert.True(t, VerifyAggregated(agg3, pubs3, msg1))
	assert.True(t, VerifyAggregated(agg2, pubs2, msg1))
	assert.False(t, VerifyAggregated(agg3, pubs2, msg1))
}

func TestCrossAggregation(t *testing.T) {
	pub1, prv1 := GenerateTestKeyPair()
	pub2, prv2 := GenerateTestKeyPair()

	msg1 := []byte("zarb")

	sig1 := prv1.Sign(msg1).(*Signature)
	sig2 := prv2.Sign(msg1).(*Signature)

	agg1 := Aggregate([]*Signature{sig1, sig2})
	agg2 := Aggregate([]*Signature{sig2, sig1})

	assert.Equal(t, agg1.RawBytes(), agg2.RawBytes())

	pubs1 := []*PublicKey{pub1, pub2}
	pubs2 := []*PublicKey{pub2, pub1}

	assert.True(t, VerifyAggregated(agg1, pubs2, msg1))
	assert.True(t, VerifyAggregated(agg2, pubs1, msg1))
}

func TestDuplicatedAggregation(t *testing.T) {
	pub1, prv1 := GenerateTestKeyPair()
	pub2, prv2 := GenerateTestKeyPair()

	msg1 := []byte("zarb")

	sig1 := prv1.Sign(msg1).(*Signature)
	sig2 := prv2.Sign(msg1).(*Signature)

	agg1 := Aggregate([]*Signature{sig1, sig2, sig1})

	pubs1 := []*PublicKey{pub1, pub2}

	assert.False(t, VerifyAggregated(agg1, pubs1, msg1))
}

func TestZeroKeys(t *testing.T) {
	prv, _ := PrivateKeyFromString("0000000000000000000000000000000000000000000000000000000000000000")
	msg := []byte("test")
	sig := prv.Sign(msg)
	assert.False(t, prv.PublicKey().Verify(msg, sig))

	_, err := PublicKeyFromString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.Error(t, err)
	_, err = SignatureFromString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.Error(t, err)

}
