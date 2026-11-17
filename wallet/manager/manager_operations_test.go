package manager_test

import (
	"testing"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/wallet/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testWallet   = "main"
	testPassword = "password"
)

type managerTest struct {
	*testsuite.TestSuite

	conf     *manager.Config
	provider *provider.MockWalletProvider
	mgr      manager.WalletManager
}

// setupManager creates a manager with one wallet. The provider answers that each account
// has a balance of 10 PAC and each validator a stake of 20 PAC.
func setupManager(t *testing.T) *managerTest {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	conf := &manager.Config{
		WalletsDir: util.TempDirPath(),
		ChainType:  genesis.Mainnet,
	}

	prov := provider.NewMockWalletProvider(ts.MockController())
	prov.EXPECT().LastBlockHeight().Return(types.Height(100), nil).AnyTimes()
	prov.EXPECT().GetAccount(gomock.Any()).DoAndReturn(func(string) (*account.Account, error) {
		acc := account.NewAccount(1)
		acc.AddToBalance(10e9)

		return acc, nil
	}).AnyTimes()
	prov.EXPECT().GetValidator(gomock.Any()).DoAndReturn(func(string) (*validator.Validator, error) {
		pub, _ := ts.RandBLSKeyPair()
		val := validator.NewValidator(pub, 1)
		val.AddToStake(20e9)

		return val, nil
	}).AnyTimes()
	prov.EXPECT().Close().Return(nil).AnyTimes()

	mgr, err := manager.NewManager(t.Context(), conf, prov, pipeline.New[any](t.Context()))
	require.NoError(t, err)
	t.Cleanup(mgr.Stop)

	_, err = mgr.CreateWallet(testWallet, testPassword)
	require.NoError(t, err)

	return &managerTest{TestSuite: ts, conf: conf, provider: prov, mgr: mgr}
}

func TestManagerStartsAndStops(t *testing.T) {
	suite := setupManager(t)

	require.NoError(t, suite.mgr.Start())
	suite.mgr.Stop()
}

// Each operation on a wallet that is not loaded returns ErrWalletNotLoaded.
func TestManagerUnknownWallet(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr
	addr := suite.RandAccAddress().String()
	unknown := "unknown"

	_, err := mgr.NewAddress(unknown, crypto.AddressTypeBLSAccount, "label")
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.PrivateKey(unknown, testPassword, addr)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.Mnemonic(unknown, testPassword)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.AddressLabel(unknown, addr)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	require.ErrorIs(t, mgr.SetAddressLabel(unknown, addr, "label"), manager.ErrWalletNotLoaded)
	_, err = mgr.Balance(unknown, addr)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.Stake(unknown, addr)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	require.ErrorIs(t, mgr.SetDefaultFee(unknown, 1), manager.ErrWalletNotLoaded)
	_, err = mgr.MakeTransferTx(unknown, addr, addr, 1)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.MakeBondTx(unknown, addr, addr, "", 1)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.MakeUnbondTx(unknown, addr)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.MakeWithdrawTx(unknown, addr, addr, 1)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	require.ErrorIs(t, mgr.SignTransaction(unknown, testPassword, suite.GenerateTestTransferTx()),
		manager.ErrWalletNotLoaded)
	_, err = mgr.BroadcastTransaction(unknown, suite.GenerateTestTransferTx())
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.TotalBalance(unknown)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.TotalStake(unknown)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, _, err = mgr.SignRawTransaction(unknown, testPassword, nil)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.SignMessage(unknown, testPassword, addr, "msg")
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.AddressInfo(unknown, addr)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.WalletInfo(unknown)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.ListAddresses(unknown)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	_, err = mgr.ListTransactions(unknown)
	require.ErrorIs(t, err, manager.ErrWalletNotLoaded)
	require.ErrorIs(t, mgr.UpdatePassword(unknown, testPassword, "new"), manager.ErrWalletNotLoaded)
}

func TestManagerCreateAndRestore(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr

	_, err := mgr.CreateWallet(testWallet, testPassword)
	require.ErrorIs(t, err, manager.ErrWalletAlreadyExists)

	mnemonic, err := mgr.Mnemonic(testWallet, testPassword)
	require.NoError(t, err)
	_, err = mgr.Mnemonic(testWallet, "wrong")
	require.Error(t, err)

	// A wallet that is restored from the same seed phrase derives the same addresses.
	require.NoError(t, mgr.RestoreWallet("restored", mnemonic, testPassword))
	require.ErrorIs(t, mgr.RestoreWallet("restored", mnemonic, testPassword), manager.ErrWalletAlreadyExists)
	require.Error(t, mgr.RestoreWallet("invalid", "not a seed phrase", testPassword))

	first, err := mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "first")
	require.NoError(t, err)
	restored, err := mgr.NewAddress("restored", crypto.AddressTypeBLSAccount, "first")
	require.NoError(t, err)
	assert.Equal(t, first.Address, restored.Address)

	wallets, err := mgr.ListWallets()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{testWallet, "restored"}, wallets)

	info, err := mgr.WalletInfo(testWallet)
	require.NoError(t, err)
	assert.Equal(t, genesis.Mainnet, info.Network)
	assert.True(t, info.Encrypted)
}

