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

type store struct {
	lk deadlock.RWMutex

	config         *Config
	blockStore     *blockStore
	txStore        *txStore
	accountStore   *accountStore
	validatorStore *validatorStore
}

func NewStore(conf *Config) (Store, error) {
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
	return &store{
		config:         conf,
		blockStore:     blockStore,
		txStore:        txStore,
		accountStore:   accountStore,
		validatorStore: validatorStore,
	}, nil
}
func (s *store) Close() error {
	if err := s.blockStore.close(); err != nil {
		return err
	}
	if err := s.txStore.close(); err != nil {
		return err
	}
	if err := s.accountStore.close(); err != nil {
		return err
	}
	if err := s.validatorStore.close(); err != nil {
		return err
	}

	return nil
}

func (s *store) SaveBlock(block block.Block, height int) error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.saveBlock(block, height)
}

func (s *store) Block(height int) (*block.Block, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.block(height)
}

func (s *store) BlockHeight(hash crypto.Hash) (int, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.blockHeight(hash)
}

func (s *store) SaveTransaction(ctrx tx.CommittedTx) {
	s.lk.Lock()
	defer s.lk.Unlock()

	if err := s.txStore.saveTx(ctrx); err != nil {
		logger.Panic("Error on saving a transaction: %v", err)
	}
}

func (s *store) Transaction(hash crypto.Hash) (*tx.CommittedTx, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.tx(hash)
}

func (s *store) HasAccount(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.hasAccount(addr)
}

func (s *store) Account(addr crypto.Address) (*account.Account, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.account(addr)
}

func (s *store) TotalAccounts() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.total
}

func (s *store) IterateAccounts(consumer func(*account.Account) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.iterateAccounts(consumer)
}

func (s *store) UpdateAccount(acc *account.Account) {
	s.lk.Lock()
	defer s.lk.Unlock()

	if err := s.accountStore.updateAccount(acc); err != nil {
		logger.Panic("Error on updating an account: %v", err)
	}
}

func (s *store) HasValidator(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.hasValidator(addr)
}

func (s *store) Validator(addr crypto.Address) (*validator.Validator, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.validator(addr)
}

func (s *store) ValidatorByNumber(num int) (*validator.Validator, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.validatorByNumber(num)
}

func (s *store) TotalValidators() int {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.total
}

func (s *store) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.validatorStore.iterateValidators(consumer)
}

func (s *store) UpdateValidator(acc *validator.Validator) {
	s.lk.Lock()
	defer s.lk.Unlock()

	if err := s.validatorStore.updateValidator(acc); err != nil {
		logger.Panic("Error on updating a validator: %v", err)
	}
}

func (s *store) HasAnyBlock() bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.hasAnyBlock()
}
