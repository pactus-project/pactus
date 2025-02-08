package state

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/state/lastinfo"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/state/score"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/persistentmerkle"
	"github.com/pactus-project/pactus/util/simplemerkle"
)

type state struct {
	lk sync.RWMutex

	valKeys         []*bls.ValidatorKey
	genDoc          *genesis.Genesis
	store           store.Store
	params          *param.Params
	txPool          txpool.TxPool
	committee       committee.Committee
	totalPower      int64
	lastInfo        *lastinfo.LastInfo
	accountMerkle   *persistentmerkle.Tree
	validatorMerkle *persistentmerkle.Tree
	scoreMgr        *score.Manager
	logger          *logger.SubLogger
	eventCh         chan<- any
}

func LoadOrNewState(
	genDoc *genesis.Genesis,
	valKeys []*bls.ValidatorKey,
	store store.Store,
	txPool txpool.TxPool, eventCh chan<- any,
) (Facade, error) {
	state := &state{
		valKeys:         valKeys,
		genDoc:          genDoc,
		txPool:          txPool,
		params:          param.FromGenesis(genDoc.Params()),
		store:           store,
		lastInfo:        lastinfo.NewLastInfo(),
		accountMerkle:   persistentmerkle.New(),
		validatorMerkle: persistentmerkle.New(),
		eventCh:         eventCh,
	}
	state.logger = logger.NewSubLogger("_state", state)
	state.store = store

	// Check if the number of accounts is greater than the genesis time;
	// this indicates we are not at the genesis height anymore.
	if store.TotalAccounts() > int32(len(genDoc.Accounts())) {
		err := state.tryLoadLastInfo()
		if err != nil {
			return nil, err
		}
	} else {
		// We are at the genesis height.
		err := state.makeGenesisState(genDoc)
		if err != nil {
			return nil, err
		}
	}

	state.totalPower = state.retrieveTotalPower()

	state.loadMerkels()

	txPool.SetNewSandboxAndRecheck(state.concreteSandbox())

	// Restoring score manager
	state.logger.Info("calculating the availability scores...")
	scoreWindow := uint32(60000)
	startHeight := uint32(2)
	endHeight := state.lastInfo.BlockHeight()
	if endHeight > scoreWindow {
		startHeight = endHeight - scoreWindow
	}

	scoreMgr := score.NewScoreManager(scoreWindow)
	for h := startHeight; h <= endHeight; h++ {
		cBlk, err := state.store.Block(h)
		if err != nil {
			return nil, err
		}
		// This code decodes the block certificate from the block data
		// without decoding the header and transactions.
		r := bytes.NewReader(cBlk.Data[138:]) // Block header is 138 bytes
		cert := new(certificate.BlockCertificate)
		err = cert.Decode(r)
		if err != nil {
			return nil, err
		}
		scoreMgr.SetCertificate(cert)
	}
	state.scoreMgr = scoreMgr

	for _, num := range state.committee.Committers() {
		state.logger.Debug("availability score", "val", num, "score", state.scoreMgr.AvailabilityScore(num))
	}

	state.logger.Debug("last info", "committers", state.committee.Committers(), "state_root", state.stateRoot())

	return state, nil
}

func (st *state) concreteSandbox() sandbox.Sandbox {
	return sandbox.NewSandbox(st.lastInfo.BlockHeight(),
		st.store, st.params, st.committee, st.totalPower)
}

func (st *state) tryLoadLastInfo() error {
	logger.Debug("try to restore the last state")
	committeeInstance, err := st.lastInfo.RestoreLastInfo(st.store, st.params.CommitteeSize)
	if err != nil {
		return err
	}
	st.committee = committeeInstance

	err = st.checkXeggexState()
	if err != nil {
		return err
	}

	logger.Info("last state restored",
		"last height", st.lastInfo.BlockHeight(),
		"last block time", st.lastInfo.BlockTime())

	return nil
}

