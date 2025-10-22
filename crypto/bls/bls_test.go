package bls_test

import (
	"encoding/hex"
	"testing"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigning(t *testing.T) {
	msg := []byte("pactus")
	prv, _ := bls.PrivateKeyFromString(
		"SECRET1P9QAUKRJAU7SQ7AT6ZZ6HXHYLMKPQSQYTGDL2VMH5Q5N0P5Q2QW0QL45AY3")
	pub, _ := bls.PublicKeyFromString(
		"public1p5dwsgfwmacjpuhaxhy0522j87qc5390v56ndh92f7flxge7vt3zfuxlvuwpnk7tdeed4s4l2r5nj" +
			"5zuyjfh0uzjmvrauf4t5xfvff5cpljvpqqpk7pzhv0hxfhf9gt5896vnllsf89ux8kc7anqlu7nxvvxcclw7")
	sig, _ := bls.SignatureFromString(
		"8c3ba687e8e4c016293a2c369493faa565065987544a59baba7aadae3f17ada07883552b6c7d1d7eb49f46fbdf0975c4")
	accAddr, _ := crypto.AddressFromString("pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf")
	valAddr, _ := crypto.AddressFromString("pc1p0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpf5uf6u5")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig.Bytes(), sig1.Bytes())
	assert.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, pub, prv.PublicKey())
	assert.Equal(t, valAddr, pub.ValidatorAddress())
	assert.Equal(t, accAddr, pub.AccountAddress())
}

func TestSignatureAggregate(t *testing.T) {
	msg := []byte("pactus")
	prv1, _ := bls.PrivateKeyFromString(
		"SECRET1P9QAUKRJAU7SQ7AT6ZZ6HXHYLMKPQSQYTGDL2VMH5Q5N0P5Q2QW0QL45AY3")
	prv2, _ := bls.PrivateKeyFromString(
		"SECRET1PVJHEKQ3F4NX5CA9L69CSLLNWMYWPAXDQ64ZLEQHFSV4JLFGXMXWQPDPHR0")

	sig1 := prv1.SignNative(msg)
	sig2 := prv2.SignNative(msg)
	agg, _ := bls.SignatureAggregate(sig1, sig2)
	aggExpected, _ := bls.SignatureFromString(
		"a74f05102c6217d06527cfcd1854ba6c38f4047f75a74958ad01fe66a5120c77c5416bfd875669588566670dc61f1168")

	assert.True(t, agg.EqualsTo(aggExpected))
}

