package store

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/syndtr/goleveldb/leveldb"
)

// publicKeyError converts the error of a public key lookup for the given address.
// Only a missing key is reported as PublicKeyNotFoundError; any other failure
// (for example a read or decoding error) is wrapped, so it is not mistaken for an absent key.
func publicKeyError(addr crypto.Address, err error) error {
	var notFoundErr PublicKeyNotFoundError
	if errors.As(err, &notFoundErr) {
		return err
	}

	if errors.Is(err, leveldb.ErrNotFound) || errors.Is(err, ErrNotFound) {
		return PublicKeyNotFoundError{Address: addr}
	}

	return fmt.Errorf("unable to read the public key of %s: %w", addr, err)
}

// CommittedBlock is a lazy reference to a committed block stored in the database.
// It carries raw serialized bytes and metadata, but does not parse or deserialize
// the block eagerly. The full [block.Block] is constructed on demand by calling
// [CommittedBlock.ToBlock], which deserializes the data and restores any stripped
// public keys by looking them up through the backing [Store].
type CommittedBlock struct {
	store Store // back-reference to the store, used by ToBlock to restore public keys.

	BlockHash hash.Hash    // hash of this block.
	Height    types.Height // height at which this block was committed.
	Data      []byte       // raw serialized block bytes, as returned by the database.
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
				return nil, publicKeyError(trx.Payload().Signer(), err)
			}
			trx.SetPublicKey(pub)
		}
	}

	return blk, nil
}

type CommittedTx struct {
	store *store

	TxID      tx.ID
	Height    types.Height
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
			return nil, publicKeyError(trx.Payload().Signer(), err)
		}
		trx.SetPublicKey(pub)
	}

	return trx, nil
}

type Reader interface {
	Block(height types.Height) (*CommittedBlock, error)
	BlockHeight(h hash.Hash) types.Height
	BlockHash(height types.Height) hash.Hash
	SortitionSeed(height types.Height) *sortition.VerifiableSeed
	Transaction(txID tx.ID) (*CommittedTx, error)
	RecentTransaction(txID tx.ID) bool
	PublicKey(addr crypto.Address) (crypto.PublicKey, error)
	HasPublicKey(addr crypto.Address) bool
	HasAccount(crypto.Address) (bool, error)
	Account(addr crypto.Address) (*account.Account, error)
	TotalAccounts() int32
	HasValidator(addr crypto.Address) bool
	ValidatorAddresses() []crypto.Address
	Validator(addr crypto.Address) (*validator.Validator, error)
	ValidatorByNumber(num int32) (*validator.Validator, error)
	IterateValidators(consumer func(*validator.Validator) (stop bool))
	IterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool))
	AnchorAddresses(skip, count uint32) ([]crypto.Address, uint32)
	TotalValidators() int32
	ActiveValidators() int32
	LastCertificate() *certificate.Certificate
	IsBanned(addr crypto.Address) bool
	IsPruned() bool
	PruningHeight() types.Height
}

type Store interface {
	Reader

	UpdateAccount(addr crypto.Address, acc *account.Account)
	UpdateValidator(val *validator.Validator)
	UpdateValidatorProtocolVersion(addr crypto.Address, ver protocol.Version)
	SaveBlock(blk *block.Block, cert *certificate.Certificate)
	Prune(callback func(pruned bool, pruningHeight types.Height) bool) error
	WriteBatch() error

	Close()
}
