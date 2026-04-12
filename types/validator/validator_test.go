package validator_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator()
	assert.False(t, val1.IsDelegated())
	assert.False(t, val1.DelegateExpired(ts.RandHeight()))

	// Round-trip serialization
	data, err := val1.Bytes()
	require.NoError(t, err)
	assert.Len(t, data, 120)

	for i := range len(data) {
		_, err := validator.FromBytes(data[:i])
		require.Error(t, err)
	}

	val2, err := validator.FromBytes(data)
	require.NoError(t, err)
	assert.Equal(t, val1.Hash(), val2.Hash())
	assert.Equal(t, val1.Address(), val2.Address())
	assert.Equal(t, val1.Number(), val2.Number())
}

func TestFromBytesDelegation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator()

	owner := ts.RandAccAddress()
	share := amount.Amount(350_000_000) // 0.35 PAC
	expiry := uint32(1000)
	val1.SetDelegation(owner, share, expiry)

	assert.True(t, val1.IsDelegated())
	assert.Equal(t, owner, val1.DelegateOwner())
	assert.Equal(t, share, val1.DelegateShare())
	assert.Equal(t, expiry, val1.DelegateExpiry())
	assert.False(t, val1.DelegateExpired(999))
	assert.True(t, val1.DelegateExpired(1000))
	assert.True(t, val1.DelegateExpired(1001))

	// Round-trip serialization with delegation
	data, err := val1.Bytes()
	require.NoError(t, err)
	assert.Len(t, data, 120+21+8+4)

	for i := range 32 {
		_, err := validator.FromBytes(data[0 : 121+i])
		require.Error(t, err)
	}

	val2, err := validator.FromBytes(data)
	require.NoError(t, err)
	assert.Equal(t, val1.Hash(), val2.Hash())
}

func TestUpdateValidator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()

	bondingHeight := ts.RandHeight()
	sortitionHeight := ts.RandHeight()
	unbondingHeight := ts.RandHeight()

	val.UpdateLastBondingHeight(bondingHeight)
	val.UpdateLastSortitionHeight(sortitionHeight)
	val.UpdateUnbondingHeight(unbondingHeight)

	assert.Equal(t, bondingHeight, val.LastBondingHeight())
	assert.Equal(t, sortitionHeight, val.LastSortitionHeight())
	assert.Equal(t, unbondingHeight, val.UnbondingHeight())
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
	assert.Equal(t, len(data), val.SerializeSize())

	assert.Equal(t, hash.CalcHash(data), val.Hash())
	expected, _ := hash.FromString("243e65ae04727f21d5f7618cea9ff8d4bc82fded1179cf8bd9e11a6b99ac42b2")
	assert.Equal(t, expected, val.Hash())
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

func TestUpdateProtocolVersion(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	assert.Equal(t, protocol.Version(0), val.ProtocolVersion())

	val.UpdateProtocolVersion(1)
	assert.Equal(t, protocol.Version(1), val.ProtocolVersion())
}
