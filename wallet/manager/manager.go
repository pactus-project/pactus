package manager

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/provider"
	"github.com/pactus-project/pactus/wallet/types"
)

var _ WalletManager = (*Manager)(nil)

type Manager struct {
	// lk protects the wallets map: the wallet service adds wallets
	// while other calls and the event pipeline look them up.
	lk sync.RWMutex

	ctx               context.Context
	wallets           map[string]*wallet.Wallet
	lockMode          bool
	chainType         genesis.ChainType
	walletDirectory   string
	DefaultWalletName string
	provider          provider.WalletProvider
	eventPipe         pipeline.Pipeline[any]
}

func NewManager(ctx context.Context, conf *Config,
	provider provider.WalletProvider,
	eventPipe pipeline.Pipeline[any],
) (WalletManager, error) {
	wallets := make(map[string]*wallet.Wallet)
	files, err := util.ListFilesInDir(conf.WalletsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		opts := []wallet.OpenWalletOption{
			wallet.WithBlockchainProvider(provider),
			wallet.WithLockMode(conf.LockMode),
		}
		wlt, err := wallet.Open(ctx, file, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to open wallet %s: %w", file, err)
		}

		wallets[filepath.Base(file)] = wlt
	}

	mgr := &Manager{
		ctx:             ctx,
		wallets:         wallets,
		lockMode:        conf.LockMode,
		chainType:       conf.ChainType,
		walletDirectory: conf.WalletsDir,
		provider:        provider,
		eventPipe:       eventPipe,
	}

	// The receivers of a pipeline can't be registered once it is running.
	// The manager registers one receiver and passes the events to the wallets,
	// including the wallets that are created later.
	if eventPipe != nil {
		eventPipe.RegisterReceiver(mgr.processEvent)
	}

	return mgr, nil
}

func (wm *Manager) processEvent(event any) {
	wm.lk.RLock()
	wallets := make([]*wallet.Wallet, 0, len(wm.wallets))
	for _, wlt := range wm.wallets {
		wallets = append(wallets, wlt)
	}
	wm.lk.RUnlock()

	for _, wlt := range wallets {
		wlt.ProcessEvent(event)
	}
}

func (wm *Manager) getWallet(walletName string) (*wallet.Wallet, bool) {
	wm.lk.RLock()
	defer wm.lk.RUnlock()

	wlt, ok := wm.wallets[walletName]

	return wlt, ok
}

func (*Manager) Start() error {
	return nil
}

func (wm *Manager) Stop() {
	wm.lk.Lock()
	defer wm.lk.Unlock()

	for _, wlt := range wm.wallets {
		wlt.Close()
	}
}

func (wm *Manager) getWalletPath(walletName string) (string, error) {
	return util.SecureJoinPath(wm.walletDirectory, walletName)
}

func (wm *Manager) createWalletWithMnemonic(
	walletName, mnemonic, password string,
) error {
	walletPath, err := wm.getWalletPath(walletName)
	if err != nil {
		return err
	}

	wm.lk.Lock()
	defer wm.lk.Unlock()

	if isExists := util.PathExists(walletPath); isExists {
		return ErrWalletAlreadyExists
	}

	wlt, err := wallet.Create(wm.ctx, walletPath, mnemonic, password, wm.chainType,
		wallet.WithBlockchainProvider(wm.provider),
		wallet.WithLockMode(wm.lockMode))
	if err != nil {
		return err
	}

	wm.wallets[walletName] = wlt

	return nil
}

// Deprecated: Move it to the utils service.
func (*Manager) GetValidatorAddress(
	publicKey string,
) (string, error) {
	pubKey, err := bls.PublicKeyFromString(publicKey)
	if err != nil {
		return "", err
	}

	return pubKey.ValidatorAddress().String(), nil
}

