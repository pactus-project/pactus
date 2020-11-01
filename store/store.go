package store

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type Store struct {
	lk deadlock.RWMutex

	config         *Config
	blockStore     *blockStore
	txStore        *txStore
	accountStore   *accountStore
	validatorStore *validatorStore
	logger         *logger.Logger
}

func NewStore(conf *Config, logger *logger.Logger) (*Store, error) {
	blockStore, err := newBlockStore(conf.BlockStorePath(), logger)
	if err != nil {
		return nil, err
	}
	txStore, err := newTxStore(conf.TxStorePath(), logger)
	if err != nil {
		return nil, err
	}
	accountStore, err := newAccountStore(conf.AccountStorePath())
	if err != nil {
		return nil, err
	}
	validatorStore, err := newValidatorStore(conf.ValidatorStorePath())
	if err != nil {
		return nil, err
	}
	return &Store{
		config:         conf,
		blockStore:     blockStore,
		txStore:        txStore,
		accountStore:   accountStore,
		validatorStore: validatorStore,
		logger:         logger,
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

	height, err := s.blockStore.blockHeight(hash)
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

	return s.blockStore.blockHeight(hash)
}

func (s *Store) SaveTx(tx tx.Tx, receipt tx.Receipt) error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.SaveTx(tx, receipt)
}

func (s *Store) Tx(hash crypto.Hash) (*tx.Tx, *tx.Receipt, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.Tx(hash)
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
