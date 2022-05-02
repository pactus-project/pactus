package wallet

import (
	_ "embed"
	"encoding/json"
	"errors"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/tx/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/wallet/vault"
)

type Network uint8

const (
	NetworkMainNet = Network(0)
	NetworkTestNet = Network(1)
)

type Wallet struct {
	*store

	path   string
	client *grpcClient
}

type serverInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}
type servers = map[string][]serverInfo

// GenerateMnemonic is a wrapper for `vault.GenerateMnemonic`
func GenerateMnemonic() string {
	return vault.GenerateMnemonic()
}

//go:embed servers.json
var serversJSON []byte

/// OpenWallet tries to open a wallet at given path
func OpenWallet(path string) (*Wallet, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	store := new(store)
	err = store.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	return newWallet(path, store, true)
}

/// FromMnemonic creates a wallet from mnemonic (seed phrase)
func FromMnemonic(path, mnemonic, password string, net Network) (*Wallet, error) {
	path = util.MakeAbs(path)
	if util.PathExists(path) {
		return nil, NewErrWalletExits(path)
	}
	keyInfo := []byte{}
	if net == NetworkTestNet {
		keyInfo = []byte{1}
	}
	vault, err := vault.CreateVaultFromMnemonic(mnemonic, password, keyInfo)
	if err != nil {
		return nil, err
	}
	store := &store{
		data: storeData{
			Version:   1,
			UUID:      uuid.New(),
			CreatedAt: time.Now().Round(time.Second).UTC(),
			Network:   net,
			Vault:     vault,
		},
	}

	return newWallet(path, store, false)
}

func newWallet(path string, store *store, online bool) (*Wallet, error) {
	if store.data.Network == NetworkTestNet {
		crypto.DefaultHRP = "tzc"
	}

	w := &Wallet{
		store: store,
		path:  path,
	}

	err := w.connectToRandomServer()
	if err != nil {
		return nil, err
	}

	return w, nil
}
func (w *Wallet) Name() string {
	return path.Base(w.path)
}

func (w *Wallet) UpdatePassword(old, new string) error {
	return w.store.UpdatePassword(old, new)
}

func (w *Wallet) connectToRandomServer() error {
	serversInfo := servers{}
	err := json.Unmarshal(serversJSON, &serversInfo)
	if err != nil {
		return err
	}

	var netServers []serverInfo
	switch w.store.data.Network {
	case NetworkMainNet:
		{ // mainnet
			netServers = serversInfo["mainnet"]
		}
	case NetworkTestNet:
		{ // testnet
			netServers = serversInfo["testnet"]
		}

	default:
		{
			return ErrInvalidNetwork
		}
	}

	for i := 0; i < 3; i++ {
		n := util.RandInt32(int32(len(netServers)))
		serverInfo := netServers[n]
		client, err := gewGRPCClient(serverInfo.IP)
		if err == nil {
			w.client = client
			return nil
		}
	}

	return errors.New("unable to connect to the servers")
}

func (w *Wallet) Path() string {
	return w.path
}

func (w *Wallet) Save() error {
	bs, err := w.store.MarshalJSON()
	if err != nil {
		return err
	}

	return util.WriteFile(w.path, bs)
}

// Balance returns the account balance amount
func (w *Wallet) Balance(addrStr string) (int64, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return 0, err
	}

	balance, _ := w.client.getAccountBalance(addr)
	//exitOnErr(err)

	return balance, nil
}

// Stake returns the validator stake amount
func (w *Wallet) Stake(addrStr string) (int64, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return 0, err
	}

	stake, _ := w.client.getValidatorStake(addr)
	//exitOnErr(err)

	return stake, nil
}

// MakeSendTx creates a new send transaction based on the given parameters
func (w *Wallet) MakeSendTx(sender, receiver string, amount int64,
	options ...TxOption) (*tx.Tx, error) {
	maker, err := newTxMaker(w.client, options...)
	if err != nil {
		return nil, err
	}
	err = maker.setFromAddr(sender)
	if err != nil {
		return nil, err
	}
	err = maker.setToAddress(sender)
	if err != nil {
		return nil, err
	}
	maker.amount = amount
	maker.typ = payload.PayloadTypeSend

	return maker.build()
}

// MakeBondTx creates a new bond transaction based on the given parameters
func (w *Wallet) MakeBondTx(sender, receiver, pubKey string, amount int64,
	options ...TxOption) (*tx.Tx, error) {
	maker, err := newTxMaker(w.client, options...)
	if err != nil {
		return nil, err
	}
	err = maker.setFromAddr(sender)
	if err != nil {
		return nil, err
	}
	err = maker.setToAddress(sender)
	if err != nil {
		return nil, err
	}
	if pubKey != "" {
		maker.pub, err = bls.PublicKeyFromString(pubKey)
		if err != nil {
			return nil, err
		}
	}
	maker.amount = amount
	maker.typ = payload.PayloadTypeBond

	return maker.build()
}

// MakeUnbondTx creates a new unbond transaction based on the given parameters
func (w *Wallet) MakeUnbondTx(addr string, options ...TxOption) (*tx.Tx, error) {
	maker, err := newTxMaker(w.client, options...)
	if err != nil {
		return nil, err
	}
	err = maker.setFromAddr(addr)
	if err != nil {
		return nil, err
	}
	maker.typ = payload.PayloadTypeUnbond

	return maker.build()
}

// TODO: write tests for me by mocking grpc server
// MakeWithdrawTx creates a new unbond transaction based on the given parameters
func (w *Wallet) MakeWithdrawTx(sender, receiver string, amount int64,
	options ...TxOption) (*tx.Tx, error) {
	maker, err := newTxMaker(w.client, options...)
	if err != nil {
		return nil, err
	}
	err = maker.setFromAddr(sender)
	if err != nil {
		return nil, err
	}
	err = maker.setToAddress(sender)
	if err != nil {
		return nil, err
	}
	maker.amount = amount
	maker.typ = payload.PayloadTypeWithdraw

	return maker.build()
}

func (w *Wallet) SignAndBroadcast(password string, tx *tx.Tx) (string, error) {
	prvStr, err := w.PrivateKey(password, tx.Payload().Signer().String())
	if err != nil {
		return "", err
	}

	prv, err := bls.PrivateKeyFromString(prvStr)
	if err != nil {
		return "", err
	}

	signer := crypto.NewSigner(prv)
	signer.SignMsg(tx)
	b, err := tx.Bytes()
	if err != nil {
		return "", err
	}

	return w.client.sendTx(b)
}
