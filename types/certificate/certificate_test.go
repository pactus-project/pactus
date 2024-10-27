package certificate_test

import (
	"fmt"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestCertificateCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestBlockCertificate(ts.RandHeight())
	bz1, err := cbor.Marshal(cert1)
	assert.NoError(t, err)
	var cert2 certificate.BlockCertificate
	err = cbor.Unmarshal(bz1, &cert2)
	assert.NoError(t, err)
	assert.NoError(t, cert2.BasicCheck())
	assert.Equal(t, cert1.Hash(), cert1.Hash())
	assert.True(t, cert1.Signature().EqualsTo(cert2.Signature()))

	err = cbor.Unmarshal([]byte{1}, &cert2)
	assert.Error(t, err)
}

func TestInvalidCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid height", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(0, 0)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("Invalid round", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(1, -1)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "round is negative: -1",
		})
	})

	t.Run("Committers is nil", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature(nil, []int32{1}, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "committers is missing",
		})
	})

	t.Run("Absentees is nil", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{1, 2, 3, 4}, nil, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "absentees is missing",
		})
	})

	t.Run("Signature is nil", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{1, 2, 3, 4, 5, 6}, []int32{1}, nil)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "signature is missing",
		})
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{11, 2, 3, 4, 5, 6}, []int32{66}, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), []int32{66}),
		})
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		cert := certificate.NewBlockCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{1, 2, 3, 4, 5, 6}, []int32{2, 1}, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), []int32{2, 1}),
		})
	})
}

func TestEncodingCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestBlockCertificate(ts.RandHeight())
	length := cert1.SerializeSize()

	for i := 0; i < length; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, cert1.Encode(w), "encode test %v failed", i)
	}
	writer := util.NewFixedWriter(length)
	assert.NoError(t, cert1.Encode(writer))

	for i := 0; i < length; i++ {
		cert := new(certificate.BlockCertificate)
		r := util.NewFixedReader(i, writer.Bytes())
		assert.Error(t, cert.Decode(r), "decode test %v failed", i)
	}

	cert2 := new(certificate.BlockCertificate)
	reader := util.NewFixedReader(length, writer.Bytes())
	assert.NoError(t, cert2.Decode(reader))
	assert.Equal(t, cert1.Hash(), cert2.Hash())
}

func TestAddSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewBlockCertificate(height, round)
	signBytes := cert.SignBytes(blockHash)
	committers := ts.RandSlice(6)
	sigs := []*bls.Signature{}
	validators := []*validator.Validator{}

	for _, committer := range committers {
		valKey := ts.RandValKey()
		val := validator.NewValidator(valKey.PublicKey(), committer)
		sig := valKey.Sign(signBytes)

		validators = append(validators, val)
		sigs = append(sigs, sig)
	}

	absentees := committers[4:]
	aggSig := bls.SignatureAggregate(sigs[:4]...)
	cert.SetSignature(committers, absentees, aggSig)

	err := cert.Validate(validators, blockHash)
	assert.NoError(t, err)

	numAbsentees := len(cert.Absentees())

	t.Run("Add an existing signature", func(t *testing.T) {
		cert.AddSignature(validators[0].Number(), sigs[0])

		err := cert.Validate(validators, blockHash)
		assert.NoError(t, err)

		assert.Len(t, cert.Absentees(), numAbsentees)
	})

	t.Run("Add non existing signature", func(t *testing.T) {
		cert.AddSignature(validators[5].Number(), sigs[5])

		err := cert.Validate(validators, blockHash)
		assert.NoError(t, err)

		assert.Len(t, cert.Absentees(), numAbsentees-1)
	})
}

func TestClone(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestBlockCertificate(ts.RandHeight())
	cert2 := cert1.Clone()

	cert2.AddSignature(cert2.Absentees()[0], ts.RandBLSSignature())
	assert.NotEqual(t, cert1.Absentees(), cert2.Absentees())
	assert.NotEqual(t, cert1, cert2)
}
