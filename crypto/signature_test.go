package crypto

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
)

func TestSignatureMarshaling(t *testing.T) {
	_, _, priv := RandomKeyPair()
	sig1 := priv.Sign(util.IntToSlice(util.RandInt(9999999999)))

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

	require.True(t, sig1.EqualsTo(*sig4))
	require.NoError(t, sig1.SanityCheck())
}

func TestSignatureFromBytes(t *testing.T) {
	_, err := SignatureFromRawBytes(nil)
	assert.Error(t, err)
	_, _, priv := RandomKeyPair()
	sig1 := priv.Sign(util.IntToSlice(util.RandInt(9999999999)))
	sig2, err := SignatureFromRawBytes(sig1.RawBytes())
	assert.NoError(t, err)
	require.True(t, sig1.EqualsTo(sig2))

	inv, _ := hex.DecodeString(strings.Repeat("ff", SignatureSize))
	_, err = SignatureFromRawBytes(inv)
	assert.Error(t, err)
}

func TestSignatureFromString(t *testing.T) {
	_, _, priv := RandomKeyPair()
	sig1 := priv.Sign(util.IntToSlice(util.RandInt(9999999999)))
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

	_, pb1, pv1 := RandomKeyPair()
	_, pb2, pv2 := RandomKeyPair()
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
	priv, err := PrivateKeyFromString("d0c6a560de2e60b6ac55386defefdf93b0c907290c2ad1b4dbd3338186bfdc68")
	assert.NoError(t, err)
	pub, err := PublicKeyFromString("37bfe636693eac0b674ae6603442192ef0432ad84384f0cec8bea5f63c9f45c29bf085b8b9b7f069ae873ccefe61a50a59ad3fefd729af5d63e9cb2325a8f064ab2514b3f846dbfded53234800603a9e752422ad48b99f835bcd95df945aac93")
	assert.NoError(t, err)
	sig, err := SignatureFromString("76da6c523c4abac463aad1ead5b7a042f143e354c346f6921a4975cc16959559e9b738fa197ab4df123f580a553b1596")
	assert.NoError(t, err)
	addr, err := AddressFromString("f6edd7e1d53d730a3ae0d44e6b6ce5dc102c0b63")
	assert.NoError(t, err)

	sig1 := priv.Sign(msg)
	assert.Equal(t, sig1.RawBytes(), sig.RawBytes())
	assert.True(t, pub.Verify(msg, sig))
	assert.Equal(t, pub.Address(), addr)
}
