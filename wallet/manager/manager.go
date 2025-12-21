package manager

import (
	"fmt"
	"path/filepath"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/storage"
)

var _ IManager = (*Manager)(nil)

type Manager struct {
	wallets         map[string]*wallet.Wallet
	chainType       genesis.ChainType
	walletDirectory string
}

func NewManager(conf *Config) IManager {
	return &Manager{
		wallets:         make(map[string]*wallet.Wallet),
		chainType:       conf.ChainType,
		walletDirectory: conf.WalletsDir,
	}
}

func (wm *Manager) getWalletPath(walletName string) string {
	return util.MakeAbs(filepath.Join(wm.walletDirectory, walletName))
}

func (wm *Manager) WalletPath(walletName string) string {
	return wm.getWalletPath(walletName)
}

func (wm *Manager) createWalletWithMnemonic(
	walletName, mnemonic, password string,
) error {
	walletPath := wm.getWalletPath(walletName)
	if isExists := util.PathExists(walletPath); isExists {
		return ErrWalletAlreadyExists
	}

	_, err := wallet.Create(walletPath, mnemonic, password, wm.chainType)
	if err != nil {
		return err
	}

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

func (wm *Manager) LoadWallet(walletName, serverAddr string) error {
	if _, ok := wm.wallets[walletName]; ok {
		return ErrWalletAlreadyLoaded
	}

	walletPath := wm.getWalletPath(walletName)
	wlt, err := wallet.Open(walletPath, true, wallet.WithCustomServers([]string{serverAddr}))
	if err != nil {
		return err
	}

	wm.wallets[walletName] = wlt

	return nil
}

func (wm *Manager) NewAddress(
	walletName string,
	addressType crypto.AddressType,
	label string,
	opts ...wallet.NewAddressOption,
) (*storage.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.NewAddress(addressType, label, opts...)
}

func (wm *Manager) PrivateKey(walletName, password, addr string) (crypto.PrivateKey, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.PrivateKey(password, addr)
}

func (wm *Manager) Mnemonic(walletName, password string) (string, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.Mnemonic(password)
}

func (wm *Manager) ListAccountAddresses(walletName string) ([]storage.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.ListAccountAddresses(), nil
}

func (wm *Manager) ListValidatorAddresses(walletName string) ([]storage.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.ListValidatorAddresses(), nil
}

func (wm *Manager) AddressLabel(walletName, addr string) (string, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.AddressLabel(addr), nil
}

func (wm *Manager) SetAddressLabel(walletName, addr, label string) error {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.SetAddressLabel(addr, label)
}

func (wm *Manager) Balance(walletName, addr string) (amount.Amount, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.Balance(addr)
}

func (wm *Manager) Stake(walletName, addr string) (amount.Amount, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.Stake(addr)
}

func (wm *Manager) SetDefaultFee(walletName string, fee amount.Amount) error {
	wlt, ok := wm.wallets[walletName]
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
	wlt, ok := wm.wallets[walletName]
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
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeBondTx(sender, receiver, publicKey, amt, opts...)
}

func (wm *Manager) MakeUnbondTx(walletName, validator string, opts ...wallet.TxOption) (*tx.Tx, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeUnbondTx(validator, opts...)
}

func (wm *Manager) MakeWithdrawTx(
	walletName, sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.MakeWithdrawTx(sender, receiver, amt, opts...)
}

func (wm *Manager) SignTransaction(walletName, password string, trx *tx.Tx) error {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.SignTransaction(password, trx)
}

func (wm *Manager) BroadcastTransaction(walletName string, trx *tx.Tx) (string, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.BroadcastTransaction(trx)
}

func (wm *Manager) UnloadWallet(
	walletName string,
) error {
	if _, ok := wm.wallets[walletName]; !ok {
		return ErrWalletNotLoaded
	}

	delete(wm.wallets, walletName)

	return nil
}

func (wm *Manager) TotalBalance(walletName string) (amount.Amount, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.TotalBalance()
}

func (wm *Manager) TotalStake(walletName string) (amount.Amount, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return 0, ErrWalletNotLoaded
	}

	return wlt.TotalStake()
}

func (wm *Manager) SignRawTransaction(
	walletName, password string, rawTx []byte,
) (txID, data []byte, err error) {
	wlt, ok := wm.wallets[walletName]
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

func (wm *Manager) GetNewAddress(
	walletName, label, password string,
	addressType crypto.AddressType,
) (*storage.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	info, err := wlt.NewAddress(addressType, label, wallet.WithPassword(password))
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (wm *Manager) SignMessage(walletName, password, addr, msg string) (string, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return "", ErrWalletNotLoaded
	}

	return wlt.SignMessage(password, addr, msg)
}

func (wm *Manager) AddressInfo(walletName, address string) (*storage.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.AddressInfo(address), nil
}

func (wm *Manager) WalletInfo(walletName string) (*wallet.Info, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.Info(), nil
}

func (wm *Manager) ListWallets() ([]string, error) {
	wallets := make([]string, 0)

	files, err := util.ListFilesInDir(wm.walletDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		_, err = wallet.Open(file, true)
		if err != nil {
			logger.Warn(fmt.Sprintf("file %s is not wallet", file))

			continue
		}

		wallets = append(wallets, filepath.Base(file))
	}

	return wallets, nil
}

func (wm *Manager) ListAddresses(walletName string) ([]storage.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, ErrWalletNotLoaded
	}

	return wlt.ListAddresses(), nil
}

func (wm *Manager) UpdatePassword(walletName, oldPassword, newPassword string) error {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return ErrWalletNotLoaded
	}

	return wlt.UpdatePassword(oldPassword, newPassword)
}
