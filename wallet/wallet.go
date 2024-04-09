package wallet

import (
	_ "embed"
	"encoding/json"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Wallet struct {
	store *store
	path  string

	lazyClient *grpcClient
}

//go:embed servers.json
var serversJSON []byte

// GenerateMnemonic is a wrapper for `vault.GenerateMnemonic`.
func GenerateMnemonic(entropy int) (string, error) {
	return vault.GenerateMnemonic(entropy)
}

// CheckMnemonic is a wrapper for `vault.CheckMnemonic`.
func CheckMnemonic(mnemonic string) error {
	return vault.CheckMnemonic(mnemonic)
}

// Open tries to open a wallet at the given path.
// If the wallet doesn’t exist on this path, it returns an error.
// A wallet can be opened in offline or online modes.
// Offline wallet doesn’t have any connection to any node.
// Online wallet has a connection to one of the pre-defined servers.
func Open(walletPath string, offline bool) (*Wallet, error) {
	data, err := util.ReadFile(walletPath)
	if err != nil {
		return nil, err
	}

	store := new(store)
	err = store.Save(data)
	if err != nil {
		return nil, err
	}

	return newWallet(walletPath, store, offline)
}

// Create creates a wallet from mnemonic (seed phrase) and save it at the
// given path.
func Create(walletPath, mnemonic, password string, chain genesis.ChainType) (*Wallet, error) {
	walletPath = util.MakeAbs(walletPath)
	if util.PathExists(walletPath) {
		return nil, ExitsError{
			Path: walletPath,
		}
	}

	var coinType uint32
	switch chain {
	case genesis.Mainnet:
		coinType = 21888
	case genesis.Testnet, genesis.Localnet:
		coinType = 21777
	default:
		return nil, ErrInvalidNetwork
	}

	store := &store{
		Version:   1,
		UUID:      uuid.New(),
		CreatedAt: time.Now().Round(time.Second).UTC(),
		Network:   chain,
		Vault:     nil,
	}
	wallet, err := newWallet(walletPath, store, true)
	if err != nil {
		return nil, err
	}
	vlt, err := vault.CreateVaultFromMnemonic(mnemonic, coinType)
	if err != nil {
		return nil, err
	}
	err = vlt.UpdatePassword("", password)
	if err != nil {
		return nil, err
	}
	wallet.store.Vault = vlt

	return wallet, nil
}

func newWallet(walletPath string, store *store, offline bool) (*Wallet, error) {
	if !store.Network.IsMainnet() {
		crypto.AddressHRP = "tpc"
		crypto.PublicKeyHRP = "tpublic"
		crypto.PrivateKeyHRP = "tsecret"
		crypto.XPublicKeyHRP = "txpublic"
		crypto.XPrivateKeyHRP = "txsecret"
	}

	w := &Wallet{
		store: store,
		path:  walletPath,
	}

	if !offline {
		serversInfo := map[string][]string{}
		err := json.Unmarshal(serversJSON, &serversInfo)
		if err != nil {
			return nil, err
		}

		var netServers []string
		switch w.store.Network {
		case genesis.Mainnet:
			// mainnet
			netServers = serversInfo["mainnet"]

		case genesis.Testnet:
			// testnet
			netServers = serversInfo["testnet"]

		case genesis.Localnet:
			// localnet
			netServers = []string{"localhost:50052"}

		default:
			return nil, ErrInvalidNetwork
		}

		w.lazyClient = newGRPCClient(netServers)
	}

	return w, nil
}

func (w *Wallet) SetServerAddr(addr string) {
	w.lazyClient = newGRPCClient([]string{addr})
}

func (w *Wallet) Name() string {
	return path.Base(w.path)
}

func (w *Wallet) IsOffline() bool {
	return w.lazyClient == nil
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

// Balance returns balance of the account associated with the address..
func (w *Wallet) Balance(addrStr string) (amount.Amount, error) {
	if w.IsOffline() {
		return 0, ErrOffline
	}

	acc, err := w.lazyClient.getAccount(addrStr)
	if err != nil {
		return 0, err
	}

	return amount.Amount(acc.Balance), nil
}

// Stake returns stake of the validator associated with the address..
func (w *Wallet) Stake(addrStr string) (amount.Amount, error) {
	if w.IsOffline() {
		return 0, ErrOffline
	}

	val, err := w.lazyClient.getValidator(addrStr)
	if err != nil {
		return 0, err
	}

	return amount.Amount(val.Stake), nil
}

// TotalBalance return the total available balance of the wallet.
func (w *Wallet) TotalBalance() amount.Amount {
	totalBalance := int64(0)
	infos := w.store.Vault.AllAccountAddresses()
	for _, info := range infos {
		acc, _ := w.lazyClient.getAccount(info.Address)
		if acc != nil {
			totalBalance += acc.Balance
		}
	}

	return amount.Amount(totalBalance)
}

// MakeTransferTx creates a new transfer transaction based on the given parameters.
func (w *Wallet) MakeTransferTx(sender, receiver string, amt amount.Amount,
	options ...TxOption,
) (*tx.Tx, error) {
	maker, err := newTxBuilder(w.lazyClient, options...)
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
	maker.amount = amt
	maker.typ = payload.TypeTransfer

	return maker.build()
}

// MakeBondTx creates a new bond transaction based on the given parameters.
func (w *Wallet) MakeBondTx(sender, receiver, pubKey string, amt amount.Amount,
	options ...TxOption,
) (*tx.Tx, error) {
	maker, err := newTxBuilder(w.lazyClient, options...)
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
	if pubKey == "" {
		// Let's check if we can get public key from the wallet
		info := w.store.Vault.AddressInfo(receiver)
		if info != nil {
			pubKey = info.PublicKey
		}
	}
	if pubKey != "" {
		maker.pub, err = bls.PublicKeyFromString(pubKey)
		if err != nil {
			return nil, err
		}
	}
	maker.amount = amt
	maker.typ = payload.TypeBond

	return maker.build()
}

// MakeUnbondTx creates a new unbond transaction based on the given parameters.
func (w *Wallet) MakeUnbondTx(addr string, opts ...TxOption) (*tx.Tx, error) {
	maker, err := newTxBuilder(w.lazyClient, opts...)
	if err != nil {
		return nil, err
	}
	err = maker.setFromAddr(addr)
	if err != nil {
		return nil, err
	}
	maker.typ = payload.TypeUnbond

	return maker.build()
}

// MakeWithdrawTx creates a new withdraw transaction based on the given
// parameters.
func (w *Wallet) MakeWithdrawTx(sender, receiver string, amt amount.Amount,
	options ...TxOption,
) (*tx.Tx, error) {
	maker, err := newTxBuilder(w.lazyClient, options...)
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
	maker.amount = amt
	maker.typ = payload.TypeWithdraw

	return maker.build()
}

func (w *Wallet) SignTransaction(password string, trx *tx.Tx) error {
	prv, err := w.PrivateKey(password, trx.Payload().Signer().String())
	if err != nil {
		return err
	}

	sig := prv.Sign(trx.SignBytes())
	trx.SetSignature(sig)
	trx.SetPublicKey(prv.PublicKey())

	return nil
}

func (w *Wallet) BroadcastTransaction(trx *tx.Tx) (string, error) {
	if w.IsOffline() {
		return "", ErrOffline
	}

	id, err := w.lazyClient.sendTx(trx)
	if err != nil {
		return "", err
	}

	d, _ := trx.Bytes()
	w.store.History.addPending(trx.Payload().Signer().String(), trx.Payload().Value(), id, d)

	return id.String(), nil
}

func (w *Wallet) CalculateFee(amt amount.Amount, payloadType payload.Type) (amount.Amount, error) {
	return w.lazyClient.getFee(amt, payloadType)
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

func (w *Wallet) AddressInfos() []vault.AddressInfo {
	return w.store.Vault.AddressInfos()
}

// AddressCount returns the number of addresses inside the wallet.
func (w *Wallet) AddressCount() int {
	return w.store.Vault.AddressCount()
}

func (w *Wallet) AllValidatorAddresses() []vault.AddressInfo {
	return w.store.Vault.AllValidatorAddresses()
}

func (w *Wallet) AllAccountAddresses() []vault.AddressInfo {
	return w.store.Vault.AllAccountAddresses()
}

func (w *Wallet) AddressFromPath(p string) *vault.AddressInfo {
	return w.store.Vault.AddressFromPath(p)
}

func (w *Wallet) ImportPrivateKey(password string, prv *bls.PrivateKey) error {
	return w.store.Vault.ImportPrivateKey(password, prv)
}

func (w *Wallet) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	keys, err := w.store.Vault.PrivateKeys(password, []string{addr})
	if err != nil {
		return nil, err
	}

	return keys[0], nil
}

func (w *Wallet) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	return w.store.Vault.PrivateKeys(password, addrs)
}

// NewBLSAccountAddress create a new BLS-based account address and
// associates it with the given label.
func (w *Wallet) NewBLSAccountAddress(label string) (*vault.AddressInfo, error) {
	return w.store.Vault.NewBLSAccountAddress(label)
}

// NewValidatorAddress creates a new BLS validator address and
// associates it with the given label.
func (w *Wallet) NewValidatorAddress(label string) (*vault.AddressInfo, error) {
	return w.store.Vault.NewValidatorAddress(label)
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

	trxRes, err := w.lazyClient.getTransaction(id)
	if err != nil {
		return err
	}

	var sender string
	var receiver *string
	switch pld := trxRes.Transaction.Payload.(type) {
	case *pactus.TransactionInfo_Transfer:
		sender = pld.Transfer.Sender
		receiver = &pld.Transfer.Receiver
	case *pactus.TransactionInfo_Bond:
		sender = pld.Bond.Sender
		receiver = &pld.Bond.Receiver
		// TODO: complete me!
	// case *pactus.TransactionInfo_Unbond:
	// 	sender = pld.Unbond.Validator
	// 	receiver = nil
	// case *payload.WithdrawPayload:
	// 	sender = pld.Withdraw.From
	// 	receiver = &pld.Withdraw.To
	case *pactus.TransactionInfo_Sortition:
		sender = pld.Sortition.Address
		receiver = nil
	}

	if w.store.Vault.Contains(sender) {
		amt := amount.Amount(-(trxRes.Transaction.Fee + trxRes.Transaction.Value))
		w.store.History.addActivity(sender, amt, trxRes)
	}

	if receiver != nil {
		if w.store.Vault.Contains(*receiver) {
			amt := amount.Amount(trxRes.Transaction.Value)
			w.store.History.addActivity(*receiver, amt, trxRes)
		}
	}

	return nil
}

func (w *Wallet) GetHistory(addr string) []HistoryInfo {
	return w.store.History.getAddrHistory(addr)
}
