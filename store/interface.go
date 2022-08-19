package store

import (
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/validator"
)

// TODO: store blocks inside flat files (to reduce the size of levelDB)
// Bitcoin impl:
// https://github.com/btcsuite/btcd/blob/0886f1e5c1fd28ad24aaca4dbccc5f4ab85e58ca/database/ffldb/blockio.go
// https://bitcoindev.network/understanding-the-data/

// TODO: How to undo or rollback at least for last 21 blocks

type StoredBlock struct {
	height uint32
	data   []byte
}

func (s *StoredBlock) Height() uint32 {
	return s.height
}

func (s *StoredBlock) ToFullBlock() (*block.Block, error) {
	b, err := block.FromBytes(s.data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Reader interface {
	Block(hash hash.Hash) (*StoredBlock, error)
	BlockHash(height uint32) hash.Hash
	// It only remembers most recent stamps
	FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool)
	// It only remembers most recent stamps
	FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool)
	Transaction(id tx.ID) (*tx.Tx, error)
	HasAccount(crypto.Address) bool
	Account(addr crypto.Address) (*account.Account, error)
	TotalAccounts() int32
	HasValidator(crypto.Address) bool
	Validator(addr crypto.Address) (*validator.Validator, error)
	ValidatorByNumber(num int32) (*validator.Validator, error)
	IterateValidators(consumer func(*validator.Validator) (stop bool))
	IterateAccounts(consumer func(*account.Account) (stop bool))
	TotalValidators() int32
	LastCertificate() (uint32, *block.Certificate)
}

type Store interface {
	Reader

	UpdateAccount(acc *account.Account)
	UpdateValidator(val *validator.Validator)
	SaveBlock(height uint32, block *block.Block, cert *block.Certificate)
	WriteBatch() error
	Close() error
}
