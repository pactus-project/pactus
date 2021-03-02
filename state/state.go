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
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

// baseSubsidy, one tenth of Bitcoin network
const baseSubsidy = 5 * 1e8

type state struct {
	lk deadlock.RWMutex

	config            *Config
	signer            crypto.Signer
	mintbaseAddr      crypto.Address
	genDoc            *genesis.Genesis
	store             *store.Store
	params            param.Params
	txPool            txpool.TxPool
	txPoolSandbox     *sandbox.SandboxConcrete
	execution         *execution.Execution
	executionSandbox  *sandbox.SandboxConcrete
	validatorSet      *validator.ValidatorSet
	sortition         *sortition.Sortition
	lastSortitionSeed sortition.Seed
	lastBlockHeight   int
	lastBlockHash     crypto.Hash
	lastReceiptsHash  crypto.Hash
	lastCommit        *block.Commit
	lastBlockTime     time.Time
	logger            *logger.Logger
}

func LoadOrNewState(
	conf *Config,
	genDoc *genesis.Genesis,
	signer crypto.Signer,
	txPool txpool.TxPool) (State, error) {

	var mintbaseAddr crypto.Address
	if conf.MintbaseAddress != "" {
		addr, err := crypto.AddressFromString(conf.MintbaseAddress)
		if err != nil {
			return nil, err
		}
		mintbaseAddr = addr
	} else {
		mintbaseAddr = signer.Address()
	}

	st := &state{
		config:       conf,
		genDoc:       genDoc,
		txPool:       txPool,
		params:       genDoc.Params(),
		signer:       signer,
		mintbaseAddr: mintbaseAddr,
		sortition:    sortition.NewSortition(),
	}
	st.logger = logger.NewLogger("_state", st)

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
	li, err := st.loadLastInfo()
	if err != nil {
		return err
	}
	// Make sure genesis hash is same
	//
	// This check is not important because genesis state is committed.
	// But it is good to have it to make sure genesis doc hasn't changed
	genHash := st.calculateGenesisStateHashFromGenesisDoc()
	blockOne, err := st.store.Block(1)
	if err != nil {
		return err
	}
	if !genHash.EqualsTo(blockOne.Header().StateHash()) {
		return fmt.Errorf("Invalid genesis doc")
	}

	logger.Info("Try to load the last state info", "height", li.LastHeight)

	b, err := st.store.Block(li.LastHeight)
	if err != nil {
		return err
	}
	st.lastBlockHeight = li.LastHeight
	st.lastBlockHash = b.Header().Hash()
	st.lastCommit = &li.LastCommit
	st.lastBlockTime = b.Header().Time()
	st.lastSortitionSeed = b.Header().SortitionSeed()
	st.lastReceiptsHash = li.LastReceiptHash

	vals := make([]*validator.Validator, len(st.lastCommit.Committers()))
	for i, num := range li.Committee {
		val, err := st.store.ValidatorByNumber(num)
		if err != nil {
			return fmt.Errorf("Unknown committee member: %v", err)
		}
		vals[i] = val
	}
	st.validatorSet, err = validator.NewValidatorSet(vals, st.params.CommitteeSize, li.NextProposer)
	if err != nil {
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

	valSet, err := validator.NewValidatorSet(vals, st.params.CommitteeSize, vals[0].Address())
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
	st.lk.RLock()
	defer st.lk.RUnlock()

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

	return st.params.BlockTime()
}

func (st *state) UpdateLastCommit(commit *block.Commit) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	// Check if commit has more signers ...
	if commit.Signers() > st.lastCommit.Signers() {
		if err := st.validateCommitForPreviousHeight(commit); err != nil {
			st.logger.Warn("Try to update last commit, but it's invalid", "err", err)
			return err
		}

		st.lastCommit = commit
	}
	return nil
}

func (st *state) createSubsidyTx(fee int64) *tx.Tx {
	acc, err := st.store.Account(crypto.TreasuryAddress)
	if err != nil {
		return nil
	}
	stamp := st.lastBlockHash
	seq := acc.Sequence() + 1
	amt := calcBlockSubsidy(st.lastBlockHeight+1, st.params.SubsidyReductionInterval)

	tx := tx.NewMintbaseTx(stamp, seq, st.mintbaseAddr, amt+fee, "")
	return tx
}

