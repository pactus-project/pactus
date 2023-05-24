package wallet

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"path"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var tWallet *Wallet
var tPassword string
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

func setup(t *testing.T) {
	tPassword := ""
	walletPath := util.TempFilePath()
	mnemonic := GenerateMnemonic(128)
	w, err := Create(walletPath, mnemonic, tPassword, genesis.Mainnet)
	assert.NoError(t, err)
	assert.False(t, w.IsEncrypted())
	assert.Equal(t, w.Path(), walletPath)
	assert.Equal(t, w.Name(), path.Base(walletPath))

	tWallet = w

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}

	client := &grpcClient{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
	}

	tWallet.client = client

	assert.False(t, w.IsEncrypted())
	assert.False(t, tWallet.IsOffline())
}

func TestOpenWallet(t *testing.T) {
	setup(t)

	t.Run("Ok", func(t *testing.T) {
		assert.NoError(t, tWallet.Save())
		_, err := Open(tWallet.path, true)
		assert.NoError(t, err)
	})

	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := Open(util.TempFilePath(), true)
		assert.Error(t, err)
	})

	t.Run("Invalid crc", func(t *testing.T) {
		tWallet.store.VaultCRC = 0
		bs, _ := json.Marshal(tWallet.store)
		assert.NoError(t, util.WriteFile(tWallet.path, bs))

		_, err := Open(tWallet.path, true)
		assert.ErrorIs(t, err, ErrInvalidCRC)
	})

	t.Run("Invalid json", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(tWallet.path, []byte("invalid_json")))

		_, err := Open(tWallet.path, true)
		assert.Error(t, err)
	})
}

func TestRecoverWallet(t *testing.T) {
	setup(t)

	mnemonic, _ := tWallet.Mnemonic(tPassword)
	password := ""
	t.Run("Wallet exists", func(t *testing.T) {
		// Save the test wallet first then
		// try to recover a wallet at the same place
		assert.NoError(t, tWallet.Save())

		_, err := Create(tWallet.path, mnemonic, password, 0)
		assert.ErrorIs(t, err, NewErrWalletExits(tWallet.path))
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
	setup(t)

	t.Run("Invalid path", func(t *testing.T) {
		tWallet.path = "/"
		assert.Error(t, tWallet.Save())
	})
}
func TestInvalidAddress(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress().String()
	_, err := tWallet.PrivateKey(tPassword, addr)
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	setup(t)

	_, prv := bls.GenerateTestKeyPair()
	assert.NoError(t, tWallet.ImportPrivateKey(tPassword, prv))

	addr := prv.PublicKey().Address().String()
	assert.True(t, tWallet.Contains(addr))
}

func TestTestKeyInfo(t *testing.T) {
	mnemonic := GenerateMnemonic(128)
	w1, err := Create(util.TempFilePath(), mnemonic, tPassword,
		genesis.Mainnet)
	assert.NoError(t, err)
	addrStr1, _ := w1.DeriveNewAddress("")
	prv1, _ := w1.PrivateKey("", addrStr1)

	w2, err := Create(util.TempFilePath(), mnemonic, tPassword,
		genesis.Testnet)
	assert.NoError(t, err)
	addrStr2, _ := w2.DeriveNewAddress("")
	prv2, _ := w2.PrivateKey("", addrStr2)

	assert.NotEqual(t, prv1.Bytes(), prv2.Bytes(),
		"Should generate different private key for the testnet")
}

func TestBalance(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	tAccountRequest = &pactus.GetAccountRequest{Address: addr.String()}
	tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Balance: 1}}
	amt, err := tWallet.Balance(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, amt, int64(1))
}

func TestStake(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	tValidatorRequest = &pactus.GetValidatorRequest{Address: addr.String()}
	tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Stake: 1}}
	amt, err := tWallet.Stake(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, amt, int64(1))
}

