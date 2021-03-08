package validator

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestMarshaling(t *testing.T) {
	val1, _ := GenerateTestValidator(util.RandInt(1000))
	assert.Equal(t, val1.PublicKey().Address(), val1.Address())
	val1.AddToStake(1)
	val1.IncSequence()
	val1.UpdateLastJoinedHeight(100)

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
	bs, _ := hex.DecodeString("A6015860DFF46FBCE5AE1BA4837DE551206176C0A74DEB5DFCA803228F570F7C9BA093EA109700559B72FE1D385492F0D5A10F17A4CEC41EB2E552F51E1F7F48AB311D4E195B1563C1FCBA8EE201173E4E6362CABEDACCEE541F9EFC9C4140D9FB268102021901B4031902F7041A2AF78F210514061864")
	val := new(Validator)
	err := val.Decode(bs)
	require.NoError(t, err)
	assert.Equal(t, val.Stake(), int64(720867105))
	assert.Equal(t, val.Sequence(), 759)
	assert.Equal(t, val.BondingHeight(), 20)
	assert.Equal(t, val.LastJoinedHeight(), 100)
	bs2, _ := val.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), crypto.HashH(bs))
	expected, _ := crypto.HashFromString("24118cc654fdc5333c222b40a932fccf0a058e3c0045a1e34d9298df4c128fba")
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
}
