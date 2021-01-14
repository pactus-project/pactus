package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func getAccount(t *testing.T, addr crypto.Address) *account.Account {
	for i := 0; i < 20; i++ {
		res := tCapnpServer.GetAccount(tCtx, func(p capnp.ZarbServer_getAccount_Params) error {
			p.SetAddress([]byte(addr.String()))
			return nil
		}).Result()

		st, err := res.Struct()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		d, _ := st.Data()
		acc := new(account.Account)
		assert.NoError(t, acc.Decode(d))
		return acc
	}
	require.NoError(t, fmt.Errorf("timeout"))
	return nil
}

func TestTreasuryAccount(t *testing.T) {

	res := getAccount(t, crypto.TreasuryAddress)
	require.NotNil(t, res)
	assert.Equal(t, tGenDoc.Accounts()[0].Balance(), res.Balance()+int64(res.Sequence()*500000000))
}
