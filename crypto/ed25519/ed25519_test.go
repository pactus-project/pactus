package ed25519_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/stretchr/testify/assert"
)

func TestSigning(t *testing.T) {
	msg := []byte("pactus")
	prv, _ := ed25519.PrivateKeyFromString(
		"SECRET1RYY62A96X25ZAL4DPL5Z63G83GCSFCCQ7K0CMQD3MFNLYK3A6R26QUUK3Y0")
	pub, _ := ed25519.PublicKeyFromString(
		"public1rvqxnpfph8tnc3ck55z85w285t5jetylmmktr9wlzs0zvx7kx500szxfudh")
	sig, _ := ed25519.SignatureFromString(
		"361aaa09c408bfcf7e79dd90c583eeeaefe7c732ca5643cfb2ea7a6d22105b87" +
			"4a412080525a855bbd5df94110a7d0083d6e386e016ecf8b7f522c339f79d305")
	addr, _ := crypto.AddressFromString("pc1r7jkvfnegf0rf5ua05fzu9krjhjxcrrygl3v4nl")

	sig1 := prv.Sign(msg)
	assert.Equal(t, sig.Bytes(), sig1.Bytes())
	assert.NoError(t, pub.Verify(msg, sig))
	assert.Equal(t, pub, prv.PublicKey())
	assert.Equal(t, addr, pub.AccountAddress())
}
