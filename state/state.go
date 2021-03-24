package state

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/genesis"
	merkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/state/last_info"
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

	config       *Config
	signer       crypto.Signer
	mintbaseAddr crypto.Address
	genDoc       *genesis.Genesis
	store        store.Store
	params       param.Params
	txPool       txpool.TxPool
	committee    *committee.Committee
	sortition    *sortition.Sortition
	lastInfo     *last_info.LastInfo
	logger       *logger.Logger
}

func LoadOrNewState(
	conf *Config,
	genDoc *genesis.Genesis,
	signer crypto.Signer,
	store store.Store,
	txPool txpool.TxPool) (StateFacade, error) {

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
		store:        store,
		mintbaseAddr: mintbaseAddr,
		sortition:    sortition.NewSortition(),
		lastInfo:     last_info.NewLastInfo(),
	}
	st.logger = logger.NewLogger("_state", st)
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

	txPool.SetNewSandboxAndRecheck(st.makeSandbox())

	return st, nil
}

func (st *state) makeSandbox() *sandbox.SandboxConcrete {
	return sandbox.NewSandbox(st.store, st.params, st.lastInfo.BlockHeight(), st.sortition, st.committee)
}

func (st *state) tryLoadLastInfo() error {
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

	li, err := st.loadLastInfo()
	if err != nil {
		return err
	}
	logger.Info("Try to load the last state info", "height", li.LastHeight)

	b, err := st.store.Block(li.LastHeight)
	if err != nil {
		return err
	}
	st.lastInfo.SetBlockHeight(li.LastHeight)
	st.lastInfo.SetBlockHash(b.Header().Hash())
	st.lastInfo.SetCertificate(li.LastCertificate)
	st.lastInfo.SetBlockTime(b.Header().Time())
	st.lastInfo.SetSortitionSeed(b.Header().SortitionSeed())
	st.lastInfo.SetReceiptsHash(li.LastReceiptHash)

	vals := make([]*validator.Validator, len(st.lastInfo.Certificate().Committers()))
	for i, num := range li.Committee {
		val, err := st.store.ValidatorByNumber(num)
		if err != nil {
			return fmt.Errorf("Unknown committee member: %v", err)
		}
		vals[i] = val
	}
	st.committee, err = committee.NewCommittee(vals, st.params.CommitteeSize, li.NextProposer)
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

	committee, err := committee.NewCommittee(vals, st.params.CommitteeSize, vals[0].Address())
	if err != nil {
		return err
	}
	st.committee = committee
	st.lastInfo.SetBlockTime(genDoc.GenesisTime())
	st.sortition.SetTotalStake(totalStake)
	return nil
}

func (st *state) Close() error {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.store.Close()
}

func (st *state) GenesisHash() crypto.Hash {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.genDoc.Hash()
}

func (st *state) LastBlockHeight() int {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.BlockHeight()
}

func (st *state) LastBlockHash() crypto.Hash {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.BlockHash()
}

func (st *state) LastBlockTime() time.Time {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.BlockTime()
}

func (st *state) LastCertificate() *block.Certificate {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.Certificate()
}

func (st *state) BlockTime() time.Duration {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.params.BlockTime()
}

func (st *state) UpdateLastCertificate(cert *block.Certificate) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	// Check if commit has more signers ...
	if len(cert.Absences()) < len(st.lastInfo.Certificate().Absences()) {
		if err := st.validateCertificateForPreviousHeight(cert); err != nil {
			st.logger.Warn("Try to update last commit, but it's invalid", "err", err)
			return err
		}
		st.lastInfo.SetCertificate(cert)
	}

	return nil
}

func (st *state) createSubsidyTx(fee int64) *tx.Tx {
	acc, err := st.store.Account(crypto.TreasuryAddress)
	if err != nil {
		return nil
	}
	stamp := st.lastInfo.BlockHash()
	seq := acc.Sequence() + 1
	amt := calcBlockSubsidy(st.lastInfo.BlockHeight()+1, st.params.SubsidyReductionInterval)

	tx := tx.NewMintbaseTx(stamp, seq, st.mintbaseAddr, amt+fee, "")
	return tx
}

