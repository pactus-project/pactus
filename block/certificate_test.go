package block

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestNilCertificateHash(t *testing.T) {
	var cert Certificate
	assert.Equal(t, cert.Hash(), hash.UndefHash)
}

func TestCertificateJSONMarshaling(t *testing.T) {
	c1 := GenerateTestCertificate(hash.UndefHash)
	bz, err := c1.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, bz)
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
	d, _ := hex.DecodeString("a50158203e29108301d725a3db2e79eae52bfa148cc52bf3de32026f034b37c84f0a2bc4020403840a120c1004811205583096f6eef7da613939caebba051f2c4bc362b35b931378b23de620f7dcc3dddd286165463a949efbf56c143673a3ba3eab")
	cert := new(Certificate)
	err := cbor.Unmarshal(d, cert)
	assert.NoError(t, err)
	d2, err := cbor.Marshal(cert)
	assert.NoError(t, err)
	assert.Equal(t, d, d2)

	expected1 := hash.CalcHash(d)
	assert.Equal(t, cert.Hash(), expected1)
}

func TestInvalidCertificate(t *testing.T) {
	cert1 := GenerateTestCertificate(hash.UndefHash)
	assert.Error(t, cert1.SanityCheck())

	c2 := GenerateTestCertificate(hash.GenerateTestHash())
	c2.data.Round = -1
	assert.Error(t, c2.SanityCheck())

	c3 := GenerateTestCertificate(hash.GenerateTestHash())
	c3.data.Committers = nil
	assert.Error(t, c3.SanityCheck())

	c4 := GenerateTestCertificate(hash.GenerateTestHash())
	c4.data.Absentees = nil
	assert.Error(t, c4.SanityCheck())

	c6 := GenerateTestCertificate(hash.GenerateTestHash())
	c6.data.Absentees = append(c6.data.Absentees, -1)
	assert.Error(t, c6.SanityCheck())

	c7 := GenerateTestCertificate(hash.GenerateTestHash())
	c7.data.Absentees = []int{2, 1}
	assert.Error(t, c7.SanityCheck())
}

func TestCertificateersHash(t *testing.T) {
	temp := GenerateTestCertificate(hash.GenerateTestHash())

	cert1 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{}, temp.Signature())
	assert.Equal(t, cert1.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert1.Absentees(), []int{})
	assert.NoError(t, cert1.SanityCheck())

	cert2 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{2, 6}, temp.Signature())
	assert.Equal(t, cert2.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert2.Absentees(), []int{2, 6})
	assert.NoError(t, cert2.SanityCheck())

	cert3 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{18}, temp.Signature())
	assert.Equal(t, cert3.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert3.Absentees(), []int{18})
	assert.NoError(t, cert3.SanityCheck())
}