// checkXeggexState checks the state of Xeggex account based on PIP-38.
func (st *state) checkXeggexState() error {
	xeggexAcc := st.store.XeggexAccount()
	if st.lastInfo.BlockHeight() < xeggexAcc.FreezeHeight {
		return nil
	}

	acc, _ := st.store.Account(st.store.XeggexAccount().DepositAddrs)
	if acc.Hash() != xeggexAcc.AccountHash {
		if st.store.HasPublicKey(xeggexAcc.WatcherAddrs) {
			// Unfrozen state

			return nil
		}

		watcherAcc, _ := st.store.Account(xeggexAcc.WatcherAddrs)
		if watcherAcc == nil || watcherAcc.Balance() < xeggexAcc.Balance {
			return errors.New("please re-sync your blockchain to remove the illegitimate Xeggex transaction")
		}
	}

	return nil
}

func (st *state) makeGenesisState(genDoc *genesis.Genesis) error {
	accs := genDoc.Accounts()
	for addr, acc := range accs {
		st.store.UpdateAccount(addr, acc)
	}

	vals := genDoc.Validators()
	for _, val := range vals {
		st.store.UpdateValidator(val)
	}

	err := st.store.WriteBatch()
	if err != nil {
		return err
	}

	cmt, err := committee.NewCommittee(vals, st.params.CommitteeSize, vals[0].Address())
	if err != nil {
		return err
	}
	st.committee = cmt
	st.lastInfo.UpdateBlockTime(genDoc.GenesisTime())

	return nil
}

func (st *state) loadMerkels() {
	totalAccount := st.store.TotalAccounts()
	st.store.IterateAccounts(func(_ crypto.Address, acc *account.Account) bool {
		// Let's keep this check, even we have tested it
		if acc.Number() >= totalAccount {
			panic(fmt.Sprintf(
				"Account number is out of range: %v >= %v", acc.Number(), totalAccount))
		}
		st.accountMerkle.SetHash(acc.Number(), acc.Hash())

		return false
	})

	totalValidator := st.store.TotalValidators()
	st.store.IterateValidators(func(val *validator.Validator) bool {
		// Let's keep this check, even we have tested it
		if val.Number() >= totalValidator {
			panic(fmt.Sprintf(
				"Validator number is out of range: %v >= %v", val.Number(), totalValidator))
		}
		st.validatorMerkle.SetHash(val.Number(), val.Hash())

		return false
	})
}

func (st *state) retrieveTotalPower() int64 {
	totalPower := int64(0)
	st.store.IterateValidators(func(val *validator.Validator) bool {
		totalPower += val.Power()

		return false
	})

	return totalPower
}

func (st *state) stateRoot() hash.Hash {
	accRoot := st.accountMerkle.Root()
	valRoot := st.validatorMerkle.Root()

	stateRoot := simplemerkle.HashMerkleBranches(&accRoot, &valRoot)

	return *stateRoot
}

func (st *state) Close() {
	st.lk.RLock()
	defer st.lk.RUnlock()

	st.store.Close()
}

func (st *state) Genesis() *genesis.Genesis {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.genDoc
}

func (st *state) LastBlockHeight() uint32 {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.BlockHeight()
}

func (st *state) LastBlockHash() hash.Hash {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.BlockHash()
}

func (st *state) LastBlockTime() time.Time {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.BlockTime()
}

func (st *state) LastCertificate() *certificate.BlockCertificate {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.Certificate()
}

func (st *state) UpdateLastCertificate(vte *vote.Vote) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	lastCert := st.lastInfo.Certificate()
	if vte.Type() != vote.VoteTypePrecommit ||
		vte.Height() != lastCert.Height() ||
		vte.Round() != lastCert.Round() {
		return InvalidVoteForCertificateError{
			Vote: vte,
		}
	}

	val, err := st.store.Validator(vte.Signer())
	if err != nil {
		return err
	}

	if !util.Contains(lastCert.Absentees(), val.Number()) {
		return InvalidVoteForCertificateError{
			Vote: vte,
		}
	}

	err = vte.Verify(val.PublicKey())
	if err != nil {
		return err
	}

	// prevent race condition
	cloneLastCert := lastCert.Clone()

	cloneLastCert.AddSignature(val.Number(), vte.Signature())
	st.lastInfo.UpdateCertificate(cloneLastCert)

	st.logger.Debug("certificate updated", "validator", val.Address(), "power", val.Power())

	return nil
}

func (st *state) createSubsidyTx(rewardAddr crypto.Address, accumulatedFee amount.Amount) *tx.Tx {
	lockTime := st.lastInfo.BlockHeight() + 1
	transaction := tx.NewSubsidyTx(lockTime, rewardAddr, st.params.BlockReward+accumulatedFee)

	return transaction
}

