package wallet

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/vault"
)

type Network uint8

const (
	NetworkMainNet = Network(0)
	NetworkTestNet = Network(1)
)

type Wallet struct {
	store  *store
	path   string
	client *grpcClient
}

//go:embed servers.json
var serversJSON []byte

type serverInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}
type servers = map[string][]serverInfo

// GenerateMnemonic is a wrapper for `vault.GenerateMnemonic.
func GenerateMnemonic(entropy int) string {
	return vault.GenerateMnemonic(entropy)
}

// OpenWallet tries to open a wallet at the given path.
// If the wallet doesn’t exist on this path, it returns an error.
// A wallet can be opened in offline or online modes.
// Offline wallet doesn’t have any connection to any node.
// Online wallet has a connection to one of the pre-defined servers.
func OpenWallet(path string, offline bool) (*Wallet, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	store := new(store)
	err = store.Save(data)
	if err != nil {
		return nil, err
	}

	return newWallet(path, store, offline)
}

// Create creates a wallet from mnemonic (seed phrase) and save it at the
// given path.
func Create(path, mnemonic, password string, net Network) (*Wallet, error) {
	path = util.MakeAbs(path)
	if util.PathExists(path) {
		return nil, NewErrWalletExits(path)
	}
	coinType := uint32(21888)
	if net == NetworkTestNet {
		coinType = uint32(21777)
	}
	store := &store{
		Version:   1,
		UUID:      uuid.New(),
		CreatedAt: time.Now().Round(time.Second).UTC(),
		Network:   net,
		Vault:     nil,
	}
	wallet, err := newWallet(path, store, true)
	if err != nil {
		return nil, err
	}
	vault, err := vault.CreateVaultFromMnemonic(mnemonic, coinType)
	if err != nil {
		return nil, err
	}
	err = vault.UpdatePassword("", password)
	if err != nil {
		return nil, err
	}
	wallet.store.Vault = vault

	return wallet, nil
}

func newWallet(path string, store *store, offline bool) (*Wallet, error) {
	if store.Network == NetworkTestNet {
		crypto.AddressHRP = "tpc"
		crypto.PublicKeyHRP = "tpublic"
		crypto.PrivateKeyHRP = "tsecret"
		crypto.XPublicKeyHRP = "txpublic"
		crypto.XPrivateKeyHRP = "txsecret"
	}

	w := &Wallet{
		store: store,
		path:  path,
	}

	if !offline {
		client, err := w.connectToRandomServer()
		if err != nil {
			return nil, err
		}
		w.client = client
	}

	return w, nil
}

func (w *Wallet) Connect(addr string) error {
	client, err := newGRPCClient(addr)
	if err != nil {
		return err
	}

	w.client = client
	return nil
}

func (w *Wallet) Name() string {
	return path.Base(w.path)
}

func (w *Wallet) IsOffline() bool {
	return w.client == nil
}

func (w *Wallet) connectToRandomServer() (*grpcClient, error) {
	serversInfo := servers{}
	err := json.Unmarshal(serversJSON, &serversInfo)
	if err != nil {
		return nil, err
	}

	var netServers []serverInfo
	switch w.store.Network {
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
			return nil, ErrInvalidNetwork
		}
	}

	for i := 0; i < 3; i++ {
		n := util.RandInt32(int32(len(netServers)))
		serverInfo := netServers[n]
		client, err := newGRPCClient(serverInfo.IP)
		if err == nil {
			return client, nil
		}
	}

	return nil, errors.New("unable to connect to the servers")
}

func (w *Wallet) Path() string {
	return w.path
}

func (w *Wallet) Save() error {
	bs, err := w.store.Load()
	if err != nil {
		return err
	}

	return util.WriteFile(w.path, bs)
}

// Balance returns the account balance amount.
func (w *Wallet) Balance(addrStr string) (int64, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return 0, err
	}

	if w.client == nil {
		return 0, ErrOffline
	}

	acc, _ := w.client.getAccount(addr)
	if acc != nil {
		return acc.Balance, nil
	}

	return 0, nil
}

// Stake returns the validator stake amount.
func (w *Wallet) Stake(addrStr string) (int64, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return 0, err
	}

	if w.client == nil {
		return 0, ErrOffline
	}

	val, _ := w.client.getValidator(addr)
	if val != nil {
		return val.Stake, nil
	}

	return 0, nil
}

// MakeSendTx creates a new send transaction based on the given parameters.
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
	err = maker.setToAddress(receiver)
	if err != nil {
		return nil, err
	}
	maker.amount = amount
	maker.typ = payload.PayloadTypeSend

	return maker.build()
}

