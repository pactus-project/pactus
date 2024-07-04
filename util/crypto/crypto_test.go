package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignMessageWithPrivateKey(t *testing.T) {
	msg := "pactus"
	prvStr :=
		"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67"
	sigStr := 
		"923d67a8624cbb7972b29328e15ec76cc846076ccf00a9e94d991c677846f334ae4ba4551396fbcd6d1cab7593baf3b7"

		sig, err := SignMessageWithPrivateKey(prvStr, msg)
		assert.Nil(t, err)
		assert.Equal(t, sig, sigStr)
}

func TestVerifyMessage(t *testing.T) {
	msg := "pactus"
	pubStr := 
		"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
			"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx47a"
	sigStr := 
		"923d67a8624cbb7972b29328e15ec76cc846076ccf00a9e94d991c677846f334ae4ba4551396fbcd6d1cab7593baf3b7"

		ok, err := VerifyMessage(sigStr,pubStr,msg)
		assert.Nil(t, err)
		assert.True(t, ok)
}