func (st *state) ProposeBlock(valKey *bls.ValidatorKey, rewardAddr crypto.Address) (*block.Block, error) {
	st.lk.Lock()
	defer st.lk.Unlock()

	// Create new sbx and execute transactions
	sbx := st.concreteSandbox()

	// Re-check all transactions strictly and remove invalid ones
	txs := st.txPool.PrepareBlockTransactions()
	txs = util.Trim(txs, st.params.MaxTransactionsPerBlock-1)
	for i := 0; i < txs.Len(); i++ {
		// Only one subsidy transaction per blk
		if txs[i].IsSubsidyTx() {
			st.logger.Error("found duplicated subsidy transaction", "tx", txs[i])
			txs.Remove(i)
			i--

			continue
		}

		if err := execution.CheckAndExecute(txs[i], sbx, true); err != nil {
			st.logger.Debug("found invalid transaction", "tx", txs[i], "error", err)
			txs.Remove(i)
			i--
		}
	}

	subsidyTx := st.createSubsidyTx(rewardAddr, sbx.AccumulatedFee())
	if subsidyTx == nil {
		// probably the node is shutting down.
		st.logger.Error("no subsidy transaction")

		return nil, ErrInvalidSubsidyTransaction
	}
	txs.Prepend(subsidyTx)
	prevSeed := st.lastInfo.SortitionSeed()

	blk := block.MakeBlock(
		st.params.BlockVersion,
		st.proposeNextBlockTime(),
		txs,
		st.lastInfo.BlockHash(),
		st.stateRoot(),
		st.lastInfo.Certificate(),
		prevSeed.GenerateNext(valKey.PrivateKey()),
		valKey.Address())

	return blk, nil
}

func (st *state) ValidateBlock(blk *block.Block, round int16) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if err := st.validateBlock(blk, round); err != nil {
		return err
	}

	t := blk.Header().Time()
	if err := st.validateBlockTime(t); err != nil {
		return err
	}

	sb := st.concreteSandbox()

	return st.executeBlock(blk, sb, true)
}

func (st *state) CommitBlock(blk *block.Block, cert *certificate.BlockCertificate) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	height := cert.Height()
	if height != st.lastInfo.BlockHeight()+1 {
		st.logger.Debug("block is committed before", "height", height)

		return nil
	}

	err := st.validateCurCertificate(cert, blk.Hash())
	if err != nil {
		return err
	}

	// There are two modules that can commit a block: Consensus and Sync.
	// The Consensus engine is ours, we have full control over that, and we know when
	// and why a block should be committed.
	// On the other hand, Sync module receives new blocks from the network and
	// tries to commit them.
	// We should never have a fork in our blockchain.
	// But if it happens, here we can catch it.
	if blk.Header().PrevBlockHash() != st.lastInfo.BlockHash() {
		st.logger.Panic("a possible fork is detected",
			"our hash", st.lastInfo.BlockHash(),
			"block hash", blk.Header().PrevBlockHash())
	}

	err = st.validateBlock(blk, cert.Round())
	if err != nil {
		return err
	}

	// -----------------------------------
	// Execute block
	sbx := st.concreteSandbox()
	if err := st.executeBlock(blk, sbx, false); err != nil {
		return err
	}

	// -----------------------------------
	// Commit block
	st.lastInfo.UpdateBlockHash(blk.Hash())
	st.lastInfo.UpdateBlockTime(blk.Header().Time())
	st.lastInfo.UpdateSortitionSeed(blk.Header().SortitionSeed())
	st.lastInfo.UpdateCertificate(cert)
	st.lastInfo.UpdateValidators(st.committee.Validators())

	// Commit and update the committee
	st.commitSandbox(sbx, cert.Round())

	st.store.SaveBlock(blk, cert)

	if err := st.store.WriteBatch(); err != nil {
		st.logger.Panic("unable to update state", "error", err)
	}

	// Remove transactions from pool and update consumption
	st.txPool.HandleCommittedBlock(blk)

	st.logger.Info("new block committed", "block", blk, "round", cert.Round())

	st.evaluateSortition()

	// -----------------------------------
	// At this point we can assign a new sandbox to tx pool
	st.txPool.SetNewSandboxAndRecheck(st.concreteSandbox())

	// -----------------------------------
	// Updating the score manager:
	// This code updates the availability scores.
	// To enhance syncing process, only blocks with timestamps from the last 10 days are considered.
	if blk.Header().Time().After(time.Now().AddDate(0, 0, -10)) {
		prevCert := blk.PrevCertificate()
		if prevCert != nil {
			st.scoreMgr.SetCertificate(prevCert)
		}
	}

	// publish committed block to event channel zeromq
	st.publishEvent(blk)

	return nil
}

