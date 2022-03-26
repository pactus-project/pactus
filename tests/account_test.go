package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func getAccount(t *testing.T, addr crypto.Address) *account.Account {
	res := tCapnpServer.GetAccount(tCtx, func(p capnp.ZarbServer_getAccount_Params) error {
		assert.NoError(t, p.SetAddress([]byte(addr.String())))
		return nil
	}).Result()

	st, err := res.Struct()
	if err != nil {
		return nil
	}

	d, _ := st.Data()
	acc, err := account.FromBytes(d)
	assert.NoError(t, err)
	return acc
}

func TestGetAccount(t *testing.T) {
	acc := getAccount(t, crypto.TreasuryAddress)
	require.NotNil(t, acc)
	assert.Equal(t, acc.Number(), int32(0))
}
