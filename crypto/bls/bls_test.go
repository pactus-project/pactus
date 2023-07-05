package bls_test

import (
	"encoding/hex"
	"testing"

	bls12381 "github.com/kilic/bls12-381"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestSigning(t *testing.T) {
	msg := []byte("zarb")
	prv, _ := bls.PrivateKeyFromString(
		"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67")
	pub, _ := bls.PublicKeyFromString(
		"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
			"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx47a")
	sig, _ := bls.SignatureFromString(
		"ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792")
	addr, _ := crypto.AddressFromString("pc1p5x2a0lkt5nrrdqe0rkcv6r4pfkmdhrr3xk73tq")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig1.Bytes(), sig.Bytes())
	assert.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, pub.Address(), addr)
}

func TestSignatureAggregate(t *testing.T) {
	msg := []byte("zarb")
	prv1, _ := bls.PrivateKeyFromString(
		"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67")
	prv2, _ := bls.PrivateKeyFromString(
		"SECRET1PDUV97560CWDGW2DR453YPUT84REN04G0DZFAPJQL5DV0CKDAN75QCJEV6F")
	agg, _ := bls.SignatureFromString(
		"a390ffec7061827b7e89193a26841dd9e3537b5db0af55661b624e8b93b855e9f65278850002ea72fb3098e674220eca")
	sig1 := prv1.Sign(msg).(*bls.Signature)
	sig2 := prv2.Sign(msg).(*bls.Signature)

	assert.True(t, bls.SignatureAggregate([]*bls.Signature{sig1, sig2}).EqualsTo(agg))
	assert.False(t, prv1.EqualsTo(prv2))
}

func TestAggregateFailed(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, prv1 := ts.RandomBLSKeyPair()
	pub2, prv2 := ts.RandomBLSKeyPair()
	pub3, prv3 := ts.RandomBLSKeyPair()
	pub4, prv4 := ts.RandomBLSKeyPair()
	msg1 := []byte("zarb")
	msg2 := []byte("zarb0")

	sig1 := prv1.Sign(msg1).(*bls.Signature)
	sig11 := prv1.Sign(msg2).(*bls.Signature)
	sig2 := prv2.Sign(msg1).(*bls.Signature)
	sig3 := prv3.Sign(msg1).(*bls.Signature)
	sig4 := prv4.Sign(msg1).(*bls.Signature)

	agg1 := bls.SignatureAggregate([]*bls.Signature{sig1, sig2, sig3})
	agg2 := bls.SignatureAggregate([]*bls.Signature{sig1, sig2, sig4})
	agg3 := bls.SignatureAggregate([]*bls.Signature{sig11, sig2, sig3})
	agg4 := bls.SignatureAggregate([]*bls.Signature{sig1, sig2})
	agg5 := bls.SignatureAggregate([]*bls.Signature{sig3, sig2, sig1})

	pubs1 := []*bls.PublicKey{pub1, pub2, pub3}
	pubs2 := []*bls.PublicKey{pub1, pub2, pub4}
	pubs3 := []*bls.PublicKey{pub1, pub2}
	pubs4 := []*bls.PublicKey{pub3, pub2, pub1}

	pubAgg1 := bls.PublicKeyAggregate(pubs1)
	pubAgg2 := bls.PublicKeyAggregate(pubs2)
	pubAgg3 := bls.PublicKeyAggregate(pubs3)
	pubAgg4 := bls.PublicKeyAggregate(pubs4)

	assert.NoError(t, pub1.Verify(msg1, sig1))
	assert.NoError(t, pub2.Verify(msg1, sig2))
	assert.NoError(t, pub3.Verify(msg1, sig3))
	assert.Error(t, pub2.Verify(msg1, sig1))
	assert.Error(t, pub3.Verify(msg1, sig1))
	assert.Error(t, pub1.Verify(msg1, agg1))
	assert.Error(t, pub2.Verify(msg1, agg1))
	assert.Error(t, pub3.Verify(msg1, agg1))

	assert.True(t, bls.VerifyAggregated(agg1, pubs1, msg1))
	assert.False(t, bls.VerifyAggregated(agg1, pubs1, msg2))
	assert.False(t, bls.VerifyAggregated(agg2, pubs1, msg1))
	assert.False(t, bls.VerifyAggregated(agg1, pubs2, msg1))
	assert.True(t, bls.VerifyAggregated(agg2, pubs2, msg1))
	assert.False(t, bls.VerifyAggregated(agg2, pubs2, msg2))
	assert.False(t, bls.VerifyAggregated(agg3, pubs1, msg1))
	assert.False(t, bls.VerifyAggregated(agg3, pubs1, msg2))
	assert.False(t, bls.VerifyAggregated(agg4, pubs1, msg1))
	assert.False(t, bls.VerifyAggregated(agg1, pubs3, msg1))
	assert.True(t, bls.VerifyAggregated(agg5, pubs1, msg1))
	assert.True(t, bls.VerifyAggregated(agg1, pubs4, msg1))

	assert.Nil(t, pubAgg1.Verify(msg1, agg1))
	assert.NotNil(t, pubAgg1.Verify(msg2, agg1))
	assert.NotNil(t, pubAgg1.Verify(msg1, agg2))
	assert.NotNil(t, pubAgg2.Verify(msg1, agg1))
	assert.Nil(t, pubAgg2.Verify(msg1, agg2))
	assert.NotNil(t, pubAgg2.Verify(msg2, agg2))
	assert.NotNil(t, pubAgg1.Verify(msg1, agg3))
	assert.NotNil(t, pubAgg1.Verify(msg2, agg3))
	assert.NotNil(t, pubAgg1.Verify(msg1, agg4))
	assert.NotNil(t, pubAgg3.Verify(msg1, agg1))
	assert.Nil(t, pubAgg1.Verify(msg1, agg5))
	assert.Nil(t, pubAgg4.Verify(msg1, agg1))
}

