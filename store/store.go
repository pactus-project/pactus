package store

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type StoreReader interface {
	BlockByHeight(height int) (*block.Block, error)
	BlockByHash(hash crypto.Hash) (*block.Block, int, error)
	BlockHeight(hash crypto.Hash) (int, error)
}

type Store struct {
	lk deadlock.RWMutex

	config         *config.Config
	blockStore     *blockStore
	txStore        *txStore
	accountStore   *accountStore
	validatorStore *validatorStore
}

func NewStore(config *config.Config) (*Store, error) {
	blockStore, err := newBlockStore(config.Store.BlockStorePath())
	if err != nil {
		return nil, err
	}
	txStore, err := newTxStore(config.Store.TxStorePath())
	if err != nil {
		return nil, err
	}
	accountStore, err := newAccountStore(config.Store.AccountStorePath())
	if err != nil {
		return nil, err
	}
	validatorStore, err := newValidatorStore(config.Store.ValidatorStorePath())
	if err != nil {
		return nil, err
	}
	return &Store{
		config:         config,
		blockStore:     blockStore,
		txStore:        txStore,
		accountStore:   accountStore,
		validatorStore: validatorStore,
	}, nil
}

func (s *Store) SaveBlock(block block.Block, height int) error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.SaveBlock(block, height)
}

func (s *Store) BlockByHeight(height int) (*block.Block, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.RetrieveBlock(height)
}

func (s *Store) BlockByHash(hash crypto.Hash) (*block.Block, int, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	height, err := s.blockStore.RetrieveBlockHeight(hash)
	if err != nil {
		return nil, -1, err
	}
	block, err := s.blockStore.RetrieveBlock(height)
	if err != nil {
		return nil, -1, err
	}
	return block, height, nil
}

func (s *Store) BlockHeight(hash crypto.Hash) (int, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.RetrieveBlockHeight(hash)
}

func (s *Store) SaveTx(tx tx.Tx, receipt tx.Receipt) error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.SaveTx(tx, receipt)
}

func (s *Store) RetrieveTx(hash crypto.Hash) (*tx.Tx, *tx.Receipt, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.RetrieveTx(hash)
}

func (s *Store) HasAccount(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.HasAccount(addr)
}

func (s *Store) RetrieveAccount(addr crypto.Address) *account.Account {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.RetrieveAccount(addr)
}

func (s *Store) AccountCount() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.AccountCount()
}

func (s *Store) IterateAccounts(consumer func(*account.Account) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.IterateAccounts(consumer)
}

func (s *Store) UpdateAccount(acc *account.Account) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.UpdateAccount(acc)
}

func (s *Store) HasValidator(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.HasValidator(addr)
}

func (s *Store) RetrieveValidator(addr crypto.Address) *validator.Validator {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.RetrieveValidator(addr)
}

func (s *Store) ValidatorCount() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.ValidatorCount()
}

func (s *Store) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.validatorStore.IterateValidators(consumer)
}

func (s *Store) UpdateValidator(acc *validator.Validator) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.validatorStore.UpdateValidator(acc)
}
