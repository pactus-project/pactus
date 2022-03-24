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

func TestCertificateMarshaling(t *testing.T) {
	/*
		{
			1: h'3E29108301D725A3DB2E79EAE52BFA148CC52BF3DE32026F034B37C84F0A2BC4',
			2: 4,
			3: [10, 18, 12, 16],
			4: [18],
			5: h'96F6EEF7DA613939CAEBBA051F2C4BC362B35B931378B23DE620F7DCC3DDDD286165463A949EFBF56C143673A3BA3EAB'
		}
	*/
	// d, _ := hex.DecodeString("a50158203e29108301d725a3db2e79eae52bfa148cc52bf3de32026f034b37c84f0a2bc4020403840a120c1004811205583096f6eef7da613939caebba051f2c4bc362b35b931378b23de620f7dcc3dddd286165463a949efbf56c143673a3ba3eab")
	// cert1 := new(Certificate)
	// cert2 := new(Certificate)
	// assert.NoError(t, cert1.Decode(d))
	// d2, err := cert1.Encode()
	// assert.NoError(t, err)
	// assert.Equal(t, d, d2)
	// assert.NoError(t, cert2.Decode(d))

	// expected1 := hash.CalcHash(d)
	// assert.Equal(t, cert1.Hash(), expected1)
	// assert.Equal(t, cert1.Hash(), expected1)
	// assert.Equal(t, cert2.Hash(), expected1)
}

func TestInvalidCertificate(t *testing.T) {
	cert := GenerateTestCertificate(hash.GenerateTestHash())
	cert.data.Committers = nil
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
