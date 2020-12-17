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
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
)

// baseSubsidy, one tenth of Bitcoin network
const baseSubsidy = 5 * 1e8

type StateReader interface {
	StoreReader() store.StoreReader
	ValidatorSet() validator.ValidatorSetReader
	LastBlockHeight() int
	GenesisHash() crypto.Hash
	LastBlockHash() crypto.Hash
	LastBlockTime() time.Time
	LastCommit() *block.Commit
	BlockTime() time.Duration
	UpdateLastCommit(blockHash crypto.Hash, commit block.Commit)
	Fingerprint() string
}

type State interface {
	StateReader

	Close() error
	ProposeBlock() block.Block
	ValidateBlock(block block.Block) error
	ApplyBlock(height int, block block.Block, commit block.Commit) error
}

type state struct {
	lk deadlock.RWMutex

	config           *Config
	proposer         crypto.Address
	genDoc           *genesis.Genesis
	store            *store.Store
	params           param.Params
	txPool           txpool.TxPool
	txPoolSandbox    *sandbox.SandboxConcrete
	execution        *execution.Execution
	executionSandbox *sandbox.SandboxConcrete
	validatorSet     *validator.ValidatorSet
	sortition        *sortition.Sortition
	lastBlockHeight  int
	lastBlockHash    crypto.Hash
	lastReceiptsHash crypto.Hash
	lastCommit       *block.Commit
	lastBlockTime    time.Time
	logger           *logger.Logger
}

func LoadOrNewState(
	conf *Config,
	genDoc *genesis.Genesis,
	signer crypto.Signer,
	txPool txpool.TxPool) (State, error) {

	st := &state{
		config:    conf,
		genDoc:    genDoc,
		txPool:    txPool,
		params:    param.NewParams(),
		proposer:  signer.Address(),
		sortition: sortition.NewSortition(signer),
	}
	st.logger = logger.NewLogger("_state", st)

	st.params.BlockTime = genDoc.BlockTime()
	st.params.MaximumPower = genDoc.MaximumPower()
	st.params.MaximumMemoLength = genDoc.MaximumMemoLength()
	st.params.FeeFraction = genDoc.FeeFraction()
	st.params.MinimumFee = genDoc.MinimumFee()
	st.params.TransactionToLiveInterval = genDoc.TTL()

	store, err := store.NewStore(conf.Store)
	if err != nil {
		return nil, err
	}
	st.store = store

	if store.HasAnyBlock() {
		err := st.tryLoadLastInfo()
		if err != nil {
			return nil, err
		}
	} else {
		err := st.makeGenesisState(genDoc)
		if err != nil {
			return nil, err
		}
	}

	st.txPoolSandbox, err = sandbox.NewSandbox(store, st.params, st.lastBlockHeight, st.sortition, st.validatorSet)
	if err != nil {
		return nil, err
	}
	st.executionSandbox, err = sandbox.NewSandbox(store, st.params, st.lastBlockHeight, st.sortition, st.validatorSet)
	if err != nil {
		return nil, err
	}
	st.txPool.SetSandbox(st.txPoolSandbox)
	st.execution = execution.NewExecution(st.executionSandbox)

	return st, nil
}

func (st *state) tryLoadLastInfo() error {
	height, commit, receiptHash, err := st.loadLastInfo()
	if err != nil {
		return err
	}
	logger.Info("Try to load the last state info", "height", height)

	b, err := st.store.Block(height)
	if err != nil {
		return err
	}
	st.lastBlockHeight = height
	st.lastBlockHash = b.Header().Hash()
	st.lastCommit = commit
	st.lastBlockTime = b.Header().Time()
	st.lastReceiptsHash = *receiptHash

	vals := make([]*validator.Validator, len(commit.Committers()))
	for i, c := range commit.Committers() {
		val, err := st.store.Validator(c.Address)
		if err != nil {
			return fmt.Errorf("Last commit has unknown validator: %v", err)
		}
		vals[i] = val
	}
	st.validatorSet, err = validator.NewValidatorSet(vals, st.params.MaximumPower, b.Header().ProposerAddress())
	if err != nil {
		return err
	}
	if err := st.validatorSet.MoveToNextHeight(0, nil); err != nil {
		return err
	}

	totalStake := int64(0)
	st.store.IterateValidators(func(val *validator.Validator) (stop bool) {
		totalStake += val.Stake()
		return false
	})

	st.sortition.SetTotalStake(totalStake)

	return nil
}

