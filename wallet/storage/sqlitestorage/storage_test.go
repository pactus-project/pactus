package sqlitestorage

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:dupword // duplicated seed phrase words
var testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon cactus"

type testData struct {
	*testsuite.TestSuite

	storage *Storage
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	vlt, err := vault.CreateVaultFromMnemonic(testMnemonic, addresspath.CoinTypePactusTestnet)
	require.NoError(t, err)

	path := util.TempDirPath()
	strg, err := Create(t.Context(), path, genesis.Testnet, vlt)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = strg.Close()
	})

	return &testData{
		TestSuite: ts,
		storage:   strg,
	}
}

func (td *testData) RandomAddressInfo(t *testing.T) *types.AddressInfo {
	t.Helper()

	return &types.AddressInfo{
		Address:   td.RandAccAddress().String(),
		PublicKey: td.RandString(16),
		Label:     td.RandString(16),
		Path:      td.RandString(16),
	}
}

func (td *testData) RandomTransactionInfo(t *testing.T) *types.TransactionInfo {
	t.Helper()

	trx := td.GenerateTestTransferTx()
	txInfos, err := types.MakeTransactionInfos(trx, types.TransactionStatusPending, 0)
	require.NoError(t, err)

	return txInfos[0]
}

func TestWalletInfo(t *testing.T) {
	td := setup(t)

	info := td.storage.WalletInfo()
	assert.Equal(t, VersionLatest, info.Version)
	assert.Equal(t, genesis.Testnet, info.Network)
	assert.Equal(t, "SQLite", info.Driver)
	assert.NotEmpty(t, info.UUID)
	assert.Equal(t, amount.Amount(10_000_000), info.DefaultFee)
	assert.False(t, info.CreatedAt.IsZero())
}

func TestVault(t *testing.T) {
	td := setup(t)

	vlt := td.storage.Vault()
	require.NotNil(t, vlt)
	assert.False(t, vlt.IsEncrypted())
	assert.False(t, vlt.IsNeutered())
}

func TestUpdateVault(t *testing.T) {
	td := setup(t)

	vlt1 := td.storage.Vault()
	require.NotNil(t, vlt1)

	// Create a new address to modify the vault
	_, err := vlt1.NewValidatorAddress("Test Validator")
	require.NoError(t, err)

	// Update vault
	err = td.storage.UpdateVault(vlt1)
	require.NoError(t, err)

	// Verify update
	vlt2 := td.storage.Vault()
	assert.Equal(t, uint32(1), vlt2.Purposes.PurposeBLS.NextValidatorIndex)
}

func TestSetDefaultFee(t *testing.T) {
	td := setup(t)

	newFee := amount.Amount(20_000_000)
	err := td.storage.SetDefaultFee(newFee)
	require.NoError(t, err)

	info := td.storage.WalletInfo()
	assert.Equal(t, newFee, info.DefaultFee)
}

func TestAddressOperations(t *testing.T) {
	td := setup(t)

	// Test InsertAddress
	addrInfo := td.RandomAddressInfo(t)

	err := td.storage.InsertAddress(addrInfo)
	require.NoError(t, err)

	// Test AllAddresses
	addresses := td.storage.AllAddresses()
	require.Len(t, addresses, 1)
	assert.Equal(t, addrInfo.Address, addresses[0].Address)
	assert.Equal(t, addrInfo.PublicKey, addresses[0].PublicKey)
	assert.Equal(t, addrInfo.Label, addresses[0].Label)
	assert.Equal(t, addrInfo.Path, addresses[0].Path)

	// Test UpdateAddress
	addrInfo.Label = "Updated Label"
	err = td.storage.UpdateAddress(addrInfo)
	require.NoError(t, err)

	addresses = td.storage.AllAddresses()
	assert.Equal(t, "Updated Label", addresses[0].Label)
}

