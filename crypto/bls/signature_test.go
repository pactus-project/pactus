package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

func TestSignatureCBORMarshaling(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.Uint64ToSlice(util.RandUint64(0)))
	sig2 := new(Signature)

	bs, err := sig1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))
	assert.NoError(t, sig1.SanityCheck())

	inv, _ := hex.DecodeString(strings.Repeat("ff", SignatureSize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, sig2.UnmarshalCBOR(data))
}

func TestSignatureJSONMarshaling(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig := prv.Sign(util.Uint64ToSlice(util.RandUint64(0)))
	js, err := sig.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(js), sig.String())
}

func TestSignatureEncoding(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig := prv.Sign(util.Uint64ToSlice(util.RandUint64(0)))
	w1 := util.NewFixedWriter(20)
	assert.Error(t, sig.Encode(w1))

	w2 := util.NewFixedWriter(SignatureSize)
	assert.NoError(t, sig.Encode(w2))

	r1 := util.NewFixedReader(20, w2.Bytes())
	assert.Error(t, sig.Decode(r1))

	r2 := util.NewFixedReader(SignatureSize, w2.Bytes())
	assert.NoError(t, sig.Decode(r2))
}

func TestSignatureFromString(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.Uint64ToSlice(util.RandUint64(0)))
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
	assert.Empty(t, sig1.Bytes())

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

	assert.NotEqual(t, sig1, sig2)
	assert.NoError(t, pb1.Verify(msg, sig1))
	assert.NoError(t, pb2.Verify(msg, sig2))
	assert.Equal(t, errors.Code(pb1.Verify(msg, sig2)), errors.ErrInvalidSignature)
	assert.Equal(t, errors.Code(pb2.Verify(msg, sig1)), errors.ErrInvalidSignature)
	assert.Equal(t, errors.Code(pb1.Verify(msg[1:], sig1)), errors.ErrInvalidSignature)
}

func TestSigning(t *testing.T) {
	msg := []byte("zarb")
	prv, _ := PrivateKeyFromString("68dcbf868133d3dbb4d12a0c2907c9b093dfefef6d3855acb6602ede60a5c6d0")
	pub, _ := PublicKeyFromString("af0f74917f5065af94727ae9541b0ddcfb5b828a9e016b02498f477ed37fb44d5d882495afb6fd4f9773e4ea9deee436030c4d61c6e3a1151585e1d838cae1444a438d089ce77e10c492a55f6908125c5be9b236a246e4082d08de564e111e65")
	sig, _ := SignatureFromString("ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792")
	addr, _ := crypto.AddressFromString("zc15x2a0lkt5nrrdqe0rkcv6r4pfkmdhrr39g6klh")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig1.Bytes(), sig.Bytes())
	assert.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, pub.Address(), addr)
}

func TestSignatureSanityCheck(t *testing.T) {
	sig, err := SignatureFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, sig.SanityCheck())
}
