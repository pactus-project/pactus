package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// AddressInfo represents an address entry in the wallet.
type AddressInfo struct {
	Address   string    `json:"address"`
	PublicKey string    `json:"public_key"`
	Label     string    `json:"label"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HistoryInfo struct {
	TxID        string
	Time        *time.Time
	PayloadType string
	Desc        string
	Amount      amount.Amount
}

type IStorage interface {
	Version() int
	CreatedAt() time.Time
	UUID() uuid.UUID
	Network() genesis.ChainType
	CoinType() addresspath.CoinType
	IsEncrypted() bool
	Mnemonic(password string) (string, error)

	Neuter(path string) error

	DefaultFee() amount.Amount
	SetDefaultFee(fee amount.Amount)

	ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error
	ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) error
	PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error)

	AddressCount() int
	HasAddress(addr string) bool
	ListAddresses() []AddressInfo
	ListValidatorAddresses() []AddressInfo
	ListAccountAddresses() []AddressInfo
	AddressInfo(address string) *AddressInfo
	AddressByPath(path string) *AddressInfo
	AddressLabel(address string) string
	SetAddressLabel(address, label string) error
	NewValidatorAddress(label string) (*AddressInfo, error)
	NewBLSAccountAddress(label string) (*AddressInfo, error)
	NewEd25519AccountAddress(label, password string) (*AddressInfo, error)

	Save() error
	Upgrade() error
	UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error
	RecoverAddresses(ctx context.Context, password string, eventFunc func(addr string) (bool, error)) error

	AddPending(addr string, amt amount.Amount, txID tx.ID, data []byte)
	AddActivity(addr string, amt amount.Amount, trx *pactus.GetTransactionResponse)
	HasTransaction(id string) bool
	GetAddrHistory(addr string) []HistoryInfo
}
