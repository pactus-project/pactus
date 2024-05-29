package execution

import (
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type Executor interface {
	Execute(trx *tx.Tx, sb sandbox.Sandbox) error
}
type Execution struct {
	executors map[payload.Type]Executor
	strict    bool
}

func newExecution(strict bool) *Execution {
	execs := make(map[payload.Type]Executor)
	execs[payload.TypeTransfer] = executor.NewTransferExecutor(strict)
	execs[payload.TypeBond] = executor.NewBondExecutor(strict)
	execs[payload.TypeSortition] = executor.NewSortitionExecutor(strict)
	execs[payload.TypeUnbond] = executor.NewUnbondExecutor(strict)
	execs[payload.TypeWithdraw] = executor.NewWithdrawExecutor(strict)

	return &Execution{
		executors: execs,
		strict:    strict,
	}
}

func NewExecutor() *Execution {
	return newExecution(true)
}

func NewChecker() *Execution {
	return newExecution(false)
}

func (exe *Execution) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	if sb.IsBanned(trx.Payload().Signer()) {
		return SignerBannedError{
			addr: trx.Payload().Signer(),
		}
	}

	// TODO: remove me later.
	if trx.Payload().Type() == payload.TypeTransfer {
		if trx.Payload().Receiver() != nil && trx.Payload().Signer() == *trx.Payload().Receiver() {
			return SpamTxError{
				addr: trx.Payload().Signer(),
			}
		}
	}

	if exists := sb.AnyRecentTransaction(trx.ID()); exists {
		return TransactionCommittedError{
			ID: trx.ID(),
		}
	}

	if err := exe.checkLockTime(trx, sb); err != nil {
		return err
	}

	if err := exe.checkFee(trx); err != nil {
		return err
	}

	e, ok := exe.executors[trx.Payload().Type()]
	if !ok {
		return UnknownPayloadTypeError{
			PayloadType: trx.Payload().Type(),
		}
	}

	if err := e.Execute(trx, sb); err != nil {
		return err
	}

	sb.CommitTransaction(trx)

	return nil
}

func (exe *Execution) checkLockTime(trx *tx.Tx, sb sandbox.Sandbox) error {
	interval := sb.Params().TransactionToLiveInterval

	if trx.IsSubsidyTx() {
		interval = 0
	} else if trx.IsSortitionTx() {
		interval = sb.Params().SortitionInterval
	}

	if sb.CurrentHeight() > interval {
		if trx.LockTime() < sb.CurrentHeight()-interval {
			return PastLockTimeError{
				LockTime: trx.LockTime(),
			}
		}
	}

	if exe.strict {
		// In strict mode, transactions with future lock times are rejected.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if trx.LockTime() > sb.CurrentHeight() {
			return FutureLockTimeError{
				LockTime: trx.LockTime(),
			}
		}
	}

	return nil
}

func (exe *Execution) checkFee(trx *tx.Tx) error {
	if trx.IsFreeTx() {
		if trx.Fee() != 0 {
			return InvalidFeeError{
				Fee:      trx.Fee(),
				Expected: 0,
			}
		}
	}

	return nil
}