func TestAccountSequence(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	tAccountRequest = &pactus.GetAccountRequest{Address: addr.String()}
	tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Sequence: 123}}
	seq, err := tWallet.AccountSequence(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, seq, int32(123))
}

func TestValidatorSequence(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	tValidatorRequest = &pactus.GetValidatorRequest{Address: addr.String()}
	tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Sequence: 123}}
	seq, err := tWallet.ValidatorSequence(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, seq, int32(123))
}

func TestSigningTx(t *testing.T) {
	setup(t)

	sender, _ := tWallet.DeriveNewAddress("testing addr")
	receiver := crypto.GenerateTestAddress()
	amount := util.RandInt64(10000)
	seq := util.RandInt32(10000)

	opts := []TxOption{
		OptionStamp(hash.GenerateTestStamp().String()),
		OptionFee(util.CoinToChange(10)),
		OptionSequence(seq),
		OptionMemo("test"),
	}

	trx, err := tWallet.MakeSendTx(sender, receiver.String(), amount, opts...)
	assert.NoError(t, err)
	err = tWallet.SignTransaction(tPassword, trx)
	assert.NoError(t, err)
	assert.NotNil(t, trx.Signature())
	assert.NoError(t, trx.SanityCheck())

	id, err := tWallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)
}

func TestMakeSendTx(t *testing.T) {
	setup(t)

	sender, _ := tWallet.DeriveNewAddress("testing addr")
	receiver := crypto.GenerateTestAddress()
	amount := util.RandInt64(10000)
	seq := util.RandInt32(10000)
	lastBlockHsh := hash.GenerateTestHash().Bytes()

	tAccountRequest = &pactus.GetAccountRequest{Address: sender}
	tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Sequence: seq}}
	tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHsh}

	t.Run("query parameters from the node", func(t *testing.T) {
		trx, err := tWallet.MakeSendTx(sender, receiver.String(), amount)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.Equal(t, trx.Payload().Value(), amount)
		assert.Equal(t, trx.Fee(), tWallet.CalculateFee(amount))
	})

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := tWallet.MakeSendTx(sender, receiver.String(), amount, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Equal(t, trx.Fee(), util.CoinToChange(10))
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("sender address doesn't exist", func(t *testing.T) {
		_, err := tWallet.MakeSendTx(crypto.GenerateTestAddress().String(), receiver.String(), amount)
		assert.Error(t, err)
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := tWallet.MakeSendTx("invalid_addr_string", receiver.String(), amount)
		assert.Error(t, err)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := tWallet.MakeSendTx(sender, "invalid_addr_string", amount)
		assert.Error(t, err)
	})
}

func TestMakeBondTx(t *testing.T) {
	setup(t)

	sender, _ := tWallet.DeriveNewAddress("testing addr")
	receiver := bls.GenerateTestSigner()
	amount := util.RandInt64(10000)
	seq := util.RandInt32(10000)
	lastBlockHsh := hash.GenerateTestHash().Bytes()

	tAccountRequest = &pactus.GetAccountRequest{Address: sender}
	tAccountResponse = &pactus.GetAccountResponse{Account: &pactus.AccountInfo{Sequence: seq}}
	tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHsh}

	t.Run("query parameters from the node", func(t *testing.T) {
		trx, err := tWallet.MakeBondTx(sender, receiver.Address().String(), receiver.PublicKey().String(), amount)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.True(t, trx.Payload().(*payload.BondPayload).PublicKey.EqualsTo(receiver.PublicKey()))
		assert.Equal(t, trx.Payload().Value(), amount)
		assert.Equal(t, trx.Fee(), tWallet.CalculateFee(amount))
	})

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := tWallet.MakeBondTx(sender, receiver.Address().String(),
			receiver.PublicKey().String(), amount, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Equal(t, trx.Fee(), util.CoinToChange(10))
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("sender address doesn't exist", func(t *testing.T) {
		_, err := tWallet.MakeBondTx(crypto.GenerateTestAddress().String(), receiver.Address().String(), "", amount)
		assert.Error(t, err)
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := tWallet.MakeBondTx("invalid_addr_string", receiver.Address().String(), "", amount)
		assert.Error(t, err)
	})

	t.Run("invalid receiver address", func(t *testing.T) {
		_, err := tWallet.MakeBondTx(sender, "invalid_addr_string", "", amount)
		assert.Error(t, err)
	})

	t.Run("invalid public key", func(t *testing.T) {
		_, err := tWallet.MakeBondTx(sender, receiver.Address().String(), "invalid-pub-key", amount)
		assert.Error(t, err)
	})
}

