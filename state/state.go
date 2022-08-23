package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/state/lastinfo"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/genesis"
	"github.com/zarbchain/zarb-go/types/param"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/errors"
	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/www/zmq/event"
)

type state struct {
	lk sync.RWMutex

	config        *Config
	signer        crypto.Signer
	rewardAddress crypto.Address
	genDoc        *genesis.Genesis
	store         store.Store
	params        param.Params
	txPool        txpool.TxPool
	committee     committee.Committee
	lastInfo      *lastinfo.LastInfo
	logger        *logger.Logger
	eventCh       chan event.Event
}

func LoadOrNewState(
	conf *Config,
	genDoc *genesis.Genesis,
	signer crypto.Signer,
	store store.Store,
	txPool txpool.TxPool, eventCh chan event.Event) (Facade, error) {
	// Block rewards goes to the reward address
	// If it is set inside config, we use that address
	// otherwise, it will be the signer address
	var rewardAddr crypto.Address
	if conf.RewardAddress != "" {
		addr, _ := crypto.AddressFromString(conf.RewardAddress)
		rewardAddr = addr
	} else {
		rewardAddr = signer.Address()
	}

	st := &state{
		config:        conf,
		genDoc:        genDoc,
		txPool:        txPool,
		params:        genDoc.Params(),
		signer:        signer,
		store:         store,
		rewardAddress: rewardAddr,
		lastInfo:      lastinfo.NewLastInfo(store),
		eventCh:       eventCh,
	}
	st.logger = logger.NewLogger("_state", st)
	st.store = store

	// The first account is Treasury Account at the genesis time.
	// So if we have more account, we are not in the genesis height anymore.
	if store.TotalAccounts() > 1 {
		err := st.tryLoadLastInfo()
		if err != nil {
			return nil, err
		}
	} else {
		// We are at the genesis height
		err := st.makeGenesisState(genDoc)
		if err != nil {
			return nil, err
		}
	}

	txPool.SetNewSandboxAndRecheck(st.concreteSandbox())

	return st, nil
}

func (st *state) concreteSandbox() sandbox.Sandbox {
	return sandbox.NewSandbox(st.store, st.params, st.committee)
}

func (st *state) tryLoadLastInfo() error {
	// Make sure genesis hash is same
	//
	// This check is not important because genesis state is committed.
	// But it is good to have it to make sure genesis doc hasn't changed
	genStateRoot := st.calculateGenesisStateRootFromGenesisDoc()
	blockOneHash := st.store.BlockHash(1)
	blockOneInfo, err := st.store.Block(blockOneHash)
	if err != nil {
		return err
	}
	blockOne, _ := blockOneInfo.ToFullBlock()
	if !genStateRoot.EqualsTo(blockOne.Header().StateRoot()) {
		return fmt.Errorf("invalid genesis doc")
	}

	logger.Info("try to load the last state info")
	committee, err := st.lastInfo.RestoreLastInfo(st.params.CommitteeSize)
	if err != nil {
		return err
	}

	st.committee = committee

	st.logger.Debug("last info", "committers", committee.Committers(), "state_root", st.stateRoot().Fingerprint())

	return nil
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

	err := st.store.WriteBatch()
	if err != nil {
		return err
	}

	committee, err := committee.NewCommittee(vals, st.params.CommitteeSize, vals[0].Address())
	if err != nil {
		return err
	}
	st.committee = committee
	st.lastInfo.SetBlockTime(genDoc.GenesisTime())

	return nil
}

func (st *state) Close() error {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.store.Close()
}

func (st *state) GenesisHash() hash.Hash {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.genDoc.Hash()
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

	// Check if certificate has more signers ...
	if len(cert.Absentees()) < len(st.lastInfo.Certificate().Absentees()) {
		if err := st.validateCertificateForPreviousHeight(st.lastInfo.BlockHash(), cert); err != nil {
			st.logger.Warn("try to update last certificate, but it's invalid", "err", err)
			return err
		}
		st.lastInfo.SetCertificate(cert)
	}

	return nil
}

func (st *state) createSubsidyTx(fee int64) *tx.Tx {
	acc, err := st.store.Account(crypto.TreasuryAddress)
	if err != nil {
		// TODO: This can happen when a node is shutting down
		// We can prevent it by using global context.
		// Then we can close state before closing store.
		return nil
	}
	stamp := st.lastInfo.BlockHash().Stamp()
	seq := acc.Sequence() + 1
	tx := tx.NewSubsidyTx(stamp, seq, st.rewardAddress, st.params.BlockReward+fee, "")
	return tx
}

