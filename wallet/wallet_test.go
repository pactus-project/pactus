package wallet_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/provider"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	wallet       *wallet.Wallet
	password     string
	testVault    *vault.Vault
	mockStorage  *storage.MockIStorage
	mockProvider *provider.MockIBlockchainProvider
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	mockStorage := storage.NewMockIStorage(ts.Ctrl)
	mockProvider := provider.NewMockIBlockchainProvider(ts.Ctrl)

	mnemonic1, _ := wallet.GenerateMnemonic(128)
	testVault, _ := vault.CreateVaultFromMnemonic(mnemonic1, addresspath.CoinTypePactusMainnet)
	mockStorage.EXPECT().Vault().Return(testVault).AnyTimes()

	var wlt *wallet.Wallet

	t.Cleanup(func() {
		mockProvider.EXPECT().Close().Times(1)
		mockStorage.EXPECT().Close().Times(1)
		if wlt != nil {
			wlt.Close()
		}
	})

	wlt, err := wallet.New(mockStorage, wallet.WithBlockchainProvider(mockProvider))
	assert.NoError(t, err)

	return &testData{
		TestSuite:    ts,
		testVault:    testVault,
		mockStorage:  mockStorage,
		mockProvider: mockProvider,
		wallet:       wlt,
		password:     "",
	}
}

func (td *testData) RandMemo() string {
	return td.RandString(32)
}

func TestCheckMnemonic(t *testing.T) {
	for _, entropy := range []int{128, 160, 192, 224, 256} {
		mnemonic, _ := wallet.GenerateMnemonic(entropy)
		assert.NoError(t, wallet.CheckMnemonic(mnemonic))
	}
}

func TestOpenWallet(t *testing.T) {
	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := wallet.Open(t.Context(), util.TempFilePath())
		assert.Error(t, err)
	})

	t.Run("Invalid data", func(t *testing.T) {
		path := util.TempFilePath()
		assert.NoError(t, util.WriteFile(path, []byte("invalid_data")))

		_, err := wallet.Open(t.Context(), path)
		assert.Error(t, err)
	})
}

func TestCreateWallet(t *testing.T) {
	mnemonic, _ := wallet.GenerateMnemonic(256)
	password := ""
	t.Run("Wallet exists", func(t *testing.T) {
		path := util.TempFilePath()
		err := util.WriteFile(path, []byte("something-here"))
		require.NoError(t, err)

		_, err = wallet.Create(t.Context(), path, mnemonic, password, genesis.Mainnet)
		assert.Error(t, err, wallet.ExitsError{Path: path})
	})

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := wallet.Create(t.Context(), util.TempFilePath(), "invalid mnemonic", password, genesis.Mainnet)
		assert.Error(t, err)
	})

	t.Run("Invalid path", func(t *testing.T) {
		_, err := wallet.Create(t.Context(), "\x00", mnemonic, password, genesis.Mainnet)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		walletPath := util.TempFilePath()
		_, err := wallet.Create(t.Context(), walletPath, mnemonic, password, genesis.Mainnet)
		assert.NoError(t, err)
	})
}

func TestSetDefaultFee(t *testing.T) {
	td := setup(t)

	fee := td.RandFee()
	td.mockStorage.EXPECT().SetDefaultFee(fee).Return(nil)

	err := td.wallet.SetDefaultFee(fee)
	require.NoError(t, err)
}

func TestMnemonic(t *testing.T) {
	td := setup(t)

	mnemonic1, err := td.testVault.Mnemonic(td.password)
	require.NoError(t, err)

	mnemonic2, err := td.wallet.Mnemonic(td.password)
	require.NoError(t, err)
	assert.Equal(t, mnemonic1, mnemonic2)
}

func TestSignMessage(t *testing.T) {
	td := setup(t)

	msg := "pactus"
	expectedSig := "8c3ba687e8e4c016293a2c369493faa565065987544a59baba7aadae3f17ada07883552b6c7d1d7eb49f46fbdf0975c4"
	prv, err := bls.PrivateKeyFromString("SECRET1P9QAUKRJAU7SQ7AT6ZZ6HXHYLMKPQSQYTGDL2VMH5Q5N0P5Q2QW0QL45AY3")
	require.NoError(t, err)

	_, accInfo, err := td.testVault.ImportBLSPrivateKey(td.password, prv)
	assert.NoError(t, err)

	t.Run("Unexpected Error", func(t *testing.T) {
		td.mockStorage.EXPECT().AddressInfo(accInfo.Address).Return(nil, errors.New("unexpected error"))

		_, err := td.wallet.SignMessage(td.password, "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf", msg)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		td.mockStorage.EXPECT().AddressInfo(accInfo.Address).Return(accInfo, nil)

		sig, err := td.wallet.SignMessage(td.password, "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf", msg)
		assert.NoError(t, err)
		assert.Equal(t, expectedSig, sig)
	})
}

