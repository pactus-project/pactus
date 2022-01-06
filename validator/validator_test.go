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

func TestMarshaling(t *testing.T) {
	val1, _ := GenerateTestValidator(util.RandInt(1000))
	assert.Equal(t, val1.PublicKey().Address(), val1.Address())
	val1.AddToStake(1)
	val1.IncSequence()
	val1.UpdateLastJoinedHeight(50)

	bs, err := val1.Encode()
	require.NoError(t, err)
	val2 := new(Validator)
	err = val2.Decode(bs)
	require.NoError(t, err)
	assert.Equal(t, val1.Hash(), val2.Hash())

	val3 := new(Validator)
	err = val3.Decode(bs)
	require.NoError(t, err)
	assert.Equal(t, val2, val3)

	/// test json marshaing/unmarshaling
	js, err := json.Marshal(val1)
	require.NoError(t, err)
	fmt.Println(string(js))
	val4 := new(Validator)
	require.NoError(t, json.Unmarshal(js, val4))

	assert.Equal(t, val3, val4)

	/// should fail
	val5 := new(Validator)
	err = val5.Decode([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestMarshalingRawData(t *testing.T) {
	bs, _ := hex.DecodeString("a701586090bc71f9422cd2f7274f00bbf2bb74ad601593ddd093b4e879172c780d02e948021694c0c8c09d2e022d54c335f63ad813efedbbc6bdd7a8934b7743aa62e3bc2ffd2071e452eb802ef456acd79ec3313bec5da04233128087cad3dca185ed4c0219016f031902f7041a2af78f210514061864071832")
	val := new(Validator)
	err := val.Decode(bs)
	require.NoError(t, err)
	assert.Equal(t, val.Stake(), int64(720867105))
	assert.Equal(t, val.Sequence(), 759)
	assert.Equal(t, val.LastBondingHeight(), 20)
	assert.Equal(t, val.UnbondingHeight(), 100)
	assert.Equal(t, val.LastJoinedHeight(), 50)
	bs2, _ := val.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("2d5e898cd60320512fa70c362b812e5dee0c42af6391d2cbfca32af4d5b1bbf5")
	assert.Equal(t, val.Hash(), expected)
}

func TestIncSequence(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt(1000))
	seq := val.Sequence()
	val.IncSequence()
	assert.Equal(t, val.Sequence(), seq+1)
}

func TestNumber(t *testing.T) {
	val, _ := GenerateTestValidator(5)
	assert.Equal(t, val.Number(), 5)
}

func TestPower(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt(1000))
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
