package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
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
	assert.Error(t, err)
	assert.Empty(t, pub1.String())
	assert.Empty(t, pub1.Bytes())

	var pub2 PublicKey
	err = pub2.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestPublicKeySanityCheck(t *testing.T) {
	pub, err := PublicKeyFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, pub.SanityCheck())
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, _ := GenerateTestKeyPair()

	assert.True(t, pub1.VerifyAddress(pub1.Address()))
	assert.False(t, pub1.VerifyAddress(pub2.Address()))
}
