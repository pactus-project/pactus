package validator_test

import (
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
	bs := []byte{
		0x95, 0x16, 0x7c, 0x2a, 0x0d, 0x86, 0xec, 0x36, 0x04, 0x07, 0xbc, 0xe8, // public key
		0x9b, 0x30, 0x46, 0x16, 0xe1, 0xd0, 0xf8, 0x3d, 0xbc, 0x20, 0x06, 0x42,
		0xab, 0xea, 0x84, 0x05, 0xe1, 0x83, 0x83, 0x12, 0xfb, 0x82, 0x90, 0xb1,
		0x23, 0x0e, 0xbe, 0x43, 0x69, 0xcf, 0x1b, 0x7f, 0x55, 0x69, 0x06, 0xc6,
		0x10, 0xae, 0x92, 0xbc, 0xee, 0x54, 0x4a, 0x1a, 0xf7, 0x9e, 0x25, 0x99,
		0x96, 0xe3, 0x68, 0xb1, 0x48, 0x51, 0xa1, 0xf8, 0x84, 0x42, 0x74, 0x69,
		0x0b, 0x10, 0xdf, 0x98, 0x3b, 0xc2, 0x77, 0x6a, 0xb1, 0x0c, 0xc3, 0x7e,
		0x49, 0xe1, 0x75, 0xbc, 0x7a, 0xe1, 0x7a, 0xc9, 0x19, 0xb8, 0xc3, 0x4c,
		0x1, 0x0, 0x0, 0x0, // number
		0x2, 0x0, 0x0, 0x0, // sequence
		0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // balance
		0x4, 0x0, 0x0, 0x0, // LastBondingHeight
		0x5, 0x0, 0x0, 0x0, // UnbondingHeight
		0x6, 0x0, 0x0, 0x0, // LastSortitionHeight
	}

	val, err := validator.FromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, val.Number(), int32(1))
	assert.Equal(t, val.Sequence(), int32(2))
	assert.Equal(t, val.Stake(), int64(3))
	assert.Equal(t, val.LastBondingHeight(), uint32(4))
	assert.Equal(t, val.UnbondingHeight(), uint32(5))
	assert.Equal(t, val.LastSortitionHeight(), uint32(6))
	bs2, _ := val.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("76fea239a4586e8d9c2df9062b1958703341e3ece0f665c714da850101b61185")
	assert.Equal(t, val.Hash(), expected)
	pub, _ := bls.PublicKeyFromBytes(bs[:96])
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
