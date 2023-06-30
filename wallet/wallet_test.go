package wallet

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"path"
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testData struct {
	*testsuite.TestSuite

	wallet   *Wallet
	password string
}

var tListener *bufconn.Listener

func init() {
	const bufSize = 1024 * 1024

	tListener = bufconn.Listen(bufSize)

	s := grpc.NewServer()
	blockchainServer := &blockchainServer{}
	transactionServer := &transactionServer{}

	pactus.RegisterBlockchainServer(s, blockchainServer)
	pactus.RegisterTransactionServer(s, transactionServer)

	go func() {
		if err := s.Serve(tListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return tListener.Dial()
}

func setup(t *testing.T) *testData {
	ts := testsuite.NewTestSuite(t)

	password := ""
	walletPath := util.TempFilePath()
	mnemonic := GenerateMnemonic(128)
	wallet, err := Create(walletPath, mnemonic, password, genesis.Mainnet)
	assert.NoError(t, err)
	assert.False(t, wallet.IsEncrypted())
	assert.Equal(t, wallet.Path(), walletPath)
	assert.Equal(t, wallet.Name(), path.Base(walletPath))

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}

	client := &grpcClient{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
	}

	wallet.client = client

	assert.False(t, wallet.IsEncrypted())
	assert.False(t, wallet.IsOffline())

	tBlockchainInfoResponse = nil
	tAccountResponse = nil
	tValidatorResponse = nil

	return &testData{
		TestSuite: ts,
		wallet:    wallet,
		password:  password,
	}
}

func TestOpenWallet(t *testing.T) {
	td := setup(t)

	t.Run("Ok", func(t *testing.T) {
		assert.NoError(t, td.wallet.Save())
		_, err := Open(td.wallet.path, true)
		assert.NoError(t, err)
	})

	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := Open(util.TempFilePath(), true)
		assert.Error(t, err)
	})

	t.Run("Invalid crc", func(t *testing.T) {
		td.wallet.store.VaultCRC = 0
		bs, _ := json.Marshal(td.wallet.store)
		assert.NoError(t, util.WriteFile(td.wallet.path, bs))

		_, err := Open(td.wallet.path, true)
		assert.ErrorIs(t, err, ErrInvalidCRC)
	})

	t.Run("Invalid json", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(td.wallet.path, []byte("invalid_json")))

		_, err := Open(td.wallet.path, true)
		assert.Error(t, err)
	})
}

func TestRecoverWallet(t *testing.T) {
	td := setup(t)

	mnemonic, _ := td.wallet.Mnemonic(td.password)
	password := ""
	t.Run("Wallet exists", func(t *testing.T) {
		// Save the test wallet first then
		// try to recover a wallet at the same place
		assert.NoError(t, td.wallet.Save())

		_, err := Create(td.wallet.path, mnemonic, password, 0)
		assert.ErrorIs(t, err, NewErrWalletExits(td.wallet.path))
	})

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := Create(util.TempFilePath(),
			"invali mnemonic phrase seed", password, 0)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		path := util.TempFilePath()
		recovered, err := Create(path, mnemonic, password, genesis.Mainnet)
		assert.NoError(t, err)

		addr1, err := recovered.DeriveNewAddress("addr-1")
		assert.NoError(t, err)

		assert.NoFileExists(t, path)
		assert.NoError(t, recovered.Save())

		assert.FileExists(t, path)
		assert.True(t, recovered.Contains(addr1))
	})
}

func TestSaveWallet(t *testing.T) {
	td := setup(t)

	t.Run("Invalid path", func(t *testing.T) {
		td.wallet.path = "/"
		assert.Error(t, td.wallet.Save())
	})
}
func TestInvalidAddress(t *testing.T) {
	td := setup(t)

	addr := td.RandomAddress().String()
	_, err := td.wallet.PrivateKey(td.password, addr)
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	td := setup(t)

	_, prv := td.RandomBLSKeyPair()
	assert.NoError(t, td.wallet.ImportPrivateKey(td.password, prv))

	addr := prv.PublicKey().Address().String()
	assert.True(t, td.wallet.Contains(addr))
}

func TestTestKeyInfo(t *testing.T) {
	td := setup(t)

	mnemonic := GenerateMnemonic(128)
	w1, err := Create(util.TempFilePath(), mnemonic, td.password,
		genesis.Mainnet)
	assert.NoError(t, err)
	addrStr1, _ := w1.DeriveNewAddress("")
	prv1, _ := w1.PrivateKey("", addrStr1)

	w2, err := Create(util.TempFilePath(), mnemonic, td.password,
		genesis.Testnet)
	assert.NoError(t, err)
	addrStr2, _ := w2.DeriveNewAddress("")
	prv2, _ := w2.PrivateKey("", addrStr2)

	assert.NotEqual(t, prv1.Bytes(), prv2.Bytes(),
		"Should generate different private key for the testnet")
}