func (st *state) evaluateSortition() bool {
	evaluated := false
	for _, key := range st.valKeys {
		val, _ := st.store.Validator(key.Address())
		if val == nil {
			// We are not a validator
			continue
		}

		if st.lastInfo.BlockHeight()-val.LastBondingHeight() < st.params.BondInterval {
			// Bonding period
			continue
		}

		if val.UnbondingHeight() > 0 {
			// we have Unbonded
			continue
		}

		ok, proof := sortition.EvaluateSortition(st.lastInfo.SortitionSeed(), key.PrivateKey(), st.totalPower, val.Power())
		if ok {
			trx := tx.NewSortitionTx(st.lastInfo.BlockHeight(), val.Address(), proof)
			sig := key.Sign(trx.SignBytes())
			trx.SetSignature(sig)
			trx.SetPublicKey(key.PublicKey())

			err := st.txPool.AppendTxAndBroadcast(trx)
			if err == nil {
				st.logger.Info("sortition transaction broadcasted",
					"address", key.Address(), "power", val.Power(), "tx", trx)

				evaluated = true
			} else {
				st.logger.Error("our sortition transaction is invalid!",
					"address", key.Address(), "power", val.Power(), "tx", trx, "error", err)
			}
		}
	}

	return evaluated
}

func (st *state) String() string {
	return fmt.Sprintf("{#%d ⌘ %v 🕣 %v}",
		st.lastInfo.BlockHeight(),
		st.lastInfo.BlockHash().ShortString(),
		st.lastInfo.BlockTime().Format("15.04.05"))
}

func (st *state) commitSandbox(sbx sandbox.Sandbox, round int16) {
	joiningCommittee := make([]*validator.Validator, 0)
	sbx.IterateValidators(func(val *validator.Validator, _ bool, joined bool) {
		if joined {
			st.logger.Debug("new validator joined", "address", val.Address(), "power", val.Power())

			joiningCommittee = append(joiningCommittee, val)
		}
	})
	st.committee.Update(round, joiningCommittee)

	sbx.IterateAccounts(func(addr crypto.Address, acc *account.Account, updated bool) {
		if updated {
			st.store.UpdateAccount(addr, acc)
			st.accountMerkle.SetHash(acc.Number(), acc.Hash())
		}
	})

	sbx.IterateValidators(func(val *validator.Validator, updated bool, _ bool) {
		if updated {
			st.store.UpdateValidator(val)
			st.validatorMerkle.SetHash(val.Number(), val.Hash())
		}
	})

	st.totalPower += sbx.PowerDelta()
}

func (st *state) validateBlockTime(blockTime time.Time) error {
	if blockTime.Second()%st.params.BlockIntervalInSecond != 0 {
		return InvalidBlockTimeError{
			Reason: fmt.Sprintf("block time (%s) is not rounded",
				blockTime.String()),
		}
	}
	if blockTime.Before(st.lastInfo.BlockTime()) {
		return InvalidBlockTimeError{
			Reason: fmt.Sprintf("block time (%s) is before the last block time (%s)",
				blockTime.String(), st.lastInfo.BlockTime()),
		}
	}
	if blockTime.Equal(st.lastInfo.BlockTime()) {
		return InvalidBlockTimeError{
			Reason: fmt.Sprintf("block time (%s) is same as the last block time",
				blockTime.String()),
		}
	}
	proposeTime := st.proposeNextBlockTime()
	threshold := st.params.BlockInterval()
	if blockTime.After(proposeTime.Add(threshold)) {
		return InvalidBlockTimeError{
			Reason: fmt.Sprintf("block time (%s) is more than threshold (%s)",
				blockTime.String(), proposeTime.String()),
		}
	}

	return nil
}

