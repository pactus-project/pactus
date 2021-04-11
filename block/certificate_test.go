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

func TestCertificateMarshaling(t *testing.T) {
	d, _ := hex.DecodeString("a5015820e809498c15e07dd5a99edfe01639193531137d603ac6f43bc59d9c4f5d42d047020603840001020304810008583057bf20c441010af1bc25fc49ce5d12066042d83f2202d650fcc57df40b9591de57670fef2d5237907fd4b5a9b26cb788")
	cert := new(Certificate)
	err := cbor.Unmarshal(d, cert)
	assert.NoError(t, err)
	d2, err := cbor.Marshal(cert)
	assert.NoError(t, err)
	assert.Equal(t, d, d2)
	expected1, _ := crypto.HashFromString("fd36b2597b028652ad4430b34a67094ba93ed84bd3abe5cd27f675bf431add48")
	assert.Equal(t, cert.CommitteeHash(), expected1)
	assert.Equal(t, cert.CommitteeHash(), crypto.HashH([]byte{0x84, 0x00, 0x01, 0x02, 03}))
	expected2, _ := crypto.HashFromString("3c14a9f114da38708cd4385865eafd1d90a5a0214f92cb8378ae2ed186e68ffe")
	assert.Equal(t, cert.Hash(), expected2)
	expected3, _ := hex.DecodeString("a2015820e809498c15e07dd5a99edfe01639193531137d603ac6f43bc59d9c4f5d42d0470206")
	assert.Equal(t, cert.SignBytes(), expected3)
	assert.NoError(t, cert.SanityCheck())
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
	expected2 := temp.CommitteeHash()

	cert1 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{}, temp.Signature())
	assert.Equal(t, cert1.CommitteeHash(), expected2)
	assert.Equal(t, cert1.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert1.Absences(), []int{})
	assert.NoError(t, cert1.SanityCheck())

	cert2 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{2, 6}, temp.Signature())
	assert.Equal(t, cert2.CommitteeHash(), cert1.CommitteeHash())
	assert.Equal(t, cert2.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert2.Absences(), []int{2, 6})
	assert.NoError(t, cert2.SanityCheck())

	cert3 := NewCertificate(temp.BlockHash(), temp.Round(), []int{10, 18, 2, 6}, []int{18}, temp.Signature())
	assert.Equal(t, cert3.CommitteeHash(), cert1.CommitteeHash())
	assert.Equal(t, cert3.Committers(), []int{10, 18, 2, 6})
	assert.Equal(t, cert3.Absences(), []int{18})
	assert.NoError(t, cert3.SanityCheck())
}
