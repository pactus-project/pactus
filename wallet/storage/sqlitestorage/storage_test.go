package sqlitestorage

import (
	"testing"

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
		PublicKey: td.RandString(32),
		Label:     td.RandString(16),
		Path:      td.RandString(16),
	}
}

func (td *testData) InsertRandomAddressInfo(t *testing.T) *types.AddressInfo {
	t.Helper()

	addrInfo := td.RandomAddressInfo(t)
	err := td.storage.InsertAddress(addrInfo)
	require.NoError(t, err)

	return addrInfo
}

func (td *testData) RandomDirection(t *testing.T) types.TxDirection {
	t.Helper()

	return types.TxDirection(td.RandInt(2) + 1)
}

func (td *testData) RandomTransactionInfo(t *testing.T, direction types.TxDirection,
	opts ...testsuite.TransactionMakerOption,
) *types.TransactionInfo {
	t.Helper()

	trx := td.GenerateTestTransferTx(opts...)
	txInfos, err := types.MakeTransactionInfos(trx, types.TransactionStatusPending, 0)
	require.NoError(t, err)

	txInfos[0].Direction = direction

	return txInfos[0]
}

func (td *testData) InsertRandomTransactionInfo(t *testing.T, direction types.TxDirection,
	opts ...testsuite.TransactionMakerOption,
) *types.TransactionInfo {
	t.Helper()

	txInfo := td.RandomTransactionInfo(t, direction, opts...)
	err := td.storage.InsertTransaction(txInfo)
	require.NoError(t, err)

	return txInfo
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

	addrInfo1 := td.InsertRandomAddressInfo(t)
	addrInfo2 := td.InsertRandomAddressInfo(t)

	// Test AllAddresses
	addresses := td.storage.AllAddresses()
	require.Len(t, addresses, 2)

	assert.Equal(t, 2, td.storage.AddressCount())

	assert.True(t, td.storage.HasAddress(addrInfo1.Address))
	assert.True(t, td.storage.HasAddress(addrInfo2.Address))

	// Test GetAddress
	retrieved1, err := td.storage.AddressInfo(addrInfo1.Address)
	require.NoError(t, err)
	assert.Equal(t, addrInfo1.Address, retrieved1.Address)
	assert.Equal(t, addrInfo1.PublicKey, retrieved1.PublicKey)
	assert.Equal(t, addrInfo1.Path, retrieved1.Path)
	assert.Equal(t, addrInfo1.Label, retrieved1.Label)

	// Test UpdateAddress
	addrInfo2.Label = "Updated Label"
	err = td.storage.UpdateAddress(addrInfo2)
	require.NoError(t, err)

	retrieved2, err := td.storage.AddressInfo(addrInfo2.Address)
	require.NoError(t, err)
	assert.Equal(t, "Updated Label", retrieved2.Label)

	// Test GetAddress not found
	_, err = td.storage.AddressInfo("non_existent")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestInsertAddress(t *testing.T) {
	td := setup(t)

	addrInfo1 := td.InsertRandomAddressInfo(t)
	assert.True(t, td.storage.HasAddress(addrInfo1.Address))

	addrInfo2 := td.InsertRandomAddressInfo(t)
	assert.True(t, td.storage.HasAddress(addrInfo2.Address))

	assert.Equal(t, 2, td.storage.AddressCount())
}

func TestTransactionOperations(t *testing.T) {
	td := setup(t)

	txInfo := td.InsertRandomTransactionInfo(t, td.RandomDirection(t))

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

func TestQueryTransactions_Range(t *testing.T) {
	td := setup(t)

	// Insert multiple transactions
	receiver := td.RandAccAddress()
	for i := 0; i < 5; i++ {
		td.InsertRandomTransactionInfo(t, td.RandomDirection(t), testsuite.TransactionWithReceiver(receiver))
	}

	// Test QueryTransactions with pagination
	transactions, err := td.storage.QueryTransactions(storage.QueryParams{
		Address:   receiver.String(),
		Direction: types.TxDirectionAny,
		Count:     3,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 3)

	// Test with skip
	transactions, err = td.storage.QueryTransactions(storage.QueryParams{
		Address:   receiver.String(),
		Direction: types.TxDirectionAny,
		Count:     3,
		Skip:      3,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 2)

	// Test with different receiver
	transactions, err = td.storage.QueryTransactions(storage.QueryParams{
		Address:   "other_receiver",
		Direction: types.TxDirectionAny,
		Count:     10,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 0)
}

func TestQueryTransactions_Direction(t *testing.T) {
	td := setup(t)

	// create three tx rows:
	// 1) sender A -> receiver B
	// 2) sender A -> receiver C
	// 3) sender D -> receiver B
	pubA, signerA := td.RandEd25519KeyPair()
	pubD, signerD := td.RandEd25519KeyPair()
	addrB := td.RandAccAddress()
	addrC := td.RandAccAddress()

	td.InsertRandomTransactionInfo(t, types.TxDirectionOutgoing,
		testsuite.TransactionWithSigner(signerA), testsuite.TransactionWithReceiver(addrB))
	td.InsertRandomTransactionInfo(t, types.TxDirectionOutgoing,
		testsuite.TransactionWithSigner(signerA), testsuite.TransactionWithReceiver(addrC))
	td.InsertRandomTransactionInfo(t, types.TxDirectionIncoming,
		testsuite.TransactionWithSigner(signerD), testsuite.TransactionWithReceiver(addrB))

	// Wildcard both: all 3
	txs, err := td.storage.QueryTransactions(storage.QueryParams{
		Address:   "*",
		Direction: types.TxDirectionAny,
		Count:     10,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 3)

	// Filter sender only (A): rows 1 and 2
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Address:   pubA.AccountAddress().String(),
		Direction: types.TxDirectionOutgoing,
		Count:     10,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 2)

	// Filter sender only (D): rows 3
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Address:   pubD.AccountAddress().String(),
		Direction: types.TxDirectionOutgoing,
		Count:     10,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 0)
	// Filter receiver only (B): rows with incoming to B (row 3)
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Address:   addrB.String(),
		Direction: types.TxDirectionIncoming,
		Count:     10,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 1)

	// Outgoing from A: rows 1 and 2
	txs, err = td.storage.QueryTransactions(storage.QueryParams{
		Address:   pubA.AccountAddress().String(),
		Direction: types.TxDirectionOutgoing,
		Count:     10,
		Skip:      0,
	})
	require.NoError(t, err)
	assert.Len(t, txs, 2)
}

func TestGetPendingTransactions(t *testing.T) {
	td := setup(t)

	// Insert multiple transactions
	txInfo1 := td.InsertRandomTransactionInfo(t, types.TxDirectionOutgoing)
	txInfo2 := td.InsertRandomTransactionInfo(t, types.TxDirectionOutgoing)
	txInfo3 := td.InsertRandomTransactionInfo(t, types.TxDirectionIncoming)

	pendings, err := td.storage.GetPendingTransactions()
	require.NoError(t, err)
	require.Len(t, pendings, 3)

	_ = td.storage.UpdateTransactionStatus(txInfo1.ID, types.TransactionStatusConfirmed, td.RandHeight())
	_ = td.storage.UpdateTransactionStatus(txInfo2.ID, types.TransactionStatusFailed, 0)

	pendings, err = td.storage.GetPendingTransactions()
	require.NoError(t, err)
	require.Len(t, pendings, 1)

	assert.Contains(t, pendings, txInfo3.ID)
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
	assert.Equal(t, originalInfo.DefaultFee, clonedInfo.DefaultFee)
	assert.Equal(t, originalInfo.Encrypted, clonedInfo.Encrypted)
	assert.Equal(t, originalInfo.Driver, clonedInfo.Driver)
	assert.Equal(t, td.storage.Vault(), cloned.Vault())
	assert.Equal(t, td.storage.AllAddresses(), cloned.AllAddresses())
}
