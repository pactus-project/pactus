package store

import (
	"bytes"

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

type CommittedBlock struct {
	BlockHash hash.Hash
	Height    uint32
	Data      []byte
}

func (s *CommittedBlock) ToBlock() *block.Block {
	b, err := block.FromBytes(s.Data)
	if err != nil {
		panic(err)
	}
	if !b.Hash().EqualsTo(s.BlockHash) {
		panic("invalid data. block hash does not match")
	}
	return b
}

type CommittedTx struct {
	TxID      tx.ID
	Height    uint32
	BlockTime uint32
	Data      []byte
}

func (s *CommittedTx) ToTx() *tx.Tx {
	trx := new(tx.Tx)
	r := bytes.NewReader(s.Data)
	if err := trx.Decode(r); err != nil {
		panic(err)
	}

	// TODO: we can set public key, if it doesn't set
	// pubKey, found := store.PublicKey(trx.Payload().Signer())
	// if !found {
	// 	panic("unable to find the public key")
	// }
	// trx.SetPublicKey(pubKey)

	return trx
}

type Reader interface {
	Block(height uint32) (*CommittedBlock, error)
	BlockHeight(hash hash.Hash) uint32
	BlockHash(height uint32) hash.Hash
	RecentBlockByStamp(stamp hash.Stamp) (uint32, *block.Block)
	Transaction(id tx.ID) (*CommittedTx, error)
	PublicKey(addr crypto.Address) (crypto.PublicKey, bool)
	HasAccount(crypto.Address) bool
	Account(addr crypto.Address) (*account.Account, error)
	AccountByNumber(number int32) (*account.Account, error)
	TotalAccounts() int32
	HasValidator(addr crypto.Address) bool
	ValidatorAddresses() []crypto.Address
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
