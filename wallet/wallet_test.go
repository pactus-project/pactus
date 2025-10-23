package wallet_test

import (
	"context"
	"path"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	server    *grpc.Server
	wallet    *wallet.Wallet
	mockState *state.MockState
	password  string
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	password := ""
	walletPath := util.TempFilePath()
	mnemonic, _ := wallet.GenerateMnemonic(128)

	grpcConf := &grpc.Config{
		Enable: true,
		Listen: "[::]:0",
	}
	walletMgrConf := &wallet.Config{
		WalletsDir: util.TempDirPath(),
		ChainType:  genesis.Mainnet,
	}
	mockState := state.MockingState(ts)
	gRPCServer := grpc.NewServer(context.Background(),
		grpcConf, mockState,
		nil, nil, nil,
		wallet.NewWalletManager(walletMgrConf), nil)

	assert.NoError(t, gRPCServer.StartServer())

	wlt, err := wallet.Create(walletPath, mnemonic, password, genesis.Mainnet,
		wallet.WithCustomServers([]string{gRPCServer.Address()}))
	assert.NoError(t, err)
	assert.False(t, wlt.IsEncrypted())
	assert.Equal(t, walletPath, wlt.Path())
	assert.Equal(t, path.Base(walletPath), wlt.Name())

	return &testData{
		TestSuite: ts,
		mockState: mockState,
		server:    gRPCServer,
		wallet:    wlt,
		password:  password,
	}
}

func (td *testData) Close() {
	td.server.StopServer()
	// TODO:  close client (wallet.Close() ??)
}

func TestOpenWallet(t *testing.T) {
	td := setup(t)
	defer td.Close()

	t.Run("Save the wallet", func(t *testing.T) {
		assert.NoError(t, td.wallet.Save())
	})

	t.Run("Re-open the wallet", func(t *testing.T) {
		_, err := wallet.Open(td.wallet.Path(), true)
		assert.NoError(t, err)
	})

	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := wallet.Open(util.TempFilePath(), true)
		assert.Error(t, err)
	})

	t.Run("Invalid crc", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(td.wallet.Path(), []byte("{}")))

		_, err := wallet.Open(td.wallet.Path(), true)
		assert.ErrorIs(t, err, wallet.UnsupportedVersionError{
			WalletVersion:    0,
			SupportedVersion: wallet.VersionLatest,
		})
	})

	t.Run("Invalid json", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(td.wallet.Path(), []byte("invalid_json")))

		_, err := wallet.Open(td.wallet.Path(), true)
		assert.Error(t, err)
	})
}

func TestRecoverWallet(t *testing.T) {
	td := setup(t)
	defer td.Close()

	mnemonic, _ := td.wallet.Mnemonic(td.password)
	password := ""
	t.Run("Wallet exists", func(t *testing.T) {
		// Save the test wallet first then
		// try to recover a wallet at the same place
		assert.NoError(t, td.wallet.Save())

		_, err := wallet.Create(td.wallet.Path(), mnemonic, password, 0)
		assert.ErrorIs(t, err, wallet.ExitsError{
			Path: td.wallet.Path(),
		})
	})

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := wallet.Create(util.TempFilePath(),
			"invalid mnemonic phrase seed", password, 0)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		walletPath := util.TempFilePath()
		recovered, err := wallet.Create(walletPath, mnemonic, password, genesis.Mainnet)
		assert.NoError(t, err)

		addrInfo1, err := recovered.NewBLSAccountAddress("addr-1")
		assert.NoError(t, err)

		assert.NoFileExists(t, walletPath)
		assert.NoError(t, recovered.Save())

		assert.FileExists(t, walletPath)
		assert.True(t, recovered.Contains(addrInfo1.Address))
	})
}

