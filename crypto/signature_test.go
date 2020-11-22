package crypto

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalingEmptySignature(t *testing.T) {
	sig1 := Signature{}

	js, err := json.Marshal(sig1)
	assert.Error(t, err)
	sig2 := new(Signature)
	err = json.Unmarshal(js, &sig2)
	assert.Error(t, err)

	bs, err := sig1.MarshalCBOR()
	assert.Error(t, err)

	sig3 := new(Signature)
	err = sig3.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestMarshalingSignature(t *testing.T) {
	_, _, privKey := RandomKeyPair()
	sig1 := privKey.Sign([]byte("Test message"))
	sig11 := privKey.Sign([]byte("Test message"))
	require.Equal(t, sig1, sig11)

	bs, err := sig1.MarshalText()
	fmt.Println(string(bs))
	require.NoError(t, err)

	var sig2 Signature
	err = sig2.UnmarshalText(bs)
	require.NoError(t, err)
	require.True(t, sig1.EqualsTo(sig2))

	bs, err = sig2.MarshalCBOR()
	assert.NoError(t, err)

	var sig3 Signature
	assert.NoError(t, sig3.UnmarshalCBOR(bs))
	require.True(t, sig3.EqualsTo(sig2))
}

func TestVerifyingSignature(t *testing.T) {
	msg := []byte("zaeb")

	_, pb1, pv1 := RandomKeyPair()
	_, pb2, pv2 := RandomKeyPair()
	sig1 := pv1.Sign(msg)
	sig2 := pv2.Sign(msg)

	require.NotEqual(t, sig1, sig2)
	require.True(t, pb1.Verify(msg, sig1))
	require.True(t, pb2.Verify(msg, sig2))
	require.False(t, pb1.Verify(msg, sig2))
	require.False(t, pb2.Verify(msg, sig1))
	require.False(t, pb1.Verify(msg[1:], sig1))
}
