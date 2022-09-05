package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGenerateMnemonic(t *testing.T) {
	conn, client := callWalletSerer(t)

	t.Run("Should return mnemonic", func(t *testing.T) {
		res, err := client.GenerateMnemonic(tCtx, &zarb.GenerateMnemonicRequest{Language: "english", Entropy: 128})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Mnemonic)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