func TestBalance(t *testing.T) {
	td := setup(t)

	addr := td.RandomAddress()
	tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Balance: 1}}
	amt, err := td.wallet.Balance(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, amt, int64(1))
}

func TestStake(t *testing.T) {
	td := setup(t)

	addr := td.RandomAddress()
	tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Stake: 1}}
	amt, err := td.wallet.Stake(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, amt, int64(1))
}

func TestAccountSequence(t *testing.T) {
	td := setup(t)

	addr := td.RandomAddress()
	tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Sequence: 123}}
	seq, err := td.wallet.AccountSequence(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, seq, int32(123))
}

func TestValidatorSequence(t *testing.T) {
	td := setup(t)

	addr := td.RandomAddress()
	tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Sequence: 123}}
	seq, err := td.wallet.ValidatorSequence(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, seq, int32(123))
}

func TestSigningTx(t *testing.T) {
	td := setup(t)

	sender, _ := td.wallet.DeriveNewAddress("testing addr")
	receiver := td.RandomAddress()
	amount := td.RandInt64(10000)
	seq := td.RandInt32(10000)

	opts := []TxOption{
		OptionStamp(td.RandomStamp().String()),
		OptionFee(util.CoinToChange(10)),
		OptionSequence(seq),
		OptionMemo("test"),
	}

	trx, err := td.wallet.MakeTransferTx(sender, receiver.String(), amount, opts...)
	assert.NoError(t, err)
	err = td.wallet.SignTransaction(td.password, trx)
	assert.NoError(t, err)
	assert.NotNil(t, trx.Signature())
	assert.NoError(t, trx.SanityCheck())

	id, err := td.wallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)
}

func TestMakeTransferTx(t *testing.T) {
	td := setup(t)

	sender, _ := td.wallet.DeriveNewAddress("testing addr")
	receiver := td.RandomAddress()
	amount := td.RandInt64(10000)
	seq := td.RandInt32(10000)

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := td.RandomStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := td.wallet.MakeTransferTx(sender, receiver.String(), amount, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Equal(t, trx.Fee(), util.CoinToChange(10))
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		lastBlockHash := td.RandomHash()
		tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHash.Bytes()}
		tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Sequence: seq}}

		trx, err := td.wallet.MakeTransferTx(sender, receiver.String(), amount)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.Equal(t, trx.Payload().Value(), amount)
		assert.Equal(t, trx.Fee(), td.wallet.CalculateFee(amount))
		assert.Equal(t, trx.Stamp(), lastBlockHash.Stamp())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeTransferTx("invalid_addr_string", receiver.String(), amount)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := td.wallet.MakeTransferTx(sender, "invalid_addr_string", amount)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("sender account doesn't exist", func(t *testing.T) {
		tAccountResponse = nil

		_, err := td.wallet.MakeTransferTx(td.RandomAddress().String(), receiver.String(), amount)
		assert.Equal(t, errors.Code(err), errors.ErrGeneric)
	})
}

