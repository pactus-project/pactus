package bls

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestSignatureMarshaling(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.IntToSlice(util.RandInt(9999999999)))
	sig2 := new(Signature)

	bs, err := sig1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))
	assert.NoError(t, sig1.SanityCheck())

	js, err := sig1.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(js), sig1.String())

	inv, _ := hex.DecodeString(strings.Repeat("ff", PublicKeySize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, sig2.UnmarshalCBOR(data))
}

func TestSignatureFromString(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.IntToSlice(util.RandInt(9999999999)))
	sig2, err := SignatureFromString(sig1.String())
	assert.NoError(t, err)
	assert.True(t, sig1.EqualsTo(sig2))

	_, err = SignatureFromString("")
	assert.Error(t, err)

	_, err = SignatureFromString("inv")
	assert.Error(t, err)

	_, err = SignatureFromString("00")
	assert.Error(t, err)
}

func TestSignatureEmpty(t *testing.T) {
	sig1 := Signature{}

	bs, err := sig1.MarshalCBOR()
	assert.Error(t, err)
	assert.Empty(t, sig1.String())
	assert.Empty(t, sig1.RawBytes())

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

	assert.NotEqual(t, sig1, sig2)
	assert.True(t, pb1.Verify(msg, sig1))
	assert.True(t, pb2.Verify(msg, sig2))
	assert.False(t, pb1.Verify(msg, sig2))
	assert.False(t, pb2.Verify(msg, sig1))
	assert.False(t, pb1.Verify(msg[1:], sig1))
}

func TestSigning(t *testing.T) {
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

func TestSignatureSanityCheck(t *testing.T) {
	sig, err := SignatureFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, sig.SanityCheck())
}