// MakeBondTx creates a new bond transaction based on the given parameters.
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
	err = maker.setToAddress(receiver)
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

// MakeUnbondTx creates a new unbond transaction based on the given parameters.
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
// MakeWithdrawTx creates a new withdraw transaction based on the given
// parameters.
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
	err = maker.setToAddress(receiver)
	if err != nil {
		return nil, err
	}
	maker.amount = amount
	maker.typ = payload.PayloadTypeWithdraw

	return maker.build()
}

func (w *Wallet) SignTransaction(password string, trx *tx.Tx) error {
	prv, err := w.store.Vault.PrivateKey(password, trx.Payload().Signer().String())
	if err != nil {
		return err
	}

	signer := crypto.NewSigner(prv)
	signer.SignMsg(trx)
	if err != nil {
		return err
	}
	return nil
}

func (w *Wallet) BroadcastTransaction(trx *tx.Tx) (string, error) {
	if w.client == nil {
		return "", ErrOffline
	}

	id, err := w.client.sendTx(trx)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// TODO: query fee from grpc client
func (w *Wallet) CalculateFee(amount int64) int64 {
	return util.Max64(amount/10000, 10000)
}

func (w *Wallet) UpdatePassword(oldPassword, newPassword string) error {
	return w.store.Vault.UpdatePassword(oldPassword, newPassword)
}

func (w *Wallet) IsEncrypted() bool {
	return w.store.Vault.IsEncrypted()
}

func (w *Wallet) AddressInfo(addr string) *vault.AddressInfo {
	return w.store.Vault.AddressInfo(addr)
}

func (w *Wallet) AddressLabels() []vault.AddressInfo {
	return w.store.Vault.AddressLabels()
}

// AddressCount returns the number of addresses inside the wallet.
func (w *Wallet) AddressCount() int {
	return w.store.Vault.AddressCount()
}

func (w *Wallet) ImportPrivateKey(password string, prv crypto.PrivateKey) error {
	return w.store.Vault.ImportPrivateKey(password, prv)
}

func (w *Wallet) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	return w.store.Vault.PrivateKey(password, addr)
}

func (w *Wallet) DeriveNewAddress(label string) (string, error) {
	return w.store.Vault.DeriveNewAddress(label, vault.PurposeBLS12381)
}

func (w *Wallet) Contains(addr string) bool {
	return w.store.Vault.Contains(addr)
}

func (w *Wallet) Mnemonic(password string) (string, error) {
	return w.store.Vault.Mnemonic(password)
}

// Label returns label of addr.
func (w *Wallet) Label(addr string) string {
	return w.store.Vault.Label(addr)
}

// SetLabel sets label for addr.
func (w *Wallet) SetLabel(addr, label string) error {
	return w.store.Vault.SetLabel(addr, label)
}

func (w *Wallet) AddTransaction(id tx.ID) error {
	idStr := id.String()
	if w.store.History.hasTransaction(idStr) {
		return ErrHistoryExists
	}

	trxInfo, err := w.client.getTransaction(id)
	if err != nil {
		return err
	}

	trx, err := tx.FromBytes(trxInfo.Data)
	if err != nil {
		return err
	}

	blockHash, err := hash.FromBytes(trxInfo.BlockHash)
	if err != nil {
		return err
	}

	var sender crypto.Address
	var receiver *crypto.Address
	switch pld := trx.Payload().(type) {
	case *payload.SendPayload:
		sender = pld.Sender
		receiver = &pld.Receiver
	case *payload.BondPayload:
		sender = pld.Sender
		receiver = &pld.Receiver
	case *payload.UnbondPayload:
		sender = pld.Validator
		receiver = nil
	case *payload.WithdrawPayload:
		sender = pld.From
		receiver = &pld.To
	case *payload.SortitionPayload:
		sender = pld.Address
		receiver = nil
	}

	transaction := Transaction{
		BlockHash: blockHash.String(),
		Data:      hex.EncodeToString(trxInfo.Data),
	}
	activity := Activity{
		TxID:        trx.ID().String(),
		Status:      "confirmed",
		PayloadType: trx.Payload().Type().String(),
		BlockTime:   trxInfo.BlockTime,
	}

	if w.store.Vault.Contains(sender.String()) {
		activity.Amount = -(trx.Fee() + trx.Payload().Value())
		w.store.History.addTransaction(sender.String(),
			activity, transaction)
	}

	if receiver != nil {
		if w.store.Vault.Contains(receiver.String()) {
			activity.Amount = trx.Payload().Value()
			w.store.History.addTransaction(receiver.String(),
				activity, transaction)
		}
	}

	return nil
}

func (w *Wallet) GetHistory(addr string) []Activity {
	return w.store.History.Activities[addr]
}
