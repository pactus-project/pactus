package tests

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAccount(t *testing.T, addr crypto.Address) *pactus.AccountInfo {
	t.Helper()

	res, err := tBlockchainClient.GetAccount(tCtx,
		&pactus.GetAccountRequest{Address: addr.String()})
	if err != nil {
		return nil
	}

	return res.Account
}

func TestGetAccount(t *testing.T) {
	acc := getAccount(t, crypto.TreasuryAddress)
	require.NotNil(t, acc)
	assert.Equal(t, int32(0), acc.Number)
}
