package validator

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestFromBytes(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt32(1000000))
	val.UpdateLastBondingHeight(util.RandInt32(1000000))
	val.UpdateLastJoinedHeight(util.RandInt32(1000000))
	val.UpdateUnbondingHeight(util.RandInt32(1000000))
	bs, err := val.Bytes()
	require.NoError(t, err)
	require.Equal(t, val.SerializeSize(), len(bs))
	val2, err := ValidatorFromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, val.Address(), val2.Address())
	assert.Equal(t, val.Sequence(), val2.Sequence())
	assert.Equal(t, val.Number(), val2.Number())
	assert.Equal(t, val.Stake(), val2.Stake())
	assert.Equal(t, val.LastBondingHeight(), val2.LastBondingHeight())
	assert.Equal(t, val.LastJoinedHeight(), val2.LastJoinedHeight())
	assert.Equal(t, val.UnbondingHeight(), val2.UnbondingHeight())
}

func TestJSONMarshaling(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt32(10000))

	js, err := json.Marshal(val)
	require.NoError(t, err)
	fmt.Println(string(js))
}

func TestInvalidData(t *testing.T) {
	_, err := ValidatorFromBytes([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestMarshalingRawData(t *testing.T) {
	bs, _ := hex.DecodeString("95167c2a0d86ec360407bce89b304616e1d0f83dbc200642abea8405e1838312fb8290b1230ebe4369cf1b7f556906c610ae92bcee544a1af79e259996e368b14851a1f8844274690b10df983bc2776ab10cc37e49e175bc7ae17ac919b8c34c01000000020000000300000000000000040000000500000006000000")
	val, err := ValidatorFromBytes(bs)
	require.NoError(t, err)
	fmt.Println(val)
	bs2, _ := val.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("76fea239a4586e8d9c2df9062b1958703341e3ece0f665c714da850101b61185")
	assert.Equal(t, val.Hash(), expected)
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

func TestSubtractStake(t *testing.T) {
	val, _ := GenerateTestValidator(100)
	bal := val.Stake()
	val.SubtractFromStake(1)
	assert.Equal(t, val.Stake(), bal-1)

}
