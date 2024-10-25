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
		"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67")
	pub, _ := bls.PublicKeyFromString(
		"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
			"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx47a")
	sig, _ := bls.SignatureFromString(
		"923d67a8624cbb7972b29328e15ec76cc846076ccf00a9e94d991c677846f334ae4ba4551396fbcd6d1cab7593baf3b7")
	addr, _ := crypto.AddressFromString("pc1p5x2a0lkt5nrrdqe0rkcv6r4pfkmdhrr3xk73tq")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig.Bytes(), sig1.Bytes())
	assert.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, pub, prv.PublicKey())
	assert.Equal(t, addr, pub.ValidatorAddress())
}

func TestSignatureAggregate(t *testing.T) {
	msg := []byte("pactus")
	prv1, _ := bls.PrivateKeyFromString(
		"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67")
	prv2, _ := bls.PrivateKeyFromString(
		"SECRET1PDUV97560CWDGW2DR453YPUT84REN04G0DZFAPJQL5DV0CKDAN75QCJEV6F")
	agg, _ := bls.SignatureFromString(
		"ad747172697127cb08dda29a386e106eb24ab0edfbc044014c3bd7a5f583cc38b3a223ff2c1df9c0b4df110630e6946b")
	sig1 := prv1.Sign(msg).(*bls.Signature)
	sig2 := prv2.Sign(msg).(*bls.Signature)

	assert.True(t, bls.SignatureAggregate(sig1, sig2).EqualsTo(agg))
	assert.False(t, prv1.EqualsTo(prv2))
}

func TestAggregateFailed(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, prv1 := ts.RandBLSKeyPair()
	pub2, prv2 := ts.RandBLSKeyPair()
	pub3, prv3 := ts.RandBLSKeyPair()
	pub4, prv4 := ts.RandBLSKeyPair()
	msg1 := []byte("pactus")
	msg2 := []byte("pactus0")

	sig1 := prv1.Sign(msg1).(*bls.Signature)
	sig11 := prv1.Sign(msg2).(*bls.Signature)
	sig2 := prv2.Sign(msg1).(*bls.Signature)
	sig3 := prv3.Sign(msg1).(*bls.Signature)
	sig4 := prv4.Sign(msg1).(*bls.Signature)

	agg1 := bls.SignatureAggregate(sig1, sig2, sig3)
	agg2 := bls.SignatureAggregate(sig1, sig2, sig4)
	agg3 := bls.SignatureAggregate(sig11, sig2, sig3)
	agg4 := bls.SignatureAggregate(sig1, sig2)
	agg5 := bls.SignatureAggregate(sig3, sig2, sig1)

	pubs1 := []*bls.PublicKey{pub1, pub2, pub3}
	pubs2 := []*bls.PublicKey{pub1, pub2, pub4}
	pubs3 := []*bls.PublicKey{pub1, pub2}
	pubs4 := []*bls.PublicKey{pub3, pub2, pub1}

	pubAgg1 := bls.PublicKeyAggregate(pubs1...)
	pubAgg2 := bls.PublicKeyAggregate(pubs2...)
	pubAgg3 := bls.PublicKeyAggregate(pubs3...)
	pubAgg4 := bls.PublicKeyAggregate(pubs4...)

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
	agg1 := bls.SignatureAggregate(sig1)

	assert.True(t, agg1.EqualsTo(sig1))
}

func TestAggregateOnlyOnePublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	agg1 := bls.PublicKeyAggregate(pub1)

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

	agg1 := bls.SignatureAggregate(sig1, sig2, sig1)
	agg2 := bls.SignatureAggregate(sig1, sig2)
	assert.False(t, agg1.EqualsTo(agg2))

	pubs1 := []*bls.PublicKey{pub1, pub2}
	pubs2 := []*bls.PublicKey{pub1, pub2, pub1}
	pubAgg1 := bls.PublicKeyAggregate(pubs1...)
	pubAgg2 := bls.PublicKeyAggregate(pubs2...)
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
