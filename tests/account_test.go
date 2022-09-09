package tests

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/www/capnp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAccount(t *testing.T, addr crypto.Address) *account.Account {
	res := tCapnpServer.GetAccount(tCtx, func(p capnp.PactusServer_getAccount_Params) error {
		assert.NoError(t, p.SetAddress(addr.String()))
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
