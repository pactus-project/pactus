package manager

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/vault"
)

// IManager defines the public API of the wallet manager.
type IManager interface {
	GetValidatorAddress(publicKey string) (string, error)
	CreateWallet(walletName, password string) (string, error)
	RestoreWallet(walletName, mnemonic, password string) error

	LoadWallet(walletName, serverAddr string) error
	UnloadWallet(walletName string) error
	ListWallet() ([]string, error)
	WalletPath(walletName string) string
	WalletInfo(walletName string) (*wallet.Info, error)
	UpdatePassword(walletName, oldPassword, newPassword string) error
	TotalBalance(walletName string) (amount.Amount, error)
	TotalStake(walletName string) (amount.Amount, error)

	SignRawTransaction(walletName, password string, rawTx []byte) (txID, data []byte, err error)
	SignMessage(walletName, password, addr, msg string) (string, error)
	PrivateKey(walletName, password, addr string) (crypto.PrivateKey, error)
	Mnemonic(walletName, password string) (string, error)

	// Address management
	// NOTE: Accepting wallet.NewAddressOption keeps UI and other callers compatible
	// with wallet.WithPassword(...) etc.
	NewAddress(
		walletName string,
		addressType crypto.AddressType,
		label string,
		opts ...wallet.NewAddressOption,
	) (*vault.AddressInfo, error)
	GetNewAddress(walletName, label, password string, addressType crypto.AddressType) (*vault.AddressInfo, error)
	AddressHistory(walletName, address string) ([]wallet.HistoryInfo, error)
	GetAddressInfo(walletName, address string) (*vault.AddressInfo, error)
	SetAddressLabel(walletName, address, label string) error
	ListAddress(walletName string) ([]vault.AddressInfo, error)
	AllAccountAddresses(walletName string) ([]vault.AddressInfo, error)
	AllValidatorAddresses(walletName string) ([]vault.AddressInfo, error)
	Label(walletName, addr string) (string, error)
	SetLabel(walletName, addr, label string) error

	// Balance and fee
	Balance(walletName, addr string) (amount.Amount, error)
	Stake(walletName, addr string) (amount.Amount, error)
	SetDefaultFee(walletName string, fee amount.Amount) error

	// Transaction creation / signing / broadcast
	MakeTransferTx(walletName, sender, receiver string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	MakeBondTx(walletName, sender, receiver, publicKey string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	MakeUnbondTx(walletName, validator string, opts ...wallet.TxOption) (*tx.Tx, error)
	MakeWithdrawTx(walletName, sender, receiver string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	SignTransaction(walletName, password string, trx *tx.Tx) error
	BroadcastTransaction(walletName string, trx *tx.Tx) (string, error)
}