func (st *state) makeGenesisState(genDoc *genesis.Genesis) error {
	accs := genDoc.Accounts()
	for _, acc := range accs {
		st.store.UpdateAccount(acc)
	}

	totalStake := int64(0)
	vals := genDoc.Validators()
	for _, val := range vals {
		st.store.UpdateValidator(val)
		totalStake += val.Stake()
	}

	valSet, err := validator.NewValidatorSet(vals, st.params.MaximumPower, vals[0].Address())
	if err != nil {
		return err
	}
	st.validatorSet = valSet
	st.lastBlockTime = genDoc.GenesisTime()
	st.sortition.SetTotalStake(totalStake)
	return nil
}

func (st *state) Close() error {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.store.Close()
}

func (st *state) StoreReader() store.StoreReader {
	return st.store
}

func (st *state) ValidatorSet() validator.ValidatorSetReader {
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

func (st *state) LastCommit() *block.Commit {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastCommit
}

func (st *state) BlockTime() time.Duration {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.params.BlockTime
}

func (st *state) UpdateLastCommit(blockHash crypto.Hash, commit block.Commit) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if err := st.validateCommitForCurrentHeight(commit, blockHash); err != nil {
		st.logger.Warn("Try to update last commit, but it's invalid", "error", err)
		return
	}

	st.lastCommit = &commit
}

func (st *state) createMintbaseTx() *tx.Tx {
	acc, _ := st.store.Account(crypto.MintbaseAddress)
	stamp := st.lastBlockHash
	seq := acc.Sequence() + 1
	amt := calcBlockSubsidy(st.lastBlockHeight+1, st.params.SubsidyReductionInterval)
	tx := tx.NewMintbaseTx(stamp, seq, st.proposer, amt, "")
	return tx
}

