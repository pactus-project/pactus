package storage

import (
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	wtypes "github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

// QueryParams specifies filters for querying stored transactions.
type QueryParams struct {
	Address   string
	Direction wtypes.TxDirection
	Count     int
	Skip      int
}

type IStorage interface {
	WalletInfo() *wtypes.WalletInfo
	Vault() *vault.Vault
	UpdateVault(vault *vault.Vault) error
	SetDefaultFee(fee amount.Amount) error

	AllAddresses() []wtypes.AddressInfo
	AddressInfo(address string) (*wtypes.AddressInfo, error)
	HasAddress(address string) bool
	AddressCount() int
	InsertAddress(info *wtypes.AddressInfo) error
	UpdateAddress(info *wtypes.AddressInfo) error

	InsertTransaction(info *wtypes.TransactionInfo) error
	GetPendingTransactions() (map[string]*wtypes.TransactionInfo, error)
	UpdateTransactionStatus(no int64, status wtypes.TransactionStatus, blockHeight types.Height) error
	GetTransaction(no int64) (*wtypes.TransactionInfo, error)
	HasTransaction(txID string) bool
	QueryTransactions(params QueryParams) ([]*wtypes.TransactionInfo, error)

	Close() error
	Clone(path string) (IStorage, error)
	IsLegacy() bool
}
