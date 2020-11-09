package crypto

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalingEmptyAddress(t *testing.T) {
	addr1 := Address{}

	js, err := json.Marshal(addr1)
	assert.NoError(t, err)
	var addr2 Address
	err = json.Unmarshal(js, &addr2)
	assert.NoError(t, err) /// No error
	assert.Equal(t, addr1, addr2)

	bs, err := addr1.MarshalCBOR()
	assert.NoError(t, err)
	var addr3 Address
	err = addr3.UnmarshalCBOR(bs)
	assert.NoError(t, err) /// No error
	assert.Equal(t, addr1, addr3)
}

func TestMarshalingAddress(t *testing.T) {
	addrs := []string{
		"0123456789ABCDEF0123456789ABCDEF01234567",
		"7777777777777777777777777777777777777777",
		"B03DD2C47852775208A56FA10A49875ABC507343",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
	}

	for _, addr := range addrs {
		bs, _ := hex.DecodeString(addr)
		ac1, err := AddressFromRawBytes(bs)
		assert.NoError(t, err)
		fmt.Println(ac1.String())

		jac, err := json.Marshal(&ac1)
		assert.NoError(t, err)
		fmt.Println(string(jac))

		var ac2 Address
		assert.NoError(t, json.Unmarshal(jac, &ac2))
		require.Equal(t, ac1, ac2)

		bac, err := ac1.MarshalCBOR()
		assert.NoError(t, err)
		fmt.Println(string(jac))

		var ac3 Address
		assert.NoError(t, ac3.UnmarshalCBOR(bac))

		require.Equal(t, ac1, ac2)
	}
}