func (st *state) ProposeBlock() block.Block {
	st.lk.Lock()
	defer st.lk.Unlock()

	timestamp := st.lastBlockTime.Add(st.params.BlockTime)
	now := time.Now()
	if now.After(timestamp) {
		timestamp = now
	}

	rewardTx := st.createMintbaseTx()
	if err := st.txPool.AppendTxAndBroadcast(*rewardTx); err != nil {
		st.logger.Panic("Our mintbase transaction is invalid", "err", err)
	}

	txIDs := block.NewTxIDs()
	txIDs.Append(rewardTx.ID())
	stateHash := st.stateHash()
	committersHash := st.validatorSet.CommittersHash()
	block := block.MakeBlock(
		timestamp,
		txIDs,
		st.lastBlockHash,
		committersHash,
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

	st.executionSandbox.Clear()

	_, err := st.executeBlock(block)
	if err != nil {
		return err
	}

	return nil
}

func (st *state) ApplyBlock(height int, block block.Block, commit block.Commit) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if height != st.lastBlockHeight && height != st.lastBlockHeight+1 {
		return errors.Errorf(errors.ErrInvalidBlock, "We are not expecting a block for this height: %v", height)
	}

	/// There are two modules can commit a block: Consensus and Syncer.
	/// Consensus engine is ours, we have full control over that and we know when and why a block should be committed.
	/// In the other hand, Syncer module receives new blocks from other peers and if we are behind them, he tries to commit them.
	/// We should never have a fork in our blockchain. but if it happens here we can catch it.

	if st.lastBlockHeight == height {
		if block.Hash().EqualsTo(st.lastBlockHash) {
			st.logger.Trace("We have committed this block before", "hash", block.Hash())
			return nil
		}

		st.logger.Error("A possible fork is detected", "our hash", st.lastBlockHash, "block hash", block.Hash())
		return errors.Error(errors.ErrInvalidBlock)
	}

	err := st.validateBlock(block)
	if err != nil {
		return err
	}

	err = st.validateCommitForCurrentHeight(commit, block.Hash())
	if err != nil {
		return err
	}

	st.txPoolSandbox.Clear()
	st.executionSandbox.Clear()

	// Execute block
	ctrxs, err := st.executeBlock(block)
	if err != nil {
		return err
	}

	if err := st.store.SaveBlock(block, st.lastBlockHeight+1); err != nil {
		return err
	}

	// Save txs and receipts
	receiptsHashes := make([]crypto.Hash, len(ctrxs))
	for i, ctrx := range ctrxs {
		st.txPool.RemoveTx(ctrx.Tx.ID())
		st.store.SaveTransaction(ctrx)

		receiptsHashes[i] = ctrx.Receipt.Hash()
	}

	// Commit changes and move proposer index
	st.commitSandbox(commit.Round())

	receiptsMerkle := merkle.NewTreeFromHashes(receiptsHashes)

	st.lastBlockHeight++
	st.lastBlockHash = block.Hash()
	st.lastBlockTime = block.Header().Time()
	st.lastReceiptsHash = receiptsMerkle.Root()
	st.lastCommit = &commit

	st.logger.Info("New block is committed", "block", block, "round", commit.Round())

	st.executionSandbox.AppendNewBlock(st.lastBlockHash, st.lastBlockHeight)
	st.txPoolSandbox.AppendNewBlock(st.lastBlockHash, st.lastBlockHeight)
	st.saveLastInfo(st.lastBlockHeight, st.lastCommit, &st.lastReceiptsHash)

	st.EvaluateSortition()

	return nil
}

func (st *state) EvaluateSortition() {
	if st.validatorSet.Contains(st.proposer) {
		// We are in the validator set right now
		return
	}

	val, _ := st.store.Validator(st.proposer)
	if val == nil {
		// We are not a validator
		return
	}
	//
	trx := st.sortition.EvaluateTransaction(st.lastBlockHash, val)
	if trx != nil {
		st.logger.Info("👏 This validator is chosen to be in set", "address", st.proposer, "stake", val.Stake(), "tx", trx)
		if err := st.txPool.AppendTxAndBroadcast(*trx); err != nil {
			st.logger.Error("Our sortition transaction is invalid. Why?", "address", st.proposer, "stake", val.Stake(), "tx", trx)
		}
	}
}

func calcBlockSubsidy(height int, subsidyReductionInterval int) int64 {
	if subsidyReductionInterval == 0 {
		return baseSubsidy
	}

	// Equivalent to: baseSubsidy / 2^(height/subsidyHalvingInterval)
	return baseSubsidy >> uint(height/subsidyReductionInterval)
}

func (st *state) Fingerprint() string {
	return fmt.Sprintf("{#%d ⌘ %v 🕣 %v}",
		st.lastBlockHeight,
		st.lastBlockHash.Fingerprint(),
		st.lastBlockTime.Format("15.04.05"))
}

// TODO: add tests for me
func (st *state) commitSandbox(round int) {
	joined := make([]*validator.Validator, 0)
	st.executionSandbox.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.AddToSet {
			joined = append(joined, &vs.Validator)
		}
	})

	if err := st.validatorSet.MoveToNextHeight(0, joined); err != nil {
		//
		// We should panic here before modifying state store
		//
		logger.Panic("An error occurred", "err", err)
	}

	st.executionSandbox.IterateAccounts(func(as *sandbox.AccountStatus) {
		if as.Updated {
			st.store.UpdateAccount(&as.Account)
		}
	})

	st.executionSandbox.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.Updated {
			st.store.UpdateValidator(&vs.Validator)
		}
	})
}
