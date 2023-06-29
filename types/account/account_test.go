package account_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, _ := ts.GenerateTestAccount(ts.RandInt32(10000))
	bs, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, acc.SerializeSize(), len(bs))
	acc2, err := account.FromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, acc, acc2)

	_, err = account.FromBytes([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestDecoding(t *testing.T) {
	bs := []byte{
		0x1, 0x0, 0x0, 0x0, // number
		0x2, 0x0, 0x0, 0x0, // sequence
		0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // balance
	}

	acc, err := account.FromBytes(bs)
	require.NoError(t, err)
	bs2, _ := acc.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, acc.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("74280903e6b73b79e56b1f15cee24c444776cfeee3bea9476b549b660176f773")
	assert.Equal(t, acc.Hash(), expected)
}

func TestIncSequence(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, _ := ts.GenerateTestAccount(100)
	seq := acc.Sequence()
	acc.IncSequence()
	assert.Equal(t, acc.Sequence(), seq+1)
	assert.Equal(t, acc.Number(), int32(100))
}

func TestAddToBalance(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, _ := ts.GenerateTestAccount(100)
	bal := acc.Balance()
	acc.AddToBalance(1)
	assert.Equal(t, acc.Balance(), bal+1)
}

func TestSubtractFromBalance(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, _ := ts.GenerateTestAccount(100)
	bal := acc.Balance()
	acc.SubtractFromBalance(1)
	assert.Equal(t, acc.Balance(), bal-1)
}

func TestClone(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, _ := ts.GenerateTestAccount(100)
	cloned := acc.Clone()
	cloned.IncSequence()

	assert.NotEqual(t, acc.Sequence(), cloned.Sequence())
}
