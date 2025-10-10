package validator_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	val.UpdateLastBondingHeight(ts.RandHeight())
	val.UpdateLastSortitionHeight(ts.RandHeight())
	val.UpdateUnbondingHeight(ts.RandHeight())
	data, err := val.Bytes()
	require.NoError(t, err)
	require.Equal(t, val.SerializeSize(), len(data))
	val2, err := validator.FromBytes(data)
	require.NoError(t, err)
	assert.Equal(t, val.Address(), val2.Address())
	assert.Equal(t, val.Number(), val2.Number())
	assert.Equal(t, val.Stake(), val2.Stake())
	assert.Equal(t, val.LastBondingHeight(), val2.LastBondingHeight())
	assert.Equal(t, val.LastSortitionHeight(), val2.LastSortitionHeight())
	assert.Equal(t, val.UnbondingHeight(), val2.UnbondingHeight())

	_, err = validator.FromBytes([]byte("asdfghjkl"))
	require.Error(t, err)

	data = data[:len(data)-1]
	_, err = validator.FromBytes(data)
	require.Error(t, err)
}

func TestDecoding(t *testing.T) {
	data, _ := hex.DecodeString(
		"8d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba1591370778846681b7b558e061" + // PublicKey
			"1776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc1" +
			"01000000" + // Number
			"0200000000000000" + // Stake
			"03000000" + // LastBondingHeight
			"04000000" + // UnbondingHeight
			"05000000") // LastSortitionHeight

	val, err := validator.FromBytes(data)
	require.NoError(t, err)
	assert.Equal(t, int32(1), val.Number())
	assert.Equal(t, amount.Amount(2), val.Stake())
	assert.Equal(t, uint32(3), val.LastBondingHeight())
	assert.Equal(t, uint32(4), val.UnbondingHeight())
	assert.Equal(t, uint32(5), val.LastSortitionHeight())
	d2, _ := val.Bytes()
	assert.Equal(t, data, d2)
	assert.Equal(t, hash.CalcHash(data), val.Hash())
	expected, _ := hash.FromString("243e65ae04727f21d5f7618cea9ff8d4bc82fded1179cf8bd9e11a6b99ac42b2")
	assert.Equal(t, expected, val.Hash())
	pub, _ := bls.PublicKeyFromBytes(data[:96])
	assert.True(t, val.PublicKey().EqualsTo(pub))
	assert.Equal(t, len(data), val.SerializeSize())
}

func TestPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	val.SubtractFromStake(val.Stake())
	assert.Equal(t, amount.Amount(0), val.Stake())
	assert.Equal(t, int64(1), val.Power())
	val.AddToStake(1)
	assert.Equal(t, amount.Amount(1), val.Stake())
	assert.Equal(t, int64(1), val.Power())
	val.UpdateUnbondingHeight(1)
	assert.Equal(t, amount.Amount(1), val.Stake())
	assert.Equal(t, int64(0), val.Power())
}

func TestAddToStake(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	stake := val.Stake()
	val.AddToStake(1)
	assert.Equal(t, stake+1, val.Stake())
}

func TestSubtractFromStake(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	stake := val.Stake()
	val.SubtractFromStake(1)
	assert.Equal(t, stake-1, val.Stake())
}

func TestClone(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	cloned := val.Clone()
	cloned.AddToStake(1)

	assert.Equal(t, val.Number(), cloned.Number())
	assert.Equal(t, val.PublicKey(), cloned.PublicKey())
	assert.NotEqual(t, val.Stake(), cloned.Stake())
}

func TestIsUnbonded(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	assert.False(t, val.IsUnbonded())

	val.UpdateUnbondingHeight(ts.RandHeight())
	assert.True(t, val.IsUnbonded())
}
