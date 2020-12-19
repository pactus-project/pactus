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
	bs, _ := hex.DecodeString("A6015860741CCD3D7F9CD9FEE7C9A412EE73E248F4522018EF6EF7FCCF456E02A2DFB39EB3CEBAE083BDC2CF10203D3BDD52A90D4EA878FEBAB2D87AD7323A5F82BAFA32A4A158FE8C96EAE6B8C30440256D691B147D864313E77E181FAC6AAA2206DE980218A803190339041A2AD6EC630518570600")
	val := new(Validator)
	err := val.Decode(bs)
	require.NoError(t, err)
	fmt.Println(val)
	bs2, _ := val.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), crypto.HashH(bs))
	expected, _ := crypto.HashFromString("be773d118bc63b3ca299aabae969c6d9964185df381a7b733c119c5c4c5defb0")
	assert.Equal(t, val.Hash(), expected)
}

func TestIncSequence(t *testing.T) {
	val, _ := GenerateTestValidator(util.RandInt(1000))
	seq := val.Sequence()
	val.IncSequence()
	assert.Equal(t, val.Sequence(), seq+1)
}
