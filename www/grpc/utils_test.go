package grpc

import (
	"context"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestSignMessageWithPrivateKey(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)
	conn, client := td.utilClient(t)

	msg := "pactus"
	prvStr := "SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67"
	invalidPrvStr := "INVSECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67"
	expectedSig := "923d67a8624cbb7972b29328e15ec76cc846076ccf00a9e94d991c677846f334ae4ba4551396fbcd6d1cab7593baf3b7"

	t.Run("", func(t *testing.T) {
		res, err := client.SignMessageWithPrivateKey(context.Background(),
			&pactus.SignMessageWithPrivateKeyRequest{
				Message:    msg,
				PrivateKey: prvStr,
			})

		assert.Nil(t, err)
		assert.Equal(t, expectedSig, res.Signature)
	})

	t.Run("", func(t *testing.T) {
		res, err := client.SignMessageWithPrivateKey(context.Background(),
			&pactus.SignMessageWithPrivateKeyRequest{
				Message:    msg,
				PrivateKey: invalidPrvStr,
			})

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestVerifyMessage(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)
	conn, client := td.utilClient(t)

	msg := "pactus"
	pubStr := "public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
		"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx47a"
	sigStr := "923d67a8624cbb7972b29328e15ec76cc846076ccf00a9e94d991c677846f334ae4ba4551396fbcd6d1cab7593baf3b7"
	invalidSigStr := "113d67a8624cbb7972b29328e15ec76cc846076ccf00a9e94d991c677846f334ae4ba4551396fbcd6d1cab7593baf3c9"

	t.Run("valid message", func(t *testing.T) {
		res, err := client.VerifyMessage(context.Background(),
			&pactus.VerifyMessageRequest{
				Message:   msg,
				Signature: sigStr,
				PublicKey: pubStr,
			})
		assert.Nil(t, err)
		assert.True(t, res.IsValid)
	})

	t.Run("invalid message", func(t *testing.T) {
		res, err := client.VerifyMessage(context.Background(),
			&pactus.VerifyMessageRequest{
				Message:   msg,
				Signature: invalidSigStr,
				PublicKey: pubStr,
			})

		assert.Nil(t, err)
		assert.False(t, res.IsValid)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestBLSPublicKeyAggregate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	conf := testConfig()
	td := setup(t, conf)
	conn, client := td.utilClient(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()
	pub3, _ := ts.RandBLSKeyPair()
	aggPub := bls.PublicKeyAggregate(pub1, pub2, pub3)
	invalidPub := "invalidpub"

	t.Run("zero public keys", func(t *testing.T) {
		res, err := client.BLSPublicKeyAggregate(context.Background(),
			&pactus.BLSPublicKeyAggregateRequest{
				PublicKeys: []string{},
			})

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("only one public key", func(t *testing.T) {
		res, err := client.BLSPublicKeyAggregate(context.Background(),
			&pactus.BLSPublicKeyAggregateRequest{
				PublicKeys: []string{pub1.String()},
			})

		assert.Nil(t, err)
		assert.Equal(t, pub1.String(), res.PublicKey)
	})

	t.Run("invalid public key", func(t *testing.T) {
		res, err := client.BLSPublicKeyAggregate(context.Background(),
			&pactus.BLSPublicKeyAggregateRequest{
				PublicKeys: []string{pub1.String(), pub2.String(), invalidPub, pub3.String()},
			})

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("valid public keys", func(t *testing.T) {
		res, err := client.BLSPublicKeyAggregate(context.Background(),
			&pactus.BLSPublicKeyAggregateRequest{
				PublicKeys: []string{pub1.String(), pub2.String(), pub3.String()},
			})

		assert.Nil(t, err)
		assert.Equal(t, aggPub.String(), res.PublicKey)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestBLSSignatureAggregate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	conf := testConfig()
	td := setup(t, conf)
	conn, client := td.utilClient(t)

	sig1 := ts.RandBLSSignature()
	sig2 := ts.RandBLSSignature()
	sig3 := ts.RandBLSSignature()
	aggSig := bls.SignatureAggregate(sig1, sig2, sig3)
	invalidSig := "invalidsig"

	t.Run("zero signatures", func(t *testing.T) {
		res, err := client.BLSSignatureAggregate(context.Background(),
			&pactus.BLSSignatureAggregateRequest{
				Signatures: []string{},
			})

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("only one signature", func(t *testing.T) {
		res, err := client.BLSSignatureAggregate(context.Background(),
			&pactus.BLSSignatureAggregateRequest{
				Signatures: []string{sig1.String()},
			})

		assert.Nil(t, err)
		assert.Equal(t, sig1.String(), res.Signature)
	})

	t.Run("invalid signature", func(t *testing.T) {
		res, err := client.BLSSignatureAggregate(context.Background(),
			&pactus.BLSSignatureAggregateRequest{
				Signatures: []string{sig1.String(), sig2.String(), invalidSig, sig3.String()},
			})

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("valid signatures", func(t *testing.T) {
		res, err := client.BLSSignatureAggregate(context.Background(),
			&pactus.BLSSignatureAggregateRequest{
				Signatures: []string{sig1.String(), sig2.String(), sig3.String()},
			})

		assert.Nil(t, err)
		assert.Equal(t, aggSig.String(), res.Signature)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
