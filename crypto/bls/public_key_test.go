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

func TestPublicKeyCBORMarshaling(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2 := new(PublicKey)

	bs, err := pub1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, pub2.UnmarshalCBOR(bs))
	assert.True(t, pub1.EqualsTo(pub2))
	assert.NoError(t, pub1.SanityCheck())

	assert.Error(t, pub2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", PublicKeySize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, pub2.UnmarshalCBOR(data))
}

func TestPublicKeyEncoding(t *testing.T) {
	pub, _ := GenerateTestKeyPair()
	w1 := util.NewFixedWriter(20)
	assert.Error(t, pub.Encode(w1))

	w2 := util.NewFixedWriter(PublicKeySize)
	assert.NoError(t, pub.Encode(w2))

	r1 := util.NewFixedReader(20, w2.Bytes())
	assert.Error(t, pub.Decode(r1))

	r2 := util.NewFixedReader(PublicKeySize, w2.Bytes())
	assert.NoError(t, pub.Decode(r2))
}

func TestPublicKeyFromString(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, err := PublicKeyFromString(pub1.String())
	assert.NoError(t, err)
	assert.True(t, pub1.EqualsTo(pub2))

	_, err = PublicKeyFromString("")
	assert.Error(t, err)

	_, err = PublicKeyFromString("inv")
	assert.Error(t, err)

	_, err = PublicKeyFromString("00")
	assert.Error(t, err)
}

func TestPublicKeyEmpty(t *testing.T) {
	pub1 := PublicKey{}

	bs, err := pub1.MarshalCBOR()
	assert.NoError(t, err)
	assert.Error(t, pub1.SanityCheck())

	var pub2 PublicKey
	err = pub2.UnmarshalCBOR(bs)
	assert.NoError(t, err)
	assert.Error(t, pub2.SanityCheck())
}

func TestPublicKeySanityCheck(t *testing.T) {
	pub, err := PublicKeyFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, pub.SanityCheck())
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, _ := GenerateTestKeyPair()

	assert.NoError(t, pub1.VerifyAddress(pub1.Address()))
	assert.Equal(t, errors.Code(pub1.VerifyAddress(pub2.Address())), errors.ErrInvalidAddress)
}

func TestNilPublicKey(t *testing.T) {
	pub := &PublicKey{}
	assert.Error(t, pub.VerifyAddress(crypto.GenerateTestAddress()))
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &Signature{}))
}

func TestNilSignature(t *testing.T) {
	pub, _ := GenerateTestKeyPair()
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &Signature{}))
}

func TestInfinitySignature(t *testing.T) {
	pub, err := PublicKeyFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	sig, err := SignatureFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, pub.Verify(nil, sig))
}
