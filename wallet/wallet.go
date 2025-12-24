package wallet

import (
	"cmp"
	"context"
	_ "embed"
	"encoding/json"
	"slices"
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
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Wallet struct {
	storage    storage.IStorage
	grpcClient *grpcClient

	addressMap map[string]types.AddressInfo
}

// GenerateMnemonic is a wrapper for `vault.GenerateMnemonic`.
func GenerateMnemonic(entropy int) (string, error) {
	return vault.GenerateMnemonic(entropy)
}

// CheckMnemonic is a wrapper for `vault.CheckMnemonic`.
func CheckMnemonic(mnemonic string) error {
	return vault.CheckMnemonic(mnemonic)
}

// Create creates a wallet from mnemonic (seed phrase) and save it at the
// given path.
func Create(walletPath, mnemonic, password string, chain genesis.ChainType, opts ...OpenWalletOption) (*Wallet, error) {
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
		crypto.ToTestnetHRP()

		coinType = addresspath.CoinTypePactusTestnet
	default:
		return nil, ErrInvalidNetwork
	}

	vlt, err := vault.CreateVaultFromMnemonic(mnemonic, coinType)
	if err != nil {
		return nil, err
	}

	err = vlt.UpdatePassword("", password)
	if err != nil {
		return nil, err
	}

	storage, err := jsonstorage.Create(walletPath, chain, *vlt)
	if err != nil {
		return nil, err
	}

	return openWallet(walletPath, storage, opts...)
}

// Open tries to open a wallet at the given path.
// If the wallet doesn’t exist on this path, it returns an error.
// A wallet can be opened in offline or online modes.
// Offline wallet doesn’t have any connection to any node.
// Online wallet has a connection to one of the pre-defined servers.
func Open(walletPath string, opts ...OpenWalletOption) (*Wallet, error) {
	err := jsonstorage.Upgrade(walletPath)
	if err != nil {
		return nil, err
	}

	storage, err := jsonstorage.Open(walletPath)
	if err != nil {
		return nil, err
	}

	return openWallet(walletPath, storage, opts...)
}

type openWalletConfig struct {
	timeout time.Duration
	servers []string
	offline bool
}

var defaultOpenWalletConfig = openWalletConfig{
	timeout: 5 * time.Second,
	servers: make([]string, 0),
	offline: false,
}

type OpenWalletOption func(*openWalletConfig)

func WithTimeout(timeout time.Duration) OpenWalletOption {
	return func(cfg *openWalletConfig) {
		cfg.timeout = timeout
	}
}

func WithCustomServers(servers []string) OpenWalletOption {
	return func(cfg *openWalletConfig) {
		cfg.servers = servers
	}
}

func WithOfflineMode() OpenWalletOption {
	return func(cfg *openWalletConfig) {
		cfg.offline = true
	}
}

func openWallet(walletPath string, storage storage.IStorage, opts ...OpenWalletOption) (*Wallet, error) {
	cfg := defaultOpenWalletConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	addressMap := make(map[string]types.AddressInfo)
	addressList, err := storage.AllAddresses()
	if err != nil {
		return nil, err
	}

	for _, addrInfo := range addressList {
		addressMap[addrInfo.Address] = addrInfo
	}

	client := newGrpcClient(cfg.timeout, cfg.servers)

	wlt := &Wallet{
		storage:    storage,
		grpcClient: client,
		addressMap: addressMap,
	}

	if !cfg.offline {
		serversData := map[string][]ServerInfo{}
		err := json.Unmarshal(serversJSON, &serversData)
		if err != nil {
			return nil, err
		}

		var netServers []string
		switch wlt.storage.WalletInfo().Network {
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

func (w *Wallet) IsOffline() bool {
	return len(w.grpcClient.servers) == 0
}

func (w *Wallet) Version() int {
	return w.storage.WalletInfo().Version
}

func (w *Wallet) Info() *types.WalletInfo {
	return w.storage.WalletInfo()
}

func (w *Wallet) Path() string {
	return w.storage.WalletInfo().Path
}

func (w *Wallet) IsEncrypted() bool {
	return w.storage.WalletInfo().Encrypted
}

// Neuter clones the wallet and neuters it and saves it at the given path.
func (w *Wallet) Neuter(path string) error {
	cloned, err := w.storage.Clone(path)
	if err != nil {
		return err
	}

	vault := cloned.Vault()
	vault.Neuter()

	return cloned.UpdateVault(vault)
}

// RecoveryAddresses recovers active addresses in the wallet.
func (w *Wallet) RecoveryAddresses(ctx context.Context, password string,
	eventFunc func(addr string),
) error {
	vault := w.storage.Vault()
	//nolint:contextcheck // client manages timeout internally, external context would interfere
	recovered, err := vault.RecoverAddresses(ctx, password, func(addr string) (bool, error) {
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

	for _, info := range recovered {
		err := w.storage.InsertAddress(&info)
		if err != nil {
			return err
		}
	}

	return w.storage.UpdateVault(vault)
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
	infos := w.ListAddresses(OnlyAccountAddresses())
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

	infos := w.ListAddresses(OnlyValidatorAddresses())
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
		info, exists := w.addressMap[receiver]
		if exists {
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
	err = w.storage.AddPending(trx.Payload().Signer().String(), trx.Payload().Value(), txID, data)
	if err != nil {
		return "", err
	}

	return txID.String(), nil
}

func (w *Wallet) UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error {
	vault := w.storage.Vault()
	err := vault.UpdatePassword(oldPassword, newPassword, opts...)
	if err != nil {
		return err
	}

	return w.storage.UpdateVault(vault)
}

func (w *Wallet) AddressInfo(addr string) *types.AddressInfo {
	info, exists := w.addressMap[addr]
	if !exists {
		return nil
	}
	return &info
}

// listAddressConfig contains options for filtering addresses.
type listAddressConfig struct {
	addressTypes []crypto.AddressType
}

var defaultListAddressConfig = listAddressConfig{
	addressTypes: []crypto.AddressType{},
}

// ListAddressOption is a functional option for ListAddresses.
type ListAddressOption func(*listAddressConfig)

// WithAddressTypes filters addresses by the specified type.
func WithAddressTypes(addressTypes []crypto.AddressType) ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = addressTypes
	}
}

// WithAddressType filters addresses by the specified type.
func WithAddressType(addressType crypto.AddressType) ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = []crypto.AddressType{addressType}
	}
}

// OnlyValidatorAddresses filters to show only validator addresses.
func OnlyValidatorAddresses() ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = []crypto.AddressType{crypto.AddressTypeValidator}
	}
}