func TestBalance(t *testing.T) {
	td := setup(t)

	t.Run("existing account", func(t *testing.T) {
		acc, addr := td.GenerateTestAccount()
		td.mockProvider.EXPECT().GetAccount(addr.String()).Return(acc, nil)

		amt, err := td.wallet.Balance(addr.String())
		assert.NoError(t, err)
		assert.Equal(t, amt, acc.Balance())
	})

	t.Run("non-existing account", func(t *testing.T) {
		addr := td.RandAccAddress()
		td.mockProvider.EXPECT().GetAccount(addr.String()).Return(nil, errors.New("account not found"))

		amt, err := td.wallet.Balance(addr.String())
		assert.Error(t, err)
		assert.Zero(t, amt)
	})
}

func TestStake(t *testing.T) {
	td := setup(t)

	t.Run("existing validator", func(t *testing.T) {
		val := td.GenerateTestValidator()
		td.mockProvider.EXPECT().GetValidator(val.Address().String()).Return(val, nil)

		amt, err := td.wallet.Stake(val.Address().String())
		assert.NoError(t, err)
		assert.Equal(t, amt, val.Stake())
	})

	t.Run("non-existing validator", func(t *testing.T) {
		addr := td.RandValAddress()
		td.mockProvider.EXPECT().GetValidator(addr.String()).Return(nil, errors.New("validator not found"))

		amt, err := td.wallet.Stake(addr.String())
		assert.Error(t, err)
		assert.Zero(t, amt)
	})
}

func TestSigningTxWithBLS(t *testing.T) {
	td := setup(t)

	senderInfo, err := td.testVault.NewBLSAccountAddress("test")
	require.NoError(t, err)
	receiver := td.RandAccAddress()
	amt := td.RandAmount()
	fee := td.RandFee()
	lockTime := td.RandHeight()
	memo := td.RandMemo()

	opts := []wallet.TxOption{
		wallet.OptionFee(fee.String()),
		wallet.OptionLockTime(lockTime),
		wallet.OptionMemo(memo),
	}

	td.mockStorage.EXPECT().WalletInfo().Return(&types.WalletInfo{DefaultFee: td.RandFee()})
	td.mockStorage.EXPECT().AddressInfo(senderInfo.Address).Return(senderInfo, nil)

	trx, err := td.wallet.MakeTransferTx(senderInfo.Address, receiver.String(), amt, opts...)
	assert.NoError(t, err)
	err = td.wallet.SignTransaction(td.password, trx)
	assert.NoError(t, err)
	assert.NotNil(t, trx.Signature())
	assert.NoError(t, trx.BasicCheck())

	td.mockProvider.EXPECT().SendTx(trx).Return(trx.ID().String(), nil)
	td.mockStorage.EXPECT().InsertTransaction(gomock.Any()).Return(nil)

	id, err := td.wallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)
	assert.Equal(t, fee, trx.Fee())
	assert.Equal(t, lockTime, trx.LockTime())
	assert.Equal(t, memo, trx.Memo())
}

func TestSigningTxWithEd25519(t *testing.T) {
	td := setup(t)

	senderInfo, err := td.testVault.NewEd25519AccountAddress("testing addr", td.password)
	require.NoError(t, err)
	receiver := td.RandAccAddress()
	amt := td.RandAmount()
	fee := td.RandFee()
	lockTime := td.RandHeight()
	memo := td.RandMemo()

	opts := []wallet.TxOption{
		wallet.OptionFee(fee.String()),
		wallet.OptionLockTime(lockTime),
		wallet.OptionMemo(memo),
	}

	td.mockStorage.EXPECT().WalletInfo().Return(&types.WalletInfo{DefaultFee: td.RandFee()})
	td.mockStorage.EXPECT().AddressInfo(senderInfo.Address).Return(senderInfo, nil)
	trx, err := td.wallet.MakeTransferTx(senderInfo.Address, receiver.String(), amt, opts...)
	assert.NoError(t, err)

	err = td.wallet.SignTransaction(td.password, trx)
	assert.NoError(t, err)
	assert.NotNil(t, trx.Signature())
	assert.NoError(t, trx.BasicCheck())

	td.mockProvider.EXPECT().SendTx(trx).Return(trx.ID().String(), nil)
	td.mockStorage.EXPECT().InsertTransaction(gomock.Any()).Return(nil)

	id, err := td.wallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)
	assert.Equal(t, fee, trx.Fee())
	assert.Equal(t, lockTime, trx.LockTime())
	assert.Equal(t, memo, trx.Memo())
}

