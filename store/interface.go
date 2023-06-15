package store

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
)

// TODO: store blocks inside flat files (to reduce the size of levelDB)
// Bitcoin impl:
// https://github.com/btcsuite/btcd/blob/0886f1e5c1fd28ad24aaca4dbccc5f4ab85e58ca/database/ffldb/blockio.go
// https://bitcoindev.network/understanding-the-data/

// TODO: How to undo or rollback at least for last 21 blocks

type StoredBlock struct {
	BlockHash hash.Hash
	Height    uint32
	Data      []byte
}

func (s *StoredBlock) ToBlock() *block.Block {
	b, err := block.FromBytes(s.Data)
	if err != nil {
		panic(err)
	}
	if !b.Hash().EqualsTo(s.BlockHash) {
		panic("invalid data. block hash does not match")
	}
	return b
}

type StoredTx struct {
	TxID      tx.ID
	Height    uint32
	BlockTime uint32
	Data      []byte
}

func (s *StoredTx) ToTx() *tx.Tx {
	trx, err := tx.FromBytes(s.Data)
	if err != nil {
		panic(err)
	}
	if !trx.ID().EqualsTo(s.TxID) {
		panic("invalid data. transaction id does not match")
	}
	return trx
}

type Reader interface {
	Block(height uint32) (*StoredBlock, error)
	BlockHeight(hash hash.Hash) uint32
	BlockHash(height uint32) hash.Hash
	// It only remembers most recent stamps
	FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool)
	// It only remembers most recent stamps
	FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool)
	Transaction(id tx.ID) (*StoredTx, error)
	HasAccount(crypto.Address) bool
	Account(addr crypto.Address) (*account.Account, error)
	AccountByNumber(number int32) (*account.Account, error)
	TotalAccounts() int32
	HasValidator(crypto.Address) bool
	Validator(addr crypto.Address) (*validator.Validator, error)
	ValidatorByNumber(num int32) (*validator.Validator, error)
	IterateValidators(consumer func(*validator.Validator) (stop bool))
	IterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool))
	TotalValidators() int32
	LastCertificate() (uint32, *block.Certificate)
}

type Store interface {
	Reader

	UpdateAccount(addr crypto.Address, acc *account.Account)
	UpdateValidator(val *validator.Validator)
	SaveBlock(height uint32, block *block.Block, cert *block.Certificate)
	WriteBatch() error
	Close() error
}