func TestInvalidAddress(t *testing.T) {
	td := setup(t)
	defer td.Close()

	addr := td.RandAccAddress().String()
	_, err := td.wallet.PrivateKey(td.password, addr)
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	td := setup(t)
	defer td.Close()

	_, prv := td.RandBLSKeyPair()
	assert.NoError(t, td.wallet.ImportBLSPrivateKey(td.password, prv))

	pub := prv.PublicKeyNative()
	accAddr := pub.AccountAddress().String()
	valAddr := pub.AccountAddress().String()

	assert.True(t, td.wallet.Contains(accAddr))
	assert.True(t, td.wallet.Contains(valAddr))

	accAddrInfo := td.wallet.AddressInfo(accAddr)
	valAddrInfo := td.wallet.AddressInfo(accAddr)

	assert.Equal(t, pub.String(), accAddrInfo.PublicKey)
	assert.Equal(t, pub.String(), valAddrInfo.PublicKey)
}

func TestSignMessage(t *testing.T) {
	td := setup(t)
	defer td.Close()

	msg := "pactus"
	expectedSig := "8c3ba687e8e4c016293a2c369493faa565065987544a59baba7aadae3f17ada07883552b6c7d1d7eb49f46fbdf0975c4"
	prv, err := bls.PrivateKeyFromString("SECRET1P9QAUKRJAU7SQ7AT6ZZ6HXHYLMKPQSQYTGDL2VMH5Q5N0P5Q2QW0QL45AY3")

	require.NoError(t, err)

	err = td.wallet.ImportBLSPrivateKey(td.password, prv)
	assert.NoError(t, err)

	sig, err := td.wallet.SignMessage(td.password, "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf", msg)
	assert.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestBalance(t *testing.T) {
	td := setup(t)
	defer td.Close()

	t.Run("existing account", func(t *testing.T) {
		addr, acc := td.mockState.TestStore.AddTestAccount()
		amt, err := td.wallet.Balance(addr.String())
		assert.NoError(t, err)
		assert.Equal(t, amt, acc.Balance())
	})

	t.Run("non-existing account", func(t *testing.T) {
		amt, err := td.wallet.Balance(
			td.RandAccAddress().String())
		assert.Error(t, err)
		assert.Zero(t, amt)
	})
}

func TestStake(t *testing.T) {
	td := setup(t)
	defer td.Close()

	t.Run("existing validator", func(t *testing.T) {
		val := td.mockState.TestStore.AddTestValidator()
		amt, err := td.wallet.Stake(val.Address().String())
		assert.NoError(t, err)
		assert.Equal(t, amt, val.Stake())
	})

	t.Run("non-existing validator", func(t *testing.T) {
		amt, err := td.wallet.Stake(
			td.RandValAddress().String())
		assert.Error(t, err)
		assert.Zero(t, amt)
	})
}

func TestSigningTxWithBLS(t *testing.T) {
	td := setup(t)
	defer td.Close()

	senderInfo, _ := td.wallet.NewBLSAccountAddress("testing addr")
	receiver := td.RandAccAddress()
	amt := td.RandAmount()
	fee := td.RandFee()
	lockTime := td.RandHeight()

	opts := []wallet.TxOption{
		wallet.OptionFee(fee.String()),
		wallet.OptionLockTime(lockTime),
		wallet.OptionMemo("test"),
	}

	trx, err := td.wallet.MakeTransferTx(senderInfo.Address, receiver.String(), amt, opts...)
	assert.NoError(t, err)
	err = td.wallet.SignTransaction(td.password, trx)
	assert.NoError(t, err)
	assert.NotNil(t, trx.Signature())
	assert.NoError(t, trx.BasicCheck())

	id, err := td.wallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)
	assert.Equal(t, fee, trx.Fee())
}

func TestSigningTxWithEd25519(t *testing.T) {
	td := setup(t)
	defer td.Close()

	senderInfo, _ := td.wallet.NewEd25519AccountAddress("testing addr", td.password)
	receiver := td.RandAccAddress()
	amt := td.RandAmount()
	fee := td.RandFee()
	lockTime := td.RandHeight()

	opts := []wallet.TxOption{
		wallet.OptionFee(fee.String()),
		wallet.OptionLockTime(lockTime),
		wallet.OptionMemo("test"),
	}

	trx, err := td.wallet.MakeTransferTx(senderInfo.Address, receiver.String(), amt, opts...)
	assert.NoError(t, err)
	err = td.wallet.SignTransaction(td.password, trx)
	assert.NoError(t, err)
	assert.NotNil(t, trx.Signature())
	assert.NoError(t, trx.BasicCheck())

	id, err := td.wallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)
	assert.Equal(t, fee, trx.Fee())
}

