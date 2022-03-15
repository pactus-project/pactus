package store

import (
	"fmt"
	"sync"

	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

const lasteStoreVersion = 1

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
	Height int       `cbor:"1,keyasint"`
	Hash   hash.Hash `cbor:"2,keyasint"`
}

type lastInfo struct {
	// Version keeps the store version and helps us to upgrade the store, if needed
	Version int                `cbor:"1,keyasint"`
	Height  int                `cbor:"2,keyasint"`
	Cert    *block.Certificate `cbor:"3,keyasint"`
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

	for height := lastHeight - stampLookupCapacity; height <= lastHeight; height++ {
		if height > 0 {
			hash := s.BlockHash(height)
			s.appendStamp(hash, height)
		}
	}

	return s, nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) appendStamp(hash hash.Hash, height int) {
	pair := &hashPair{
		Height: height,
		Hash:   hash,
	}
	s.stampLookup.PushBack(hash.Stamp(), pair)
}

func (s *store) SaveBlock(height int, block *block.Block, cert *block.Certificate) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.blockStore.saveBlock(s.batch, height, block, cert)

	for index, trx := range block.Transactions() {
		pos := &txPos{
			Height: height,
			Index:  index,
		}
		s.txStore.saveTx(s.batch, trx.ID(), pos)
	}

	// Save last certificate
	lc := lastInfo{
		Version: lasteStoreVersion,
		Height:  height,
		Cert:    cert,
	}
	lastCertData, err := cbor.Marshal(lc)
	if err != nil {
		logger.Panic("unable to encode last certificate: %v", err)
	}
	s.batch.Put(lastInfoKey, lastCertData)

	// Update stamp to height lookup
	s.appendStamp(block.Hash(), height)
}

func (s *store) Block(hash hash.Hash) (*StoreBlock, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	bi, err := s.blockStore.block(hash)
	if err != nil {
		return nil, err
	}

	header := new(block.Header)
	if err := header.Decode(bi.HeaderData); err != nil {
		return nil, err
	}
	var cert *block.Certificate
	if bi.PrevCertData != nil {
		cert = new(block.Certificate)
		if err := cert.Decode(bi.PrevCertData); err != nil {
			return nil, err
		}
	}
	txs := make([]*tx.Tx, len(bi.TransactionsData))
	for i, d := range bi.TransactionsData {
		tx := new(tx.Tx)
		if err := tx.Decode(d); err != nil {
			return nil, err
		}
		txs[i] = tx
	}

	b := block.NewBlock(*header, cert, txs)
	if err := b.SanityCheck(); err != nil {
		return nil, err
	}

	return &StoreBlock{
		Block:      b,
		Height:     bi.Height,
		HeaderData: bi.HeaderData,
	}, nil
}

func (s *store) BlockHash(height int) hash.Hash {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.BlockHash(height)
}

func (s *store) BlockHashByStamp(stamp hash.Stamp) hash.Hash {
	s.lk.Lock()
	defer s.lk.Unlock()

	if stamp.EqualsTo(hash.UndefHash.Stamp()) {
		return hash.UndefHash
	}

	v, ok := s.stampLookup.Get(stamp)
	if ok {
		return v.(*hashPair).Hash
	}
	return hash.UndefHash
}
func (s *store) BlockHeightByStamp(stamp hash.Stamp) int {
	s.lk.Lock()
	defer s.lk.Unlock()

	if stamp.EqualsTo(hash.UndefHash.Stamp()) {
		return 0
	}

	v, ok := s.stampLookup.Get(stamp)
	if ok {
		return v.(*hashPair).Height
	}
	return -1
}

func (s *store) Transaction(id tx.ID) (*tx.Tx, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	pos, err := s.txStore.tx(id)
	if err != nil {
		return nil, err
	}
	blockHash := s.blockStore.BlockHash(pos.Height)
	block, err := s.blockStore.block(blockHash)
	if err != nil {
		return nil, err
	}
	if pos.Index >= len(block.TransactionsData) {
		return nil, fmt.Errorf("index is out of range") // TODO: Shall we panic here?
	}
	tx := new(tx.Tx)
	err = tx.Decode(block.TransactionsData[pos.Index])
	if err != nil {
		return nil, err
	}
	if tx.ID() != id {
		return nil, fmt.Errorf("transaction id is not matched") // TODO: Shall we panic here?
	}
	return tx, nil
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

	s.validatorStore.updateValidator(s.batch, acc)
}

func (s *store) LastCertificate() (int, *block.Certificate) {
	s.lk.Lock()
	defer s.lk.Unlock()

	data, _ := tryGet(s.db, lastInfoKey)
	if data == nil {
		// Genesis block
		return 0, nil
	}
	lc := new(lastInfo)
	err := cbor.Unmarshal(data, lc)
	if err != nil {
		// TODO: should panic here?
		return -1, nil
	}
	return lc.Height, lc.Cert

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