func TestAggregateFailed(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, prv1 := ts.RandBLSKeyPair()
	pub2, prv2 := ts.RandBLSKeyPair()
	pub3, prv3 := ts.RandBLSKeyPair()
	pub4, prv4 := ts.RandBLSKeyPair()
	msg1 := ts.RandBytes(14)
	msg2 := ts.RandBytes(16)

	sig1 := prv1.Sign(msg1).(*bls.Signature)
	sig11 := prv1.Sign(msg2).(*bls.Signature)
	sig2 := prv2.Sign(msg1).(*bls.Signature)
	sig3 := prv3.Sign(msg1).(*bls.Signature)
	sig4 := prv4.Sign(msg1).(*bls.Signature)

	agg1, _ := bls.SignatureAggregate(sig1, sig2, sig3)
	agg2, _ := bls.SignatureAggregate(sig1, sig2, sig4)
	agg3, _ := bls.SignatureAggregate(sig11, sig2, sig3)
	agg4, _ := bls.SignatureAggregate(sig1, sig2)
	agg5, _ := bls.SignatureAggregate(sig3, sig2, sig1)

	pubs1 := []*bls.PublicKey{pub1, pub2, pub3}
	pubs2 := []*bls.PublicKey{pub1, pub2, pub4}
	pubs3 := []*bls.PublicKey{pub1, pub2}
	pubs4 := []*bls.PublicKey{pub3, pub2, pub1}

	pubAgg1, _ := bls.PublicKeyAggregate(pubs1...)
	pubAgg2, _ := bls.PublicKeyAggregate(pubs2...)
	pubAgg3, _ := bls.PublicKeyAggregate(pubs3...)
	pubAgg4, _ := bls.PublicKeyAggregate(pubs4...)

	assert.NoError(t, pub1.Verify(msg1, sig1))
	assert.NoError(t, pub2.Verify(msg1, sig2))
	assert.NoError(t, pub3.Verify(msg1, sig3))
	assert.Error(t, pub2.Verify(msg1, sig1))
	assert.Error(t, pub3.Verify(msg1, sig1))
	assert.Error(t, pub1.Verify(msg1, agg1))
	assert.Error(t, pub2.Verify(msg1, agg1))
	assert.Error(t, pub3.Verify(msg1, agg1))

	assert.NoError(t, pubAgg1.Verify(msg1, agg1))
	assert.Error(t, pubAgg1.Verify(msg2, agg1))
	assert.Error(t, pubAgg1.Verify(msg1, agg2))
	assert.Error(t, pubAgg2.Verify(msg1, agg1))
	assert.NoError(t, pubAgg2.Verify(msg1, agg2))
	assert.Error(t, pubAgg2.Verify(msg2, agg2))
	assert.Error(t, pubAgg1.Verify(msg1, agg3))
	assert.Error(t, pubAgg1.Verify(msg2, agg3))
	assert.Error(t, pubAgg1.Verify(msg1, agg4))
	assert.Error(t, pubAgg3.Verify(msg1, agg1))
	assert.NoError(t, pubAgg1.Verify(msg1, agg5))
	assert.NoError(t, pubAgg4.Verify(msg1, agg1))
}

func TestAggregateOnlyOneSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandBLSKeyPair()
	msg1 := []byte("pactus")
	sig1 := prv1.Sign(msg1).(*bls.Signature)
	agg1, _ := bls.SignatureAggregate(sig1)

	assert.True(t, agg1.EqualsTo(sig1))
}

func TestAggregateOnlyOnePublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	agg1, _ := bls.PublicKeyAggregate(pub1)

	assert.True(t, agg1.EqualsTo(pub1))
}

// TODO: should we check for duplication here?
func TestDuplicatedAggregate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, prv1 := ts.RandBLSKeyPair()
	pub2, prv2 := ts.RandBLSKeyPair()

	msg1 := []byte("pactus")

	sig1 := prv1.Sign(msg1).(*bls.Signature)
	sig2 := prv2.Sign(msg1).(*bls.Signature)

	agg1, _ := bls.SignatureAggregate(sig1, sig2, sig1)
	agg2, _ := bls.SignatureAggregate(sig1, sig2)
	assert.False(t, agg1.EqualsTo(agg2))

	pubs1 := []*bls.PublicKey{pub1, pub2}
	pubs2 := []*bls.PublicKey{pub1, pub2, pub1}
	pubAgg1, _ := bls.PublicKeyAggregate(pubs1...)
	pubAgg2, _ := bls.PublicKeyAggregate(pubs2...)
	assert.False(t, pubAgg1.EqualsTo(pubAgg2))
}