func (st *state) ProposeBlock(round int16) (*block.Block, error) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if !st.committee.IsProposer(st.signer.Address(), round) {
		return nil, errors.Errorf(errors.ErrGeneric, "we are not propser for this round")
	}

	// Create new sandbox and execute transactions
	sb := st.concreteSandbox()
	exe := execution.NewExecutor()

	// Re-chaeck all transactions strictly and remove invalid ones
	txs := st.txPool.PrepareBlockTransactions()
	for i := 0; i < txs.Len(); i++ {
		// Only one subsidy transaction per block
		if txs[i].IsSubsidyTx() {
			st.logger.Error("found duplicated subsidy transaction", "tx", txs[i])
			st.txPool.RemoveTx(txs[i].ID())
			txs.Remove(i)
			i--
			continue
		}

		if err := exe.Execute(txs[i], sb); err != nil {
			st.logger.Debug("found invalid transaction", "tx", txs[i], "err", err)
			txs.Remove(i)
			i--
		}
		// Maximum 1000 transactions per block
		if txs.Len() >= 1000 {
			break
		}
	}

	subsidyTx := st.createSubsidyTx(exe.AccumulatedFee())
	if subsidyTx == nil {
		// probably the node is shutting down.
		st.logger.Error("no subsidy transaction")
		return nil, errors.Errorf(errors.ErrInvalidBlock, "no subsidy transaction")
	}
	txs.Prepend(subsidyTx)
	seed := st.lastInfo.SortitionSeed()

	block := block.MakeBlock(
		st.params.BlockVersion,
		st.proposeNextBlockTime(),
		txs,
		st.lastInfo.BlockHash(),
		st.stateRoot(),
		st.lastInfo.Certificate(),
		seed.Generate(st.signer),
		st.signer.Address())

	return block, nil
}

func (st *state) ValidateBlock(block *block.Block) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if err := st.validateBlock(block); err != nil {
		return err
	}

	t := block.Header().Time()
	if err := st.validateBlockTime(t); err != nil {
		return err
	}

	sb := st.concreteSandbox()
	if err := st.executeBlock(block, sb); err != nil {
		return err
	}

	return nil
}

func (st *state) CommitBlock(height uint32, block *block.Block, cert *block.Certificate) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if height != st.lastInfo.BlockHeight()+1 {
		// Returning error here will cause so many error logs during syncing blockchain
		// Syncing is asynchronous job and we might receive blocks not in order
		st.logger.Debug("unexpected block height", "height", height)
		return nil
	}

	err := st.validateCertificate(block.Hash(), cert)
	if err != nil {
		return err
	}

	// There are two modules that can commit a block: Consensus and Sync.
	// Consensus engine is ours, we have full control over that and we know when
	// and why a block should be committed.
	// On the other hand, Sync module receives new blocks from the network and
	// tries to commit them.
	// We should never have a fork in our blockchain.
	// But if it happens, here we can catch it.
	if !block.Header().PrevBlockHash().EqualsTo(st.lastInfo.BlockHash()) {
		st.logger.Panic("a possible fork is detected",
			"our hash", st.lastInfo.BlockHash(),
			"block hash", block.Header().PrevBlockHash())
		return errors.Error(errors.ErrInvalidBlock)
	}

	err = st.validateBlock(block)
	if err != nil {
		return err
	}

	// Verify proposer
	proposer := st.committee.Proposer(cert.Round())
	if !proposer.Address().EqualsTo(block.Header().ProposerAddress()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid proposer, expected %s, got %s", proposer.Address(), block.Header().ProposerAddress())
	}
	// Validate sortition seed
	seed := block.Header().SortitionSeed()
	if !seed.Verify(proposer.PublicKey(), st.lastInfo.SortitionSeed()) {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid sortition seed")
	}

	// -----------------------------------
	// Execute block
	sb := st.concreteSandbox()
	if err := st.executeBlock(block, sb); err != nil {
		return err
	}

	// -----------------------------------
	// Commit block
	st.lastInfo.SetBlockHeight(height)
	st.lastInfo.SetBlockHash(block.Hash())
	st.lastInfo.SetBlockTime(block.Header().Time())
	st.lastInfo.SetSortitionSeed(block.Header().SortitionSeed())
	st.lastInfo.SetCertificate(cert)

	// Commit and update the committee
	st.commitSandbox(sb, cert.Round())

	st.store.SaveBlock(height, block, cert)

	// Remove transactions from pool
	for _, trx := range block.Transactions() {
		st.txPool.RemoveTx(trx.ID())
	}

	if err := st.store.WriteBatch(); err != nil {
		st.logger.Panic("unable to update state", "err", err)
	}

	st.logger.Info("new block is committed", "block", block, "round", cert.Round())

	// -----------------------------------
	// Evaluate sortition before updating the committee
	if st.evaluateSortition() {
		st.logger.Info("👏 this validator is chosen to be in the committee", "address", st.signer.Address())
	}

	// -----------------------------------
	// At this point we can assign new sandbox to tx pool
	st.txPool.SetNewSandboxAndRecheck(st.concreteSandbox())
	// -----------------------------------
	// Publishing the events to the zmq
	st.publishEvents(height, block)
	return nil
}

