package offline_test

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/provider/offline"
	"github.com/stretchr/testify/require"
)

// An offline wallet answers every question about the blockchain with ErrOffline.
func TestOfflineProvider(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	prov := offline.NewOfflineBlockchainProvider()

	_, err := prov.LastBlockHeight()
	require.ErrorIs(t, err, offline.ErrOffline)

	_, err = prov.GetAccount(ts.RandAccAddress().String())
	require.ErrorIs(t, err, offline.ErrOffline)

	_, err = prov.GetValidator(ts.RandValAddress().String())
	require.ErrorIs(t, err, offline.ErrOffline)

	_, _, err = prov.GetTransaction(ts.RandHash().String())
	require.ErrorIs(t, err, offline.ErrOffline)

	_, err = prov.SendTx(ts.GenerateTestTransferTx())
	require.ErrorIs(t, err, offline.ErrOffline)

	require.ErrorIs(t, prov.CheckTransaction(nil), offline.ErrOffline)
	require.NoError(t, prov.Close())
}
