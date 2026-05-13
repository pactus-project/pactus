package secp256k1_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncoding uses PIP-53 vectors for secp256k1 key/address encoding and signing.
func TestEncoding(t *testing.T) {
	prvData, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	pubData, _ := hex.DecodeString("036d6caac248af96f6afa7f904f550253a0f3ef3f5aa2fe6838a95b216691468e2")
	addrData, _ := hex.DecodeString("042bc1db7e0797c45b918dc401093c9257c6012b4c")

	prvStr := "SECRET1YQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SPVXU8Z"
	pubStr := "public1yqdkke2kzfzheda405lusfa2sy5aq70hn7k4zle5r322my9nfz35wyfamrfs"
	addrStr := "pc1y90qakls8jlz9hyvdcsqsj0yj2lrqz26vqu7l0z"

	prv, _ := secp256k1.PrivateKeyFromString(prvStr)
	pub, _ := secp256k1.PublicKeyFromString(pubStr)
	addr, _ := crypto.AddressFromString(addrStr)

	assert.Equal(t, prvData, prv.Bytes())
	assert.Equal(t, pubData, pub.Bytes())
	assert.Equal(t, addrData, addr.Bytes())

	msg := []byte("pactus")
	sig, _ := secp256k1.SignatureFromString("16e6f8bcdb92964a35773aae200628a5b470b6488d42ceef6538da0b4ffd3b42098dd821eea96f66ba02c9c4473443ab51c411ab78adfbb90d53b07ca1d6862b")

	require.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, sig.Bytes(), prv.Sign(msg).Bytes())
	assert.True(t, pub.EqualsTo(prv.PublicKey()))
	assert.Equal(t, addr, pub.AccountAddress())
}
