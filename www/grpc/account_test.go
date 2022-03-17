package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetAccount(t *testing.T) {
	conn, client := callServer(t)
	acc := tMockState.TestStore.AddTestAccount()

	t.Run("Should return error for non-parsable address ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{
			Address: nil,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{
			Address: crypto.GenerateTestAddress().RawBytes(),
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return account details", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{
			Address: acc.Address().RawBytes(),
		})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Account.Address, acc.Address().RawBytes())
		assert.Equal(t, res.Account.Balance, acc.Balance())
		assert.Equal(t, int(res.Account.Number), acc.Number())
		assert.Equal(t, int(res.Account.Sequence), acc.Sequence())
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}
