package store

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/linkedmap"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const lastStoreVersion = int32(1)

// TODO: add cache for me

var (
	lastInfoKey       = []byte{0x00}
	blockPrefix       = []byte{0x01}
	txPrefix          = []byte{0x03}
	accountPrefix     = []byte{0x05}
	validatorPrefix   = []byte{0x07}
	blockHeightPrefix = []byte{0x09}
)

func tryGet(db *leveldb.DB, key []byte) ([]byte, error) {
	data, err := db.Get(key, nil)
	if err != nil {
		// Probably key doesn't exist in database
		logger.Trace("database error", "err", err, "key", key)
		return nil, err
	}
	return data, nil
}

type blockHeightPair struct {
	Height uint32
	Block  *block.Block
}

type store struct {
	lk sync.RWMutex

	config         *Config
	db             *leveldb.DB
	batch          *leveldb.Batch
	blockStore     *blockStore
	txStore        *txStore
	accountStore   *accountStore
	validatorStore *validatorStore
	stampLookup    *linkedmap.LinkedMap[hash.Stamp, blockHeightPair]
}

func NewStore(conf *Config, stampLookupCapacity int) (Store, error) {
	options := &opt.Options{
		Strict:      opt.DefaultStrict,
		Compression: opt.NoCompression,
	}
	db, err := leveldb.OpenFile(conf.StorePath(), options)
	if err != nil {
		return nil, err
	}
	stampLookup := linkedmap.NewLinkedMap[hash.Stamp, blockHeightPair](stampLookupCapacity)
	s := &store{
		config:         conf,
		db:             db,
		batch:          new(leveldb.Batch),
		blockStore:     newBlockStore(db),
		txStore:        newTxStore(db),
		accountStore:   newAccountStore(db),
		validatorStore: newValidatorStore(db),
		stampLookup:    stampLookup,
	}

	lastHeight, _ := s.LastCertificate()
	height := uint32(1)
	if lastHeight > uint32(stampLookupCapacity) {
		height = lastHeight - uint32(stampLookupCapacity)
	}
	for ; height <= lastHeight; height++ {
		storedBlock, _ := s.Block(height)
		s.updateStampLookup(height, storedBlock.ToBlock())
	}

	return s, nil
}

func (s *store) Close() error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.db.Close()
}

func (s *store) updateStampLookup(height uint32, block *block.Block) {
	pair := blockHeightPair{
		Height: height,
		Block:  block,
	}
	s.stampLookup.PushBack(block.Stamp(), pair)
}

func (s *store) SaveBlock(height uint32, block *block.Block, cert *block.Certificate) {
	s.lk.Lock()
	defer s.lk.Unlock()

	reg := s.blockStore.saveBlock(s.batch, height, block)
	for i, trx := range block.Transactions() {
		s.txStore.saveTx(s.batch, trx.ID(), &reg[i])
	}

	// Save last certificate
	w := bytes.NewBuffer(make([]byte, 0, 8+cert.SerializeSize()))
	err := encoding.WriteElements(w, lastStoreVersion, height)
	if err != nil {
		panic(err)
	}
	err = cert.Encode(w)
	if err != nil {
		panic(err)
	}

	s.batch.Put(lastInfoKey, w.Bytes())

	// Update stamp lookup
	s.updateStampLookup(height, block)
}

func (s *store) Block(height uint32) (*StoredBlock, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	data, err := s.blockStore.block(height)
	if err != nil {
		return nil, err
	}

	blockHash, err := hash.FromBytes(data[0:hash.HashSize])
	if err != nil {
		return nil, err
	}

	return &StoredBlock{
		BlockHash: blockHash,
		Height:    height,
		Data:      data[hash.HashSize:],
	}, nil
}

func (s *store) BlockHeight(hash hash.Hash) uint32 {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.blockHeight(hash)
}

func (s *store) BlockHash(height uint32) hash.Hash {
	s.lk.Lock()
	defer s.lk.Unlock()

	data, err := s.blockStore.block(height)
	if err == nil {
		blockHash, _ := hash.FromBytes(data[0:hash.HashSize])
		return blockHash
	}

	return hash.UndefHash
}

func (s *store) RecentBlockByStamp(stamp hash.Stamp) (uint32, *block.Block) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	n := s.stampLookup.GetNode(stamp)
	if n != nil {
		return n.Data.Value.Height, n.Data.Value.Block
	}
	return 0, nil
}

func (s *store) Transaction(id tx.ID) (*StoredTx, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	pos, err := s.txStore.tx(id)
	if err != nil {
		return nil, err
	}
	data, err := s.blockStore.block(pos.height)
	if err != nil {
		return nil, err
	}
	start := pos.offset
	end := pos.offset + pos.length
	if end > uint32(len(data)) {
		return nil, fmt.Errorf("offset is out of range") // TODO: Shall we panic here?
	}
	blockTime := util.SliceToUint32(data[hash.HashSize+1 : hash.HashSize+5])

	return &StoredTx{
		TxID:      id,
		Height:    pos.height,
		BlockTime: blockTime,
		Data:      data[start:end],
	}, nil
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

func (s *store) AccountByNumber(number int32) (*account.Account, error) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.accountStore.accountByNumber(number)
}

func (s *store) TotalAccounts() int32 {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.accountStore.total
}

func (s *store) IterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool)) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.iterateAccounts(consumer)
}

func (s *store) UpdateAccount(addr crypto.Address, acc *account.Account) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.accountStore.updateAccount(s.batch, addr, acc)
}

func (s *store) HasValidator(addr crypto.Address) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.hasValidator(addr)
}

func (s *store) ValidatorAddresses() []crypto.Address {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.validatorStore.ValidatorAddresses()
}

func (s *store) Validator(addr crypto.Address) (*validator.Validator, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.validator(addr)
}

func (s *store) ValidatorByNumber(num int32) (*validator.Validator, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.validatorStore.validatorByNumber(num)
}

func (s *store) TotalValidators() int32 {
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

	s.validatorStore.updateValidator(s.batch, acc)
}

func (s *store) LastCertificate() (uint32, *block.Certificate) {
	s.lk.Lock()
	defer s.lk.Unlock()

	data, _ := tryGet(s.db, lastInfoKey)
	if data == nil {
		// Genesis block
		return 0, nil
	}
	r := bytes.NewReader(data)
	version := int32(0)
	height := uint32(0)
	cert := new(block.Certificate)
	err := encoding.ReadElements(r, &version, &height)
	if err != nil {
		return 0, nil
	}
	err = cert.Decode(r)
	if err != nil {
		return 0, nil
	}
	return height, cert
}

func (s *store) WriteBatch() error {
	s.lk.Lock()
	defer s.lk.Unlock()

	if err := s.db.Write(s.batch, nil); err != nil {
		return err
	}
	s.batch.Reset()
	return nil
}