// A new manager on the same directory loads the wallets that were created.
func TestManagerLoadsTheWalletsOfTheDirectory(t *testing.T) {
	suite := setupManager(t)

	created, err := suite.mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "saved")
	require.NoError(t, err)
	suite.mgr.Stop()

	mgr, err := manager.NewManager(t.Context(), suite.conf, suite.provider, pipeline.New[any](t.Context()))
	require.NoError(t, err)
	t.Cleanup(mgr.Stop)

	wallets, err := mgr.ListWallets()
	require.NoError(t, err)
	assert.Equal(t, []string{testWallet}, wallets)

	info, err := mgr.AddressInfo(testWallet, created.Address)
	require.NoError(t, err)
	assert.Equal(t, "saved", info.Label)
}

func TestManagerAddresses(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr

	accInfo, err := mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "account")
	require.NoError(t, err)
	valInfo, err := mgr.NewAddress(testWallet, crypto.AddressTypeValidator, "validator")
	require.NoError(t, err)
	edInfo, err := mgr.NewAddress(testWallet, crypto.AddressTypeEd25519Account, "ed25519",
		wallet.WithPassword(testPassword))
	require.NoError(t, err)

	label, err := mgr.AddressLabel(testWallet, accInfo.Address)
	require.NoError(t, err)
	assert.Equal(t, "account", label)

	require.NoError(t, mgr.SetAddressLabel(testWallet, accInfo.Address, "renamed"))
	info, err := mgr.AddressInfo(testWallet, accInfo.Address)
	require.NoError(t, err)
	assert.Equal(t, "renamed", info.Label)
	assert.Equal(t, accInfo.PublicKey, info.PublicKey)

	_, err = mgr.AddressInfo(testWallet, suite.RandAccAddress().String())
	require.Error(t, err)

	all, err := mgr.ListAddresses(testWallet)
	require.NoError(t, err)
	assert.Len(t, all, 3)

	validators, err := mgr.ListAddresses(testWallet, wallet.OnlyValidatorAddresses())
	require.NoError(t, err)
	require.Len(t, validators, 1)
	assert.Equal(t, valInfo.Address, validators[0].Address)

	// The validator address of a public key doesn't need a wallet.
	valAddr, err := mgr.GetValidatorAddress(valInfo.PublicKey)
	require.NoError(t, err)
	assert.Equal(t, valInfo.Address, valAddr)
	_, err = mgr.GetValidatorAddress("invalid")
	require.Error(t, err)

	// The private key of an address matches its public key.
	prv, err := mgr.PrivateKey(testWallet, testPassword, edInfo.Address)
	require.NoError(t, err)
	assert.Equal(t, edInfo.PublicKey, prv.PublicKey().String())
	_, err = mgr.PrivateKey(testWallet, "wrong", edInfo.Address)
	require.Error(t, err)
}

func TestManagerBalances(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr

	accInfo, err := mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "account")
	require.NoError(t, err)
	_, err = mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "account 2")
	require.NoError(t, err)
	valInfo, err := mgr.NewAddress(testWallet, crypto.AddressTypeValidator, "validator")
	require.NoError(t, err)

	balance, err := mgr.Balance(testWallet, accInfo.Address)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(10e9), balance)

	stake, err := mgr.Stake(testWallet, valInfo.Address)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(20e9), stake)

	total, err := mgr.TotalBalance(testWallet)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(20e9), total, "two accounts of 10 PAC")

	totalStake, err := mgr.TotalStake(testWallet)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(20e9), totalStake, "one validator of 20 PAC")
}

