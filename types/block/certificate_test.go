package block_test

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestCertificateCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	c1 := ts.GenerateTestCertificate(ts.RandomHash())
	bz1, err := cbor.Marshal(c1)
	assert.NoError(t, err)
	var c2 block.Certificate
	err = cbor.Unmarshal(bz1, &c2)
	assert.NoError(t, err)
	assert.NoError(t, c2.BasicCheck())
	assert.Equal(t, c1.Hash(), c1.Hash())

	assert.Equal(t, c1.Hash(), c2.Hash())

	err = cbor.Unmarshal([]byte{1}, &c2)
	assert.Error(t, err)
}

func TestCertificateSignBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h := ts.RandomHash()
	c1 := ts.GenerateTestCertificate(h)
	bz := block.CertificateSignBytes(h, c1.Round())
	assert.NotEqual(t, bz, block.CertificateSignBytes(h, c1.Round()+1))
	assert.NotEqual(t, bz, block.CertificateSignBytes(ts.RandomHash(), c1.Round()))
}

func TestInvalidCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert0 := ts.GenerateTestCertificate(ts.RandomHash())

	t.Run("Invalid round", func(t *testing.T) {
		cert := block.NewCertificate(-1, cert0.Committers(), cert0.Absentees(), cert0.Signature())

		err := cert.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidRound)
	})

	t.Run("Committers is nil", func(t *testing.T) {
		cert := block.NewCertificate(cert0.Round(), nil, cert0.Absentees(), cert0.Signature())

		err := cert.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidBlock)
	})

	t.Run("Absentees is nil", func(t *testing.T) {
		cert := block.NewCertificate(cert0.Round(), cert0.Committers(), nil, cert0.Signature())

		err := cert.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidBlock)
	})

	t.Run("Signature is nil", func(t *testing.T) {
		cert := block.NewCertificate(cert0.Round(), cert0.Committers(), cert0.Absentees(), nil)

		err := cert.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		abs := cert0.Absentees()
		abs = append(abs, 0)
		cert := block.NewCertificate(cert0.Round(), cert0.Committers(), abs, cert0.Signature())

		err := cert.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidBlock)
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		abs := []int32{2, 1}
		cert := block.NewCertificate(cert0.Round(), cert0.Committers(), abs, cert0.Signature())

		err := cert.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidBlock)
	})
}

func TestCertificateHash(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	temp := ts.GenerateTestCertificate(ts.RandomHash())

	cert1 := block.NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{}, temp.Signature())
	assert.Equal(t, cert1.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert1.Absentees(), []int32{})
	assert.NoError(t, cert1.BasicCheck())

	cert2 := block.NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{2, 6}, temp.Signature())
	assert.Equal(t, cert2.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert2.Absentees(), []int32{2, 6})
	assert.NoError(t, cert2.BasicCheck())

	cert3 := block.NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{18}, temp.Signature())
	assert.Equal(t, cert3.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert3.Absentees(), []int32{18})
	assert.NoError(t, cert3.BasicCheck())
}