func (st *state) ProposeBlock(round int) (*block.Block, error) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if !st.committee.IsProposer(st.signer.Address(), round) {
		return nil, errors.Errorf(errors.ErrInvalidAddress, "We are not propser for this round")
	}

	// Create new sandbox and execute transactions
	sb := st.makeSandbox()
	exe := execution.NewExecution()

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

		if err := exe.Execute(trx, sb); err != nil {
			st.logger.Debug("Found invalid transaction", "tx", trx, "err", err)
			st.txPool.RemoveTx(trx.ID())
		} else {
			txIDs.Append(trx.ID())

			if txIDs.Len() >= st.params.MaximumTransactionPerBlock {
				break
			}
		}
	}

	subsidyTx := st.createSubsidyTx(exe.AccumulatedFee())
	if subsidyTx == nil {
		st.logger.Error("Probably the node is shutting down.")
		return nil, errors.Errorf(errors.ErrInvalidBlock, "No subsidy transaction")
	}
	if err := st.txPool.AppendTxAndBroadcast(subsidyTx); err != nil {
		st.logger.Error("Our subsidy transaction is invalid. Why?", "err", err)
		return nil, err
	}
	txIDs.Prepend(subsidyTx.ID())

	stateHash := st.stateHash()
	committeeHash := st.committee.CommitteeHash()
	timestamp := st.proposeNextBlockTime()
	newSortitionSeed := st.lastInfo.SortitionSeed().Generate(st.signer)

	block := block.MakeBlock(
		st.params.BlockVersion,
		timestamp,
		txIDs,
		st.lastInfo.BlockHash(),
		committeeHash,
		stateHash,
		st.lastInfo.ReceiptsHash(),
		st.lastInfo.Certificate(),
		newSortitionSeed,
		st.signer.Address())

	return block, nil
}

func (st *state) ValidateBlock(block *block.Block) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	t := block.Header().Time()
	if err := st.validateBlockTime(t); err != nil {
		return err
	}

	if err := st.validateBlock(block); err != nil {
		return err
	}

	sb := st.makeSandbox()
	_, err := st.executeBlock(block, sb)
	if err != nil {
		return err
	}

	return nil
}

func (st *state) CommitBlock(height int, block *block.Block, cert *block.Certificate) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if height != st.lastInfo.BlockHeight() && height != st.lastInfo.BlockHeight()+1 {
		/// Returning error here will cause so many error logs during syncing blockchain
		/// Syncing is asynchronous job and we might receive blocks not in order
		st.logger.Debug("Unexpected block height", "height", height)
		return nil
	}

	/// There are two modules that can commit a block: Consensus and Syncer.
	/// Consensus engine is ours, we have full control over that and we know when and why a block should be committed.
	/// In the other hand, Syncer module receives new blocks from other peers and if we are behind them, it tries to commit them.
	/// We should never have a fork in our blockchain. but if it happens here we can catch it.
	if st.lastInfo.BlockHeight() == height {
		if block.Hash().EqualsTo(st.lastInfo.BlockHash()) {
			st.logger.Debug("This block committed before", "hash", block.Hash())
			return nil
		}

		st.logger.Error("A possible fork is detected", "our hash", st.lastInfo.BlockHash(), "block hash", block.Hash())
		return errors.Error(errors.ErrInvalidBlock)
	}

	err := st.validateBlock(block)
	if err != nil {
		return err
	}

	err = st.validateCertificateForCurrentHeight(cert, block.Hash())
	if err != nil {
		return err
	}

	// Verify proposer
	proposer := st.committee.Proposer(cert.Round())
	if !proposer.Address().EqualsTo(block.Header().ProposerAddress()) {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid proposer. Expected %s, got %s", proposer.Address(), block.Header().ProposerAddress())
	}
	// Validate sortition seed
	if !block.Header().SortitionSeed().Validate(proposer.PublicKey(), st.lastInfo.SortitionSeed()) {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid sortition seed.")
	}

	sb := st.makeSandbox()
	ctrxs, err := st.executeBlock(block, sb)
	if err != nil {
		return err
	}

	// Commit and update the validator set
	st.commitSandbox(sb, cert.Round())

	if err := st.store.SaveBlock(st.lastInfo.BlockHeight()+1, block); err != nil {
		return err
	}

	// Save txs and receipts
	receiptsHashes := make([]crypto.Hash, len(ctrxs))
	for i, ctrx := range ctrxs {
		st.txPool.RemoveTx(ctrx.Tx.ID())
		st.store.SaveTransaction(ctrx)

		receiptsHashes[i] = ctrx.Receipt.Hash()
	}
	receiptsMerkle := merkle.NewTreeFromHashes(receiptsHashes)

	st.lastInfo.SetBlockHeight(st.lastInfo.BlockHeight() + 1)
	st.lastInfo.SetBlockHash(block.Hash())
	st.lastInfo.SetBlockTime(block.Header().Time())
	st.lastInfo.SetSortitionSeed(block.Header().SortitionSeed())
	st.lastInfo.SetReceiptsHash(receiptsMerkle.Root())
	st.lastInfo.SetCertificate(cert)

	// Evaluate sortition before updating the validator set
	if st.evaluateSortition() {
		st.logger.Info("👏 This validator is chosen to be in the set", "address", st.signer.Address())
	}

	st.logger.Info("New block is committed", "block", block, "round", cert.Round())

	st.saveLastInfo(st.lastInfo.BlockHeight(), cert, st.lastInfo.ReceiptsHash(),
		st.committee.Committers(),
		st.committee.Proposer(0).Address())

	// At this point we can assign new sandbox to tx pool
	st.txPool.SetNewSandboxAndRecheck(st.makeSandbox())

	return nil
}

