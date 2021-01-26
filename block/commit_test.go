package block

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestNilCommitHash(t *testing.T) {
	var c Commit
	assert.Equal(t, c.Hash(), crypto.UndefHash)
}

func TestCommitMarshaling(t *testing.T) {
	d, _ := hex.DecodeString("a40158207e7f2aceae7e3bfa6a023cd3221a093d27fbd6e05a7686538b83d9d6308b1f2502060384a201000200a201010201a201020201a201030201045830e9d9e964b40713024ac0b28798795db48cb07122513c782fe7aae9f7ed5c984b944bfe90cf2dd03ab9929170dcbf2090")
	c := new(Commit)
	err := cbor.Unmarshal(d, c)
	assert.NoError(t, err)
	d2, err := cbor.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, d, d2)
	expected1, _ := crypto.HashFromString("fd36b2597b028652ad4430b34a67094ba93ed84bd3abe5cd27f675bf431add48")
	assert.Equal(t, c.CommitteeHash(), expected1)
	assert.Equal(t, c.CommitteeHash(), crypto.HashH([]byte{0x84, 0x00, 0x01, 0x02, 03}))
	expected2, _ := crypto.HashFromString("9e954e738f696a49ae6aac4fb837ec1fff2757b36d4ec0647aacb90cca180bd1")
	assert.Equal(t, c.Hash(), expected2)
	expected3, _ := hex.DecodeString("a20158207e7f2aceae7e3bfa6a023cd3221a093d27fbd6e05a7686538b83d9d6308b1f250206")
	assert.Equal(t, c.SignBytes(), expected3)
}

func TestCommitSanityCheck(t *testing.T) {
	c := GenerateTestCommit(crypto.GenerateTestHash())
	assert.NoError(t, c.SanityCheck())
	c.data.Committers[1].Status = 0 // not signed
	assert.False(t, c.Committers()[0].HasSigned())
	// Not enough signer
	assert.Error(t, c.SanityCheck())
	assert.False(t, c.HasTwoThirdThreshold())
	c.data.Committers[0].Status = 1 // signed
	assert.NoError(t, c.SanityCheck())
	assert.True(t, c.HasTwoThirdThreshold())
	c.data.Committers[2].Status = 2 // invalid status
	assert.Error(t, c.SanityCheck())
	assert.False(t, c.HasTwoThirdThreshold())

	c2 := GenerateTestCommit(crypto.UndefHash)
	assert.Error(t, c2.SanityCheck())

	c3 := GenerateTestCommit(crypto.GenerateTestHash())
	c3.data.Round = -1
	assert.Error(t, c3.SanityCheck())
}

func TestThreshold(t *testing.T) {
	c := GenerateTestCommit(crypto.GenerateTestHash())

	assert.Equal(t, c.Threshold(), 75) // 3÷4=0.75
	assert.True(t, c.HasTwoThirdThreshold())
	c.data.Committers = append(c.data.Committers, Committer{5, CommitNotSigned})
	assert.Equal(t, c.Threshold(), 60) // 3÷5=0.6
	assert.False(t, c.HasTwoThirdThreshold())
	c.data.Committers = append(c.data.Committers, Committer{6, CommitSigned})
	assert.False(t, c.HasTwoThirdThreshold())
	c.data.Committers = append(c.data.Committers, Committer{7, CommitSigned})
	assert.Equal(t, c.Threshold(), 71) //6÷8=0.71
	assert.True(t, c.HasTwoThirdThreshold())
}

func TestCommitersHash(t *testing.T) {
	temp := GenerateTestCommit(crypto.GenerateTestHash())
	expected2 := temp.CommitteeHash()
	c1 := NewCommit(temp.BlockHash(), temp.Round(), []Committer{
		{0, CommitSigned},
		{1, CommitSigned},
		{2, CommitSigned},
		{3, CommitSigned},
	}, temp.Signature())
	assert.Equal(t, c1.CommitteeHash(), expected2)

	c2 := NewCommit(temp.BlockHash(), temp.Round(), []Committer{
		{0, CommitSigned},
		{1, CommitSigned},
		{2, CommitNotSigned},
		{3, CommitNotSigned},
	}, temp.Signature())
	assert.Equal(t, c2.CommitteeHash(), expected2)

	c3 := NewCommit(temp.BlockHash(), temp.Round(), []Committer{
		{1, CommitSigned},
		{2, CommitSigned},
		{3, CommitSigned},
		{0, CommitNotSigned},
	}, temp.Signature())
	assert.NotEqual(t, c3.CommitteeHash(), expected2)
}
