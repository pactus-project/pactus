package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestCertificateJSONMarshaling(t *testing.T) {
	c1 := GenerateTestCertificate(hash.UndefHash)
	bz, err := c1.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, bz)
}

func TestCertificateSignBytes(t *testing.T) {
	h := hash.GenerateTestHash()
	c1 := GenerateTestCertificate(h)
	bz := CertificateSignBytes(h, c1.Round())
	assert.NotEqual(t, bz, CertificateSignBytes(h, c1.Round()+1))
	assert.NotEqual(t, bz, CertificateSignBytes(hash.GenerateTestHash(), c1.Round()))
}

func TestInvalidCertificate(t *testing.T) {
	cert := GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Committers = nil
	assert.Error(t, cert.SanityCheck())

	cert = GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Round = -1
	assert.Error(t, cert.SanityCheck())

	cert = GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Absentees = nil
	assert.Error(t, cert.SanityCheck())

	cert = GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Absentees = append(cert.data.Absentees, 0)
	assert.Error(t, cert.SanityCheck())

	cert = GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Absentees = []int32{2, 1}
	assert.Error(t, cert.SanityCheck())

	cert = GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Signature = nil
	assert.Error(t, cert.SanityCheck())

	cert = GenerateTestCertificate(hash.GenerateTestHash())
	sig, _ := bls.SignatureFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	cert.data.Signature = sig
	assert.Error(t, cert.SanityCheck())
}

func TestCertificateersHash(t *testing.T) {
	temp := GenerateTestCertificate(hash.GenerateTestHash())

	cert1 := NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{}, temp.Signature())
	assert.Equal(t, cert1.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert1.Absentees(), []int32{})
	assert.NoError(t, cert1.SanityCheck())

	cert2 := NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{2, 6}, temp.Signature())
	assert.Equal(t, cert2.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert2.Absentees(), []int32{2, 6})
	assert.NoError(t, cert2.SanityCheck())

	cert3 := NewCertificate(temp.Round(), []int32{10, 18, 2, 6}, []int32{18}, temp.Signature())
	assert.Equal(t, cert3.Committers(), []int32{10, 18, 2, 6})
	assert.Equal(t, cert3.Absentees(), []int32{18})
	assert.NoError(t, cert3.SanityCheck())
}
