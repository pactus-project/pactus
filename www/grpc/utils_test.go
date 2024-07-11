package grpc

import (
	"context"
	"testing"

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
		assert.Equal(t, res.Signature, expectedSig)
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

	t.Run("", func(t *testing.T) {
		res, err := client.VerifyMessage(context.Background(),
			&pactus.VerifyMessageRequest{
				Message:   msg,
				Signature: sigStr,
				PublicKey: pubStr,
			})
		assert.Nil(t, err)
		assert.True(t, res.IsValid)
	})

	t.Run("", func(t *testing.T) {
		_, err := client.VerifyMessage(context.Background(),
			&pactus.VerifyMessageRequest{
				Message:   msg,
				Signature: invalidSigStr,
				PublicKey: pubStr,
			})

		assert.NotNil(t, err)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
