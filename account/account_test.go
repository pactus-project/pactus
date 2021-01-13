package account

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
	acc1, _ := GenerateTestAccount(util.RandInt(10000))
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

func TestMarshalingRawData(t *testing.T) {
	bs, _ := hex.DecodeString("A401540C9819C4D4B1EDB7B70E6665287D4CE95401A37702191BD7031823041A007F5535")
	acc := new(Account)
	err := acc.Decode(bs)
	require.NoError(t, err)
	fmt.Println(acc)
	bs2, _ := acc.Encode()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, acc.Hash(), crypto.HashH(bs))
	expected, _ := crypto.HashFromString("8a010969c3bf76f764bcf6cf230b1fbc55e5c07b954dbf5595e59b854e97ce7c")
	assert.Equal(t, acc.Hash(), expected)
}

func TestIncSequence(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	seq := acc.Sequence()
	acc.IncSequence()
	assert.Equal(t, acc.Sequence(), seq+1)
	assert.Equal(t, acc.Number(), 100)
}

func TestSubtractBalance(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	bal := acc.Balance()
	acc.SubtractFromBalance(1)
	assert.Equal(t, acc.Balance(), bal-1)

}
