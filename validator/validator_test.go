package validator

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestMarshaling(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val1.AddToStake(1)
	val1.IncSequence()

	bs, err := val1.Encode()
	fmt.Printf("%X\n", bs)
	fmt.Printf("%X", val1.Address().RawBytes())
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
	bs, _ := hex.DecodeString("A5015860765981A07FC9C14FDF07963689049A224D775CC680982AB8F8F966A0F0079E53A5077037A315FAAA9970476E493D3D096B067BD4E4D7833FB3660895DF820363E63738B0EBE6E2DD3808670EDA8D5F782D80D8F68CB23825A436A2BD3881B70F02182A031874041A24C5455E05191FBC")
	val := new(Validator)
	err := val.Decode(bs)
	require.NoError(t, err)
	fmt.Println(val)
	bs2, _ := val.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), crypto.HashH(bs))
	expected, _ := crypto.HashFromString("74d2080fcd05f48956ef5035fff25d312463221649359655754e08f52fdb6c1c")
	assert.Equal(t, val.Hash(), expected)
}

func TestAddToStake(t *testing.T) {
	val, _ := GenerateTestValidator()
	amt := val.Stake()

	assert.Error(t, val.AddToStake(-1))
	assert.NoError(t, val.AddToStake(1))
	assert.Error(t, val.SubtractFromStake(-2))
	assert.NoError(t, val.SubtractFromStake(2))
	assert.Error(t, val.SubtractFromStake(amt))
	assert.Equal(t, val.Stake(), amt-1)
}

func TestIncSequence(t *testing.T) {
	val, _ := GenerateTestValidator()
	seq := val.Sequence()
	val.IncSequence()
	assert.Equal(t, val.Sequence(), seq+1)
}
