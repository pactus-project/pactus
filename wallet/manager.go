package wallet

import (
	"path/filepath"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Manager struct {
	wallets         map[string]*Wallet
	chainType       genesis.ChainType
	walletDirectory string
}

func NewWalletManager(conf *Config) *Manager {
	return &Manager{
		wallets:         make(map[string]*Wallet),
		chainType:       conf.ChainType,
		walletDirectory: conf.WalletsDir,
	}
}

func (w *Manager) getWalletPath(walletName string) string {
	return util.MakeAbs(filepath.Join(w.walletDirectory, walletName))
}

func (w *Manager) createWalletWithMnemonic(
	walletName, mnemonic, password string,
) error {
	walletPath := w.getWalletPath(walletName)
	wlt, err := Create(walletPath, mnemonic, password, w.chainType)
	if err != nil {
		return err
	}

	return wlt.Save()
}

func (*Manager) GetValidatorAddress(
	publicKey string,
) (string, error) {
	pubKey, err := bls.PublicKeyFromString(publicKey)
	if err != nil {
		return "", err
	}

	return pubKey.ValidatorAddress().String(), nil
}

func (w *Manager) CreateWallet(
	walletName, password string,
) (string, error) {
	mnemonic, err := GenerateMnemonic(128)
	if err != nil {
		return "", err
	}

	walletPath := w.getWalletPath(walletName)
	if isExists := util.PathExists(walletPath); isExists {
		return "", status.Errorf(codes.AlreadyExists, "wallet already exists")
	}

	if err := w.createWalletWithMnemonic(walletName, mnemonic, password); err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (w *Manager) RestoreWallet(walletName, mnemonic, password string) error {
	walletPath := w.getWalletPath(walletName)
	if isExists := util.PathExists(walletPath); isExists {
		return status.Errorf(codes.AlreadyExists, "wallet already exists")
	}

	return w.createWalletWithMnemonic(walletName, mnemonic, password)
}

func (w *Manager) LoadWallet(walletName string) error {
	if _, ok := w.wallets[walletName]; ok {
		// TODO: define special codes for errors
		return status.Errorf(codes.AlreadyExists, "wallet already loaded")
	}

	walletPath := util.MakeAbs(filepath.Join(w.walletDirectory, walletName))
	wlt, err := Open(walletPath, true)
	if err != nil {
		return err
	}

	w.wallets[walletName] = wlt

	return nil
}

func (w *Manager) UnloadWallet(
	walletName string,
) error {
	if _, ok := w.wallets[walletName]; !ok {
		return status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	delete(w.wallets, walletName)

	return nil
}

func (w *Manager) TotalBalance(
	walletName string,
) (int64, error) {
	wlt, ok := w.wallets[walletName]
	if !ok {
		return 0, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.TotalBalance().ToNanoPAC(), nil
}

func (w *Manager) SignRawTransaction(
	walletName, password string, rawTx []byte,
) ([]byte, []byte, error) {
	wlt, ok := w.wallets[walletName]
	if !ok {
		return nil, nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	trx, err := tx.FromBytes(rawTx)
	if err != nil {
		return nil, nil, err
	}

	if err := wlt.SignTransaction(password, trx); err != nil {
		return nil, nil, err
	}

	data, err := trx.Bytes()
	if err != nil {
		return nil, nil, err
	}

	return trx.ID().Bytes(), data, nil
}

func (w *Manager) GetNewAddress(
	walletName, label string,
	addressType crypto.AddressType,
) (*vault.AddressInfo, error) {
	wlt, ok := w.wallets[walletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	var addressInfo *vault.AddressInfo
	switch addressType {
	case crypto.AddressTypeBLSAccount:
		info, err := wlt.NewBLSAccountAddress(label)
		if err != nil {
			return nil, err
		}
		addressInfo = info

	case crypto.AddressTypeValidator:
		info, err := wlt.NewValidatorAddress(label)
		if err != nil {
			return nil, err
		}
		addressInfo = info

	case crypto.AddressTypeTreasury:
		return nil, status.Errorf(codes.InvalidArgument, "invalid address type")

	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid address type")
	}

	if err := wlt.Save(); err != nil {
		return nil, err
	}

	return addressInfo, nil
}

func (w *Manager) AddressHistory(
	walletName, address string,
) ([]HistoryInfo, error) {
	wlt, ok := w.wallets[walletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.GetHistory(address), nil
}
