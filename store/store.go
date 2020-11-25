package store

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
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
}

func NewStore(conf *Config) (*Store, error) {
	blockStore, err := newBlockStore(conf.BlockStorePath())
	if err != nil {
		return nil, err
	}
	txStore, err := newTxStore(conf.TxStorePath())
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
	}, nil
}

func (s *Store) SaveBlock(block block.Block, height int) error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.saveBlock(block, height)
}

func (s *Store) BlockByHeight(height int) (*block.Block, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.block(height)
}

func (s *Store) BlockByHash(hash crypto.Hash) (*block.Block, int, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	height, err := s.blockStore.blockHeight(hash)
	if err != nil {
		return nil, -1, err
	}
	block, err := s.blockStore.block(height)
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

	return s.txStore.saveTx(tx, receipt)
}

func (s *Store) Tx(hash crypto.Hash) (*tx.Tx, *tx.Receipt, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.tx(hash)
}

func (s *Store) HasAccount(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.hasAccount(addr)
}

func (s *Store) Account(addr crypto.Address) (*account.Account, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.account(addr)
}

func (s *Store) AccountLen() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.len()
}

func (s *Store) IterateAccounts(consumer func(*account.Account) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.iterateAccounts(consumer)
}

func (s *Store) UpdateAccount(acc *account.Account) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.updateAccount(acc)
}

func (s *Store) HasValidator(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.hasValidator(addr)
}

func (s *Store) Validator(addr crypto.Address) (*validator.Validator, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.validator(addr)
}

func (s *Store) ValidatorLen() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.len()
}

func (s *Store) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.validatorStore.iterateValidators(consumer)
}

func (s *Store) UpdateValidator(acc *validator.Validator) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.validatorStore.updateValidator(acc)
}

func (s *Store) LastBlockHeight() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.lastHeight()
}