func TestMakeTransferTx(t *testing.T) {
	td := setup(t)

	sender := td.RandAccAddress()
	receiver := td.RandAccAddress()
	amt := td.RandAmount()

	td.mockStorage.EXPECT().WalletInfo().Return(&types.WalletInfo{DefaultFee: td.RandFee()}).AnyTimes()

	t.Run("set parameters manually", func(t *testing.T) {
		fee := td.RandFee()
		lockTime := td.RandHeight()
		memo := td.RandMemo()
		opts := []wallet.TxOption{
			wallet.OptionFee(fee.String()),
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo(memo),
		}

		trx, err := td.wallet.MakeTransferTx(sender.String(), receiver.String(), amt, opts...)
		assert.NoError(t, err)
		assert.Equal(t, fee, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, memo, trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(testHeight), nil)

		trx, err := td.wallet.MakeTransferTx(sender.String(), receiver.String(), amt)
		assert.NoError(t, err)
		assert.Equal(t, testHeight+1, trx.LockTime())
		assert.Equal(t, amt, trx.Payload().Value())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeTransferTx("invalid_addr_string", receiver.String(), amt)
		assert.Error(t, err)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := td.wallet.MakeTransferTx(sender.String(), "invalid_addr_string", amt)
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(0), errors.New("not found"))

		_, err := td.wallet.MakeTransferTx(td.RandAccAddress().String(), receiver.String(), amt)
		assert.Error(t, err)
	})
}

func TestMakeBondTx(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().WalletInfo().Return(&types.WalletInfo{DefaultFee: td.RandFee()}).AnyTimes()

	sender := td.RandAccAddress()
	amt := td.RandAmount()

	t.Run("set parameters manually", func(t *testing.T) {
		receiver := td.RandValKey()

		lockTime := td.RandHeight()
		fee := td.RandFee()
		memo := td.RandMemo()
		opts := []wallet.TxOption{
			wallet.OptionFee(fee.String()),
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo(memo),
		}
		td.mockProvider.EXPECT().GetValidator(receiver.Address().String()).Return(nil, nil)

		trx, err := td.wallet.MakeBondTx(sender.String(), receiver.Address().String(),
			receiver.PublicKey().String(), amt, opts...)
		assert.NoError(t, err)
		assert.Equal(t, fee, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, memo, trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		receiver := td.RandValKey()

		testHeight := td.RandHeight()
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(testHeight), nil)
		td.mockProvider.EXPECT().GetValidator(receiver.Address().String()).Return(nil, nil)

		trx, err := td.wallet.MakeBondTx(sender.String(), receiver.Address().String(), receiver.PublicKey().String(), amt)
		assert.NoError(t, err)
		assert.Equal(t, testHeight+1, trx.LockTime())
		assert.Equal(t, amt, trx.Payload().Value())
	})

	t.Run("validator address is not stored in wallet", func(t *testing.T) {
		receiver := td.RandValKey()
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(td.RandHeight()), nil).Times(4)
		td.mockStorage.EXPECT().AddressInfo(receiver.Address().String()).Return(nil, storage.ErrNotFound).AnyTimes()

		t.Run("validator doesn't exist and public key not set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiver.Address().String()).Return(nil, errors.New("not exist"))

			trx, err := td.wallet.MakeBondTx(sender.String(), receiver.Address().String(), "", amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator doesn't exist and public key set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiver.Address().String()).Return(nil, errors.New("not exist"))

			trx, err := td.wallet.MakeBondTx(sender.String(), receiver.Address().String(), receiver.PublicKey().String(), amt)
			require.NoError(t, err)
			assert.Equal(t, receiver.PublicKey().String(), trx.Payload().(*payload.BondPayload).PublicKey.String())
		})

		t.Run("validator exists and public key not set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiver.Address().String()).Return(td.GenerateTestValidator(), nil)

			trx, err := td.wallet.MakeBondTx(sender.String(), receiver.Address().String(), "", amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator exists and public key set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiver.Address().String()).Return(td.GenerateTestValidator(), nil)

			trx, err := td.wallet.MakeBondTx(sender.String(),
				receiver.Address().String(), receiver.PublicKey().String(), amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})
	})

	t.Run("validator address stored in wallet", func(t *testing.T) {
		td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil)
		td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)
		receiverInfo, err := td.wallet.NewValidatorAddress("validator-address")
		require.NoError(t, err)

		td.mockStorage.EXPECT().AddressInfo(receiverInfo.Address).Return(receiverInfo, nil).AnyTimes()
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(td.RandHeight()), nil).Times(4)

		t.Run("validator doesn't exist and public key not set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiverInfo.Address).Return(nil, errors.New("not exist"))

			trx, err := td.wallet.MakeBondTx(sender.String(), receiverInfo.Address, "", amt)
			assert.NoError(t, err)
			assert.Equal(t, receiverInfo.PublicKey, trx.Payload().(*payload.BondPayload).PublicKey.String())
		})

		t.Run("validator doesn't exist and public key set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiverInfo.Address).Return(nil, errors.New("not exist"))

			trx, err := td.wallet.MakeBondTx(sender.String(), receiverInfo.Address, receiverInfo.PublicKey, amt)
			assert.NoError(t, err)
			assert.Equal(t, receiverInfo.PublicKey, trx.Payload().(*payload.BondPayload).PublicKey.String())
		})

		t.Run("validator exists and public key not set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiverInfo.Address).Return(td.GenerateTestValidator(), nil)

			trx, err := td.wallet.MakeBondTx(sender.String(),
				receiverInfo.Address, "", amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator exists and public key set", func(t *testing.T) {
			td.mockProvider.EXPECT().GetValidator(receiverInfo.Address).Return(td.GenerateTestValidator(), nil)

			trx, err := td.wallet.MakeBondTx(sender.String(),
				receiverInfo.Address, receiverInfo.PublicKey, amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx("invalid_addr_string", td.RandValAddress().String(), "", amt)
		assert.Error(t, err)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx(sender.String(), "invalid_addr_string", "", amt)
		assert.Error(t, err)
	})

	t.Run("invalid public key", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx(sender.String(), td.RandValAddress().String(), "invalid-pub-key", amt)
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.mockStorage.EXPECT().AddressInfo(gomock.Any()).Return(nil, storage.ErrNotFound)
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(0), errors.New("unable to get height"))
		td.mockProvider.EXPECT().GetValidator(gomock.Any()).Return(nil, errors.New("unable to get validator info")).AnyTimes()

		_, err := td.wallet.MakeBondTx(td.RandAccAddress().String(), td.RandValAddress().String(), "", amt)
		assert.Error(t, err)
	})
}