// TestHashToCurve ensures that the hash-to-curve function in kilic/bls12-381
// works as intended and is compatible with the spec.
// test vectors can be found here:
// https://datatracker.ietf.org/doc/html/rfc9380
func TestHashToCurve(t *testing.T) {
	domain := []byte("QUUX-V01-CS02-with-BLS12381G1_XMD:SHA-256_SSWU_RO_")
	tests := []struct {
		msg      string
		expected string
	}{
		{
			"",
			"052926add2207b76ca4fa57a8734416c8dc95e24501772c814278700eed6d1e4e8cf62d9c09db0fac349612b759e79a1" +
				"08ba738453bfed09cb546dbb0783dbb3a5f1f566ed67bb6be0e8c67e2e81a4cc68ee29813bb7994998f3eae0c9c6a265",
		},
		{
			"abc",
			"03567bc5ef9c690c2ab2ecdf6a96ef1c139cc0b2f284dca0a9a7943388a49a3aee664ba5379a7655d3c68900be2f6903" +
				"0b9c15f3fe6e5cf4211f346271d7b01c8f3b28be689c8429c85b67af215533311f0b8dfaaa154fa6b88176c229f2885d",
		},
		{
			"abcdef0123456789",
			"11e0b079dea29a68f0383ee94fed1b940995272407e3bb916bbf268c263ddd57a6a27200a784cbc248e84f357ce82d98" +
				"03a87ae2caf14e8ee52e51fa2ed8eefe80f02457004ba4d486d6aa1f517c0889501dc7413753f9599b099ebcbbd2d709",
		},
		{
			"q128_qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
			"15f68eaa693b95ccb85215dc65fa81038d69629f70aeee0d0f677cf22285e7bf58d7cb86eefe8f2e9bc3f8cb84fac488" +
				"1807a1d50c29f430b8cafc4f8638dfeeadf51211e1602a5f184443076715f91bb90a48ba1e370edce6ae1062f5e6dd38",
		},
		{
			"a512_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"082aabae8b7dedb0e78aeb619ad3bfd9277a2f77ba7fad20ef6aabdc6c31d19ba5a6d12283553294c1825c4b3ca2dcfe" +
				"05b84ae5a942248eea39e1d91030458c40153f3b654ab7872d779ad1e942856a20c438e8d99bc8abfbf74729ce1f7ac8",
		},
	}

	for no, tt := range tests {
		mappedPoint, _ := bls12381.HashToG1([]byte(tt.msg), domain)
		d, _ := hex.DecodeString(tt.expected)

		expectedPoint := bls12381.G1Affine{}
		err := expectedPoint.Unmarshal(d)
		require.NoError(t, err)

		assert.Equal(t, expectedPoint, mappedPoint,
			"test %v: not match", no)
	}
}

// TestSignatureAggregateErrorHandling tests error scenarios for SignatureAggregate.
func TestSignatureAggregateErrorHandling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("EmptyInput", func(t *testing.T) {
		aggSig, err := bls.SignatureAggregate()
		assert.Error(t, err)
		assert.Nil(t, aggSig)
		assert.Contains(t, err.Error(), "no signatures provided")
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		// Point at infinity
		invalidSig, err := bls.SignatureFromString(
			"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		assert.NoError(t, err)

		aggSig, err := bls.SignatureAggregate(invalidSig)
		assert.Error(t, err)
		assert.Nil(t, aggSig)
	})

	t.Run("MixedValidAndInvalid", func(t *testing.T) {
		validSig := ts.RandBLSSignature()

		invalidSig, err := bls.SignatureFromString(
			"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		assert.NoError(t, err)

		aggSig, err := bls.SignatureAggregate(validSig, invalidSig)
		assert.Error(t, err)
		assert.Nil(t, aggSig)
	})
}

// TestPublicKeyAggregateErrorHandling tests error scenarios for PublicKeyAggregate.
func TestPublicKeyAggregateErrorHandling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("EmptyInput", func(t *testing.T) {
		aggPub, err := bls.PublicKeyAggregate()
		assert.Error(t, err)
		assert.Nil(t, aggPub)
		assert.Contains(t, err.Error(), "no public keys provided")
	})

	t.Run("InvalidPublicKeyData", func(t *testing.T) {
		// Point at infinity
		invalidPub, err := bls.PublicKeyFromString(
			"public1pcqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9")
		assert.NoError(t, err)

		aggPub, err := bls.PublicKeyAggregate(invalidPub)
		assert.Error(t, err)
		assert.Nil(t, aggPub)
	})

	t.Run("MixedValidAndInvalid", func(t *testing.T) {
		validPub, _ := ts.RandBLSKeyPair()

		invalidPub, err := bls.PublicKeyFromString(
			"public1pcqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9")
		assert.NoError(t, err)

		aggPub, err := bls.PublicKeyAggregate(validPub, invalidPub)
		assert.Error(t, err)
		assert.Nil(t, aggPub)
	})
}
