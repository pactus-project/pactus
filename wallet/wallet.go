package wallet

import (
	"context"

	"github.com/ezex-io/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/provider"
	offlineprovider "github.com/pactus-project/pactus/wallet/provider/offline"
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

	provider provider.IBlockchainProvider
	storage  storage.IStorage
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

	coinType := addresspath.CoinTypePactusMainnet
	if chain != genesis.Mainnet {
		coinType = addresspath.CoinTypePactusTestnet
		crypto.ToTestnetHRP()
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

	return New(storage, opts...)
}

// Open tries to open a wallet at the given path.
// It first tries the SQLite backend; if that fails, it falls back to the legacy JSON wallet format.
// A wallet can be opened in offline or online modes.
// Offline wallet doesn’t have any connection to any node.
// Online wallet has a connection to one of the pre-defined servers.
func Open(ctx context.Context, walletPath string, opts ...OpenWalletOption) (*Wallet, error) {
	sqliteStrg, err := sqlitestorage.Open(ctx, walletPath)
	if err == nil {
		return New(sqliteStrg, opts...)
	}

	// Fallback to JSON storage for legacy wallets
	if err := jsonstorage.Upgrade(walletPath); err != nil {
		return nil, err
	}

	jsonStrg, err := jsonstorage.Open(walletPath)
	if err != nil {
		return nil, err
	}

	return New(jsonStrg, opts...)
}

type openWalletConfig struct {
	eventPipe pipeline.Pipeline[any]
	provider  provider.IBlockchainProvider
}

var defaultOpenWalletConfig = openWalletConfig{
	eventPipe: nil,
	provider:  offlineprovider.NewOfflineBlockchainProvider(),
}

type OpenWalletOption func(*openWalletConfig)

func WithEventPipe(eventPipe pipeline.Pipeline[any]) OpenWalletOption {
	return func(cfg *openWalletConfig) {
		cfg.eventPipe = eventPipe
	}
}

func WithBlockchainProvider(provider provider.IBlockchainProvider) OpenWalletOption {
	return func(cfg *openWalletConfig) {
		cfg.provider = provider
	}
}

func New(storage storage.IStorage, opts ...OpenWalletOption) (*Wallet, error) {
	if storage.Vault().CoinType != addresspath.CoinTypePactusMainnet {
		crypto.ToTestnetHRP()
	}

	cfg := defaultOpenWalletConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	wlt := &Wallet{
		addresses:    newAddresses(storage),
		transactions: newTransactions(storage, cfg.provider),
		provider:     cfg.provider,
		storage:      storage,
	}

	if cfg.eventPipe != nil {
		cfg.eventPipe.RegisterReceiver(wlt.transactions.processEvent)
	}

	return wlt, nil
}

func (w *Wallet) Close() {
	if err := w.provider.Close(); err != nil {
		logger.Warn("failed to close provider", "error", err)
	}

	if err := w.storage.Close(); err != nil {
		logger.Warn("failed to close storage", "error", err)
	}
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
	recovered, err := vault.RecoverAddresses(ctx, password, func(addr string) (bool, error) {
		_, err := w.provider.GetAccount(addr)
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
	acc, err := w.provider.GetAccount(addrStr)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}

// Stake returns stake of the validator associated with the address..
func (w *Wallet) Stake(addrStr string) (amount.Amount, error) {
	val, err := w.provider.GetValidator(addrStr)
	if err != nil {
		return 0, err
	}

	return val.Stake(), nil
}

// TotalBalance return the total available balance of the wallet.
func (w *Wallet) TotalBalance() (amount.Amount, error) {
	totalBalance := amount.Amount(0)
	infos := w.ListAddresses(OnlyAccountAddresses())
	for _, info := range infos {
		acc, err := w.provider.GetAccount(info.Address)
		if err == nil {
			totalBalance += acc.Balance()
		}
	}

	return totalBalance, nil
}

// TotalStake return total available stake of the wallet.
func (w *Wallet) TotalStake() (amount.Amount, error) {
	totalStake := amount.Amount(0)

	infos := w.ListAddresses(OnlyValidatorAddresses())
	for _, info := range infos {
		val, err := w.provider.GetValidator(info.Address)
		if err == nil {
			totalStake += val.Stake()
		}
	}

	return totalStake, nil
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
		err = maker.setPublicKey(pubKey)
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
	hash, err := w.provider.SendTx(trx)
	if err != nil {
		return "", err
	}

	txInfos, _ := types.MakeTransactionInfos(trx, types.TransactionStatusPending, 0)
	for _, info := range txInfos {
		info.Direction = types.TxDirectionOutgoing
		err := w.storage.InsertTransaction(info)
		if err != nil {
			logger.Warn("transaction broadcasted but not recorded in wallet storage", "error", err)
		}
	}

	return hash, nil
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
		provider: w.provider,
		fee:      w.storage.WalletInfo().DefaultFee,
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

// SetProvider sets the blockchain provider for the wallet.
func (w *Wallet) SetProvider(provider provider.IBlockchainProvider) {
	w.provider = provider
	w.transactions.provider = provider
}
