package manager

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
)

// IManager defines the public API of the wallet manager.
type IManager interface {
	Start() error
	Stop()

	GetValidatorAddress(publicKey string) (string, error)
	CreateWallet(walletName, password string) (string, error)
	RestoreWallet(walletName, mnemonic, password string) error

	LoadWallet(walletName string, opts ...wallet.OpenWalletOption) error
	UnloadWallet(walletName string) error
	ListWallets() ([]string, error)
	WalletInfo(walletName string) (*types.WalletInfo, error)
	UpdatePassword(walletName, oldPassword, newPassword string) error
	TotalBalance(walletName string) (amount.Amount, error)
	TotalStake(walletName string) (amount.Amount, error)
	SetDefaultFee(walletName string, fee amount.Amount) error

	SignRawTransaction(walletName, password string, rawTx []byte) (txID, data []byte, err error)
	SignMessage(walletName, password, addr, msg string) (string, error)
	PrivateKey(walletName, password, addr string) (crypto.PrivateKey, error)
	Mnemonic(walletName, password string) (string, error)

	NewAddress(walletName string, addressType crypto.AddressType, label string,
		opts ...wallet.NewAddressOption) (*types.AddressInfo, error)
	ListAddresses(walletName string, opts ...wallet.ListAddressOption) ([]types.AddressInfo, error)
	AddressInfo(walletName, address string) (*types.AddressInfo, error)
	AddressLabel(walletName, addr string) (string, error)
	SetAddressLabel(walletName, addr, label string) error
	Balance(walletName, addr string) (amount.Amount, error)
	Stake(walletName, addr string) (amount.Amount, error)

	// Transaction creation / signing / broadcast
	MakeTransferTx(walletName, sender, receiver string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	MakeBondTx(walletName, sender, receiver, publicKey string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	MakeUnbondTx(walletName, validator string, opts ...wallet.TxOption) (*tx.Tx, error)
	MakeWithdrawTx(walletName, sender, receiver string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	SignTransaction(walletName, password string, trx *tx.Tx) error
	BroadcastTransaction(walletName string, trx *tx.Tx) (string, error)
	ListTransactions(walletName string, opts ...wallet.ListTransactionsOption) ([]*types.TransactionInfo, error)
}