func (st *state) TotalPower() int64 {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.totalPower
}

func (st *state) CommitteePower() int64 {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.committee.TotalPower()
}

func (st *state) proposeNextBlockTime() time.Time {
	timestamp := st.lastInfo.BlockTime().Add(st.params.BlockInterval())

	now := time.Now()
	if now.After(timestamp.Add(10 * time.Second)) {
		st.logger.Debug("it looks the last block had delay", "delay", now.Sub(timestamp))
		timestamp = util.RoundNow(st.params.BlockIntervalInSecond)
	}

	return timestamp
}

func (st *state) CommitteeValidators() []*validator.Validator {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.committee.Validators()
}

func (st *state) TotalAccounts() int32 {
	return st.store.TotalAccounts()
}

func (st *state) TotalValidators() int32 {
	return st.store.TotalValidators()
}

func (st *state) IsInCommittee(addr crypto.Address) bool {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.committee.Contains(addr)
}

func (st *state) Proposer(round int16) *validator.Validator {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.committee.Proposer(round)
}

func (st *state) IsProposer(addr crypto.Address, round int16) bool {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.committee.IsProposer(addr, round)
}

func (st *state) IsValidator(addr crypto.Address) bool {
	return st.store.HasValidator(addr)
}

func (st *state) CommittedBlock(height uint32) *store.CommittedBlock {
	blk, err := st.store.Block(height)
	if err != nil {
		st.logger.Trace("error on retrieving block", "error", err)

		return nil
	}

	return blk
}

func (st *state) CommittedTx(txID tx.ID) *store.CommittedTx {
	transaction, err := st.store.Transaction(txID)
	if err != nil {
		st.logger.Trace("searching transaction in local store failed", "id", txID, "error", err)
	}

	return transaction
}

func (st *state) BlockHash(height uint32) hash.Hash {
	return st.store.BlockHash(height)
}

func (st *state) BlockHeight(h hash.Hash) uint32 {
	return st.store.BlockHeight(h)
}

func (st *state) AccountByAddress(addr crypto.Address) *account.Account {
	acc, err := st.store.Account(addr)
	if err != nil {
		st.logger.Trace("error on retrieving account", "error", err)
	}

	return acc
}

func (st *state) ValidatorAddresses() []crypto.Address {
	return st.store.ValidatorAddresses()
}

func (st *state) ValidatorByAddress(addr crypto.Address) *validator.Validator {
	val, err := st.store.Validator(addr)
	if err != nil {
		st.logger.Trace("error on retrieving validator", "error", err)
	}

	return val
}

// ValidatorByNumber returns validator data based on validator number.
func (st *state) ValidatorByNumber(n int32) *validator.Validator {
	val, err := st.store.ValidatorByNumber(n)
	if err != nil {
		st.logger.Trace("error on retrieving validator", "error", err)
	}

	return val
}

func (st *state) PendingTx(txID tx.ID) *tx.Tx {
	return st.txPool.PendingTx(txID)
}

func (st *state) AddPendingTx(trx *tx.Tx) error {
	return st.txPool.AppendTx(trx)
}

func (st *state) AddPendingTxAndBroadcast(trx *tx.Tx) error {
	return st.txPool.AppendTxAndBroadcast(trx)
}

func (st *state) Params() *param.Params {
	return st.params
}

func (st *state) CalculateFee(amt amount.Amount, payloadType payload.Type) amount.Amount {
	return st.txPool.EstimatedFee(amt, payloadType)
}

func (st *state) PublicKey(addr crypto.Address) (crypto.PublicKey, error) {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.store.PublicKey(addr)
}

func (st *state) AvailabilityScore(valNum int32) float64 {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.scoreMgr.AvailabilityScore(valNum)
}

func (st *state) AllPendingTxs() []*tx.Tx {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.txPool.AllPendingTxs()
}

func (st *state) IsPruned() bool {
	return st.store.IsPruned()
}

func (st *state) PruningHeight() uint32 {
	return st.store.PruningHeight()
}

func (st *state) publishEvent(msg any) {
	if st.eventCh == nil {
		return
	}

	select {
	case st.eventCh <- msg:
	default:
	}
}