func (st *state) evaluateSortition() bool {
	if st.committee.Contains(st.signer.Address()) {
		// We are in the validator set right now
		return false
	}

	val, _ := st.store.Validator(st.signer.Address())
	if val == nil {
		// We are not a validator
		return false
	}

	if st.lastInfo.BlockHeight()-val.BondingHeight() < 2*st.params.CommitteeSize {
		// Bonding period
		return false
	}

	//
	ok, proof := st.sortition.EvaluateSortition(st.lastInfo.SortitionSeed(), st.signer, val.Stake())
	if ok {
		//
		trx := tx.NewSortitionTx(st.lastInfo.BlockHash(), val.Sequence()+1, val.Address(), proof)
		st.signer.SignMsg(trx)

		if err := st.txPool.AppendTxAndBroadcast(trx); err != nil {
			st.logger.Error("Our sortition transaction is invalid. Why?", "address", st.signer.Address(), "stake", val.Stake(), "tx", trx, "err", err)
			return false
		} else {
			st.logger.Debug("Sortition transaction broadcasted", "address", st.signer.Address(), "stake", val.Stake(), "tx", trx)
			return true
		}
	}

	return false
}

func calcBlockSubsidy(height int, subsidyReductionInterval int) int64 {
	// Equivalent to: baseSubsidy / 2^(height/subsidyHalvingInterval)
	return baseSubsidy >> uint(height/subsidyReductionInterval)
}

func (st *state) Fingerprint() string {
	return fmt.Sprintf("{#%d ⌘ %v 🕣 %v}",
		st.lastInfo.BlockHeight(),
		st.lastInfo.BlockHash().Fingerprint(),
		st.lastInfo.BlockTime().Format("15.04.05"))
}

func (st *state) commitSandbox(sb *sandbox.SandboxConcrete, round int) {
	joined := make([]*validator.Validator, 0)
	sb.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.JoinedCommittee {
			st.logger.Info("New validator joined", "address", vs.Validator.Address(), "stake", vs.Validator.Stake())

			joined = append(joined, &vs.Validator)
		}
	})

	if err := st.committee.Update(round, joined); err != nil {
		//
		// We should panic here before updating the state
		//
		logger.Panic("An error occurred", "err", err)
	}

	sb.IterateAccounts(func(as *sandbox.AccountStatus) {
		if as.Updated {
			st.store.UpdateAccount(&as.Account)
		}
	})

	sb.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.Updated {
			st.store.UpdateValidator(&vs.Validator)
		}
	})

	st.sortition.AddToTotalStake(sb.TotalStakeChange())
}

func (st *state) validateBlockTime(t time.Time) error {
	if t.Second()%st.params.BlockTimeInSecond != 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Block time is not rounded")
	}
	if t.Before(st.lastInfo.BlockTime().Add(1 * time.Second)) {
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
	timestamp := st.lastInfo.BlockTime().Add(st.params.BlockTime())
	timestamp = util.RoundTime(timestamp, st.params.BlockTimeInSecond)

	now := util.Now()
	if now.After(timestamp.Add(1 * time.Second)) {
		st.logger.Debug("It looks the last commit had delay", "delay", now.Sub(timestamp))
		timestamp = util.RoundNow(st.params.BlockTimeInSecond)
	}
	return timestamp
}

func (st *state) CommitteeValidators() []*validator.Validator {
	return st.committee.Validators()
}

func (st *state) IsInCommittee(addr crypto.Address) bool {
	return st.committee.Contains(addr)
}

func (st *state) Proposer(round int) *validator.Validator {
	return st.committee.Proposer(round)
}

func (st *state) IsProposer(addr crypto.Address, round int) bool {
	return st.committee.IsProposer(addr, round)
}

func (st *state) Transaction(id tx.ID) *tx.CommittedTx {
	tx, _ := st.store.Transaction(id)
	return tx
}

func (st *state) Block(height int) *block.Block {
	b, err := st.store.Block(height)
	if err != nil {
		st.logger.Trace("Error on retrieving block", "err", err)
	}
	return b
}

func (st *state) BlockHeight(hash crypto.Hash) int {
	h, err := st.store.BlockHeight(hash)
	if err != nil {
		st.logger.Trace("Error on retrieving block height", "err", err)
	}
	return h
}

func (st *state) Account(addr crypto.Address) *account.Account {
	acc, err := st.store.Account(addr)
	if err != nil {
		st.logger.Trace("Error on retrieving account", "err", err)
	}
	return acc
}

func (st *state) Validator(addr crypto.Address) *validator.Validator {
	val, err := st.store.Validator(addr)
	if err != nil {
		st.logger.Trace("Error on retrieving validator", "err", err)
	}
	return val
}
func (st *state) PendingTx(id tx.ID) *tx.Tx {
	return st.txPool.PendingTx(id)
}
func (st *state) AddPendingTx(trx *tx.Tx) error {
	return st.txPool.AppendTx(trx)
}
func (st *state) AddPendingTxAndBroadcast(trx *tx.Tx) error {
	return st.txPool.AppendTxAndBroadcast(trx)
}
