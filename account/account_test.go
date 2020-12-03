package account

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
)

type acc struct {
	Test crypto.Hash
}
type acb struct {
	Test *crypto.Hash
}

func TestMarshaling23(t *testing.T) {
	h := crypto.HashH([]byte("a"))
	a1 := acc{}
	a2 := acc{Test: crypto.HashH([]byte("a"))}
	a3 := acb{}
	a4 := acb{Test: nil}
	a5 := acb{Test: &h}

	bz1, _ := cbor.Marshal(&a1)
	bz2, _ := cbor.Marshal(&a2)
	bz3, _ := cbor.Marshal(&a3)
	bz4, _ := cbor.Marshal(&a4)
	bz5, _ := cbor.Marshal(&a5)

	fmt.Printf("%x\n", bz1)
	fmt.Printf("%x\n", bz2)
	fmt.Printf("%x\n", bz3)
	fmt.Printf("%x\n", bz4)
	fmt.Printf("%x\n", bz5)

}

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
	bs, _ := hex.DecodeString("a401581a12ec0b329db4db42793200be69c1cd1eb4ad575b94e206aa335b020103060443020304")
	acc := new(Account)
	err := acc.Decode(bs)
	require.NoError(t, err)
	fmt.Println(acc)
}
