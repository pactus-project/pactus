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
	"github.com/pactus-project/pactus/types/account"
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
	networkServer := &networkServer{}

	pactus.RegisterBlockchainServer(s, blockchainServer)
	pactus.RegisterNetworkServer(s, networkServer)

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
	w, err := Create(walletPath, mnemonic, tPassword, NetworkMainNet)
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
		_, err := OpenWallet(tWallet.path, true)
		assert.NoError(t, err)
	})

	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := OpenWallet(util.TempFilePath(), true)
		assert.Error(t, err)
	})

	t.Run("Invalid crc", func(t *testing.T) {
		tWallet.store.VaultCRC = 0
		bs, _ := json.Marshal(tWallet.store)
		assert.NoError(t, util.WriteFile(tWallet.path, bs))

		_, err := OpenWallet(tWallet.path, true)
		assert.ErrorIs(t, err, ErrInvalidCRC)
	})

	t.Run("Invalid json", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(tWallet.path, []byte("invalid_json")))

		_, err := OpenWallet(tWallet.path, true)
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
		recovered, err := Create(path, mnemonic, password, NetworkMainNet)
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
		NetworkMainNet)
	assert.NoError(t, err)
	addrStr1, _ := w1.DeriveNewAddress("")
	prv1, _ := w1.PrivateKey("", addrStr1)

	w2, err := Create(util.TempFilePath(), mnemonic, tPassword,
		NetworkTestNet)
	assert.NoError(t, err)
	addrStr2, _ := w2.DeriveNewAddress("")
	prv2, _ := w2.PrivateKey("", addrStr2)

	assert.NotEqual(t, prv1.Bytes(), prv2.Bytes(),
		"Should generate different private key for the testnet")
}

func TestGetBalance(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	tAccountRequest = &pactus.AccountRequest{Address: addr.String()}
	tAccountResponse = &pactus.AccountResponse{Account: &pactus.AccountInfo{Balance: 1}}
	bal, err := tWallet.Balance(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, bal, int64(1))
}

func TestGetStake(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	tValidatorRequest = &pactus.ValidatorRequest{Address: addr.String()}
	tValidatorResponse = &pactus.ValidatorResponse{Validator: &pactus.ValidatorInfo{Stake: 1}}
	bal, err := tWallet.Stake(addr.String())
	assert.NoError(t, err)
	assert.Equal(t, bal, int64(1))
}

func TestMakeSendTx(t *testing.T) {
	setup(t)

	senderAcc, _ := account.GenerateTestAccount(util.RandInt32(0))
	receiver := crypto.GenerateTestAddress()
	amount := int64(1)
	lastBlockHsh := hash.GenerateTestHash().Bytes()

	tAccountRequest = &pactus.AccountRequest{Address: senderAcc.Address().String()}
	tAccountResponse = &pactus.AccountResponse{Account: &pactus.AccountInfo{Sequence: senderAcc.Sequence()}}
	tBlockchainInfoResponse = &pactus.BlockchainInfoResponse{LastBlockHash: lastBlockHsh}
	tx1, err := tWallet.MakeSendTx(senderAcc.Address().String(), receiver.String(), amount)
	assert.NoError(t, err)
	assert.Equal(t, tx1.Sequence(), senderAcc.Sequence()+1)
	assert.Equal(t, tx1.Payload().Value(), amount)

	stamp := hash.GenerateTestStamp()
	opts := []TxOption{
		OptionStamp(stamp.String()),
		OptionFee(util.CoinToChange(10)),
		OptionSequence(int32(20)),
		OptionMemo("test"),
	}

	tx2, err := tWallet.MakeSendTx(senderAcc.Address().String(), receiver.String(), amount, opts...)
	assert.NoError(t, err)
	assert.Equal(t, tx2.Stamp(), stamp)
	assert.Equal(t, tx2.Fee(), util.CoinToChange(10))
	assert.Equal(t, tx2.Sequence(), int32(20))
	assert.Equal(t, tx2.Memo(), "test")
}