func (st *state) ProposeBlock(round int) (*block.Block, error) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if !st.validatorSet.IsProposer(st.signer.Address(), round) {
		return nil, errors.Errorf(errors.ErrInvalidAddress, "We are not propser for this round")
	}

	// Reset Sandbox and clear the accululated fee
	st.executionSandbox.Clear()
	st.execution.ResetFee()

	txIDs := block.NewTxIDs()

	// Re-chaeck all transactions again, remove invalid ones
	trxs := st.txPool.AllTransactions()
	for _, trx := range trxs {
		// All subsidy transactions (probably from invalid rounds)
		// should be removed from the pool
		if trx.IsMintbaseTx() {
			st.logger.Debug("Found duplicated subsidy transaction", "tx", trx)
			st.txPool.RemoveTx(trx.ID())
			continue
		}

		if err := st.execution.Execute(trx); err != nil {
			st.logger.Debug("Found invalid transaction", "tx", trx, "err", err)
			st.txPool.RemoveTx(trx.ID())
		} else {
			txIDs.Append(trx.ID())

			if txIDs.Len() >= st.params.MaximumTransactionPerBlock {
				break
			}
		}
	}

	subsidyTx := st.createSubsidyTx(st.execution.AccumulatedFee())
	if subsidyTx == nil {
		st.logger.Error("Probably the node is shutting down.")
		return nil, errors.Errorf(errors.ErrInvalidBlock, "No subsidy transaction")
	}
	if err := st.txPool.AppendTx(subsidyTx); err != nil {
		st.logger.Error("Our subsidy transaction is invalid. Why?", "err", err)
		return nil, err
	}
	txIDs.Prepend(subsidyTx.ID())

	// Broadcast all transaction
	st.txPool.BroadcastTxs(txIDs.IDs())

	stateHash := st.stateHash()
	committeeHash := st.validatorSet.CommitteeHash()
	timestamp := st.proposeNextBlockTime()
	newSortitionSeed := st.lastSortitionSeed.Generate(st.signer)

	block := block.MakeBlock(
		st.params.BlockVersion,
		timestamp,
		txIDs,
		st.lastBlockHash,
		committeeHash,
		stateHash,
		st.lastReceiptsHash,
		st.lastCommit,
		newSortitionSeed,
		st.signer.Address())

	return &block, nil
}

func (st *state) ValidateBlock(block block.Block) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	t := block.Header().Time()
	if err := st.validateBlockTime(t); err != nil {
		return err
	}

	if err := st.validateBlock(block); err != nil {
		return err
	}

	_, err := st.executeBlock(block)
	if err != nil {
		return err
	}

	return nil
}