func TestMakeUnbondTx(t *testing.T) {
	setup(t)

	sender, _ := tWallet.DeriveNewAddress("testing addr")
	seq := util.RandInt32(10000)
	lastBlockHsh := hash.GenerateTestHash().Bytes()

	tValidatorRequest = &pactus.GetValidatorRequest{Address: sender}
	tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Sequence: seq}}
	tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHsh}

	t.Run("query parameters from the node", func(t *testing.T) {
		trx, err := tWallet.MakeUnbondTx(sender)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.Zero(t, trx.Payload().Value())
		assert.Zero(t, trx.Fee())
	})

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := tWallet.MakeUnbondTx(sender, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Zero(t, trx.Fee()) // Fee for unbond transaction is zero
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("sender address doesn't exist", func(t *testing.T) {
		_, err := tWallet.MakeUnbondTx(crypto.GenerateTestAddress().String())
		assert.Error(t, err)
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := tWallet.MakeUnbondTx("invalid_addr_string")
		assert.Error(t, err)
	})
}

func TestMakeWithdrawTx(t *testing.T) {
	setup(t)

	sender, _ := tWallet.DeriveNewAddress("testing addr")
	receiver, _ := tWallet.DeriveNewAddress("testing addr")
	amount := util.RandInt64(10000)
	seq := util.RandInt32(10000)
	lastBlockHsh := hash.GenerateTestHash().Bytes()

	tValidatorRequest = &pactus.GetValidatorRequest{Address: sender}
	tValidatorResponse = &pactus.GetValidatorResponse{Validator: &pactus.ValidatorInfo{Sequence: seq}}
	tBlockchainInfoResponse = &pactus.GetBlockchainInfoResponse{LastBlockHash: lastBlockHsh}

	t.Run("query parameters from the node", func(t *testing.T) {
		trx, err := tWallet.MakeWithdrawTx(sender, receiver, amount)
		assert.NoError(t, err)
		assert.Equal(t, trx.Sequence(), seq+1)
		assert.Equal(t, trx.Payload().Value(), amount)
		assert.Equal(t, trx.Fee(), tWallet.CalculateFee(amount))
	})

	t.Run("set parameters manually", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		opts := []TxOption{
			OptionStamp(stamp.String()),
			OptionFee(util.CoinToChange(10)),
			OptionSequence(seq),
			OptionMemo("test"),
		}

		trx, err := tWallet.MakeWithdrawTx(sender, receiver, amount, opts...)
		assert.NoError(t, err)
		assert.Equal(t, trx.Stamp(), stamp)
		assert.Equal(t, trx.Fee(), util.CoinToChange(10)) // Fee for unbond transaction is zero
		assert.Equal(t, trx.Sequence(), seq)
		assert.Equal(t, trx.Memo(), "test")
	})

	t.Run("sender address doesn't exist", func(t *testing.T) {
		_, err := tWallet.MakeWithdrawTx(crypto.GenerateTestAddress().String(), receiver, amount)
		assert.Error(t, err)
	})

	t.Run("invalid sender address", func(t *testing.T) {
		_, err := tWallet.MakeWithdrawTx("invalid_addr_string", receiver, amount)
		assert.Error(t, err)
	})
}
