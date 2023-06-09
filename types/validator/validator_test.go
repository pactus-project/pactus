package validator

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt32(1000000))
	val.UpdateLastBondingHeight(util.RandUint32(1000000))
	val.UpdateLastSortitionHeight(util.RandUint32(1000000))
	val.UpdateUnbondingHeight(util.RandUint32(1000000))
	bs, err := val.Bytes()
	require.NoError(t, err)
	require.Equal(t, val.SerializeSize(), len(bs))
	val2, err := FromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, val.Address(), val2.Address())
	assert.Equal(t, val.Sequence(), val2.Sequence())
	assert.Equal(t, val.Number(), val2.Number())
	assert.Equal(t, val.Stake(), val2.Stake())
	assert.Equal(t, val.LastBondingHeight(), val2.LastBondingHeight())
	assert.Equal(t, val.LastSortitionHeight(), val2.LastSortitionHeight())
	assert.Equal(t, val.UnbondingHeight(), val2.UnbondingHeight())

	_, err = FromBytes([]byte("asdfghjkl"))
	require.Error(t, err)

	bs = bs[:len(bs)-1]
	_, err = FromBytes(bs)
	require.Error(t, err)
}

func TestDecoding(t *testing.T) {
	bs, _ := hex.DecodeString(
		"95167c2a0d86ec360407bce89b304616e1d0f83dbc200642abea8405e1838312fb8290b1230ebe4369cf1b7f556906c610ae92bcee544a1" +
			"af79e259996e368b14851a1f8844274690b10df983bc2776ab10cc37e49e175bc7ae17ac919b8c34c01000000020000000300000000" +
			"000000040000000500000006000000")
	val, err := FromBytes(bs)
	require.NoError(t, err)
	bs2, _ := val.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("76fea239a4586e8d9c2df9062b1958703341e3ece0f665c714da850101b61185")
	assert.Equal(t, val.Hash(), expected)
	pub, _ := bls.PublicKeyFromBytes(bs[:96])
	assert.True(t, val.PublicKey().EqualsTo(pub))
}

func TestIncSequence(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt32(1000))
	seq := val.Sequence()
	val.IncSequence()
	assert.Equal(t, val.Sequence(), seq+1)
}

func TestPower(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt32(1000))
	val.data.Stake = 0
	assert.Equal(t, val.Stake(), int64(0))
	assert.Equal(t, val.Power(), int64(1))
	val.data.Stake = 1
	assert.Equal(t, val.Stake(), int64(1))
	assert.Equal(t, val.Power(), int64(1))
	val.UpdateUnbondingHeight(1)
	assert.Equal(t, val.Stake(), int64(1))
	assert.Equal(t, val.Power(), int64(0))
}
func TestAddToStake(t *testing.T) {
	val, _ := GenerateTestValidator(100)
	stake := val.Stake()
	val.AddToStake(1)
	assert.Equal(t, val.Stake(), stake+1)
}

func TestSubtractFromStake(t *testing.T) {
	val, _ := GenerateTestValidator(100)
	stake := val.Stake()
	val.SubtractFromStake(1)
	assert.Equal(t, val.Stake(), stake-1)
}
