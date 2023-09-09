package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/state/lastinfo"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/persistentmerkle"
	"github.com/pactus-project/pactus/util/simplemerkle"
	"github.com/pactus-project/pactus/www/nanomsg/event"
)

type state struct {
	lk sync.RWMutex

	signers         []crypto.Signer
	genDoc          *genesis.Genesis
	store           store.Store
	params          param.Params
	txPool          txpool.TxPool
	committee       committee.Committee
	totalPower      int64
	lastInfo        *lastinfo.LastInfo
	accountMerkle   *persistentmerkle.Tree
	validatorMerkle *persistentmerkle.Tree
	logger          *logger.SubLogger
	eventCh         chan event.Event
}

func LoadOrNewState(
	genDoc *genesis.Genesis,
	signers []crypto.Signer,
	store store.Store,
	txPool txpool.TxPool, eventCh chan event.Event,
) (Facade, error) {
	st := &state{
		signers:         signers,
		genDoc:          genDoc,
		txPool:          txPool,
		params:          genDoc.Params(),
		store:           store,
		lastInfo:        lastinfo.NewLastInfo(),
		accountMerkle:   persistentmerkle.New(),
		validatorMerkle: persistentmerkle.New(),
		eventCh:         eventCh,
	}
	st.logger = logger.NewSubLogger("_state", st)
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

	st.totalPower = st.retrieveTotalPower()

	st.loadMerkels()

	txPool.SetNewSandboxAndRecheck(st.concreteSandbox())

	st.logger.Debug("last info", "committers", st.committee.Committers(), "state_root", st.stateRoot().ShortString())

	return st, nil
}

func (st *state) concreteSandbox() sandbox.Sandbox {
	return sandbox.NewSandbox(st.store, st.params, st.committee, st.totalPower)
}

func (st *state) tryLoadLastInfo() error {
	// Make sure the genesis doc is the same as before.
	//
	// This check is not strictly necessary, since the genesis state is already committed.
	// However, it is good to perform this check to ensure that the genesis document has not been modified.
	genStateRoot := st.calculateGenesisStateRootFromGenesisDoc()
	blockOneInfo, err := st.store.Block(1)
	if err != nil {
		return err
	}

	blockOne, err := blockOneInfo.ToBlock()
	if err != nil {
		return err
	}

	if !genStateRoot.EqualsTo(blockOne.Header().StateRoot()) {
		return fmt.Errorf("invalid genesis doc")
	}

	logger.Debug("try to restore the last state")
	committee, err := st.lastInfo.RestoreLastInfo(st.store, st.params.CommitteeSize)
	if err != nil {
		return err
	}

	st.committee = committee

	logger.Info("last state restored",
		"last height", st.lastInfo.BlockHeight(),
		"last block time", st.lastInfo.BlockTime())

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

	committee, err := committee.NewCommittee(vals, st.params.CommitteeSize, vals[0].Address())
	if err != nil {
		return err
	}
	st.committee = committee
	st.lastInfo.UpdateBlockTime(genDoc.GenesisTime())

	return nil
}

func (st *state) loadMerkels() {
	totalAccount := st.store.TotalAccounts()
	st.store.IterateAccounts(func(addr crypto.Address, acc *account.Account) (stop bool) {
		// Let's keep this check, even we have tested it
		if acc.Number() >= totalAccount {
			panic("Account number is out of range")
		}
		st.accountMerkle.SetHash(int(acc.Number()), acc.Hash())

		return false
	})

	totalValidator := st.store.TotalValidators()
	st.store.IterateValidators(func(val *validator.Validator) (stop bool) {
		// Let's keep this check, even we have tested it
		if val.Number() >= totalValidator {
			panic("Validator number is out of range")
		}
		st.validatorMerkle.SetHash(int(val.Number()), val.Hash())

		return
	})
}

