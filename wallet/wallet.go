package wallet

import (
	"context"
	_ "embed"
	"encoding/json"
	"path"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/pactus-project/pactus/wallet/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Address string `json:"address"`
}

type HistoryInfo struct {
	TxID        string
	Time        *time.Time
	PayloadType string
	Desc        string
	Amount      amount.Amount
}

type Wallet struct {
	storage    storage.IStorage
	path       string
	grpcClient *grpcClient
}

type Info struct {
	WalletName string
	Version    int
	Network    string
	DefaultFee amount.Amount
	UUID       string
	Encrypted  bool
	CreatedAt  time.Time
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
func Open(walletPath string, offline bool, options ...Option) (*Wallet, error) {
	storage, err := jsonstorage.Open(walletPath)
	if err != nil {
		return nil, err
	}

	err = storage.Upgrade()
	if err != nil {
		return nil, err
	}

	opts := defaultWalletOpt
	for _, opt := range options {
		opt(opts)
	}

	return newWallet(walletPath, storage, offline, opts)
}

// Create creates a wallet from mnemonic (seed phrase) and save it at the
// given path.
func Create(walletPath, mnemonic, password string, chain genesis.ChainType,
	options ...Option,
) (*Wallet, error) {
	opts := defaultWalletOpt

	for _, opt := range options {
		opt(opts)
	}

	walletPath = util.MakeAbs(walletPath)
	if util.PathExists(walletPath) {
		return nil, ExitsError{
			Path: walletPath,
		}
	}

	var coinType addresspath.CoinType
	switch chain {
	case genesis.Mainnet:
		coinType = addresspath.CoinTypePactusMainnet
	case genesis.Testnet, genesis.Localnet:
		coinType = addresspath.CoinTypePactusTestnet
	default:
		return nil, ErrInvalidNetwork
	}

	vlt, err := vault.CreateVaultFromMnemonic(mnemonic, coinType)
	if err != nil {
		return nil, err
	}

	storage, _ := jsonstorage.Create(walletPath, version.Version5, chain, *vlt)
	wallet, err := newWallet(walletPath, storage, false, opts)
	if err != nil {
		return nil, err
	}

	err = storage.UpdatePassword("", password)
	if err != nil {
		return nil, err
	}

	err = wallet.save()
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func newWallet(walletPath string, storage storage.IStorage, offline bool, option *walletOpt) (*Wallet, error) {
	client := newGrpcClient(option.timeout, option.servers)

	wlt := &Wallet{
		storage:    storage,
		path:       walletPath,
		grpcClient: client,
	}

	if !offline {
		serversData := map[string][]ServerInfo{}
		err := json.Unmarshal(serversJSON, &serversData)
		if err != nil {
			return nil, err
		}

		var netServers []string
		switch wlt.storage.Network() {
		case genesis.Mainnet:
			for _, srv := range serversData["mainnet"] {
				netServers = append(netServers, srv.Address)
			}

		case genesis.Testnet:
			crypto.ToTestnetHRP()

			for _, srv := range serversData["testnet"] {
				netServers = append(netServers, srv.Address)
			}

		case genesis.Localnet:
			crypto.ToTestnetHRP()

			netServers = []string{"localhost:50052"}

		default:
			return nil, ErrInvalidNetwork
		}

		util.Shuffle(netServers)

		if client.servers == nil {
			client.servers = netServers
		}
	}

	return wlt, nil
}

func (w *Wallet) Name() string {
	return path.Base(w.path)
}

func (w *Wallet) CoinType() addresspath.CoinType {
	return w.storage.CoinType()
}

func (w *Wallet) IsOffline() bool {
	return len(w.grpcClient.servers) == 0
}

func (w *Wallet) Path() string {
	return w.path
}

func (w *Wallet) save() error {
	return w.storage.Save()
}

// RecoveryAddresses recovers active addresses in the wallet.
func (w *Wallet) RecoveryAddresses(ctx context.Context, password string,
	eventFunc func(addr string),
) error {
	//nolint:contextcheck // client manages timeout internally, external context would interfere
	err := w.storage.RecoverAddresses(ctx, password, func(addr string) (bool, error) {
		_, err := w.grpcClient.getAccount(addr)
		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.NotFound {
				return false, nil
			}

			return false, err
		}

		if eventFunc != nil {
			eventFunc(addr)
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	return w.save()
}

// Balance returns balance of the account associated with the address..
func (w *Wallet) Balance(addrStr string) (amount.Amount, error) {
	acc, err := w.grpcClient.getAccount(addrStr)
	if err != nil {
		return 0, err
	}

	return amount.Amount(acc.Balance), nil
}

// Stake returns stake of the validator associated with the address..
func (w *Wallet) Stake(addrStr string) (amount.Amount, error) {
	val, err := w.grpcClient.getValidator(addrStr)
	if err != nil {
		return 0, err
	}

	return amount.Amount(val.Stake), nil
}

// TotalBalance return the total available balance of the wallet.
func (w *Wallet) TotalBalance() (amount.Amount, error) {
	totalBalance := int64(0)
	infos := w.storage.ListAccountAddresses()
	for _, info := range infos {
		acc, _ := w.grpcClient.getAccount(info.Address)
		if acc != nil {
			totalBalance += acc.Balance
		}
	}

	return amount.Amount(totalBalance), nil
}

// TotalStake return total available stake of the wallet.
func (w *Wallet) TotalStake() (amount.Amount, error) {
	totalStake := int64(0)

	infos := w.storage.ListValidatorAddresses()
	for _, info := range infos {
		val, _ := w.grpcClient.getValidator(info.Address)
		if val != nil {
			totalStake += val.Stake
		}
	}

	return amount.Amount(totalStake), nil
}

// MakeTransferTx creates a new transfer transaction based on the given parameters.
func (w *Wallet) MakeTransferTx(sender, receiver string, amt amount.Amount,
	options ...TxOption,
) (*tx.Tx, error) {
	maker, err := w.makeTxBuilder(options...)
	if err != nil {
		return nil, err
	}
	err = maker.setSenderAddr(sender)
	if err != nil {
		return nil, err
	}
	err = maker.setReceiverAddress(receiver)
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
	maker, err := w.makeTxBuilder(options...)
	if err != nil {
		return nil, err
	}
	err = maker.setSenderAddr(sender)
	if err != nil {
		return nil, err
	}
	err = maker.setReceiverAddress(receiver)
	if err != nil {
		return nil, err
	}
	if pubKey == "" {
		// Let's check if we can get public key from the wallet
		info := w.storage.AddressInfo(receiver)
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
	maker, err := w.makeTxBuilder(opts...)
	if err != nil {
		return nil, err
	}
	err = maker.setSenderAddr(addr)
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
	maker, err := w.makeTxBuilder(options...)
	if err != nil {
		return nil, err
	}
	err = maker.setSenderAddr(sender)
	if err != nil {
		return nil, err
	}
	err = maker.setReceiverAddress(receiver)
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
	txID, err := w.grpcClient.sendTx(trx)
	if err != nil {
		return "", err
	}

	data, _ := trx.Bytes()
	w.storage.AddPending(trx.Payload().Signer().String(), trx.Payload().Value(), txID, data)

	err = w.save()
	if err != nil {
		return "", err
	}

	return txID.String(), nil
}

func (w *Wallet) UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error {
	err := w.storage.UpdatePassword(oldPassword, newPassword, opts...)
	if err != nil {
		return err
	}

	return w.save()
}

func (w *Wallet) IsEncrypted() bool {
	return w.storage.IsEncrypted()
}

func (w *Wallet) AddressInfo(addr string) *storage.AddressInfo {
	return w.storage.AddressInfo(addr)
}

func (w *Wallet) ListAddresses() []storage.AddressInfo {
	return w.storage.ListAddresses()
}

// AddressCount returns the number of addresses inside the wallet.
func (w *Wallet) AddressCount() int {
	return w.storage.AddressCount()
}

func (w *Wallet) ListValidatorAddresses() []storage.AddressInfo {
	return w.storage.ListValidatorAddresses()
}

func (w *Wallet) ListAccountAddresses() []storage.AddressInfo {
	return w.storage.ListAccountAddresses()
}

func (w *Wallet) AddressByPath(p string) *storage.AddressInfo {
	return w.storage.AddressByPath(p)
}

func (w *Wallet) ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error {
	err := w.storage.ImportBLSPrivateKey(password, prv)
	if err != nil {
		return err
	}

	return w.save()
}

func (w *Wallet) ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) error {
	err := w.storage.ImportEd25519PrivateKey(password, prv)
	if err != nil {
		return err
	}

	return w.save()
}

func (w *Wallet) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	keys, err := w.storage.PrivateKeys(password, []string{addr})
	if err != nil {
		return nil, err
	}

	return keys[0], nil
}

func (w *Wallet) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	return w.storage.PrivateKeys(password, addrs)
}

type addressBuilder struct {
	password string
}

type NewAddressOption func(*addressBuilder)

func WithPassword(password string) NewAddressOption {
	return func(opt *addressBuilder) {
		opt.password = password
	}
}

func (w *Wallet) NewAddress(addressType crypto.AddressType, label string,
	opts ...NewAddressOption,
) (*storage.AddressInfo, error) {
	builder := &addressBuilder{}

	for _, opt := range opts {
		opt(builder)
	}

	var info *storage.AddressInfo
	var err error
	switch addressType {
	case crypto.AddressTypeValidator:
		info, err = w.storage.NewValidatorAddress(label)
	case crypto.AddressTypeBLSAccount:
		info, err = w.storage.NewBLSAccountAddress(label)
	case crypto.AddressTypeEd25519Account:
		info, err = w.storage.NewEd25519AccountAddress(label, builder.password)
	case crypto.AddressTypeTreasury:
		return nil, jsonstorage.ErrInvalidAddressType

	default:
		return nil, jsonstorage.ErrInvalidAddressType
	}

	if err != nil {
		return nil, err
	}

	err = w.save()
	if err != nil {
		return nil, err
	}

	return info, nil
}

// NewBLSAccountAddress create a new BLS-based account address and
// associates it with the given label.
func (w *Wallet) NewBLSAccountAddress(label string) (*storage.AddressInfo, error) {
	return w.NewAddress(crypto.AddressTypeBLSAccount, label)
}

// NewEd25519AccountAddress create a new Ed25519-based account address and
// associates it with the given label.
// The password is required to access the master private key needed for address generation.
func (w *Wallet) NewEd25519AccountAddress(label, password string) (*storage.AddressInfo, error) {
	return w.NewAddress(crypto.AddressTypeEd25519Account, label, WithPassword(password))
}

// NewValidatorAddress creates a new BLS validator address and
// associates it with the given label.
func (w *Wallet) NewValidatorAddress(label string) (*storage.AddressInfo, error) {
	return w.NewAddress(crypto.AddressTypeValidator, label)
}

func (w *Wallet) HasAddress(addr string) bool {
	return w.storage.HasAddress(addr)
}

func (w *Wallet) Mnemonic(password string) (string, error) {
	return w.storage.Mnemonic(password)
}

// AddressLabel returns label of the given address.
func (w *Wallet) AddressLabel(addr string) string {
	return w.storage.AddressLabel(addr)
}

// SetAddressLabel updates the label of the given address.
func (w *Wallet) SetAddressLabel(addr, label string) error {
	err := w.storage.SetAddressLabel(addr, label)
	if err != nil {
		return err
	}

	return w.save()
}

func (w *Wallet) AddTransaction(txID tx.ID) error {
	idStr := txID.String()
	if w.storage.HasTransaction(idStr) {
		return jsonstorage.ErrHistoryExists
	}

	trxRes, err := w.grpcClient.getTransaction(txID)
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

	if w.storage.HasAddress(sender) {
		amt := amount.Amount(-(trxRes.Transaction.Fee + trxRes.Transaction.Value))
		w.storage.AddActivity(sender, amt, trxRes)
	}

	if receiver != nil {
		if w.storage.HasAddress(*receiver) {
			amt := amount.Amount(trxRes.Transaction.Value)
			w.storage.AddActivity(*receiver, amt, trxRes)
		}
	}

	return w.save()
}

func (w *Wallet) History(addr string) []storage.HistoryInfo {
	return w.storage.GetAddrHistory(addr)
}

func (w *Wallet) SignMessage(password, addr, msg string) (string, error) {
	prv, err := w.PrivateKey(password, addr)
	if err != nil {
		return "", err
	}

	return prv.Sign([]byte(msg)).String(), nil
}

func (w *Wallet) Version() int {
	return w.storage.Version()
}

func (w *Wallet) CreationTime() time.Time {
	return w.storage.CreatedAt()
}

func (w *Wallet) Network() genesis.ChainType {
	return w.storage.Network()
}

func (w *Wallet) Info() *Info {
	return &Info{
		WalletName: w.Name(),
		Version:    w.storage.Version(),
		Network:    w.storage.Network().String(),
		DefaultFee: w.storage.DefaultFee(),
		UUID:       w.storage.UUID().String(),
		Encrypted:  w.IsEncrypted(),
		CreatedAt:  w.storage.CreatedAt(),
	}
}

// Neuter clones the wallet and neuters it and saves it at the given path.
func (w *Wallet) Neuter(path string) error {
	return w.storage.Neuter(path)
}

// makeTxBuilder initializes a txBuilder with provided options, allowing for flexible configuration of the transaction.
func (w *Wallet) makeTxBuilder(options ...TxOption) (*txBuilder, error) {
	builder := &txBuilder{
		client: w.grpcClient,
		fee:    w.storage.DefaultFee(),
	}
	for _, op := range options {
		err := op(builder)
		if err != nil {
			return nil, err
		}
	}

	return builder, nil
}

func (w *Wallet) SetDefaultFee(fee amount.Amount) error {
	w.storage.SetDefaultFee(fee)

	return w.save()
}

func GetServerList(network string) ([]ServerInfo, error) {
	// Default to mainnet if network is empty
	if network == "" {
		network = "mainnet"
	}

	serversData := map[string][]ServerInfo{}
	err := json.Unmarshal(serversJSON, &serversData)
	if err != nil {
		return nil, err
	}

	servers, exists := serversData[network]
	if !exists {
		return []ServerInfo{}, nil
	}

	return servers, nil
}
