package state

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.com/zarb-chain/zarb-go/block"
	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/errors"
	"gitlab.com/zarb-chain/zarb-go/execution"
	"gitlab.com/zarb-chain/zarb-go/genesis"
	merkle "gitlab.com/zarb-chain/zarb-go/libs/merkle"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/network"
	"gitlab.com/zarb-chain/zarb-go/store"
	"gitlab.com/zarb-chain/zarb-go/tx"
	"gitlab.com/zarb-chain/zarb-go/txpool"
	"gitlab.com/zarb-chain/zarb-go/validator"
)

// TODO
// implement go-level db snapshot, and keep latest snap shots
var (
	lastBlockHeightKey = []byte{0x01}
	accountPrefix      = []byte{0x02}
	accountNumPrefix   = []byte{0x04}
	validatorPrefix    = []byte{0x08}
	validatorNumPrefix = []byte{0x10}
)

func accountKey(addr crypto.Address) []byte   { return append(accountPrefix, addr.RawBytes()...) }
func validatorKey(addr crypto.Address) []byte { return append(validatorPrefix, addr.RawBytes()...) }

type State struct {
	lk deadlock.RWMutex

	db                 *leveldb.DB
	config             *config.Config
	store              *store.Store
	txPool             *txpool.TxPool
	cache              *Cache
	executor           *execution.Executor
	validatorSet       *validator.ValidatorSet
	syncer             *synchronizer
	lastBlockHeight    int
	lastBlockHash      crypto.Hash
	lastReceiptsHash   crypto.Hash
	nextValidatorsHash crypto.Hash
	lastBlockTime      time.Time
	updateCh           chan int
	logger             *logger.Logger
}

func LoadStateOrNewState(conf *config.Config, genDoc *genesis.Genesis, net *network.Network, store *store.Store, txPool *txpool.TxPool) (*State, error) {
	db, err := leveldb.OpenFile(conf.Store.StateStorePath(), nil)
	if err != nil {
		return nil, err
	}
	st := &State{
		db:     db,
		config: conf,
		txPool: txPool,
		store:  store,
	}

	st.cache = newCache(store)
	st.executor, err = execution.NewExecutor(st.config, st.cache)
	if err != nil {
		return nil, err
	}

	err = st.loadState()
	if err != nil {
		err = st.makeGenesisState(genDoc)
	}

	if err != nil {
		return nil, err
	}

	st.logger = logger.NewLogger("State", st)

	syncer, err := newSynchronizer(conf, st, store, net, st.logger)
	if err != nil {
		return nil, err
	}
	st.syncer = syncer

	return st, nil
}

func (st *State) Start() error {
	if err := st.syncer.Start(); err != nil {
		return err
	}

	return nil
}

func (st *State) Stop() {
	st.syncer.Stop()
}

func (st *State) loadState() error {

	return fmt.Errorf("temp error")
	//return nil
}

func (st *State) makeGenesisState(genDoc *genesis.Genesis) error {
	accs := genDoc.Accounts()
	for _, acc := range accs {
		st.store.UpdateAccount(acc)
	}

	vals := genDoc.Validators()
	for _, val := range vals {
		st.store.UpdateValidator(val)
	}

	st.validatorSet = validator.NewValidatorSet(vals, len(vals))
	st.lastBlockTime = genDoc.GenesisTime()
	return nil
}

func (st *State) SetNewHeightListener(listener chan int) {
	st.updateCh = listener
}

func (st *State) stateHash() crypto.Hash {
	accRootHash := st.accountsMerkleRootHash()
	valRootHash := st.validatorsMerkleRootHash()

	rootHash := merkle.HashMerkleBranches(accRootHash, valRootHash)
	if rootHash == nil {
		logger.Panic("State hash can't be nil")
	}

	return *rootHash
}

func (st *State) ValidatorSet() *validator.ValidatorSet {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.validatorSet
}

func (st *State) LastBlockHeight() int {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastBlockHeight
}

func (st *State) LastBlockInfo() (int, crypto.Hash) {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastBlockHeight, st.lastBlockHash
}

func (st *State) LastBlockTime() time.Time {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastBlockTime
}

func (st *State) ProposeBlock(height int, proposer crypto.Address, lastCommit *block.Commit) block.Block {
	st.lk.Lock()
	defer st.lk.Unlock()

	timestamp := st.lastBlockTime.Add(st.config.BlockTime())
	now := time.Now()
	if now.After(timestamp) {
		timestamp = now
	}

	mintbaseTx := tx.NewMintbaseTx(st.lastBlockHash, proposer, 10, "Minbase transaction")

	st.txPool.AppendTx(mintbaseTx)

	txs := block.NewTxs()
	txs.Append(mintbaseTx.Hash())
	stateHash := st.stateHash()
	block := block.MakeBlock(
		timestamp,
		txs,
		st.lastBlockHash,
		crypto.UndefHash,
		stateHash,
		st.lastReceiptsHash,
		lastCommit,
		proposer)

	return block
}

func (st *State) SyncTxPool(block *block.Block) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if block == nil {
		return
	}
}

func (st *State) ValidateBlock(block *block.Block) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if block == nil {
		return errors.Error(errors.ErrInvalidBlock)
	}

	if err := st.validateBlock(block); err != nil {
		return err
	}

	st.cache.reset()
	_, err := st.executeBlock(block, st.executor)
	if err != nil {
		return err
	}

	return nil
}

func (st *State) ApplyBlock(block *block.Block, round int) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	err := st.validateBlock(block)
	if err != nil {
		st.logger.Panic("Applying block failed", "err", err)
	}

	st.cache.reset()
	receipts, err := st.executeBlock(block, st.executor)
	if err != nil {
		st.logger.Panic("Applying block failed", "err", err)
	}
	hashes := make([]crypto.Hash, len(receipts))
	for i, r := range receipts {
		hashes[i] = r.Hash()
		st.txPool.RemoveTx(r.TxHash())
	}
	receiptsMerkle := merkle.NewTreeFromHashes(hashes)
	receiptsHash := receiptsMerkle.Root()

	// Commit the changes
	st.cache.commit(nil)

	if err := st.store.SaveBlock(*block, st.lastBlockHeight+1); err != nil {
		return err
	}

	// Move validator set
	st.validatorSet.MoveProposer(round)
	st.lastBlockHeight += 1
	st.lastBlockHash = block.Hash()
	st.lastBlockTime = block.Header().Time()
	st.lastReceiptsHash = *receiptsHash

	return nil
}

func (st *State) Fingerprint() string {
	return fmt.Sprintf("{#%v %v @%v}",
		st.lastBlockHeight,
		st.lastBlockHash.Fingerprint(),
		st.lastBlockTime.Format("15.04.05"))
}