func TestMakeTransferTx(t *testing.T) {
	td := setup(t)
	defer td.Close()

	senderInfo, _ := td.wallet.NewBLSAccountAddress("testing addr")
	receiverInfo := td.RandAccAddress()
	amt := td.RandAmount()
	lockTime := td.RandHeight()

	t.Run("set parameters manually", func(t *testing.T) {
		fee := td.RandFee()
		opts := []wallet.TxOption{
			wallet.OptionFee(fee.String()),
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo("test"),
		}

		trx, err := td.wallet.MakeTransferTx(senderInfo.Address, receiverInfo.String(), amt, opts...)
		assert.NoError(t, err)
		assert.Equal(t, fee, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, "test", trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		_ = td.mockState.TestStore.AddTestBlock(testHeight)

		trx, err := td.wallet.MakeTransferTx(senderInfo.Address, receiverInfo.String(), amt)
		assert.NoError(t, err)
		assert.Equal(t, trx.LockTime(), testHeight+1)
		assert.Equal(t, amt, trx.Payload().Value())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeTransferTx("invalid_addr_string", receiverInfo.String(), amt)
		assert.Error(t, err)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := td.wallet.MakeTransferTx(senderInfo.Address, "invalid_addr_string", amt)
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.Close()

		_, err := td.wallet.MakeTransferTx(td.RandAccAddress().String(), receiverInfo.String(), amt)
		assert.Error(t, err)
	})
}

func TestMakeBondTx(t *testing.T) {
	td := setup(t)
	defer td.Close()

	senderInfo, _ := td.wallet.NewValidatorAddress("testing addr")
	receiver := td.RandValKey()
	amt := td.RandAmount()

	t.Run("set parameters manually", func(t *testing.T) {
		lockTime := td.RandHeight()
		fee := td.RandFee()
		opts := []wallet.TxOption{
			wallet.OptionFee(fee.String()),
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo("test"),
		}

		trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address().String(),
			receiver.PublicKey().String(), amt, opts...)
		assert.NoError(t, err)
		assert.Equal(t, fee, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, "test", trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		_ = td.mockState.TestStore.AddTestBlock(testHeight)

		trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address().String(), receiver.PublicKey().String(), amt)
		assert.NoError(t, err)
		assert.Equal(t, trx.LockTime(), testHeight+1)
		assert.Equal(t, amt, trx.Payload().Value())
	})

	t.Run("validator address is not stored in wallet", func(t *testing.T) {
		t.Run("validator doesn't exist and public key not set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address().String(), "", amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator doesn't exist and public key set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address().String(), receiver.PublicKey().String(), amt)
			assert.NoError(t, err)
			assert.Equal(t, trx.Payload().(*payload.BondPayload).PublicKey.String(), receiver.PublicKey().String())
		})

		t.Run("validator exists and public key not set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address().String(), "", amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator exists and public key set", func(t *testing.T) {
			val := td.mockState.TestStore.AddTestValidator()

			trx, err := td.wallet.MakeBondTx(senderInfo.Address,
				val.Address().String(), receiver.PublicKey().String(), amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})
	})

	t.Run("validator address stored in wallet", func(t *testing.T) {
		receiver, _ := td.wallet.NewValidatorAddress("validator-address")
		receiverInfo := td.wallet.AddressInfo(receiver.Address)

		t.Run("validator doesn't exist and public key not set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address, "", amt)
			assert.NoError(t, err)
			assert.Equal(t, trx.Payload().(*payload.BondPayload).PublicKey.String(), receiverInfo.PublicKey)
		})

		t.Run("validator doesn't exist and public key set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address, receiverInfo.PublicKey, amt)
			assert.NoError(t, err)
			assert.Equal(t, trx.Payload().(*payload.BondPayload).PublicKey.String(), receiverInfo.PublicKey)
		})

		t.Run("validator exists and public key not set", func(t *testing.T) {
			val := td.mockState.TestStore.AddTestValidator()

			trx, err := td.wallet.MakeBondTx(senderInfo.Address,
				val.Address().String(), "", amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator exists and public key set", func(t *testing.T) {
			val := td.mockState.TestStore.AddTestValidator()

			trx, err := td.wallet.MakeBondTx(senderInfo.Address,
				val.Address().String(), receiverInfo.PublicKey, amt)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx("invalid_addr_string", receiver.Address().String(), "", amt)
		assert.Error(t, err)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx(senderInfo.Address, "invalid_addr_string", "", amt)
		assert.Error(t, err)
	})

	t.Run("invalid public key", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx(senderInfo.Address, receiver.Address().String(), "invalid-pub-key", amt)
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.Close()

		_, err := td.wallet.MakeBondTx(td.RandAccAddress().String(), receiver.Address().String(), "", amt)
		assert.Error(t, err)
	})
}