func TestMakeBondTx(t *testing.T) {
	td := setup(t)

	sender, _ := td.wallet.DeriveNewAddress("testing addr")
	receiver := td.RandomSigner()
	amount := td.RandInt64(10000)
	seq := td.RandInt32(10000)

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := td.RandomStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := td.wallet.MakeBondTx(sender, receiver.Address().String(),
			receiver.PublicKey().String(), amount, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Equal(t, trx.Fee(), util.CoinToChange(10))
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		lastBlockHash := td.RandomHash()
		tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Sequence: seq}}
		tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHash.Bytes()}

		trx, err := td.wallet.MakeBondTx(sender, receiver.Address().String(), receiver.PublicKey().String(), amount)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.True(t, trx.Payload().(*payload.BondPayload).PublicKey.EqualsTo(receiver.PublicKey()))
		assert.Equal(t, trx.Payload().Value(), amount)
		assert.Equal(t, trx.Fee(), td.wallet.CalculateFee(amount))
		assert.Equal(t, trx.Stamp(), lastBlockHash.Stamp())
	})

	t.Run("validator address is not stored in wallet", func(t *testing.T) {
		t.Run("validator doesn't exist and public key not set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(sender, receiver.Address().String(), "", amount)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator doesn't exist and public key set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(sender, receiver.Address().String(), receiver.PublicKey().String(), amount)
			assert.NoError(t, err)
			assert.Equal(t, trx.Payload().(*payload.BondPayload).PublicKey.String(), receiver.PublicKey().String())
		})

		t.Run("validator exists and public key not set", func(t *testing.T) {
			tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{}}
			trx, err := td.wallet.MakeBondTx(sender, receiver.Address().String(), "", amount)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator exists and public key set", func(t *testing.T) {
			trx, err := td.wallet.MakeBondTx(sender, receiver.Address().String(), receiver.PublicKey().String(), amount)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})
	})

	t.Run("validator address stored in wallet", func(t *testing.T) {
		receiver, _ := td.wallet.DeriveNewAddress("validator-address")
		receiverInfo := td.wallet.AddressInfo(receiver)

		t.Run("validator doesn't exist and public key not set", func(t *testing.T) {
			tValidatorResponse = nil
			trx, err := td.wallet.MakeBondTx(sender, receiver, "", amount)
			assert.NoError(t, err)
			assert.Equal(t, trx.Payload().(*payload.BondPayload).PublicKey.String(), receiverInfo.Pub.String())
		})

		t.Run("validator doesn't exist and public key set", func(t *testing.T) {
			receiverInfo := td.wallet.AddressInfo(receiver)
			trx, err := td.wallet.MakeBondTx(sender, receiver, receiverInfo.Pub.String(), amount)
			assert.NoError(t, err)
			assert.Equal(t, trx.Payload().(*payload.BondPayload).PublicKey.String(), receiverInfo.Pub.String())
		})

		t.Run("validator exists and public key not set", func(t *testing.T) {
			tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{}}
			trx, err := td.wallet.MakeBondTx(sender, receiver, "", amount)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})

		t.Run("validator exists and public key set", func(t *testing.T) {
			receiverInfo := td.wallet.AddressInfo(receiver)
			trx, err := td.wallet.MakeBondTx(sender, receiver, receiverInfo.Pub.String(), amount)
			assert.NoError(t, err)
			assert.Nil(t, trx.Payload().(*payload.BondPayload).PublicKey)
		})
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx("invalid_addr_string", receiver.Address().String(), "", amount)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx(sender, "invalid_addr_string", "", amount)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("invalid public key", func(t *testing.T) {
		_, err := td.wallet.MakeBondTx(sender, receiver.Address().String(), "invalid-pub-key", amount)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	t.Run("sender account doesn't exist", func(t *testing.T) {
		tAccountResponse = nil

		_, err := td.wallet.MakeBondTx(td.RandomAddress().String(), receiver.Address().String(), "", amount)
		assert.Equal(t, errors.Code(err), errors.ErrGeneric)
	})
}

func TestMakeUnbondTx(t *testing.T) {
	td := setup(t)

	sender, _ := td.wallet.DeriveNewAddress("testing addr")
	seq := td.RandInt32(10000)

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := td.RandomStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := td.wallet.MakeUnbondTx(sender, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Zero(t, trx.Fee()) // Fee for unbond transaction is zero
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		lastBlockHash := td.RandomHash()
		tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Sequence: seq}}
		tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHash.Bytes()}

		trx, err := td.wallet.MakeUnbondTx(sender)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.Zero(t, trx.Payload().Value())
		assert.Zero(t, trx.Fee())
		assert.Equal(t, trx.Stamp(), lastBlockHash.Stamp())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeUnbondTx("invalid_addr_string")
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("sender account doesn't exist", func(t *testing.T) {
		tValidatorResponse = nil

		_, err := td.wallet.MakeUnbondTx(td.RandomAddress().String())
		assert.Equal(t, errors.Code(err), errors.ErrGeneric)
	})
}

func TestMakeWithdrawTx(t *testing.T) {
	td := setup(t)

	sender, _ := td.wallet.DeriveNewAddress("testing addr")
	receiver, _ := td.wallet.DeriveNewAddress("testing addr")
	amount := td.RandInt64(10000)
	seq := td.RandInt32(10000)

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := td.RandomStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := td.wallet.MakeWithdrawTx(sender, receiver, amount, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Equal(t, trx.Fee(), util.CoinToChange(10)) // Fee for unbond transaction is zero
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("query parameters from the node", func(t *testing.T) {
		lastBlockHash := td.RandomHash()
		tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Sequence: seq}}
		tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHash.Bytes()}

		trx, err := td.wallet.MakeWithdrawTx(sender, receiver, amount)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.Equal(t, trx.Payload().Value(), amount)
		assert.Equal(t, trx.Fee(), td.wallet.CalculateFee(amount))
		assert.Equal(t, trx.Stamp(), lastBlockHash.Stamp())
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := td.wallet.MakeWithdrawTx("invalid_addr_string", receiver, amount)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("sender account doesn't exist", func(t *testing.T) {
		tValidatorResponse = nil

		_, err := td.wallet.MakeWithdrawTx(td.RandomAddress().String(), receiver, amount)
		assert.Equal(t, errors.Code(err), errors.ErrGeneric)
	})
}

func TestCheckMnemonic(t *testing.T) {
	mnemonic := GenerateMnemonic(128)
	assert.NoError(t, CheckMnemonic(mnemonic))
}
