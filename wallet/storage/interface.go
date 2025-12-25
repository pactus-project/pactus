package storage

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

type IStorage interface {
	WalletInfo() *types.WalletInfo
	Vault() *vault.Vault
	UpdateVault(vault *vault.Vault) error
	SetDefaultFee(fee amount.Amount) error

	AllAddresses() []types.AddressInfo
	AddressInfo(address string) (*types.AddressInfo, error)
	HasAddress(address string) bool
	AddressCount() int
	InsertAddress(info *types.AddressInfo) error
	UpdateAddress(info *types.AddressInfo) error

	InsertTransaction(info *types.TransactionInfo) error
	UpdateTransactionStatus(id string, status types.TransactionStatus) error
	HasTransaction(id string) bool
	GetTransaction(id string) (*types.TransactionInfo, error)
	ListTransactions(receiver string, count int, skip int) ([]types.TransactionInfo, error)

	Close() error
	Clone(path string) (IStorage, error)
}