func TestMakeUnbondTx(t *testing.T) {
	td := setup(t)
	defer td.Close()

	senderInfo, _ := td.wallet.NewValidatorAddress("testing addr")

	t.Run("set parameters manually", func(t *testing.T) {
		lockTime := td.RandHeight()
		opts := []wallet.TxOption{
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo("test"),
		}

		trx, err := td.wallet.MakeUnbondTx(senderInfo.Address, opts...)
		assert.NoError(t, err)
		assert.Zero(t, trx.Fee()) // Fee for unbond transaction is zero
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, "test", trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		_ = td.mockState.TestStore.AddTestBlock(testHeight)

		trx, err := td.wallet.MakeUnbondTx(senderInfo.Address)
		assert.NoError(t, err)
		assert.Equal(t, trx.LockTime(), testHeight+1)
		assert.Zero(t, trx.Payload().Value())
		assert.Zero(t, trx.Fee())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeUnbondTx("invalid_addr_string")
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.Close()

		_, err := td.wallet.MakeUnbondTx(td.RandAccAddress().String())
		assert.Error(t, err)
	})
}

func TestMakeWithdrawTx(t *testing.T) {
	td := setup(t)
	defer td.Close()

	senderInfo, _ := td.wallet.NewBLSAccountAddress("testing addr")
	receiverInfo, _ := td.wallet.NewBLSAccountAddress("testing addr")
	amt := td.RandAmount()

	t.Run("set parameters manually", func(t *testing.T) {
		lockTime := td.RandHeight()
		fee := td.RandFee()
		opts := []wallet.TxOption{
			wallet.OptionFee(fee.String()),
			wallet.OptionLockTime(lockTime),
			wallet.OptionMemo("test"),
		}

		trx, err := td.wallet.MakeWithdrawTx(senderInfo.Address, receiverInfo.Address, amt, opts...)
		assert.NoError(t, err)
		assert.Equal(t, fee, trx.Fee())
		assert.Equal(t, lockTime, trx.LockTime())
		assert.Equal(t, "test", trx.Memo())
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		testHeight := td.RandHeight()
		_ = td.mockState.TestStore.AddTestBlock(testHeight)

		trx, err := td.wallet.MakeWithdrawTx(senderInfo.Address, receiverInfo.Address, amt)
		assert.NoError(t, err)
		assert.Equal(t, trx.LockTime(), testHeight+1)
		assert.Equal(t, amt, trx.Payload().Value())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeWithdrawTx("invalid_addr_string", receiverInfo.Address, amt)
		assert.Error(t, err)
	})

	t.Run("unable to get the blockchain info", func(t *testing.T) {
		td.Close()

		_, err := td.wallet.MakeWithdrawTx(td.RandAccAddress().String(), receiverInfo.Address, amt)
		assert.Error(t, err)
	})
}

func TestCheckMnemonic(t *testing.T) {
	mnemonic, _ := wallet.GenerateMnemonic(128)
	assert.NoError(t, wallet.CheckMnemonic(mnemonic))
}

