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
	d, _ := hex.DecodeString("a50158201c8f67440c5d2fcaec3176cde966e8b46ec744c836f643612bec96eb6a83c1fe02060383010203048100055830df65ca781a94080e2fcfb66bd443a9681f6b93985c9e696c248b778355231c44411cc2400a02a763827bc9251c553b03")
	c := new(Commit)
	err := cbor.Unmarshal(d, c)
	assert.NoError(t, err)
	d2, err := cbor.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, d, d2)
	expected1, _ := crypto.HashFromString("df5c58d8b7c13806b6d23e878526ccdf331c4fed72780e52ea2775f4aa082a44")
	assert.Equal(t, c.CommittersHash(), expected1)
	expected2, _ := crypto.HashFromString("ab25b4097fa3d9bdd70bf8910064a14a497bae8e1e715621f6e6818506c3d047")
	assert.Equal(t, c.Hash(), expected2)
	expected3, _ := hex.DecodeString("a20158201c8f67440c5d2fcaec3176cde966e8b46ec744c836f643612bec96eb6a83c1fe0206")
	assert.Equal(t, c.SignBytes(), expected3)
}

func TestCommitSanityCheck(t *testing.T) {
	c1 := GenerateTestCommit(crypto.GenerateTestHash())
	assert.NoError(t, c1.SanityCheck())
	c1.data.Missed = append(c1.data.Missed, 5)
	assert.Error(t, c1.SanityCheck())

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
	c.data.Missed = append(c.data.Missed, 5)
	assert.Equal(t, c.Threshold(), 60) // 3÷5=0.6
	assert.False(t, c.HasTwoThirdThreshold())
	c.data.Signed = append(c.data.Signed, 6)
	assert.False(t, c.HasTwoThirdThreshold())
	c.data.Signed = append(c.data.Signed, 7)
	assert.Equal(t, c.Threshold(), 71) //6÷8=0.71
	assert.True(t, c.HasTwoThirdThreshold())
}

func TestCommiters(t *testing.T) {
	temp := GenerateTestCommit(crypto.GenerateTestHash())
	expected1 := temp.Committers()
	expected2 := temp.CommittersHash()
	c1 := NewCommit(temp.BlockHash(), temp.Round(), []int{0, 1, 2, 3}, []int{}, temp.Signature())
	assert.Equal(t, c1.Committers(), expected1)
	assert.Equal(t, c1.CommittersHash(), expected2)

	c2 := NewCommit(temp.BlockHash(), temp.Round(), []int{2, 3}, []int{0, 1}, temp.Signature())
	assert.Equal(t, c2.Committers(), expected1)
	assert.Equal(t, c2.CommittersHash(), expected2)
}
