package block

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestNilCertificateHash(t *testing.T) {
	var cert Certificate
	assert.Equal(t, cert.Hash(), crypto.UndefHash)
}

func TestCertificateJSONMarshaling(t *testing.T) {
	c1 := GenerateTestCertificate(crypto.UndefHash)
	bz, err := c1.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, bz)
}

func TestCertificateMarshaling(t *testing.T) {
	/*
		{
			1: h'D846EF49A6C72390645F12970987865A795A55FA19C92DBB9CBE24D6503ECA9F',
			2: 6,
			3: [10, 18, 2, 6],
			4: [10],
			5: h'85C368E9E6DF4EA1B16E29AEBBF74A3DA45A033683E753C93130336E035C2181BF469DAB5E0448064FB64F6282B28296'
		}
	*/
	d, _ := hex.DecodeString("a5015820d846ef49a6c72390645f12970987865a795a55fa19c92dbb9cbe24d6503eca9f020603840a12020604810a05583085c368e9e6df4ea1b16e29aebbf74a3da45a033683e753c93130336e035c2181bf469dab5e0448064fb64f6282b28296")
	cert := new(Certificate)
	err := cbor.Unmarshal(d, cert)
	assert.NoError(t, err)
	d2, err := cbor.Marshal(cert)
	assert.NoError(t, err)
	assert.Equal(t, d, d2)

	expected1 := crypto.HashH(d)
	assert.Equal(t, cert.Hash(), expected1)
}

func TestInvalidCertificate(t *testing.T) {
	cert1 := GenerateTestCertificate(crypto.UndefHash)
	assert.Error(t, cert1.SanityCheck())

	c2 := GenerateTestCertificate(crypto.GenerateTestHash())
	c2.data.Round = -1
	assert.Error(t, c2.SanityCheck())

	c3 := GenerateTestCertificate(crypto.GenerateTestHash())
	c3.data.Committers = nil
	assert.Error(t, c3.SanityCheck())

	c4 := GenerateTestCertificate(crypto.GenerateTestHash())
	c4.data.Absences = nil
	assert.Error(t, c4.SanityCheck())

	c6 := GenerateTestCertificate(crypto.GenerateTestHash())
	c6.data.Absences = append(c6.data.Absences, -1)
	assert.Error(t, c6.SanityCheck())

	c7 := GenerateTestCertificate(crypto.GenerateTestHash())
	c7.data.Absences = []int{2, 1}
	assert.Error(t, c7.SanityCheck())
}

func TestCertificateersHash(t *testing.T) {
	temp := GenerateTestCertificate(crypto.GenerateTestHash())

	cert1 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{}, temp.Signature())
	assert.Equal(t, cert1.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert1.Absences(), []int{})
	assert.NoError(t, cert1.SanityCheck())

	cert2 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{2, 6}, temp.Signature())
	assert.Equal(t, cert2.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert2.Absences(), []int{2, 6})
	assert.NoError(t, cert2.SanityCheck())

	cert3 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{18}, temp.Signature())
	assert.Equal(t, cert3.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert3.Absences(), []int{18})
	assert.NoError(t, cert3.SanityCheck())
}
