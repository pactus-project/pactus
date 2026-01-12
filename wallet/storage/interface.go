package storage

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

// QueryParams specifies filters for querying stored transactions.
type QueryParams struct {
	Address   string
	Direction types.TxDirection
	Count     int
	Skip      int
}

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
	GetPendingTransactions() (map[string]*types.TransactionInfo, error)
	UpdateTransactionStatus(no int64, status types.TransactionStatus, blockHeight uint32) error
	GetTransaction(no int64) (*types.TransactionInfo, error)
	HasTransaction(txID string) bool
	QueryTransactions(params QueryParams) ([]*types.TransactionInfo, error)

	Close() error
	Clone(path string) (IStorage, error)
	IsLegacy() bool
}