func TestMakeUnbondTx(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().WalletInfo().Return(&types.WalletInfo{DefaultFee: td.RandFee()}).AnyTimes()

	sender := td.RandValAddress()

	t.Run("set parameters manually", func(t *testing.T) {
		lockTime := td.RandHeight()
		opts := []wallet.TxOption{
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo("test"),
		}

		trx, err := td.wallet.MakeUnbondTx(sender.String(), opts...)
		assert.NoError(t, err)
		assert.Zero(t, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, "test", trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(testHeight), nil)

		trx, err := td.wallet.MakeUnbondTx(sender.String())
		assert.NoError(t, err)
		assert.Equal(t, testHeight+1, trx.LockTime())
		assert.Zero(t, trx.Payload().Value())
		assert.Zero(t, trx.Fee())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeUnbondTx("invalid_addr_string")
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(0), errors.New("unable to get height"))

		_, err := td.wallet.MakeUnbondTx(td.RandAccAddress().String())
		assert.Error(t, err)
	})
}

func TestMakeWithdrawTx(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().WalletInfo().Return(&types.WalletInfo{DefaultFee: td.RandFee()}).AnyTimes()

	sender := td.RandValAddress()
	receiver := td.RandAccAddress()

	amt := td.RandAmount()

	t.Run("set parameters manually", func(t *testing.T) {
		lockTime := td.RandHeight()
		fee := td.RandFee()
		opts := []wallet.TxOption{
			wallet.OptionFee(fee.String()),
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo("test"),
		}

		trx, err := td.wallet.MakeWithdrawTx(sender.String(), receiver.String(), amt, opts...)
		assert.NoError(t, err)
		assert.Equal(t, fee, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, "test", trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(testHeight), nil)

		trx, err := td.wallet.MakeWithdrawTx(sender.String(), receiver.String(), amt)
		assert.NoError(t, err)
		assert.Equal(t, testHeight+1, trx.LockTime())
		assert.Equal(t, amt, trx.Payload().Value())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeWithdrawTx("invalid_addr_string", receiver.String(), amt)
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.mockProvider.EXPECT().LastBlockHeight().Return(block.Height(0), errors.New("unable to get height"))

		_, err := td.wallet.MakeWithdrawTx(td.RandAccAddress().String(), receiver.String(), amt)
		assert.Error(t, err)
	})
}