// OnlyAccountAddresses filters to show only account addresses (BLS and Ed25519).
func OnlyAccountAddresses() ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = []crypto.AddressType{
			crypto.AddressTypeBLSAccount,
			crypto.AddressTypeEd25519Account,
		}
	}
}

func (w *Wallet) ListAddresses(opts ...ListAddressOption) []types.AddressInfo {
	cfg := defaultListAddressConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	infos := make([]types.AddressInfo, 0)
	for _, info := range w.addressMap {
		if len(cfg.addressTypes) == 0 {
			infos = append(infos, info)
			continue
		}

		addr, err := crypto.AddressFromString(info.Address)
		if err != nil {
			return nil
		}

		for _, addrType := range cfg.addressTypes {
			if addr.Type() == addrType {
				infos = append(infos, info)
				break
			}
		}
	}

	w.sortAddressesByAddressIndex(infos...)
	w.sortAddressesByAddressType(infos...)
	w.sortAddressesByPurpose(infos...)

	return infos
}

func (w *Wallet) sortAddressesByPurpose(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.Purpose(), pathB.Purpose())
	})
}

func (w *Wallet) sortAddressesByAddressType(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressType(), pathB.AddressType())
	})
}

func (w *Wallet) sortAddressesByAddressIndex(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressIndex(), pathB.AddressIndex())
	})
}

// AddressCount returns the number of addresses inside the wallet.
func (w *Wallet) AddressCount() int {
	return len(w.addressMap)
}

func (w *Wallet) ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error {
	pub := prv.PublicKeyNative()
	accAddr := pub.AccountAddress()
	if w.HasAddress(accAddr.String()) {
		return ErrAddressExists
	}

	vault := w.storage.Vault()
	accInfo, valInfo, err := vault.ImportBLSPrivateKey(password, prv)
	if err != nil {
		return err
	}

	w.addressMap[accInfo.Address] = *accInfo
	w.addressMap[valInfo.Address] = *valInfo

	err = w.storage.InsertAddress(accInfo)
	if err != nil {
		return err
	}

	err = w.storage.InsertAddress(valInfo)
	if err != nil {
		return err
	}

	return w.storage.UpdateVault(vault)
}

func (w *Wallet) ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) error {
	pub := prv.PublicKeyNative()

	accAddr := pub.AccountAddress()
	if w.HasAddress(accAddr.String()) {
		return ErrAddressExists
	}

	vault := w.storage.Vault()
	accInfo, err := vault.ImportEd25519PrivateKey(password, prv)
	if err != nil {
		return err
	}

	w.addressMap[accInfo.Address] = *accInfo

	err = w.storage.InsertAddress(accInfo)
	if err != nil {
		return err
	}

	return w.storage.UpdateVault(vault)
}

func (w *Wallet) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	keys, err := w.PrivateKeys(password, []string{addr})
	if err != nil {
		return nil, err
	}

	return keys[0], nil
}

