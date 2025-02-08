package store

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
)

// TODO: store blocks inside flat files (to reduce the size of levelDB)
// Bitcoin impl:
// https://github.com/btcsuite/btcd/blob/0886f1e5c1fd28ad24aaca4dbccc5f4ab85e58ca/database/ffldb/blockio.go
// https://bitcoindev.network/understanding-the-data/

// TODO: How to undo or rollback at least for last 21 blocks

type CommittedBlock struct {
	store *store

	BlockHash hash.Hash
	Height    uint32
	Data      []byte
}

func (s *CommittedBlock) ToBlock() (*block.Block, error) {
	blk, err := block.FromBytes(s.Data)
	if err != nil {
		return nil, err
	}

	trxs := blk.Transactions()
	for i := 0; i < trxs.Len(); i++ {
		trx := trxs[i]
		if trx.IsPublicKeyStriped() {
			pub, err := s.store.PublicKey(trx.Payload().Signer())
			if err != nil {
				return nil, PublicKeyNotFoundError{
					Address: trx.Payload().Signer(),
				}
			}
			trx.SetPublicKey(pub)
		}
	}

	return blk, nil
}

type CommittedTx struct {
	store *store

	TxID      tx.ID
	Height    uint32
	BlockTime uint32
	Data      []byte
}

func (s *CommittedTx) ToTx() (*tx.Tx, error) {
	trx, err := tx.FromBytes(s.Data)
	if err != nil {
		return nil, err
	}

	if trx.IsPublicKeyStriped() {
		pub, err := s.store.PublicKey(trx.Payload().Signer())
		if err != nil {
			return nil, PublicKeyNotFoundError{
				Address: trx.Payload().Signer(),
			}
		}
		trx.SetPublicKey(pub)
	}

	return trx, nil
}

type Reader interface {
	Block(height uint32) (*CommittedBlock, error)
	BlockHeight(h hash.Hash) uint32
	BlockHash(height uint32) hash.Hash
	SortitionSeed(blockHeight uint32) *sortition.VerifiableSeed
	Transaction(txID tx.ID) (*CommittedTx, error)
	RecentTransaction(txID tx.ID) bool
	PublicKey(addr crypto.Address) (crypto.PublicKey, error)
	HasPublicKey(addr crypto.Address) bool
	HasAccount(crypto.Address) bool
	Account(addr crypto.Address) (*account.Account, error)
	TotalAccounts() int32
	HasValidator(addr crypto.Address) bool
	ValidatorAddresses() []crypto.Address
	Validator(addr crypto.Address) (*validator.Validator, error)
	ValidatorByNumber(num int32) (*validator.Validator, error)
	IterateValidators(consumer func(*validator.Validator) (stop bool))
	IterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool))
	TotalValidators() int32
	LastCertificate() *certificate.BlockCertificate
	IsBanned(addr crypto.Address) bool
	IsPruned() bool
	PruningHeight() uint32
	XeggexAccount() *XeggexAccount
}

type Store interface {
	Reader

	UpdateAccount(addr crypto.Address, acc *account.Account)
	UpdateValidator(val *validator.Validator)
	SaveBlock(blk *block.Block, cert *certificate.BlockCertificate)
	Prune(callback func(pruned bool, pruningHeight uint32) bool) error
	WriteBatch() error
	Close()
}
