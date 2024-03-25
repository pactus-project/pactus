package account_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
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
	d, _ := hex.DecodeString(
		"01000000" + // number
			"0200000000000000") // balance

	acc, err := account.FromBytes(d)
	require.NoError(t, err)
	assert.Equal(t, acc.Number(), int32(1))
	assert.Equal(t, acc.Balance(), amount.Amount(2))
	d2, _ := acc.Bytes()
	assert.Equal(t, d, d2)
	assert.Equal(t, acc.Hash(), hash.CalcHash(d))
	expected, _ := hash.FromString("c3b75f08e64a66cb980fdc03c3a0b78635a7b1db049096e8bbbd9a2873f3071a")
	assert.Equal(t, acc.Hash(), expected)
	assert.Equal(t, acc.SerializeSize(), len(d))
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
	cloned.AddToBalance(1)

	assert.NotEqual(t, acc.Balance(), cloned.Balance())
}
