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
	acc, _ := GenerateTestAccount(util.RandInt32(10000))
	bs, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, acc.SerializeSize(), len(bs))
	acc2, err := AccountFromBytes(bs)
	require.NoError(t, err)
	assert.Equal(t, acc, acc2)
}

func TestJSONMarshaling(t *testing.T) {
	acc, _ := GenerateTestAccount(util.RandInt32(10000))

	js, err := json.Marshal(acc)
	require.NoError(t, err)
	fmt.Println(string(js))
}

func TestInvalidData(t *testing.T) {
	_, err := AccountFromBytes([]byte("asdfghjkl"))
	require.Error(t, err)
}

func TestMarshalingRawData(t *testing.T) {
	bs, _ := hex.DecodeString("01283993000F6484BF1E148B2B27CF11A602BBB2DC01000000020000000300000000000000")
	acc, err := AccountFromBytes(bs)
	require.NoError(t, err)
	fmt.Println(acc)
	bs2, _ := acc.Bytes()
	assert.Equal(t, bs, bs2)
	assert.Equal(t, acc.Hash(), hash.CalcHash(bs))
	expected, _ := hash.FromString("a56021e105f1fd644864d7813a131b68b1b447c4abf19a7a44df54deab5b7091")
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
