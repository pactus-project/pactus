package certificate_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestCertificate(t *testing.T) {
	d, _ := hex.DecodeString(
		"04030201" + // Height
			"0100" + // Round
			"0401020304" + // Committers
			"0102" + // Absentees
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6") // Signature

	h, _ := hash.FromString("6d5fee07c7cc35384f2f1bc695f6b6afa339df6d867dec0d324a60a48803e1aa")
	r := bytes.NewReader(d)
	cert := new(certificate.Certificate)
	err := cert.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, cert.Height(), uint32(0x01020304))
	assert.Equal(t, cert.Round(), int16(0x0001))
	assert.Equal(t, cert.Committers(), []int32{1, 2, 3, 4})
	assert.Equal(t, cert.Absentees(), []int32{2})
	assert.Equal(t, cert.Hash(), h)
}

func TestCertificateCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestCertificate(ts.RandHeight())
	bz1, err := cbor.Marshal(cert1)
	assert.NoError(t, err)
	var cert2 certificate.Certificate
	err = cbor.Unmarshal(bz1, &cert2)
	assert.NoError(t, err)
	assert.NoError(t, cert2.BasicCheck())
	assert.Equal(t, cert1.Hash(), cert1.Hash())

	assert.Equal(t, cert1.Hash(), cert2.Hash())

	err = cbor.Unmarshal([]byte{1}, &cert2)
	assert.Error(t, err)
}

func TestCertificateSignBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h := ts.RandHash()
	height := ts.RandHeight()
	cert := ts.GenerateTestCertificate(height)
	bz := certificate.BlockCertificateSignBytes(h, height, cert.Round())
	assert.NotEqual(t, bz, certificate.BlockCertificateSignBytes(h, height, cert.Round()+1))
	assert.NotEqual(t, bz, certificate.BlockCertificateSignBytes(ts.RandHash(), height, cert.Round()))
}

func TestInvalidCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert0 := ts.GenerateTestCertificate(ts.RandHeight())

	t.Run("Invalid height", func(t *testing.T) {
		cert := certificate.NewCertificate(0, 0, cert0.Committers(), cert0.Absentees(), cert0.Signature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("Invalid round", func(t *testing.T) {
		cert := certificate.NewCertificate(1, -1, cert0.Committers(), cert0.Absentees(), cert0.Signature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "round is negative: -1",
		})
	})

	t.Run("Committers is nil", func(t *testing.T) {
		cert := certificate.NewCertificate(cert0.Height(), cert0.Round(), nil, cert0.Absentees(), cert0.Signature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "committers is missing",
		})
	})

	t.Run("Absentees is nil", func(t *testing.T) {
		cert := certificate.NewCertificate(cert0.Height(), cert0.Round(), cert0.Committers(), nil, cert0.Signature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "absentees is missing",
		})
	})

	t.Run("Signature is nil", func(t *testing.T) {
		cert := certificate.NewCertificate(cert0.Height(), cert0.Round(), cert0.Committers(), cert0.Absentees(), nil)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "signature is missing",
		})
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		abs := cert0.Absentees()
		abs = append(abs, 0)
		cert := certificate.NewCertificate(cert0.Height(), cert0.Round(), cert0.Committers(), abs, cert0.Signature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), abs),
		})
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		abs := []int32{2, 1}
		cert := certificate.NewCertificate(cert0.Height(), cert0.Round(), cert0.Committers(), abs, cert0.Signature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), abs),
		})
	})
}

func TestCertificateHash(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert0 := ts.GenerateTestCertificate(ts.RandHeight())

	cert1 := certificate.NewCertificate(cert0.Height(), cert0.Round(),
		[]int32{10, 18, 2, 6}, []int32{}, cert0.Signature())
	assert.Equal(t, cert1.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert1.Absentees(), []int32{})
	assert.NoError(t, cert1.BasicCheck())

	cert2 := certificate.NewCertificate(cert0.Height(), cert0.Round(),
		[]int32{10, 18, 2, 6}, []int32{2, 6}, cert0.Signature())
	assert.Equal(t, cert2.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert2.Absentees(), []int32{2, 6})
	assert.NoError(t, cert2.BasicCheck())

	cert3 := certificate.NewCertificate(cert0.Height(), cert0.Round(),
		[]int32{10, 18, 2, 6}, []int32{18}, cert0.Signature())
	assert.Equal(t, cert3.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert3.Absentees(), []int32{18})
	assert.NoError(t, cert3.BasicCheck())
}

func TestEncodingCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestCertificate(ts.RandHeight())
	length := cert1.SerializeSize()

	for i := 0; i < length; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, cert1.Encode(w), "encode test %v failed", i)
	}
	w := util.NewFixedWriter(length)
	assert.NoError(t, cert1.Encode(w))

	for i := 0; i < length; i++ {
		cert := new(certificate.Certificate)
		r := util.NewFixedReader(i, w.Bytes())
		assert.Error(t, cert.Decode(r), "decode test %v failed", i)
	}

	cert2 := new(certificate.Certificate)
	r := util.NewFixedReader(length, w.Bytes())
	assert.NoError(t, cert2.Decode(r))
	assert.Equal(t, cert1.Hash(), cert2.Hash())
}

