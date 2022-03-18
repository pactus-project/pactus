package bls

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestSignatureMarshaling(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.IntToSlice(util.RandInt(9999999999)))

	sig2 := new(Signature)
	sig3 := new(Signature)
	sig4 := new(Signature)

	js, err := json.Marshal(sig1)
	assert.NoError(t, err)
	require.Error(t, sig2.UnmarshalJSON([]byte("bad")))
	require.NoError(t, json.Unmarshal(js, sig2))

	bs, err := sig2.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, sig3.UnmarshalCBOR(bs))

	txt, err := sig2.MarshalText()
	assert.NoError(t, err)
	assert.NoError(t, sig4.UnmarshalText(txt))

	require.True(t, sig1.EqualsTo(sig4))
	require.NoError(t, sig1.SanityCheck())
}

func TestSignatureFromBytes(t *testing.T) {
	_, err := SignatureFromRawBytes(nil)
	assert.Error(t, err)
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.IntToSlice(util.RandInt(9999999999)))
	sig2, err := SignatureFromRawBytes(sig1.RawBytes())
	assert.NoError(t, err)
	require.True(t, sig1.EqualsTo(sig2))

	inv, _ := hex.DecodeString(strings.Repeat("ff", SignatureSize))
	_, err = SignatureFromRawBytes(inv)
	assert.Error(t, err)
}

func TestSignatureFromString(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.IntToSlice(util.RandInt(9999999999)))
	sig2, err := SignatureFromString(sig1.String())
	assert.NoError(t, err)
	require.True(t, sig1.EqualsTo(sig2))

	_, err = SignatureFromString("inv")
	assert.Error(t, err)
}

func TestMarshalingEmptySignature(t *testing.T) {
	sig1 := Signature{}

	js, err := json.Marshal(sig1)
	assert.NoError(t, err)
	assert.Equal(t, js, []byte{0x22, 0x22}) // ""
	sig2 := new(Signature)
	err = json.Unmarshal(js, &sig2)
	assert.Error(t, err)

	bs, err := sig1.MarshalCBOR()
	assert.Error(t, err)

	sig3 := new(Signature)
	err = sig3.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestVerifyingSignature(t *testing.T) {
	msg := []byte("zarb")

	pb1, pv1 := GenerateTestKeyPair()
	pb2, pv2 := GenerateTestKeyPair()
	sig1 := pv1.Sign(msg)
	sig2 := pv2.Sign(msg)

	fmt.Printf("%x\n", pb1.RawBytes())
	fmt.Printf("%x\n", pv1.RawBytes())
	fmt.Printf("%x\n", sig1.RawBytes())

	require.NotEqual(t, sig1, sig2)
	require.True(t, pb1.Verify(msg, sig1))
	require.True(t, pb2.Verify(msg, sig2))
	require.False(t, pb1.Verify(msg, sig2))
	require.False(t, pb2.Verify(msg, sig1))
	require.False(t, pb1.Verify(msg[1:], sig1))
}

func TestSignature(t *testing.T) {
	msg := []byte("zarb")
	prv, _ := PrivateKeyFromString("68dcbf868133d3dbb4d12a0c2907c9b093dfefef6d3855acb6602ede60a5c6d0")
	pub, _ := PublicKeyFromString("af0f74917f5065af94727ae9541b0ddcfb5b828a9e016b02498f477ed37fb44d5d882495afb6fd4f9773e4ea9deee436030c4d61c6e3a1151585e1d838cae1444a438d089ce77e10c492a55f6908125c5be9b236a246e4082d08de564e111e65")
	sig, _ := SignatureFromString("a2d06b33af2c9e7ca878da85a96b2c2346f4306d0473bdabc38be87c19dae5e67e08724a5220d0e372fb080bbd2fbde9")
	addr, _ := crypto.AddressFromString("zc15x2a0lkt5nrrdqe0rkcv6r4pfkmdhrr39g6klh")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig1.RawBytes(), sig.RawBytes())
	assert.True(t, pub.Verify(msg, sig))
	assert.Equal(t, pub.Address(), addr)
}
