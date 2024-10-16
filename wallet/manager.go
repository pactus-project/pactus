package wallet

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

func (wm *Manager) getWalletPath(walletName string) string {
	return util.MakeAbs(filepath.Join(wm.walletDirectory, walletName))
}

func (wm *Manager) createWalletWithMnemonic(
	walletName, mnemonic, password string,
) error {
	walletPath := wm.getWalletPath(walletName)
	if isExists := util.PathExists(walletPath); isExists {
		return status.Errorf(codes.AlreadyExists, "wallet already exists")
	}

	wlt, err := Create(walletPath, mnemonic, password, wm.chainType)
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

func (wm *Manager) CreateWallet(
	walletName, password string,
) (string, error) {
	mnemonic, err := GenerateMnemonic(128)
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
		return status.Errorf(codes.AlreadyExists, "wallet already loaded")
	}

	walletPath := util.MakeAbs(filepath.Join(wm.walletDirectory, walletName))
	wlt, err := Open(walletPath, true, WithCustomServers([]string{serverAddr}))
	if err != nil {
		return err
	}

	wm.wallets[walletName] = wlt

	return nil
}

func (wm *Manager) UnloadWallet(
	walletName string,
) error {
	if _, ok := wm.wallets[walletName]; !ok {
		return status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	delete(wm.wallets, walletName)

	return nil
}

func (wm *Manager) TotalBalance(
	walletName string,
) (amount.Amount, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return 0, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.TotalBalance()
}

func (wm *Manager) TotalStake(walletName string) (amount.Amount, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return 0, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.TotalStake()
}

func (wm *Manager) SignRawTransaction(
	walletName, password string, rawTx []byte,
) ([]byte, []byte, error) {
	wlt, ok := wm.wallets[walletName]
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

func (wm *Manager) GetNewAddress(
	walletName, label, password string,
	addressType crypto.AddressType,
) (*vault.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	var addressInfo *vault.AddressInfo
	switch addressType {
	case crypto.AddressTypeValidator:
		info, err := wlt.NewValidatorAddress(label)
		if err != nil {
			return nil, err
		}
		addressInfo = info

	case crypto.AddressTypeBLSAccount:
		info, err := wlt.NewBLSAccountAddress(label)
		if err != nil {
			return nil, err
		}
		addressInfo = info

	case crypto.AddressTypeEd25519Account:
		if password == "" {
			return nil, status.Errorf(codes.InvalidArgument, "password cannot be empty when address type is Ed25519")
		}

		info, err := wlt.NewEd25519AccountAddress(label, password)
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

func (wm *Manager) AddressHistory(
	walletName, address string,
) ([]HistoryInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.History(address), nil
}

func (wm *Manager) SignMessage(walletName, password, addr, msg string) (string, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return "", status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.SignMessage(password, addr, msg)
}

func (wm *Manager) GetAddressInfo(walletName, address string) (*vault.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return wlt.AddressInfo(address), nil
}

func (wm *Manager) SetAddressLabel(walletName, address, label string) error {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	err := wlt.SetLabel(address, label)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	return wlt.Save()
}

func (wm *Manager) GetWalletInfo(walletName string) (Info, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return Info{}, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return Info{
		WalletName: walletName,
		Version:    int64(wlt.store.Version),
		Network:    wlt.store.Network.String(),
		UUID:       wlt.store.UUID.String(),
		Encrypted:  wlt.IsEncrypted(),
		Crc:        wlt.store.VaultCRC,
		CreatedAt:  wlt.store.CreatedAt,
	}, nil
}

func (wm *Manager) ListWallet() ([]string, error) {
	wallets := make([]string, 0)

	files, err := util.ListFilesInDir(wm.walletDirectory)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		_, err = Open(f, true)
		if err != nil {
			logger.Warn(fmt.Sprintf("file %s is not wallet", f))

			continue
		}

		wallets = append(wallets, filepath.Base(f))
	}

	return wallets, nil
}

func (wm *Manager) ListAddress(walletName string) ([]vault.AddressInfo, error) {
	wlt, ok := wm.wallets[walletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	return append(wlt.AllValidatorAddresses(), wlt.AllAccountAddresses()...), nil
}
