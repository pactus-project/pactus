package store

import (
	"bytes"
	"errors"
	"sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrBadOffset = errors.New("offset is out of range")
)

const (
	lastStoreVersion = int32(1)
)

var (
	lastInfoKey       = []byte{0x00}
	blockPrefix       = []byte{0x01}
	txPrefix          = []byte{0x03}
	accountPrefix     = []byte{0x05}
	validatorPrefix   = []byte{0x07}
	blockHeightPrefix = []byte{0x09}
	publicKeyPrefix   = []byte{0x0b}
)

func tryGet(db *leveldb.DB, key []byte) ([]byte, error) {
	data, err := db.Get(key, nil)
	if err != nil {
		// Probably key doesn't exist in database
		logger.Trace("database `get` error", "error", err, "key", key)

		return nil, err
	}

	return data, nil
}

func tryHas(db *leveldb.DB, key []byte) bool {
	has, err := db.Has(key, nil)
	if err != nil {
		logger.Error("database `has` error", "error", err, "key", key)

		return false
	}

	return has
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
	isPruned       bool
}

func NewStore(conf *Config) (Store, error) {
	options := &opt.Options{
		Strict:      opt.DefaultStrict,
		Compression: opt.NoCompression,
	}

	db, err := leveldb.OpenFile(conf.StorePath(), options)
	if err != nil {
		return nil, err
	}
	store := &store{
		config:         conf,
		db:             db,
		batch:          new(leveldb.Batch),
		blockStore:     newBlockStore(db, conf.SeedCacheWindow, conf.PublicKeyCacheSize),
		txStore:        newTxStore(db, conf.TxCacheWindow),
		accountStore:   newAccountStore(db, conf.AccountCacheSize),
		validatorStore: newValidatorStore(db),
		isPruned:       false,
	}

	lastCert := store.lastCertificate()
	if lastCert == nil {
		return store, nil
	}

	// Check if the node is pruned by checking genesis block.
	cBlkOne, _ := store.block(1)
	if cBlkOne == nil {
		store.isPruned = true
	}

	currentHeight := lastCert.Height()
	startHeight := uint32(1)
	if currentHeight > conf.TxCacheWindow {
		startHeight = currentHeight - conf.TxCacheWindow
	}

	for height := startHeight; height < currentHeight+1; height++ {
		cBlk, err := store.block(height)
		if err != nil {
			return nil, err
		}
		blk, err := cBlk.ToBlock()
		if err != nil {
			return nil, err
		}

		txs := blk.Transactions()
		for _, transaction := range txs {
			store.txStore.addToCache(transaction.ID(), height)
		}

		sortitionSeed := blk.Header().SortitionSeed()
		store.blockStore.addToCache(height, sortitionSeed)
	}

	return store, nil
}

func (s *store) Close() {
	s.lk.Lock()
	defer s.lk.Unlock()

	err := s.db.Close()
	if err != nil {
		logger.Error("error on closing store", "error", err)
	}
}

func (s *store) SaveBlock(blk *block.Block, cert *certificate.BlockCertificate) {
	s.lk.Lock()
	defer s.lk.Unlock()

	height := cert.Height()
	regs := s.blockStore.saveBlock(s.batch, height, blk)
	s.txStore.saveTxs(s.batch, blk.Transactions(), regs)
	s.txStore.pruneCache(height)

	// Removing old block from prune node store.
	if s.isPruned && height > s.config.RetentionBlocks() {
		pruneHeight := height - s.config.RetentionBlocks()
		deleted, err := s.pruneBlock(pruneHeight)
		if err != nil {
			panic(err)
		}

		if deleted {
			// TODO: Let's use state logger in store[?].
			logger.Debug("old block is pruned", "height", pruneHeight)
		} else {
			logger.Warn("unable to prune the old block", "height", pruneHeight, "error", err)
		}
	}

	// Save last certificate: [version: 4 bytes]+[certificate: variant]
	buf := bytes.NewBuffer(make([]byte, 0, 4+cert.SerializeSize()))
	err := encoding.WriteElements(buf, lastStoreVersion)
	if err != nil {
		panic(err)
	}
	err = cert.Encode(buf)
	if err != nil {
		panic(err)
	}

	s.batch.Put(lastInfoKey, buf.Bytes())
}

func (s *store) Block(height uint32) (*CommittedBlock, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.block(height)
}

func (s *store) block(height uint32) (*CommittedBlock, error) {
	data, err := s.blockStore.block(height)
	if err != nil {
		return nil, err
	}

	blockHash, err := hash.FromBytes(data[0:hash.HashSize])
	if err != nil {
		return nil, err
	}

	return &CommittedBlock{
		store:     s,
		BlockHash: blockHash,
		Height:    height,
		Data:      data[hash.HashSize:],
	}, nil
}

