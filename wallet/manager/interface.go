package manager

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	wallet "github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/vault"
)

// WalletManager defines the public API of the wallet manager.
type IManager interface {
	GetValidatorAddress(publicKey string) (string, error)
	CreateWallet(walletName, password string) (string, error)
	RestoreWallet(walletName, mnemonic, password string) error
	LoadWallet(walletName, serverAddr string) error
	UnloadWallet(walletName string) error
	TotalBalance(walletName string) (amount.Amount, error)
	TotalStake(walletName string) (amount.Amount, error)
	SignRawTransaction(walletName, password string, rawTx []byte) (txID, data []byte, err error)
	GetNewAddress(walletName, label, password string, addressType crypto.AddressType) (*vault.AddressInfo, error)
	AddressHistory(walletName, address string) ([]wallet.HistoryInfo, error)
	SignMessage(walletName, password, addr, msg string) (string, error)
	GetAddressInfo(walletName, address string) (*vault.AddressInfo, error)
	SetAddressLabel(walletName, address, label string) error
	WalletInfo(walletName string) (*wallet.Info, error)
	ListWallet() ([]string, error)
	ListAddress(walletName string) ([]vault.AddressInfo, error)
	UpdatePassword(walletName, oldPassword, newPassword string) error
}
