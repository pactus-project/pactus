package validator_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val, _ := ts.GenerateTestValidator(ts.RandInt32(1000000))
	val.UpdateLastBondingHeight(ts.RandHeight())
	val.UpdateLastSortitionHeight(ts.RandHeight())
	val.UpdateUnbondingHeight(ts.RandHeight())
	bs, err := val.Bytes()
	require.NoError(t, err)
	require.Equal(t, val.SerializeSize(), len(bs))
	val2, err := validator.FromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, val.Address(), val2.Address())
	assert.Equal(t, val.Sequence(), val2.Sequence())
	assert.Equal(t, val.Number(), val2.Number())
	assert.Equal(t, val.Stake(), val2.Stake())
	assert.Equal(t, val.LastBondingHeight(), val2.LastBondingHeight())
	assert.Equal(t, val.LastSortitionHeight(), val2.LastSortitionHeight())
	assert.Equal(t, val.UnbondingHeight(), val2.UnbondingHeight())

	_, err = validator.FromBytes([]byte("asdfghjkl"))
	require.Error(t, err)

	bs = bs[:len(bs)-1]
	_, err = validator.FromBytes(bs)
	require.Error(t, err)
}

func TestDecoding(t *testing.T) {
	d, _ := hex.DecodeString(
		"8d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba1591370778846681b7b558e061" + // PublicKey
			"1776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc1" +
			"01000000" + // Number
			"02000000" + // Sequence
			"0300000000000000" + // Stake
			"04000000" + // LastBondingHeight
			"05000000" + // UnbondingHeight
			"06000000") // LastSortitionHeight

	val, err := validator.FromBytes(d)
	require.NoError(t, err)
	assert.Equal(t, val.Number(), int32(1))
	assert.Equal(t, val.Sequence(), int32(2))
	assert.Equal(t, val.Stake(), int64(3))
	assert.Equal(t, val.LastBondingHeight(), uint32(4))
	assert.Equal(t, val.UnbondingHeight(), uint32(5))
	assert.Equal(t, val.LastSortitionHeight(), uint32(6))
	d2, _ := val.Bytes()
	assert.Equal(t, d, d2)
	assert.Equal(t, val.Hash(), hash.CalcHash(d))
	expected, _ := hash.FromString("6a42857652fbd8a7fbbe63b3d4b177cd4c93de2c29d633e33aa9e83d69831bda")
	assert.Equal(t, val.Hash(), expected)
	pub, _ := bls.PublicKeyFromBytes(d[:96])
	assert.True(t, val.PublicKey().EqualsTo(pub))
}

func TestIncSequence(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val, _ := ts.GenerateTestValidator(ts.RandInt32(1000))
	seq := val.Sequence()
	val.IncSequence()
	assert.Equal(t, val.Sequence(), seq+1)
}

func TestPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val, _ := ts.GenerateTestValidator(ts.RandInt32(1000))
	val.SubtractFromStake(val.Stake())
	assert.Equal(t, val.Stake(), int64(0))
	assert.Equal(t, val.Power(), int64(1))
	val.AddToStake(1)
	assert.Equal(t, val.Stake(), int64(1))
	assert.Equal(t, val.Power(), int64(1))
	val.UpdateUnbondingHeight(1)
	assert.Equal(t, val.Stake(), int64(1))
	assert.Equal(t, val.Power(), int64(0))
}

func TestAddToStake(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val, _ := ts.GenerateTestValidator(100)
	stake := val.Stake()
	val.AddToStake(1)
	assert.Equal(t, val.Stake(), stake+1)
}

func TestSubtractFromStake(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val, _ := ts.GenerateTestValidator(100)
	stake := val.Stake()
	val.SubtractFromStake(1)
	assert.Equal(t, val.Stake(), stake-1)
}

func TestClone(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val, _ := ts.GenerateTestValidator(100)
	cloned := val.Clone()
	cloned.IncSequence()

	assert.NotEqual(t, val.Sequence(), cloned.Sequence())
}
