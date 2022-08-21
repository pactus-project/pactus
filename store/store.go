package store

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/linkedmap"
	"github.com/zarbchain/zarb-go/util/logger"
)

const lasteStoreVersion = int32(1)

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

type hashPair struct {
	Height uint32
	Hash   hash.Hash
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
	stampLookup    *linkedmap.LinkedMap
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

	s := &store{
		config:         conf,
		db:             db,
		batch:          new(leveldb.Batch),
		blockStore:     newBlockStore(db),
		txStore:        newTxStore(db),
		accountStore:   newAccountStore(db),
		validatorStore: newValidatorStore(db),
		stampLookup:    linkedmap.NewLinkedMap(stampLookupCapacity),
	}

	lastHeight, _ := s.LastCertificate()
	height := uint32(0)
	if lastHeight > uint32(stampLookupCapacity) {
		height = lastHeight - uint32(stampLookupCapacity)
	}
	for ; height <= lastHeight; height++ {
		hash := s.BlockHash(height)
		s.appendStamp(hash, height)
	}

	return s, nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) appendStamp(hash hash.Hash, height uint32) {
	pair := &hashPair{
		Height: height,
		Hash:   hash,
	}
	s.stampLookup.PushBack(hash.Stamp(), pair)
}

func (s *store) SaveBlock(height uint32, block *block.Block, cert *block.Certificate) {
	s.lk.Lock()
	defer s.lk.Unlock()

	txsPos := s.blockStore.saveBlock(s.batch, height, block)

	for i, trx := range block.Transactions() {
		s.txStore.saveTx(s.batch, trx.ID(), &txsPos[i])
	}

	// Save last certificate
	w := bytes.NewBuffer(make([]byte, 0, 8+cert.SerializeSize()))
	err := encoding.WriteElements(w, lasteStoreVersion, height)
	if err != nil {
		panic(err)
	}
	err = cert.Encode(w)
	if err != nil {
		panic(err)
	}

	s.batch.Put(lastInfoKey, w.Bytes())

	// Update stamp to height lookup
	s.appendStamp(block.Hash(), height)
}

func (s *store) Block(hash hash.Hash) (*StoredBlock, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	data, err := s.blockStore.block(hash)
	if err != nil {
		return nil, err
	}

	return &StoredBlock{
		height: util.SliceToUint32(data[:4]),
		data:   data[4:],
	}, nil
}

func (s *store) BlockHash(height uint32) hash.Hash {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.BlockHash(height)
}

func (s *store) FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool) {
	s.lk.Lock()
	defer s.lk.Unlock()

	if stamp.EqualsTo(hash.UndefHash.Stamp()) {
		return hash.UndefHash, true
	}

	v, ok := s.stampLookup.Get(stamp)
	if ok {
		return v.(*hashPair).Hash, true
	}
	return hash.UndefHash, false
}

func (s *store) FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool) {
	s.lk.Lock()
	defer s.lk.Unlock()

	if stamp.EqualsTo(hash.UndefHash.Stamp()) {
		return 0, true
	}

	v, ok := s.stampLookup.Get(stamp)
	if ok {
		return v.(*hashPair).Height, true
	}
	return 0, false
}

func (s *store) Transaction(id tx.ID) (*tx.Tx, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	pos, err := s.txStore.tx(id)
	if err != nil {
		return nil, err
	}
	data, err := s.blockStore.block(pos.Hash)
	if err != nil {
		return nil, err
	}
	if pos.Offset >= int32(len(data)) {
		return nil, fmt.Errorf("offset is out of range") // TODO: Shall we panic here?
	}
	r := bytes.NewReader(data[pos.Offset:])
	trx := new(tx.Tx)
	err = trx.Decode(r)
	if err != nil {
		return nil, err
	}
	if trx.ID() != id {
		return nil, fmt.Errorf("transaction id is not matched") // TODO: Shall we panic here?
	}
	return trx, nil
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

func (s *store) TotalAccounts() int32 {
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

	s.accountStore.updateAccount(s.batch, acc)
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