func TestTransactionOperations(t *testing.T) {
	td := setup(t)

	txInfo := td.RandomTransactionInfo(t)
	err := td.storage.InsertTransaction(txInfo)
	require.NoError(t, err)

	// Test HasTransaction
	assert.True(t, td.storage.HasTransaction(txInfo.ID))
	assert.False(t, td.storage.HasTransaction("non_existing"))

	// Test GetTransaction
	retrieved, err := td.storage.GetTransaction(txInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, txInfo.ID, retrieved.ID)
	assert.Equal(t, txInfo.Sender, retrieved.Sender)
	assert.Equal(t, txInfo.Receiver, retrieved.Receiver)
	assert.Equal(t, txInfo.Amount, retrieved.Amount)
	assert.Equal(t, txInfo.Fee, retrieved.Fee)
	assert.Equal(t, txInfo.Memo, retrieved.Memo)
	assert.Equal(t, txInfo.Status, retrieved.Status)
	assert.Equal(t, txInfo.BlockHeight, retrieved.BlockHeight)
	assert.Equal(t, txInfo.PayloadType, retrieved.PayloadType)
	assert.Equal(t, txInfo.Comment, retrieved.Comment)

	// Test GetTransaction not found
	_, err = td.storage.GetTransaction("non_existent")
	assert.ErrorIs(t, err, storage.ErrNotFound)

	// Test UpdateTransactionStatus
	err = td.storage.UpdateTransactionStatus(txInfo.ID, types.TransactionStatusConfirmed, 1)
	require.NoError(t, err)

	retrieved, err = td.storage.GetTransaction(txInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, types.TransactionStatusConfirmed, retrieved.Status)
}

func TestQueryTransactions(t *testing.T) {
	td := setup(t)

	// Insert multiple transactions
	receiver := td.RandAccAddress()
	for i := 0; i < 5; i++ {
		trx := td.GenerateTestTransferTx(testsuite.TransactionWithReceiver(receiver))
		txInfo, err := types.MakeTransactionInfos(trx, types.TransactionStatusPending, 0)
		require.NoError(t, err)

		err = td.storage.InsertTransaction(txInfo[0])
		require.NoError(t, err)
	}

	// Test QueryTransactions with pagination
	transactions, err := td.storage.QueryTransactions(storage.QueryParams{
		Receiver: receiver.String(),
		Count:    3,
		Skip:     0,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 3)

	// Test with skip
	transactions, err = td.storage.QueryTransactions(storage.QueryParams{
		Receiver: receiver.String(),
		Count:    3,
		Skip:     3,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 2)

	// Test with different receiver
	transactions, err = td.storage.QueryTransactions(storage.QueryParams{
		Receiver: "other_receiver",
		Count:    10,
		Skip:     0,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 0)
}

func TestQueryTransactions_WildcardAndFilters(t *testing.T) {
	td := setup(t)

	// create three tx rows:
	// 1) sender A -> receiver B
	// 2) sender A -> receiver C
	// 3) sender D -> receiver B
	pubA, signerA := td.RandEd25519KeyPair()
	pubD, signerD := td.RandEd25519KeyPair()
	addrB := td.RandAccAddress()
	addrC := td.RandAccAddress()

	makeTx := func(signer crypto.PrivateKey, receiver crypto.Address) {
		trx := td.GenerateTestTransferTx(
			testsuite.TransactionWithSigner(signer),
			testsuite.TransactionWithReceiver(receiver),
		)
		txInfos, err := types.MakeTransactionInfos(trx, types.TransactionStatusPending, 0)
		require.NoError(t, err)
		err = td.storage.InsertTransaction(txInfos[0])
		require.NoError(t, err)
	}

	makeTx(signerA, addrB)
	makeTx(signerA, addrC)
	makeTx(signerD, addrB)

	// Wildcard both: all 3
	txs, err := td.storage.QueryTransactions(storage.QueryParams{
		Sender: "*", Receiver: "*", Count: 10, Skip: 0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 3)

	// Filter sender only (A): rows 1 and 2
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Sender: pubA.AccountAddress().String(), Receiver: "*", Count: 10, Skip: 0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 2)

	// Filter receiver only (B): rows 1 and 3
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Sender: "*", Receiver: addrB.String(), Count: 10, Skip: 0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 2)

	// AND both with matching pair (A,B): only row 1
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Sender: pubA.AccountAddress().String(), Receiver: addrB.String(), Count: 10, Skip: 0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 1)

	// AND both with non-matching pair (A,D): none
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Sender: pubA.AccountAddress().String(), Receiver: pubD.AccountAddress().String(), Count: 10, Skip: 0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 0)
}

func TestClone(t *testing.T) {
	td := setup(t)

	clonePath := util.TempDirPath()
	cloned, err := td.storage.Clone(clonePath)
	require.NoError(t, err)
	defer func() { _ = cloned.Close() }()

	// Verify cloned storage has different UUID and CreatedAt
	originalInfo := td.storage.WalletInfo()
	clonedInfo := cloned.WalletInfo()
	assert.NotEqual(t, originalInfo.UUID, clonedInfo.UUID)
	assert.NotEqual(t, originalInfo.CreatedAt, clonedInfo.CreatedAt)
	assert.Equal(t, originalInfo.Network, clonedInfo.Network)
	assert.Equal(t, originalInfo.Version, clonedInfo.Version)
}
