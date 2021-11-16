package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/util"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetAccount(t *testing.T) {
	conn, client := callServer(t)
	acc1, _ := account.GenerateTestAccount(util.RandInt(10000))
	t.Run("Should return error for non-parsable address ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{Address: "NON_EXISTING_ADDRESS"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{Address: acc1.Address().String()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	tMockState.Store.UpdateAccount(acc1)

	t.Run("Should return account details", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{Address: acc1.Address().String()})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Account.Address, acc1.Address().String())
		assert.Equal(t, res.Account.Balance, acc1.Balance())
		assert.Equal(t, int(res.Account.Number), acc1.Number())
		assert.Equal(t, int(res.Account.Sequence), acc1.Sequence())
	})
	conn.Close()

}