func (st *state) CommitBlock(height int, block block.Block, commit block.Commit) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if height != st.lastBlockHeight && height != st.lastBlockHeight+1 {
		/// Returning error here will cause so many error logs during syncing blockchain
		/// Syncing is asynchronous job and we might receive blocks not in order
		st.logger.Debug("Unexpected block height", "height", height)
		return nil
	}

	/// There are two modules that can commit a block: Consensus and Syncer.
	/// Consensus engine is ours, we have full control over that and we know when and why a block should be committed.
	/// In the other hand, Syncer module receives new blocks from other peers and if we are behind them, it tries to commit them.
	/// We should never have a fork in our blockchain. but if it happens here we can catch it.
	if st.lastBlockHeight == height {
		if block.Hash().EqualsTo(st.lastBlockHash) {
			st.logger.Debug("This block committed before", "hash", block.Hash())
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

	// Verify proposer
	proposer := st.validatorSet.Proposer(commit.Round())
	if !proposer.Address().EqualsTo(block.Header().ProposerAddress()) {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid proposer. Expected %s, got %s", proposer.Address(), block.Header().ProposerAddress())
	}
	// Validate sortition seed
	if !block.Header().SortitionSeed().Validate(proposer.PublicKey(), st.lastSortitionSeed) {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid sortition seed.")
	}

	ctrxs, err := st.executeBlock(block)
	if err != nil {
		return err
	}

	if err := st.store.SaveBlock(block, st.lastBlockHeight+1); err != nil {
		return err
	}

	// Commit changes and update the validator set
	st.commitSandbox(commit.Round())

	// Save txs and receipts
	receiptsHashes := make([]crypto.Hash, len(ctrxs))
	for i, ctrx := range ctrxs {
		st.txPool.RemoveTx(ctrx.Tx.ID())
		st.store.SaveTransaction(ctrx)

		receiptsHashes[i] = ctrx.Receipt.Hash()
	}
	receiptsMerkle := merkle.NewTreeFromHashes(receiptsHashes)

	st.lastBlockHeight++
	st.lastBlockHash = block.Hash()
	st.lastBlockTime = block.Header().Time()
	st.lastSortitionSeed = block.Header().SortitionSeed()
	st.lastReceiptsHash = receiptsMerkle.Root()
	st.lastCommit = &commit

	// Evaluate sortition before updating the validator set
	if st.evaluateSortition() {
		st.logger.Info("👏 This validator is chosen to be in the set", "address", st.signer.Address())
	}

	st.logger.Info("New block is committed", "block", block, "round", commit.Round())

	st.executionSandbox.AppendNewBlock(st.lastBlockHash, st.lastBlockHeight)
	st.txPoolSandbox.AppendNewBlock(st.lastBlockHash, st.lastBlockHeight)
	st.saveLastInfo(st.lastBlockHeight, commit, st.lastReceiptsHash,
		st.validatorSet.Committee(),
		st.validatorSet.Proposer(0).Address())

	// At this point we can reset txpool sandbox
	st.txPoolSandbox.Clear()
	st.txPool.Recheck()

	return nil
}

func (st *state) evaluateSortition() bool {
	if st.validatorSet.Contains(st.signer.Address()) {
		// We are in the validator set right now
		return false
	}

	val, _ := st.store.Validator(st.signer.Address())
	if val == nil {
		// We are not a validator
		return false
	}

	if st.lastBlockHeight-val.BondingHeight() < 2*st.params.CommitteeSize {
		// Bonding period
		return false
	}

	//
	ok, proof := st.sortition.EvaluateSortition(st.lastSortitionSeed, st.signer, val.Stake())
	if ok {
		//
		trx := tx.NewSortitionTx(st.lastBlockHash, val.Sequence()+1, val.Address(), proof)
		st.signer.SignMsg(trx)

		if err := st.txPool.AppendTxAndBroadcast(trx); err != nil {
			st.logger.Error("Our sortition transaction is invalid. Why?", "address", st.signer.Address(), "stake", val.Stake(), "tx", trx, "err", err)
			return false
		}
	}

	return true
}

func calcBlockSubsidy(height int, subsidyReductionInterval int) int64 {
	// Equivalent to: baseSubsidy / 2^(height/subsidyHalvingInterval)
	return baseSubsidy >> uint(height/subsidyReductionInterval)
}

func (st *state) Fingerprint() string {
	return fmt.Sprintf("{#%d ⌘ %v 🕣 %v}",
		st.lastBlockHeight,
		st.lastBlockHash.Fingerprint(),
		st.lastBlockTime.Format("15.04.05"))
}

func (st *state) commitSandbox(round int) {
	joined := make([]*validator.Validator, 0)
	st.executionSandbox.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.AddToSet {
			st.logger.Info("New validator joined", "address", vs.Validator.Address(), "stake", vs.Validator.Stake())

			joined = append(joined, &vs.Validator)
		}
	})

	// TODO: for joined vals write tests
	if err := st.validatorSet.UpdateTheSet(round, joined); err != nil {
		//
		// We should panic here before updating state
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

	st.sortition.AddToTotalStake(st.executionSandbox.TotalStakeChange())
}

func (st *state) validateBlockTime(t time.Time) error {
	if t.Second()%st.params.BlockTimeInSecond != 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Block time is not rounded")
	}
	if t.Before(st.lastBlockTime.Add(1 * time.Second)) {
		return errors.Errorf(errors.ErrInvalidBlock, "Block time is too early")
	}
	proposeTime := st.proposeNextBlockTime()
	threshold := 2 * st.params.BlockTime()
	if t.After(proposeTime.Add(threshold)) {
		fmt.Println(t)
		fmt.Println(util.RoundNow(st.params.BlockTimeInSecond).Add(threshold))
		return errors.Errorf(errors.ErrInvalidBlock, "Block time is too far")
	}

	return nil
}

func (st *state) proposeNextBlockTime() time.Time {
	timestamp := st.lastBlockTime.Add(st.params.BlockTime())
	timestamp = util.RoundTime(timestamp, st.params.BlockTimeInSecond)

	now := util.Now()
	if now.After(timestamp.Add(1 * time.Second)) {
		st.logger.Debug("It looks the last commit had delay", "delay", now.Sub(timestamp))
		timestamp = util.RoundNow(st.params.BlockTimeInSecond)
	}
	return timestamp
}
