package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestFromBytes(t *testing.T) {
	acc, _ := GenerateTestAccount(util.RandInt32(10000))
	bs, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, acc.SerializeSize(), len(bs))
	acc2, err := FromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, acc, acc2)

	_, err = FromBytes([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestDecoding(t *testing.T) {
	bs := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
		0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, // address
		0x1, 0x0, 0x0, 0x0, // number
		0x2, 0x0, 0x0, 0x0, // sequence
		0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // balance
	}

	acc, err := FromBytes(bs)
	require.NoError(t, err)
	bs2, _ := acc.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, acc.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("33a4208262903cd1f274e760f495eca8e56b7fcc61feec0a8e6dcd0d2e57cafc")
	assert.Equal(t, acc.Hash(), expected)
	assert.Equal(t, acc.Address().Bytes(), bs[:21])
}

func TestIncSequence(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	seq := acc.Sequence()
	acc.IncSequence()
	assert.Equal(t, acc.Sequence(), seq+1)
	assert.Equal(t, acc.Number(), int32(100))
}

func TestAddToBalance(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	bal := acc.Balance()
	acc.AddToBalance(1)
	assert.Equal(t, acc.Balance(), bal+1)
}

func TestSubtractFromBalance(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	bal := acc.Balance()
	acc.SubtractFromBalance(1)
	assert.Equal(t, acc.Balance(), bal-1)
}
