package account

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
	acc1, _ := GenerateTestAccount()
	acc1.AddToBalance(1)
	acc1.IncSequence()

	bs, err := acc1.Encode()
	fmt.Printf("%X\n", bs)
	fmt.Printf("%X", acc1.Address().RawBytes())
	require.NoError(t, err)
	acc2 := new(Account)
	err = acc2.Decode(bs)
	require.NoError(t, err)
	assert.Equal(t, acc1, acc2)

	acc3 := new(Account)
	err = acc3.Decode(bs)
	require.NoError(t, err)
	assert.Equal(t, acc2, acc3)

	/// test json marshaing/unmarshaling
	js, err := json.Marshal(acc1)
	require.NoError(t, err)
	fmt.Println(string(js))
	acc4 := new(Account)
	require.NoError(t, json.Unmarshal(js, acc4))

	assert.Equal(t, acc3, acc4)

	/// should fail
	acc5 := new(Account)
	err = acc5.Decode([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestMarshaling2(t *testing.T) {
	bs, _ := hex.DecodeString("A3015427E1E8F8BBB9B5EB067EC71FBA278310173E3356021822031A0091398A")
	acc := new(Account)
	err := acc.Decode(bs)
	require.NoError(t, err)
	fmt.Println(acc)
	bs2, _ := acc.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, acc.Hash(), crypto.HashH(bs))
}

func TestAddToBalance(t *testing.T) {
	acc, _ := GenerateTestAccount()
	amt := acc.Balance()

	assert.Error(t, acc.AddToBalance(-1))
	assert.NoError(t, acc.AddToBalance(1))
	assert.Error(t, acc.SubtractFromBalance(-2))
	assert.NoError(t, acc.SubtractFromBalance(2))
	assert.Error(t, acc.SubtractFromBalance(amt))
	assert.Equal(t, acc.Balance(), amt-1)
}

func TestIncSequence(t *testing.T) {
	acc, _ := GenerateTestAccount()
	seq := acc.Sequence()
	acc.IncSequence()
	assert.Equal(t, acc.Sequence(), seq+1)
}