func TestManagerTransactions(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr

	sender, err := mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "sender")
	require.NoError(t, err)
	valInfo, err := mgr.NewAddress(testWallet, crypto.AddressTypeValidator, "validator")
	require.NoError(t, err)
	receiver := suite.RandAccAddress().String()

	require.NoError(t, mgr.SetDefaultFee(testWallet, 2e7))
	info, err := mgr.WalletInfo(testWallet)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(2e7), info.DefaultFee)

	transfer, err := mgr.MakeTransferTx(testWallet, sender.Address, receiver, 1e9)
	require.NoError(t, err)
	assert.True(t, transfer.IsTransferTx())
	assert.Equal(t, amount.Amount(2e7), transfer.Fee(), "the default fee of the wallet")
	assert.Equal(t, sender.Address, transfer.Payload().Signer().String())

	withFee, err := mgr.MakeTransferTx(testWallet, sender.Address, receiver, 1e9,
		wallet.OptionFee("0.5"), wallet.OptionMemo("memo"), wallet.OptionLockTime(50))
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(5e8), withFee.Fee())
	assert.Equal(t, "memo", withFee.Memo())
	assert.Equal(t, types.Height(50), withFee.LockTime())

	bond, err := mgr.MakeBondTx(testWallet, sender.Address, valInfo.Address, valInfo.PublicKey, 2e9)
	require.NoError(t, err)
	assert.True(t, bond.IsBondTx())

	unbond, err := mgr.MakeUnbondTx(testWallet, valInfo.Address)
	require.NoError(t, err)
	assert.True(t, unbond.IsUnbondTx())

	withdraw, err := mgr.MakeWithdrawTx(testWallet, valInfo.Address, receiver, 1e9)
	require.NoError(t, err)
	assert.True(t, withdraw.IsWithdrawTx())

	_, err = mgr.MakeTransferTx(testWallet, "invalid", receiver, 1e9)
	require.Error(t, err)

	// Signing needs the password, and the signature is valid.
	require.Error(t, mgr.SignTransaction(testWallet, "wrong", transfer))
	require.NoError(t, mgr.SignTransaction(testWallet, testPassword, transfer))
	require.NoError(t, transfer.BasicCheck())

	// A raw transaction is signed to the same transaction.
	rawTx, err := withFee.Bytes()
	require.NoError(t, err)
	txID, signedData, err := mgr.SignRawTransaction(testWallet, testPassword, rawTx)
	require.NoError(t, err)
	signed, err := tx.FromBytes(signedData)
	require.NoError(t, err)
	require.NoError(t, signed.BasicCheck())
	assert.Equal(t, withFee.ID().Bytes(), txID)
	assert.Equal(t, signed.ID().Bytes(), txID)

	_, _, err = mgr.SignRawTransaction(testWallet, testPassword, []byte("not a transaction"))
	require.Error(t, err)
	_, _, err = mgr.SignRawTransaction(testWallet, "wrong", rawTx)
	require.Error(t, err)

	// A broadcast transaction is recorded in the history of the wallet.
	suite.provider.EXPECT().SendTx(transfer).Return(transfer.ID().String(), nil)
	hash, err := mgr.BroadcastTransaction(testWallet, transfer)
	require.NoError(t, err)
	assert.Equal(t, transfer.ID().String(), hash)

	history, err := mgr.ListTransactions(testWallet)
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, transfer.ID().String(), history[0].TxID)

	// A transaction that the node refuses is not recorded.
	suite.provider.EXPECT().SendTx(signed).Return("", assert.AnError)
	_, err = mgr.BroadcastTransaction(testWallet, signed)
	require.ErrorIs(t, err, assert.AnError)
	history, err = mgr.ListTransactions(testWallet)
	require.NoError(t, err)
	assert.Len(t, history, 1)
}

func TestManagerSignMessage(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr

	info, err := mgr.NewAddress(testWallet, crypto.AddressTypeBLSAccount, "signer")
	require.NoError(t, err)

	sigStr, err := mgr.SignMessage(testWallet, testPassword, info.Address, "message")
	require.NoError(t, err)

	pub, err := bls.PublicKeyFromString(info.PublicKey)
	require.NoError(t, err)
	sig, err := bls.SignatureFromString(sigStr)
	require.NoError(t, err)
	require.NoError(t, pub.Verify([]byte("message"), sig))
	require.Error(t, pub.Verify([]byte("other"), sig))

	_, err = mgr.SignMessage(testWallet, "wrong", info.Address, "message")
	require.Error(t, err)
}

func TestManagerUpdatePassword(t *testing.T) {
	suite := setupManager(t)
	mgr := suite.mgr

	require.Error(t, mgr.UpdatePassword(testWallet, "wrong", "new"))
	require.NoError(t, mgr.UpdatePassword(testWallet, testPassword, "new"))

	_, err := mgr.Mnemonic(testWallet, testPassword)
	require.Error(t, err)
	_, err = mgr.Mnemonic(testWallet, "new")
	require.NoError(t, err)
}