func (wm *Manager) CreateWallet(
	walletName, password string,
) (string, error) {
	mnemonic, err := wallet.GenerateMnemonic(128)
	if err != nil {
		return "", err
	}

	if err := wm.createWalletWithMnemonic(walletName, mnemonic, password); err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (wm *Manager) RestoreWallet(walletName, mnemonic, password string) error {
	return wm.createWalletWithMnemonic(walletName, mnemonic, password)
}

func (wm *Manager) MigrateWallet(walletName string) error {
	walletPath, err := wm.getWalletPath(walletName)
	if err != nil {
		return err
	}

	wm.lk.Lock()
	defer wm.lk.Unlock()

	if err := wallet.Migrate(wm.ctx, walletPath); err != nil {
		return err
	}

	// Reload the migrated wallet and replace the previous in-memory instance.
	wlt, err := wallet.Open(wm.ctx, walletPath,
		wallet.WithBlockchainProvider(wm.provider),
		wallet.WithLockMode(wm.lockMode))
	if err != nil {
		return err
	}

	wm.wallets[walletName] = wlt

	return nil
}

func (wm *Manager) NewAddress(walletName string, addressType crypto.AddressType, label string,
	opts ...wallet.NewAddressOption,
) (*types.AddressInfo, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.NewAddress(addressType, label, opts...)
}

func (wm *Manager) PrivateKey(walletName, password, addr string) (crypto.PrivateKey, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.PrivateKey(password, addr)
}

func (wm *Manager) Mnemonic(walletName, password string) (string, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.Mnemonic(password)
}

func (wm *Manager) AddressLabel(walletName, addr string) (string, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.AddressLabel(addr), nil
}

func (wm *Manager) SetAddressLabel(walletName, addr, label string) error {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.SetAddressLabel(addr, label)
}

func (wm *Manager) Balance(walletName, addr string) (amount.Amount, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.Balance(addr)
}

func (wm *Manager) Stake(walletName, addr string) (amount.Amount, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.Stake(addr)
}

func (wm *Manager) SetDefaultFee(walletName string, fee amount.Amount) error {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.SetDefaultFee(fee)
}

func (wm *Manager) MakeTransferTx(
	walletName, sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeTransferTx(sender, receiver, amt, opts...)
}

func (wm *Manager) MakeBondTx(
	walletName, sender, receiver, publicKey string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeBondTx(sender, receiver, publicKey, amt, opts...)
}

func (wm *Manager) MakeUnbondTx(walletName, validatorAddr string, opts ...wallet.TxOption) (*tx.Tx, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeUnbondTx(validatorAddr, opts...)
}

func (wm *Manager) MakeWithdrawTx(
	walletName, sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeWithdrawTx(sender, receiver, amt, opts...)
}

func (wm *Manager) SignTransaction(walletName, password string, trx *tx.Tx) error {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.SignTransaction(password, trx)
}

func (wm *Manager) BroadcastTransaction(walletName string, trx *tx.Tx) (string, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.BroadcastTransaction(trx)
}

func (wm *Manager) TotalBalance(walletName string) (amount.Amount, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.TotalBalance()
}

func (wm *Manager) TotalStake(walletName string) (amount.Amount, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.TotalStake()
}

func (wm *Manager) SignRawTransaction(
	walletName, password string, rawTx []byte,
) (txID, data []byte, err error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, nil, ErrWalletNotLoaded
	}

	trx, err := tx.FromBytes(rawTx)
	if err != nil {
		return nil, nil, err
	}

	if err := wlt.SignTransaction(password, trx); err != nil {
		return nil, nil, err
	}

	data, err = trx.Bytes()
	if err != nil {
		return nil, nil, err
	}

	return trx.ID().Bytes(), data, nil
}

func (wm *Manager) SignMessage(walletName, password, addr, msg string) (string, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.SignMessage(password, addr, msg)
}

func (wm *Manager) AddressInfo(walletName, address string) (*types.AddressInfo, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.AddressInfo(address)
}

func (wm *Manager) WalletInfo(walletName string) (*types.WalletInfo, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.Info(), nil
}

func (wm *Manager) ListWallets() ([]string, error) {
	wm.lk.RLock()
	defer wm.lk.RUnlock()

	wallets := make([]string, 0, len(wm.wallets))
	for name := range wm.wallets {
		wallets = append(wallets, name)
	}

	return wallets, nil
}

func (wm *Manager) ListAddresses(walletName string, opts ...wallet.ListAddressOption) ([]*types.AddressInfo, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.ListAddresses(opts...), nil
}

func (wm *Manager) ListTransactions(walletName string,
	opts ...wallet.ListTransactionsOption,
) ([]*types.TransactionInfo, error) {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.ListTransactions(opts...), nil
}

func (wm *Manager) UpdatePassword(walletName, oldPassword, newPassword string) error {
	wlt, ok := wm.getWallet(walletName)
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.UpdatePassword(oldPassword, newPassword)
}
