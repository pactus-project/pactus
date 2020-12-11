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

func TestMarshaling2(t *testing.T) {
	bs, _ := hex.DecodeString("A40158608FC6D6D9EE584AADFAB23D374836199DA2A26731B2482D4158E43082A83F3319B6EAD1599C03751A2C6EFF187B85D504302654775A6404E3B5EBBDA5E263B3E8E42FBA631EA2A0428571EF62F3B5247B5BA8431166647860A0DC518FE3223C07021861031A004F626E0400")
	val := new(Validator)
	err := val.Decode(bs)
	require.NoError(t, err)
	fmt.Println(val)
	bs2, _ := val.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, val.Hash(), crypto.HashH(bs))
	expected,_:=crypto.HashFromString("113f6b87a935601398993d687143241e8547ea39ba991625273e8638b5a0373f")
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
