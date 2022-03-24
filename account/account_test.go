package account

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
	acc1, _ := GenerateTestAccount(util.RandInt32(10000))
	bs, err := acc1.Bytes()
	require.NoError(t, err)
	fmt.Printf("%X\n", bs)
	acc2, err := AccountFromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, acc1, acc2)
}

func TestJSONMarshaling(t *testing.T) {
	acc1, _ := GenerateTestAccount(util.RandInt32(10000))

	js, err := json.Marshal(acc1)
	require.NoError(t, err)
	fmt.Println(string(js))
}

func TestInvalidData(t *testing.T) {
	_, err := AccountFromBytes([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestMarshalingRawData(t *testing.T) {
	bs, _ := hex.DecodeString("01283993000F6484BF1E148B2B27CF11A602BBB2DC03000000020000000100000000000000")
	acc, err := AccountFromBytes(bs)
	require.NoError(t, err)
	fmt.Println(acc)
	bs2, _ := acc.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, acc.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("0e950ffaf53de12c6ea2d17e8fef96c51937f5c48ab0e3cb9af6a8af4dcf290e")
	assert.Equal(t, acc.Hash(), expected)
}

func TestIncSequence(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	seq := acc.Sequence()
	acc.IncSequence()
	assert.Equal(t, acc.Sequence(), seq+1)
	assert.Equal(t, acc.Number(), int32(100))
}

func TestSubtractBalance(t *testing.T) {
	acc, _ := GenerateTestAccount(100)
	bal := acc.Balance()
	acc.SubtractFromBalance(1)
	assert.Equal(t, acc.Balance(), bal-1)

}