func TestAggregateNil(t *testing.T) {
	assert.Nil(t, bls.SignatureAggregate(nil))
	assert.Nil(t, bls.PublicKeyAggregate(nil))
}

func TestAggregateOnlyOneSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandomBLSKeyPair()
	msg1 := []byte("zarb")
	sig1 := prv1.Sign(msg1).(*bls.Signature)
	agg1 := bls.SignatureAggregate([]*bls.Signature{sig1})

	assert.True(t, agg1.EqualsTo(sig1))
}

func TestAggregateOnlyOnePublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandomBLSKeyPair()
	agg1 := bls.PublicKeyAggregate([]*bls.PublicKey{pub1})

	assert.True(t, agg1.EqualsTo(pub1))
}

// TODO: should we check for duplication here?
func TestDuplicatedAggregate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, prv1 := ts.RandomBLSKeyPair()
	pub2, prv2 := ts.RandomBLSKeyPair()

	msg1 := []byte("zarb")

	sig1 := prv1.Sign(msg1).(*bls.Signature)
	sig2 := prv2.Sign(msg1).(*bls.Signature)

	agg1 := bls.SignatureAggregate([]*bls.Signature{sig1, sig2, sig1})
	agg2 := bls.SignatureAggregate([]*bls.Signature{sig1, sig2})
	assert.False(t, agg1.EqualsTo(agg2))

	pubs1 := []*bls.PublicKey{pub1, pub2}
	pubs2 := []*bls.PublicKey{pub1, pub2, pub1}
	pubAgg1 := bls.PublicKeyAggregate(pubs1)
	pubAgg2 := bls.PublicKeyAggregate(pubs2)
	assert.False(t, pubAgg1.EqualsTo(pubAgg2))

	assert.False(t, bls.VerifyAggregated(agg1, pubs1, msg1))
	assert.NotNil(t, pubAgg1.Verify(msg1, agg1))

	assert.True(t, bls.VerifyAggregated(agg1, pubs2, msg1))
	assert.Nil(t, pubAgg2.Verify(msg1, agg1))
}

// TestHashToCurve ensures that the hash-to-curve function in kilic/bls12-381
// works as intended and is compatible with the spec.
// test vectors can be found here:
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-hash-to-curve-16#appendix-J.9.1
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
	}

	g1 := bls12381.NewG1()
	for no, test := range tests {
		mappedPoint, _ := g1.HashToCurve([]byte(test.msg), domain)
		d, _ := hex.DecodeString(test.expected)
		expectedPoint, _ := g1.FromBytes(d)
		assert.Equal(t, mappedPoint, expectedPoint,
			"test %v: not match", no)
	}
}