func (st *state) retrieveTotalPower() int64 {
	totalPower := int64(0)
	st.store.IterateValidators(func(val *validator.Validator) (stop bool) {
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

func (st *state) calculateGenesisStateRootFromGenesisDoc() hash.Hash {
	accs := st.genDoc.Accounts()
	vals := st.genDoc.Validators()

	accHashes := make([]hash.Hash, len(accs))
	valHashes := make([]hash.Hash, len(vals))
	for _, acc := range accs {
		accHashes[acc.Number()] = acc.Hash()
	}
	for _, val := range vals {
		valHashes[val.Number()] = val.Hash()
	}

	accTree := simplemerkle.NewTreeFromHashes(accHashes)
	valTree := simplemerkle.NewTreeFromHashes(valHashes)
	accRootHash := accTree.Root()
	valRootHash := valTree.Root()

	return *simplemerkle.HashMerkleBranches(&accRootHash, &valRootHash)
}

func (st *state) Close() error {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.store.Close()
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

func (st *state) LastCertificate() *certificate.Certificate {
	st.lk.RLock()
	defer st.lk.RUnlock()

	return st.lastInfo.Certificate()
}

func (st *state) UpdateLastCertificate(cert *certificate.Certificate) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	// Check if certificate has more signers ...
	if len(cert.Absentees()) < len(st.lastInfo.Certificate().Absentees()) {
		if err := st.validatePrevCertificate(cert, st.lastInfo.BlockHash()); err != nil {
			st.logger.Warn("try to update last certificate, but it's invalid", "err", err)
			return err
		}
		st.lastInfo.UpdateCertificate(cert)
	}

	return nil
}

func (st *state) createSubsidyTx(rewardAddr crypto.Address, fee int64) *tx.Tx {
	acc, err := st.store.Account(crypto.TreasuryAddress)
	if err != nil {
		// TODO: This can happen when a node is shutting down
		// We can prevent it by using global context.
		// Then we can close state before closing store.
		return nil
	}
	stamp := st.lastInfo.BlockHash().Stamp()
	seq := acc.Sequence() + 1
	tx := tx.NewSubsidyTx(stamp, seq, rewardAddr, st.params.BlockReward+fee, "")
	return tx
}

func (st *state) ProposeBlock(signer crypto.Signer, rewardAddr crypto.Address, round int16) (*block.Block, error) {
	st.lk.Lock()
	defer st.lk.Unlock()

	if !st.committee.IsProposer(signer.Address(), round) {
		return nil, errors.Errorf(errors.ErrGeneric, "we are not proposer for this round")
	}

	// Create new sandbox and execute transactions
	sb := st.concreteSandbox()
	exe := execution.NewExecutor()

	// Re-check all transactions strictly and remove invalid ones
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

	subsidyTx := st.createSubsidyTx(rewardAddr, exe.AccumulatedFee())
	if subsidyTx == nil {
		// probably the node is shutting down.
		st.logger.Error("no subsidy transaction")
		return nil, errors.Errorf(errors.ErrInvalidBlock, "no subsidy transaction")
	}
	txs.Prepend(subsidyTx)
	preSeed := st.lastInfo.SortitionSeed()

	block := block.MakeBlock(
		st.params.BlockVersion,
		st.proposeNextBlockTime(),
		txs,
		st.lastInfo.BlockHash(),
		st.stateRoot(),
		st.lastInfo.Certificate(),
		preSeed.GenerateNext(signer),
		signer.Address())

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
	return st.executeBlock(block, sb)
}

func (st *state) CommitBlock(height uint32, block *block.Block, cert *certificate.Certificate) error {
	st.lk.Lock()
	defer st.lk.Unlock()

	if height != st.lastInfo.BlockHeight()+1 {
		// Returning error here will cause so many error logs during syncing blockchain
		// Syncing is asynchronous job and we might receive blocks not in order
		st.logger.Debug("unexpected block height", "height", height)
		return nil
	}

	err := st.validateCertificate(cert, block.Hash())
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
	st.lastInfo.UpdateBlockHeight(height)
	st.lastInfo.UpdateBlockHash(block.Hash())
	st.lastInfo.UpdateBlockTime(block.Header().Time())
	st.lastInfo.UpdateSortitionSeed(block.Header().SortitionSeed())
	st.lastInfo.UpdateCertificate(cert)
	st.lastInfo.UpdateValidators(st.committee.Validators())

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

	st.logger.Info("new block committed", "block", block, "round", cert.Round())

	st.evaluateSortition()

	// -----------------------------------
	// At this point we can assign new sandbox to tx pool
	st.txPool.SetNewSandboxAndRecheck(st.concreteSandbox())

	// -----------------------------------
	// Publishing the events to the zmq
	st.publishEvents(height, block)

	return nil
}

func (st *state) evaluateSortition() bool {
	evaluated := false
	for _, signer := range st.signers {
		val, _ := st.store.Validator(signer.Address())
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

		ok, proof := sortition.EvaluateSortition(st.lastInfo.SortitionSeed(), signer, st.totalPower, val.Power())
		if ok {
			trx := tx.NewSortitionTx(st.lastInfo.BlockHash().Stamp(), val.Sequence()+1, val.Address(), proof)
			signer.SignMsg(trx)

			err := st.txPool.AppendTxAndBroadcast(trx)
			if err == nil {
				st.logger.Info("sortition transaction broadcasted",
					"address", signer.Address(), "power", val.Power(), "tx", trx)

				evaluated = true
			} else {
				st.logger.Error("our sortition transaction is invalid!",
					"address", signer.Address(), "power", val.Power(), "tx", trx, "err", err)
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

func (st *state) commitSandbox(sb sandbox.Sandbox, round int16) {
	joiningCommittee := make([]*validator.Validator, 0)
	sb.IterateValidators(func(val *validator.Validator, _ bool, joined bool) {
		if joined {
			st.logger.Info("new validator joined", "address", val.Address(), "power", val.Power())

			joiningCommittee = append(joiningCommittee, val)
		}
	})
	st.committee.Update(round, joiningCommittee)

	sb.IterateAccounts(func(addr crypto.Address, acc *account.Account, updated bool) {
		if updated {
			st.store.UpdateAccount(addr, acc)
			st.accountMerkle.SetHash(int(acc.Number()), acc.Hash())
		}
	})

	sb.IterateValidators(func(val *validator.Validator, updated bool, _ bool) {
		if updated {
			st.store.UpdateValidator(val)
			st.validatorMerkle.SetHash(int(val.Number()), val.Hash())
		}
	})

	st.totalPower += sb.PowerDelta()
}

func (st *state) validateBlockTime(t time.Time) error {
	if t.Second()%st.params.BlockIntervalInSecond != 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is not rounded", t.String())
	}
	if t.Before(st.lastInfo.BlockTime()) {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is before the last block time", t.String())
	}
	if t.Equal(st.lastInfo.BlockTime()) {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is same as the last block time", t.String())
	}
	proposeTime := st.proposeNextBlockTime()
	threshold := st.params.BlockInterval()
	if t.After(proposeTime.Add(threshold)) {
		return errors.Errorf(errors.ErrInvalidBlock, "block time (%s) is more than threshold (%s)",
			t.String(), proposeTime.String())
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

	now := util.Now()
	if now.After(timestamp.Add(1 * time.Second)) {
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

func (st *state) MakeCommittedBlock(data []byte, height uint32, blockHash hash.Hash) *store.CommittedBlock {
	return &store.CommittedBlock{
		Store:     st.store,
		Data:      data,
		BlockHash: blockHash,
		Height:    height,
	}
}

func (st *state) CommittedBlock(height uint32) *store.CommittedBlock {
	b, err := st.store.Block(height)
	if err != nil {
		st.logger.Trace("error on retrieving block", "err", err)
		return nil
	}
	return b
}

func (st *state) CommittedTx(id tx.ID) *store.CommittedTx {
	tx, err := st.store.Transaction(id)
	if err != nil {
		st.logger.Trace("searching transaction in local store failed", "id", id, "err", err)
	}
	return tx
}

func (st *state) BlockHash(height uint32) hash.Hash {
	return st.store.BlockHash(height)
}

func (st *state) BlockHeight(hash hash.Hash) uint32 {
	return st.store.BlockHeight(hash)
}

func (st *state) AccountByAddress(addr crypto.Address) *account.Account {
	acc, err := st.store.Account(addr)
	if err != nil {
		st.logger.Trace("error on retrieving account", "err", err)
	}
	return acc
}

func (st *state) AccountByNumber(number int32) *account.Account {
	acc, err := st.store.AccountByNumber(number)
	if err != nil {
		st.logger.Trace("error on retrieving account", "err", err)
	}
	return acc
}

func (st *state) ValidatorAddresses() []crypto.Address {
	return st.store.ValidatorAddresses()
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

// publishEvents publishes block related events.
func (st *state) publishEvents(height uint32, block *block.Block) {
	if st.eventCh == nil {
		return
	}
	blockEvent := event.CreateBlockEvent(block.Hash(), height)
	st.eventCh <- blockEvent

	for i := 1; i < block.Transactions().Len(); i++ {
		tx := block.Transactions().Get(i)
		TxEvent := event.CreateNewTransactionEvent(tx.ID(), height)
		st.eventCh <- TxEvent
	}
}

func (st *state) CalculateFee(amount int64, payloadType payload.Type) (int64, error) {
	switch payloadType {
	case payload.TypeTransfer,
		payload.TypeBond,
		payload.TypeWithdraw:
		{
			return execution.CalculateFee(amount, st.params), nil
		}

	case payload.TypeUnbond,
		payload.TypeSortition:
		{
			return 0, nil
		}

	default:
		return 0, errors.Errorf(errors.ErrInvalidTx, "unexpected tx type: %v", payloadType)
	}
}