func (s *store) BlockHeight(h hash.Hash) uint32 {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.blockHeight(h)
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

func (s *store) SortitionSeed(blockHeight uint32) *sortition.VerifiableSeed {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.blockStore.sortitionSeed(blockHeight)
}

func (s *store) PublicKey(addr crypto.Address) (crypto.PublicKey, error) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.blockStore.publicKey(addr)
}

func (s *store) HasPublicKey(addr crypto.Address) bool {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return tryHas(s.db, publicKeyKey(addr))
}

func (s *store) Transaction(txID tx.ID) (*CommittedTx, error) {
	s.lk.Lock()
	defer s.lk.Unlock()

	pos, err := s.txStore.tx(txID)
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
		return nil, ErrBadOffset
	}
	blockTime := util.SliceToUint32(data[hash.HashSize+1 : hash.HashSize+5])

	return &CommittedTx{
		store:     s,
		TxID:      txID,
		Height:    pos.height,
		BlockTime: blockTime,
		Data:      data[start:end],
	}, nil
}

// RecentTransaction checks if there is a transaction with the given ID
// within the last 8640 blocks.
// The time window for recent transactions is determined by the
// TransactionToLive interval, which is part of the consensus parameters.
func (s *store) RecentTransaction(txID tx.ID) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.txStore.recentTransaction(txID)
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

func (s *store) IterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool)) {
	s.lk.RLock()
	defer s.lk.RUnlock()

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
	s.lk.RLock()
	defer s.lk.RUnlock()

	s.validatorStore.iterateValidators(consumer)
}

func (s *store) UpdateValidator(acc *validator.Validator) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.validatorStore.updateValidator(s.batch, acc)
}

func (s *store) LastCertificate() *certificate.BlockCertificate {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.lastCertificate()
}

func (s *store) lastCertificate() *certificate.BlockCertificate {
	data, _ := tryGet(s.db, lastInfoKey)
	if data == nil {
		// Genesis block
		return nil
	}
	reader := bytes.NewReader(data)
	version := int32(0)
	cert := new(certificate.BlockCertificate)
	err := encoding.ReadElements(reader, &version)
	if err != nil {
		return nil
	}
	err = cert.Decode(reader)
	if err != nil {
		return nil
	}

	return cert
}

func (s *store) WriteBatch() error {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.writeBatch()
}

func (s *store) writeBatch() error {
	if err := s.db.Write(s.batch, nil); err != nil {
		// TODO: Should we panic here?
		// The store is unreliable if the stored data does not match the cached data.
		return err
	}
	s.batch.Reset()

	return nil
}

func (s *store) IsBanned(addr crypto.Address) bool {
	return s.config.BannedAddrs[addr]
}

// IsPruned returns true if the store is in prune mode, otherwise false.
func (s *store) IsPruned() bool {
	return s.isPruned
}

// PruningHeight returns the height at which blocks will be pruned if the store is in prune mode.
// If the store is not in prune mode, it returns 0.
func (s *store) PruningHeight() uint32 {
	s.lk.RLock()
	defer s.lk.RUnlock()

	if !s.isPruned {
		return 0
	}

	// TODO: it can be optimized (and safer?) by keeping the last block height in memory.
	cert := s.lastCertificate()

	return cert.Height() - s.config.RetentionBlocks()
}

// Prune iterates over all blocks from the pruning height to the genesis block and prunes them.
// The pruning height is `LastBlockHeight - RetentionBlocks`.
// The callback function is called after each block is pruned and can cancel the process.
func (s *store) Prune(callback func(pruned bool, pruningHeight uint32) bool) error {
	s.lk.Lock()
	defer s.lk.Unlock()

	cert := s.lastCertificate()

	// Store is at the genesis height
	if cert == nil {
		return nil
	}

	retentionBlocks := s.config.RetentionBlocks()
	if cert.Height() < retentionBlocks {
		return nil
	}

	pruningHeight := cert.Height() - retentionBlocks
	for height := pruningHeight; height >= 1; height-- {
		deleted, err := s.pruneBlock(height)
		if err != nil {
			return err
		}

		if err := s.writeBatch(); err != nil {
			return err
		}

		if callback(deleted, height) {
			// canceled
			break
		}
	}

	return nil
}

// pruneBlock removes a block and all transactions inside the block from the store.
// It accepts a block height to prune, and returns a boolean that
// indicate whether the block at the specified height existed and pruned,
// or did not exist, along with any encountered errors.
func (s *store) pruneBlock(blockHeight uint32) (bool, error) {
	if !s.blockStore.hasBlock(blockHeight) {
		return false, nil
	}

	cBlock, err := s.block(blockHeight)
	if err != nil {
		return false, err
	}

	blk, err := block.FromBytes(cBlock.Data)
	if err != nil {
		return false, err
	}

	s.batch.Delete(blockHashKey(blk.Hash()))
	s.batch.Delete(blockKey(blockHeight))

	for _, t := range blk.Transactions() {
		s.batch.Delete(t.ID().Bytes())
	}

	return true, nil
}

func (s *store) XeggexAccount() *XeggexAccount {
	return &s.config.XeggexAccount
}
