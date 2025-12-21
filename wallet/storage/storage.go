package storage

import (
	"context"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type IStorage interface {
	Version() int
	Network() genesis.ChainType
	Path() string
	IsEncrypted() bool
	WalletInfo() *types.WalletInfo

	DefaultFee() amount.Amount
	SetDefaultFee(fee amount.Amount)

	Mnemonic(password string) (string, error)
	ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error
	ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) error
	PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error)
	Neuter(path string) error

	AddressCount() int
	HasAddress(addr string) bool
	ListAddresses(opts ...types.ListAddressOption) []types.AddressInfo
	AddressInfo(address string) *types.AddressInfo
	AddressLabel(address string) string
	SetAddressLabel(address, label string) error
	NewValidatorAddress(label string) (*types.AddressInfo, error)
	NewBLSAccountAddress(label string) (*types.AddressInfo, error)
	NewEd25519AccountAddress(label, password string) (*types.AddressInfo, error)

	Save() error
	Upgrade() error
	UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error
	RecoverAddresses(ctx context.Context, password string, eventFunc func(addr string) (bool, error)) error

	AddPending(addr string, amt amount.Amount, txID tx.ID, data []byte)
	AddActivity(addr string, amt amount.Amount, trx *pactus.GetTransactionResponse)
	HasTransaction(id string) bool
	GetAddrHistory(addr string) []types.HistoryInfo
}
