package wallet

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
	"github.com/pactus-project/pactus/wallet/storage/sqlitestorage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Wallet struct {
	addresses
	transactions

	storage storage.IStorage
	// grpcClient   *grpcClient
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
func Create(ctx context.Context, walletPath, mnemonic, password string,
	chain genesis.ChainType, opts ...OpenWalletOption,
) (*Wallet, error) {
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

	storage, err := sqlitestorage.Create(ctx, walletPath, chain, vlt)
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

func openWallet(_ string, storage storage.IStorage, opts ...OpenWalletOption) (*Wallet, error) {
	cfg := defaultOpenWalletConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	client := newGrpcClient(cfg.timeout, cfg.servers)

	addresses := newAddresses(storage)
	transactions := newTransactions(storage, client)

	wlt := &Wallet{
		addresses:    addresses,
		transactions: transactions,
		storage:      storage,
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
	res, err := w.grpcClient.getAccount(addrStr)
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.Account.Balance), nil
}

// Stake returns stake of the validator associated with the address..
func (w *Wallet) Stake(addrStr string) (amount.Amount, error) {
	res, err := w.grpcClient.getValidator(addrStr)
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.Validator.Stake), nil
}

// TotalBalance return the total available balance of the wallet.
func (w *Wallet) TotalBalance() (amount.Amount, error) {
	totalBalance := int64(0)
	infos := w.ListAddresses(OnlyAccountAddresses())
	for _, info := range infos {
		res, _ := w.grpcClient.getAccount(info.Address)
		if res != nil {
			totalBalance += res.Account.Balance
		}
	}

	return amount.Amount(totalBalance), nil
}

// TotalStake return total available stake of the wallet.
func (w *Wallet) TotalStake() (amount.Amount, error) {
	totalStake := int64(0)

	infos := w.ListAddresses(OnlyValidatorAddresses())
	for _, info := range infos {
		res, _ := w.grpcClient.getValidator(info.Address)
		if res != nil {
			totalStake += res.Validator.Stake
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
		info, _ := w.AddressInfo(receiver)
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

	txInfos, _ := types.MakeTransactionInfos(trx)
	for _, info := range txInfos {
		_ = w.storage.InsertTransaction(info)
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

func (w *Wallet) Mnemonic(password string) (string, error) {
	return w.storage.Vault().Mnemonic(password)
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