func (w *Wallet) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	paths := make([]addresspath.Path, len(addrs))
	for i, addr := range addrs {
		info := w.AddressInfo(addr)
		if info == nil {
			return nil, NewErrAddressNotFound(addr)
		}

		hdPath, err := addresspath.FromString(info.Path)
		if err != nil {
			return nil, err
		}

		paths[i] = hdPath
	}

	return w.storage.Vault().PrivateKeys(password, paths)
}

// newAddressConfig contains options for creating new addresses.
type newAddressConfig struct {
	password string
}

var defaultNewAddressConfig = newAddressConfig{
	password: "",
}

// NewAddressOption is a functional option for NewAddress.
type NewAddressOption func(*newAddressConfig)

// WithPassword sets the password for address creation required for Ed25519 accounts.
func WithPassword(password string) NewAddressOption {
	return func(cfg *newAddressConfig) {
		cfg.password = password
	}
}

func (w *Wallet) NewAddress(addressType crypto.AddressType, label string, opts ...NewAddressOption,
) (*types.AddressInfo, error) {
	cfg := defaultNewAddressConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	vault := w.storage.Vault()
	var info *types.AddressInfo
	var err error
	switch addressType {
	case crypto.AddressTypeValidator:
		info, err = vault.NewValidatorAddress(label)
	case crypto.AddressTypeBLSAccount:
		info, err = vault.NewBLSAccountAddress(label)
	case crypto.AddressTypeEd25519Account:
		info, err = vault.NewEd25519AccountAddress(label, cfg.password)
	case crypto.AddressTypeTreasury:
		return nil, ErrInvalidAddressType

	default:
		return nil, ErrInvalidAddressType
	}

	if err != nil {
		return nil, err
	}

	w.addressMap[info.Address] = *info

	err = w.storage.InsertAddress(info)
	if err != nil {
		return nil, err
	}

	err = w.storage.UpdateVault(vault)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// NewBLSAccountAddress create a new BLS-based account address and
// associates it with the given label.
func (w *Wallet) NewBLSAccountAddress(label string) (*types.AddressInfo, error) {
	return w.NewAddress(crypto.AddressTypeBLSAccount, label)
}

// NewEd25519AccountAddress create a new Ed25519-based account address and
// associates it with the given label.
// The password is required to access the master private key needed for address generation.
func (w *Wallet) NewEd25519AccountAddress(label, password string) (*types.AddressInfo, error) {
	return w.NewAddress(crypto.AddressTypeEd25519Account, label, WithPassword(password))
}

// NewValidatorAddress creates a new BLS validator address and
// associates it with the given label.
func (w *Wallet) NewValidatorAddress(label string) (*types.AddressInfo, error) {
	return w.NewAddress(crypto.AddressTypeValidator, label)
}

func (w *Wallet) HasAddress(addr string) bool {
	_, exists := w.addressMap[addr]

	return exists
}

func (w *Wallet) Mnemonic(password string) (string, error) {
	return w.storage.Vault().Mnemonic(password)
}

// AddressLabel returns label of the given address.
func (w *Wallet) AddressLabel(addr string) string {
	info, exists := w.addressMap[addr]
	if !exists {
		return ""
	}

	return info.Label
}

// SetAddressLabel updates the label of the given address.
func (w *Wallet) SetAddressLabel(addr, label string) error {
	info, exists := w.addressMap[addr]
	if !exists {
		return NewErrAddressNotFound(addr)
	}

	info.Label = label
	w.addressMap[addr] = info

	return w.storage.UpdateAddress(&info)
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

	if w.HasAddress(sender) {
		amt := amount.Amount(-(trxRes.Transaction.Fee + trxRes.Transaction.Value))
		w.storage.AddActivity(sender, amt, trxRes)
	}

	if receiver != nil {
		if w.HasAddress(*receiver) {
			amt := amount.Amount(trxRes.Transaction.Value)
			err := w.storage.AddActivity(*receiver, amt, trxRes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Wallet) History(addr string) []types.HistoryInfo {
	return w.storage.GetAddrHistory(addr)
}

func (w *Wallet) SignMessage(password, addr, msg string) (string, error) {
	prv, err := w.PrivateKey(password, addr)
	if err != nil {
		return "", err
	}

	return prv.Sign([]byte(msg)).String(), nil
}

// makeTxBuilder initializes a txBuilder with provided options, allowing for flexible configuration of the transaction.
func (w *Wallet) makeTxBuilder(options ...TxOption) (*txBuilder, error) {
	builder := &txBuilder{
		client: w.grpcClient,
		fee:    w.storage.WalletInfo().DefaultFee,
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
	return w.storage.SetDefaultFee(fee)
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