func (st *state) evaluateSortition() bool {
	val, _ := st.store.Validator(st.signer.Address())
	if val == nil {
		// We are not a validator
		return false
	}

	if st.lastInfo.BlockHeight()-val.LastBondingHeight() < st.params.BondInterval {
		// Bonding period
		return false
	}

	if val.UnbondingHeight() > 0 {
		// we have Unbonded
		return false
	}

	ok, proof := sortition.EvaluateSortition(st.lastInfo.SortitionSeed(), st.signer, st.totalPower(), val.Power())
	if ok {
		trx := tx.NewSortitionTx(st.lastInfo.BlockHash().Stamp(), val.Sequence()+1, val.Address(), proof)
		st.signer.SignMsg(trx)

		err := st.txPool.AppendTxAndBroadcast(trx)
		if err == nil {
			st.logger.Debug("sortition transaction broadcasted",
				"address", st.signer.Address(), "power", val.Power(), "tx", trx)
			return true
		}
		st.logger.Error("our sortition transaction is invalid. Why?",
			"address", st.signer.Address(), "power", val.Power(), "tx", trx, "err", err)
	}

	return false
}

func (st *state) Fingerprint() string {
	return fmt.Sprintf("{#%d ⌘ %v 🕣 %v}",
		st.lastInfo.BlockHeight(),
		st.lastInfo.BlockHash().Fingerprint(),
		st.lastInfo.BlockTime().Format("15.04.05"))
}

func (st *state) commitSandbox(sb sandbox.Sandbox, round int16) {
	joined := make([]*validator.Validator, 0)
	currentHeight := sb.CurrentHeight()
	sb.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.Validator.LastJoinedHeight() == currentHeight {
			st.logger.Info("new validator joined", "address", vs.Validator.Address(), "power", vs.Validator.Power())

			joined = append(joined, &vs.Validator)
		}
	})
	st.committee.Update(round, joined)

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
}

func (st *state) validateBlockTime(t time.Time) error {
	if t.Second()%st.params.BlockTimeInSecond != 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is not rounded", t.String())
	}
	if t.Before(st.lastInfo.BlockTime()) {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is before the last block time", t.String())
	}
	if t.Equal(st.lastInfo.BlockTime()) {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is same as the last block time", t.String())
	}
	proposeTime := st.proposeNextBlockTime()
	threshold := st.params.BlockTime()
	if t.After(proposeTime.Add(threshold)) {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is more than threshold (%s)",
			t.String(), proposeTime.String())
	}

	return nil
}

func (st *state) TotalPower() int64 {
	st.lk.Lock()
	defer st.lk.Unlock()

	return st.totalPower()
}

func (st *state) CommitteePower() int64 {
	st.lk.Lock()
	defer st.lk.Unlock()

	return st.committeePower()
}

// TODO: add test for me when a validator is parked (unbonded)
// TODO: Improve performance of remember total power
// TODO: sandbox has the same logic.
func (st *state) totalPower() int64 {
	p := int64(0)
	st.store.IterateValidators(func(val *validator.Validator) bool {
		p += val.Power()
		return false
	})
	return p
}

func (st *state) committeePower() int64 {
	return st.committee.TotalPower()
}

func (st *state) proposeNextBlockTime() time.Time {
	timestamp := st.lastInfo.BlockTime().Add(st.params.BlockTime())

	now := util.Now()
	if now.After(timestamp.Add(1 * time.Second)) {
		st.logger.Debug("it looks the last block had delay", "delay", now.Sub(timestamp))
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

func (st *state) Proposer(round int16) *validator.Validator {
	return st.committee.Proposer(round)
}

func (st *state) IsProposer(addr crypto.Address, round int16) bool {
	return st.committee.IsProposer(addr, round)
}

func (st *state) Transaction(id tx.ID) *tx.Tx {
	tx, err := st.store.Transaction(id)
	if err != nil {
		st.logger.Trace("searching transaction in local store failed", "id", id, "err", err)
	}
	return tx
}

func (st *state) Block(hash hash.Hash) *block.Block {
	bi, err := st.store.Block(hash)
	if err != nil {
		st.logger.Trace("error on retrieving block", "err", err)
		return nil
	}
	b, _ := bi.ToFullBlock()
	return b
}

func (st *state) BlockHash(height uint32) hash.Hash {
	return st.store.BlockHash(height)
}

func (st *state) AccountByAddress(addr crypto.Address) *account.Account {
	acc, err := st.store.Account(addr)
	if err != nil {
		st.logger.Trace("error on retrieving account", "err", err)
	}
	return acc
}

func (st *state) ValidatorByAddress(addr crypto.Address) *validator.Validator {
	val, err := st.store.Validator(addr)
	if err != nil {
		st.logger.Trace("error on retrieving validator", "err", err)
	}
	return val
}

// ValidatorByNumber returns validator data based on validator number.
func (st *state) ValidatorByNumber(n int32) *validator.Validator {
	val, err := st.store.ValidatorByNumber(n)
	if err != nil {
		st.logger.Trace("error on retrieving validator", "err", err)
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
func (st *state) Params() param.Params {
	return st.params
}
func (st *state) RewardAddress() crypto.Address {
	return st.rewardAddress
}
func (st *state) ValidatorAddress() crypto.Address {
	return st.signer.Address()
}

func (st *state) publishEvents(height uint32, block *block.Block) {
	if st.eventCh == nil {
		return
	}
	blockEvent := event.CreateBlockEvent(block.Hash(), height)
	st.eventCh <- blockEvent
}
