package block_test

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/types/block"
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
	assert.NoError(t, c2.SanityCheck())
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

// func TestInvalidCertificate(t *testing.T) {
// 	ts := testsuite.NewTestSuite(t)

// 	cert := ts.GenerateTestCertificate(ts.RandomHash())
// 	cert.data.Committers = nil
// 	assert.Error(t, cert.SanityCheck())

// 	cert = ts.GenerateTestCertificate(ts.RandomHash())
// 	cert.data.Round = -1
// 	assert.Error(t, cert.SanityCheck())

// 	cert = ts.GenerateTestCertificate(ts.RandomHash())
// 	cert.data.Absentees = nil
// 	assert.Error(t, cert.SanityCheck())

// 	cert = ts.GenerateTestCertificate(ts.RandomHash())
// 	cert.data.Absentees = append(cert.data.Absentees, 0)
// 	assert.Error(t, cert.SanityCheck())

// 	cert = ts.GenerateTestCertificate(ts.RandomHash())
// 	cert.data.Absentees = []int32{2, 1}
// 	assert.Error(t, cert.SanityCheck())

// 	cert = ts.GenerateTestCertificate(ts.RandomHash())
// 	cert.data.Signature = nil
// 	assert.Error(t, cert.SanityCheck())
// }

func TestCertificateHash(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	temp := ts.GenerateTestCertificate(ts.RandomHash())

	cert1 := block.NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{}, temp.Signature())
	assert.Equal(t, cert1.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert1.Absentees(), []int32{})
	assert.NoError(t, cert1.SanityCheck())

	cert2 := block.NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{2, 6}, temp.Signature())
	assert.Equal(t, cert2.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert2.Absentees(), []int32{2, 6})
	assert.NoError(t, cert2.SanityCheck())

	cert3 := block.NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{18}, temp.Signature())
	assert.Equal(t, cert3.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert3.Absentees(), []int32{18})
	assert.NoError(t, cert3.SanityCheck())
}

// This test ensures that committers are not part of the certificate hash
// We can remove this tests if we remove the committers from the certificate
// This test is not logical, since we have two certificate for the same block
func TestCertificateHashWithoutCommitters(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	temp := ts.GenerateTestCertificate(ts.RandomHash())
	cert1 := block.NewCertificate(temp.Round(), []int32{1, 2, 3, 4}, []int32{2}, temp.Signature())
	cert2 := block.NewCertificate(temp.Round(), []int32{1, 2, 3, 4, 5}, []int32{2}, temp.Signature())

	assert.Equal(t, cert1.Hash(), cert2.Hash())
}
