package storage

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type IStorage interface {
	WalletInfo() *types.WalletInfo
	Vault() *vault.Vault
	UpdateVault(vault *vault.Vault) error
	SetDefaultFee(fee amount.Amount) error

	AllAddresses() ([]types.AddressInfo, error)
	InsertAddress(info *types.AddressInfo) error
	UpdateAddress(info *types.AddressInfo) error

	AddPending(addr string, amt amount.Amount, txID tx.ID, data []byte) error
	AddActivity(addr string, amt amount.Amount, trx *pactus.GetTransactionResponse) error
	HasTransaction(id string) bool
	GetAddrHistory(addr string) []types.HistoryInfo
}