func TestCertificateValidation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valKey1 := ts.RandValKey()
	valKey2 := ts.RandValKey()
	valKey3 := ts.RandValKey()
	valKey4 := ts.RandValKey()
	val1 := validator.NewValidator(valKey1.PublicKey(), 1001)
	val2 := validator.NewValidator(valKey2.PublicKey(), 1002)
	val3 := validator.NewValidator(valKey3.PublicKey(), 1003)
	val4 := validator.NewValidator(valKey4.PublicKey(), 1004)

	validators := []*validator.Validator{val1, val2, val3, val4}
	committers := []int32{
		val1.Number(), val2.Number(), val3.Number(), val4.Number(),
	}
	blockHash := ts.RandHash()
	blockHeight := ts.RandHeight()
	blockRound := ts.RandRound()
	signBytes := certificate.BlockCertificateSignBytes(blockHash, blockHeight, blockRound)
	sig1 := valKey1.Sign(signBytes)
	sig3 := valKey3.Sign(signBytes)
	sig4 := valKey4.Sign(signBytes)
	aggSig := bls.SignatureAggregate(sig1, sig3, sig4)

	t.Run("Invalid height, should return error", func(t *testing.T) {
		cert := certificate.NewCertificate(blockHeight+1, blockRound, committers,
			[]int32{}, aggSig)

		err := cert.Validate(blockHeight, validators, signBytes)
		assert.ErrorIs(t, err, certificate.UnexpectedHeightError{
			Expected: blockHeight,
			Got:      blockHeight + 1,
		})
	})

	t.Run("Invalid committer, should return error", func(t *testing.T) {
		invCommitters := committers
		invCommitters = append(invCommitters, ts.Rand.Int31n(1000))
		cert := certificate.NewCertificate(blockHeight, blockRound, invCommitters,
			[]int32{}, aggSig)

		err := cert.Validate(blockHeight, validators, signBytes)

		assert.ErrorIs(t, err, certificate.UnexpectedCommittersError{
			Committers: invCommitters,
		})
	})

	t.Run("Invalid committers, should return error", func(t *testing.T) {
		invCommitters := []int32{
			ts.Rand.Int31n(1000), val2.Number(), val3.Number(), val4.Number(),
		}
		cert := certificate.NewCertificate(blockHeight, blockRound, invCommitters,
			[]int32{}, aggSig)

		err := cert.Validate(blockHeight, validators, signBytes)
		assert.ErrorIs(t, err, certificate.UnexpectedCommittersError{
			Committers: invCommitters,
		})
	})

	t.Run("Doesn't have 2/3 majority", func(t *testing.T) {
		cert := certificate.NewCertificate(blockHeight, blockRound, committers,
			[]int32{val1.Number(), val2.Number()}, aggSig)

		err := cert.Validate(blockHeight, validators, signBytes)
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		cert := certificate.NewCertificate(blockHeight, blockRound, committers,
			[]int32{val3.Number()}, aggSig)

		err := cert.Validate(blockHeight, validators, signBytes)
		assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})
	t.Run("Ok, should return no error", func(t *testing.T) {
		cert := certificate.NewCertificate(blockHeight, blockRound, committers,
			[]int32{val2.Number()}, aggSig)

		err := cert.Validate(blockHeight, validators, signBytes)
		assert.NoError(t, err)
	})
}

func TestAddSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valKey1 := ts.RandValKey()
	valKey2 := ts.RandValKey()
	valKey3 := ts.RandValKey()
	valKey4 := ts.RandValKey()
	blockHeight := ts.RandHeight()
	blockRound := ts.RandRound()
	blockHash := ts.RandHash()

	val1 := validator.NewValidator(valKey1.PublicKey(), ts.RandInt32(10000))
	val2 := validator.NewValidator(valKey2.PublicKey(), ts.RandInt32(10000))
	val3 := validator.NewValidator(valKey3.PublicKey(), ts.RandInt32(10000))
	val4 := validator.NewValidator(valKey4.PublicKey(), ts.RandInt32(10000))

	signBytes := certificate.BlockCertificateSignBytes(blockHash, blockHeight, blockRound)
	sig1 := valKey1.Sign(signBytes)
	sig2 := valKey2.Sign(signBytes)
	sig3 := valKey3.Sign(signBytes)
	sig4 := valKey4.Sign(signBytes)
	aggSig := bls.SignatureAggregate(sig1, sig2, sig3)

	cert := certificate.NewCertificate(blockHeight, blockRound,
		[]int32{val1.Number(), val2.Number(), val3.Number(), val4.Number()}, []int32{val4.Number()}, aggSig)

	assert.Equal(t, []int32{val4.Number()}, cert.Absentees())
	cert.AddSignature(val4.Number(), sig4)
	assert.Empty(t, cert.Absentees())
	assert.NoError(t, cert.Validate(blockHeight, []*validator.Validator{val1, val2, val3, val4}, signBytes))
}

func TestClone(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestCertificate(ts.RandHeight())
	cert2 := cert1.Clone()

	cert2.AddSignature(cert2.Absentees()[0], ts.RandBLSSignature())
	assert.NotEqual(t, cert1, cert2)
}
