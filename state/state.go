package state

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/genesis"
	merkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
)

type StoreReader interface {
	BlockByHeight(height int) (*block.Block, error)
	BlockByHash(hash crypto.Hash) (*block.Block, int, error)
	BlockHeight(hash crypto.Hash) (int, error)
	Tx(hash crypto.Hash) (*tx.Tx, *tx.Receipt, error)
}

type StateReader interface {
	StoreReader() StoreReader
	ValidatorSet() *validator.ValidatorSet
	LastBlockHeight() int
	GenesisHash() crypto.Hash
	LastBlockHash() crypto.Hash
	LastBlockTime() time.Time
	BlockTime() time.Duration
	UpdateLastCommit(blockHash crypto.Hash, commit block.Commit)
	Fingerprint() string
}

type State interface {
	StateReader

	ProposeBlock() block.Block
	ValidateBlock(block block.Block) error
	ApplyBlock(block block.Block, commit block.Commit) error
}

type state struct {
	lk deadlock.RWMutex

	proposer           crypto.Address
	genDoc             *genesis.Genesis
	store              *store.Store
	txPool             *txpool.TxPool
	cache              *Cache
	params             *Params
	executor           *execution.Executor
	validatorSet       *validator.ValidatorSet
	lastBlockHeight    int
	lastBlockHash      crypto.Hash
	lastReceiptsHash   crypto.Hash
	lastCommit         *block.Commit
	nextValidatorsHash crypto.Hash
	lastBlockTime      time.Time
	updateCh           chan int
	logger             *logger.Logger
}

func LoadOrNewState(
	conf *Config,
	genDoc *genesis.Genesis,
	proposer crypto.Address,
	txPool *txpool.TxPool) (State, error) {

	st := &state{
		genDoc:   genDoc,
		proposer: proposer,
		txPool:   txPool,
		params:   NewParams(),
	}
	st.logger = logger.NewLogger("_state", st)
	store, err := store.NewStore(conf.Store)
	if err != nil {
		return nil, err
	}
	st.store = store

	err = st.loadState()
	if err != nil {
		err = st.makeGenesisState(genDoc)
	}

	st.cache = newCache(store)
	st.executor, err = execution.NewExecutor(st.cache)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (st *state) loadState() error {

	return fmt.Errorf("temp error")
	//return nil
}

func (st *state) makeGenesisState(genDoc *genesis.Genesis) error {
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

func (st *state) stateHash() crypto.Hash {
	accRootHash := st.accountsMerkleRootHash()
	valRootHash := st.validatorsMerkleRootHash()

	rootHash := merkle.HashMerkleBranches(accRootHash, valRootHash)
	if rootHash == nil {
		logger.Panic("State hash can't be nil")
	}

	return *rootHash
}

func (st *state) StoreReader() StoreReader {
	return st.store
}

func (st *state) ValidatorSet() *validator.ValidatorSet {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.validatorSet
}

func (st *state) LastBlockHeight() int {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastBlockHeight
}

func (st *state) GenesisHash() crypto.Hash {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.genDoc.Hash()
}

func (st *state) LastBlockHash() crypto.Hash {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastBlockHash
}

func (st *state) LastBlockTime() time.Time {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastBlockTime
}

func (st *state) BlockTime() time.Duration {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.params.BlockTime
}

func (st *state) UpdateLastCommit(blockHash crypto.Hash, commit block.Commit) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if err := st.validateCommit(blockHash, commit); err != nil {
		st.logger.Warn("Try to update last commit, but it's invalid", "error", err)
		return
	}

	st.lastCommit = &commit
}

func (st *state) ProposeBlock() block.Block {
	st.lk.Lock()
	defer st.lk.Unlock()

	timestamp := st.lastBlockTime.Add(st.params.BlockTime)
	now := time.Now()
	if now.After(timestamp) {
		timestamp = now
	}

	mintbaseTx := tx.NewMintbaseTx(st.lastBlockHash, st.proposer, 10, "Minbase transaction")
	st.txPool.AppendTxAndBroadcast(mintbaseTx)

	txHashes := block.NewTxHashes()
	txHashes.Append(mintbaseTx.Hash())
	stateHash := st.stateHash()
	block := block.MakeBlock(
		timestamp,
		txHashes,
		st.lastBlockHash,
		crypto.UndefHash,
		stateHash,
		st.lastReceiptsHash,
		st.lastCommit,
		st.proposer)

	return block
}

func (st *state) ValidateBlock(block block.Block) error {
	st.lk.Lock()
	defer st.lk.Unlock()

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

func (st *state) ApplyBlock(block block.Block, commit block.Commit) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if !block.Header().LastBlockHash().EqualsTo(st.lastBlockHash) {
		return errors.Errorf(errors.ErrInvalidBlock, "Previous block hash does not match. Should be %v, but got %v",
			st.lastBlockHash, block.Header().LastBlockHash())
	}

	round := commit.Round()
	err := st.validateBlock(block)
	if err != nil {
		return err
	}

	err = st.validateCommit(block.Hash(), commit)
	if err != nil {
		return err
	}

	st.cache.reset()
	// Execute block
	receipts, err := st.executeBlock(block, st.executor)
	if err != nil {
		return err
	}
	// Commit the changes
	st.cache.commit(nil)

	// Save block and txs
	receiptsHashes := make([]crypto.Hash, len(receipts))
	for i, r := range receipts {
		receiptsHashes[i] = r.Hash()
		trx := st.txPool.RemoveTx(r.TxHash())
		if trx == nil {
			return errors.Errorf(errors.ErrInvalidBlock, "Saving block failed: Transaction lost")
		}
		st.store.SaveTx(*trx, *r)
	}

	if err := st.store.SaveBlock(block, st.lastBlockHeight+1); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, "Saving block failed: %v", err)
	}

	receiptsMerkle := merkle.NewTreeFromHashes(receiptsHashes)
	receiptsHash := receiptsMerkle.Root()

	// Move validator set
	st.validatorSet.MoveProposer(round)
	st.lastBlockHeight += 1
	st.lastBlockHash = block.Hash()
	st.lastBlockTime = block.Header().Time()
	st.lastReceiptsHash = *receiptsHash
	st.lastCommit = &commit

	return nil
}

func (st *state) Fingerprint() string {
	return fmt.Sprintf("{# %v ⌘ %v 🕣 %v}",
		st.lastBlockHeight,
		st.lastBlockHash.Fingerprint(),
		st.lastBlockTime.Format("15.04.05"))
}