func TestTotalBalance(t *testing.T) {
	td := setup(t)

	accInfo1, _ := td.testVault.NewBLSAccountAddress("account-1")
	accInfo2, _ := td.testVault.NewBLSAccountAddress("account-2")
	accInfo3, _ := td.testVault.NewBLSAccountAddress("account-3")

	addr1, err := crypto.AddressFromString(accInfo1.Address)
	require.NoError(t, err)
	addr3, err := crypto.AddressFromString(accInfo3.Address)
	require.NoError(t, err)

	acc1, _ := td.GenerateTestAccount(testsuite.AccountWithAddress(addr1))
	acc2 := account.NewAccount(2)
	acc3, _ := td.GenerateTestAccount(testsuite.AccountWithAddress(addr3))

	td.mockStorage.EXPECT().AllAddresses().Return([]types.AddressInfo{*accInfo1, *accInfo2, *accInfo3})
	td.mockProvider.EXPECT().GetAccount(accInfo1.Address).Return(acc1, nil)
	td.mockProvider.EXPECT().GetAccount(accInfo2.Address).Return(acc2, nil)
	td.mockProvider.EXPECT().GetAccount(accInfo3.Address).Return(acc3, nil)

	totalBalance, err := td.wallet.TotalBalance()
	assert.NoError(t, err)
	assert.Equal(t, acc1.Balance()+acc2.Balance()+acc3.Balance(), totalBalance)
}

func TestTotalStake(t *testing.T) {
	td := setup(t)

	valInfo1, _ := td.testVault.NewValidatorAddress("val-1")
	valInfo2, _ := td.testVault.NewValidatorAddress("val-2")

	pub1, err := bls.PublicKeyFromString(valInfo1.PublicKey)
	require.NoError(t, err)
	pub2, err := bls.PublicKeyFromString(valInfo2.PublicKey)
	require.NoError(t, err)

	val1 := td.GenerateTestValidator(testsuite.ValidatorWithPublicKey(pub1))
	val2 := td.GenerateTestValidator(testsuite.ValidatorWithPublicKey(pub2))

	td.mockStorage.EXPECT().AllAddresses().Return([]types.AddressInfo{*valInfo1, *valInfo2})
	td.mockProvider.EXPECT().GetValidator(valInfo1.Address).Return(val1, nil)
	td.mockProvider.EXPECT().GetValidator(valInfo2.Address).Return(val2, nil)

	stake, err := td.wallet.TotalStake()
	require.NoError(t, err)
	require.Equal(t, val1.Stake()+val2.Stake(), stake)
}

func TestNeuter(t *testing.T) {
	td := setup(t)

	path := util.TempFilePath()
	clonedStorage := storage.NewMockIStorage(td.Ctrl)

	td.mockStorage.EXPECT().Clone(path).Return(clonedStorage, nil)
	clonedStorage.EXPECT().Vault().Return(td.testVault)
	clonedStorage.EXPECT().UpdateVault(gomock.Any()).DoAndReturn(func(vlt *vault.Vault) error {
		assert.True(t, vlt.IsNeutered())
		assert.False(t, vlt.IsEncrypted())

		return nil
	})

	err := td.wallet.Neuter(path)
	require.NoError(t, err)
}

func TestTestnetWallet(t *testing.T) {
	walletPath := util.TempFilePath()

	t.Run("Create Testnet wallet", func(t *testing.T) {
		mnemonic, _ := wallet.GenerateMnemonic(128)
		wlt, err := wallet.Create(t.Context(), walletPath, mnemonic, "", genesis.Testnet)
		require.NoError(t, err)

		assert.Equal(t, genesis.Testnet, wlt.Info().Network)

		info, err := wlt.NewBLSAccountAddress("testnet-addr-1")
		require.NoError(t, err)
		assert.Equal(t, "m/12381'/21777'/2'/0", info.Path)
		assert.True(t, strings.HasPrefix(info.Address, "tpc1"))
	})

	t.Run("Open Testnet wallet", func(t *testing.T) {
		wlt, err := wallet.Open(t.Context(), walletPath)
		require.NoError(t, err)

		assert.Equal(t, genesis.Testnet, wlt.Info().Network)

		info, err := wlt.NewBLSAccountAddress("testnet-addr-2")
		require.NoError(t, err)
		assert.Equal(t, "m/12381'/21777'/2'/1", info.Path)
		assert.True(t, strings.HasPrefix(info.Address, "tpc1"))
	})
}
