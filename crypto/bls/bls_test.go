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

func TestEncoding(t *testing.T) {
	prvData, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	pubData, _ := hex.DecodeString(
		"a7290fc800d2d14f2dc5e5cb416bebf3267dfed1c6c3a79c6edc4ebd1e657d956daa06a2fcaafd42c94b65b32d4d43ea" +
			"1368f861006829c475b7d54763a502dfd717e9d51c5cc7deae2981e56090a821c9c5bcafc129b8599203ab99031f4ce7")
	valAddrData, _ := hex.DecodeString("01c40b914373d4fc9c1e4611ad0acd5f23abf58a0d")
	accAddrData, _ := hex.DecodeString("02c40b914373d4fc9c1e4611ad0acd5f23abf58a0d")

	prvStr := "SECRET1PQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SEZYD4L"
	pubStr := "public1p5u5sljqq6tg57tw9uh95z6lt7vn8mlk3cmp608rwm38t68n90k2km2sx5t724l2ze99ktvedf4p7" +
		"5ymglpssq6pfc36m0428vwjs9h7hzl5a28zucl02u2vpu4sfp2ppe8zmet7p9xu9nysr4wvsx86vuujrva2z"
	valAddrStr := "pc1pcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsdc9v8qn"
	accAddrStr := "pc1zcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsd9wu6hw"

	prv, _ := bls.PrivateKeyFromString(prvStr)
	pub, _ := bls.PublicKeyFromString(pubStr)
	valAddr, _ := crypto.AddressFromString(valAddrStr)
	accAddr, _ := crypto.AddressFromString(accAddrStr)

	assert.Equal(t, prvData, prv.Bytes())
	assert.Equal(t, pubData, pub.Bytes())
	assert.Equal(t, valAddrData, valAddr.Bytes())
	assert.Equal(t, accAddrData, accAddr.Bytes())

	msg := []byte("pactus")
	sig, _ := bls.SignatureFromString(
		"8bdda74336efdf43b428a3811d3d6867a19e20889c91261b02a6b950b130f5bb22621394667c27660bfed2a8719d9c52")

	require.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, sig.Bytes(), prv.Sign(msg).Bytes())
	assert.True(t, pub.EqualsTo(prv.PublicKey()))
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

	require.NoError(t, pub1.Verify(msg1, sig1))
	require.NoError(t, pub2.Verify(msg1, sig2))
	require.NoError(t, pub3.Verify(msg1, sig3))
	require.Error(t, pub2.Verify(msg1, sig1))
	require.Error(t, pub3.Verify(msg1, sig1))
	require.Error(t, pub1.Verify(msg1, agg1))
	require.Error(t, pub2.Verify(msg1, agg1))
	require.Error(t, pub3.Verify(msg1, agg1))

	require.NoError(t, pubAgg1.Verify(msg1, agg1))
	require.Error(t, pubAgg1.Verify(msg2, agg1))
	require.Error(t, pubAgg1.Verify(msg1, agg2))
	require.Error(t, pubAgg2.Verify(msg1, agg1))
	require.NoError(t, pubAgg2.Verify(msg1, agg2))
	require.Error(t, pubAgg2.Verify(msg2, agg2))
	require.Error(t, pubAgg1.Verify(msg1, agg3))
	require.Error(t, pubAgg1.Verify(msg2, agg3))
	require.Error(t, pubAgg1.Verify(msg1, agg4))
	require.Error(t, pubAgg3.Verify(msg1, agg1))
	require.NoError(t, pubAgg1.Verify(msg1, agg5))
	require.NoError(t, pubAgg4.Verify(msg1, agg1))
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
		require.Error(t, err)
		assert.Nil(t, aggSig)
		assert.Contains(t, err.Error(), "no signatures provided")
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		// Point at infinity
		invalidSig, err := bls.SignatureFromString(
			"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)

		aggSig, err := bls.SignatureAggregate(invalidSig)
		require.Error(t, err)
		assert.Nil(t, aggSig)
	})

	t.Run("MixedValidAndInvalid", func(t *testing.T) {
		validSig := ts.RandBLSSignature()

		invalidSig, err := bls.SignatureFromString(
			"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)

		aggSig, err := bls.SignatureAggregate(validSig, invalidSig)
		require.Error(t, err)
		assert.Nil(t, aggSig)
	})
}

// TestPublicKeyAggregateErrorHandling tests error scenarios for PublicKeyAggregate.
func TestPublicKeyAggregateErrorHandling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("EmptyInput", func(t *testing.T) {
		aggPub, err := bls.PublicKeyAggregate()
		require.Error(t, err)
		assert.Nil(t, aggPub)
		assert.Contains(t, err.Error(), "no public keys provided")
	})

	t.Run("InvalidPublicKeyData", func(t *testing.T) {
		// Point at infinity
		invalidPub, err := bls.PublicKeyFromString(
			"public1pcqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9")
		require.NoError(t, err)

		aggPub, err := bls.PublicKeyAggregate(invalidPub)
		require.Error(t, err)
		assert.Nil(t, aggPub)
	})

	t.Run("MixedValidAndInvalid", func(t *testing.T) {
		validPub, _ := ts.RandBLSKeyPair()

		invalidPub, err := bls.PublicKeyFromString(
			"public1pcqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9")
		require.NoError(t, err)

		aggPub, err := bls.PublicKeyAggregate(validPub, invalidPub)
		require.Error(t, err)
		assert.Nil(t, aggPub)
	})
}
