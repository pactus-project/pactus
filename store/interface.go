package store

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

// TODO: store blocks inside flat files (to reduce the size of levelDB)
// Bitcoin impl:
// https://github.com/btcsuite/btcd/blob/0886f1e5c1fd28ad24aaca4dbccc5f4ab85e58ca/database/ffldb/blockio.go
// https://bitcoindev.network/understanding-the-data/

// TODO: How to undo or rollback at least for last 21 blocks

type StoreBlock struct {
	Height     int
	Block      *block.Block
	HeaderData []byte
}

type Reader interface {
	Block(hash hash.Hash) (*StoreBlock, error)
	BlockHash(height int) hash.Hash
	BlockHashByStamp(stamp hash.Stamp) hash.Hash // It only remembers most recent stamps
	BlockHeightByStamp(stamp hash.Stamp) int     // It only remembers most recent stamps
	Transaction(id tx.ID) (*tx.Tx, error)
	HasAccount(crypto.Address) bool
	Account(addr crypto.Address) (*account.Account, error)
	TotalAccounts() int
	HasValidator(crypto.Address) bool
	Validator(addr crypto.Address) (*validator.Validator, error)
	ValidatorByNumber(num int) (*validator.Validator, error)
	IterateValidators(consumer func(*validator.Validator) (stop bool))
	IterateAccounts(consumer func(*account.Account) (stop bool))
	TotalValidators() int
	LastCertificate() (int, *block.Certificate)
}

type Store interface {
	Reader

	UpdateAccount(acc *account.Account)
	UpdateValidator(val *validator.Validator)
	SaveBlock(height int, block *block.Block, cert *block.Certificate)
	WriteBatch() error
	Close() error
}