func TestTotalBalance(t *testing.T) {
	td := setup(t)
	defer td.Close()

	addrInfo1, _ := td.wallet.NewBLSAccountAddress("account-1")
	_, _ = td.wallet.NewBLSAccountAddress("account-2")
	addrInfo3, _ := td.wallet.NewBLSAccountAddress("account-3")

	addr1, _ := crypto.AddressFromString(addrInfo1.Address)
	addr3, _ := crypto.AddressFromString(addrInfo3.Address)

	acc1, _ := td.GenerateTestAccount()
	acc3, _ := td.GenerateTestAccount()

	td.mockState.TestStore.Accounts[addr1] = acc1
	td.mockState.TestStore.Accounts[addr3] = acc3

	totalBalance, err := td.wallet.TotalBalance()
	assert.NoError(t, err)
	assert.Equal(t, totalBalance, acc1.Balance()+acc3.Balance())
}

func TestTotalStake(t *testing.T) {
	td := setup(t)
	defer td.Close()

	addrInfo1, _ := td.wallet.NewValidatorAddress("val-1")
	addrInfo2, _ := td.wallet.NewValidatorAddress("val-2")

	addr1, _ := crypto.AddressFromString(addrInfo1.Address)
	addr3, _ := crypto.AddressFromString(addrInfo2.Address)

	val1 := td.GenerateTestValidator()
	val2 := td.GenerateTestValidator()

	td.mockState.TestStore.Validators[addr1] = val1
	td.mockState.TestStore.Validators[addr3] = val2

	stake, err := td.wallet.TotalStake()
	require.NoError(t, err)

	require.Equal(t, stake, val1.Stake()+val2.Stake())
}

func TestCoinType(t *testing.T) {
	td := setup(t)
	defer td.Close()

	assert.Equal(t, td.wallet.CoinType(), uint32(21888))
}

func TestTestnetKeyInfo(t *testing.T) {
	td := setup(t)
	defer td.Close()

	mnemonic, _ := wallet.GenerateMnemonic(128)
	w1, err := wallet.Create(util.TempFilePath(), mnemonic, td.password, genesis.Mainnet)
	assert.NoError(t, err)
	addrInfo1, _ := w1.NewBLSAccountAddress("")
	prv1, _ := w1.PrivateKey("", addrInfo1.Address)

	w2, err := wallet.Create(util.TempFilePath(), mnemonic, td.password, genesis.Testnet)
	assert.NoError(t, err)
	addrInfo2, _ := w2.NewBLSAccountAddress("")
	prv2, _ := w2.PrivateKey("", addrInfo2.Address)

	assert.NotEqual(t, prv1.Bytes(), prv2.Bytes(),
		"Should generate different private key for the testnet")
}

func TestGetServerList(t *testing.T) {
	t.Run("Get mainnet servers", func(t *testing.T) {
		servers, err := wallet.GetServerList("mainnet")
		require.NoError(t, err)
		assert.NotEmpty(t, servers)

		// Verify structure
		for _, srv := range servers {
			assert.NotEmpty(t, srv.Address, "Server address should not be empty")
			// Name, Email, Website can be empty for some servers
		}

		// Check that at least official Pactus bootstrap servers exist
		foundBootstrap := false
		for _, srv := range servers {
			if strings.Contains(srv.Address, "bootstrap") && strings.Contains(srv.Address, "pactus.org") {
				foundBootstrap = true

				break
			}
		}
		assert.True(t, foundBootstrap, "Should contain at least one official bootstrap server")
	})

	t.Run("Get testnet servers", func(t *testing.T) {
		servers, err := wallet.GetServerList("testnet")
		require.NoError(t, err)
		assert.NotEmpty(t, servers)

		// Verify structure
		for _, srv := range servers {
			assert.NotEmpty(t, srv.Address, "Server address should not be empty")
		}
	})

	t.Run("Get servers for non-existent network", func(t *testing.T) {
		servers, err := wallet.GetServerList("nonexistent")
		require.NoError(t, err)
		assert.Empty(t, servers, "Should return empty list for unknown network")
	})

	t.Run("Get servers for localnet", func(t *testing.T) {
		servers, err := wallet.GetServerList("localnet")
		require.NoError(t, err)
		assert.Empty(t, servers, "Should return empty list for localnet")
	})
}
