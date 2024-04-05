package execution

import (
	"math"

	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
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
	if exists := sb.AnyRecentTransaction(trx.ID()); exists {
		return TransactionCommittedError{
			ID: trx.ID(),
		}
	}

	if err := exe.checkLockTime(trx, sb); err != nil {
		return err
	}

	if err := exe.checkFee(trx, sb); err != nil {
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

func (exe *Execution) checkFee(trx *tx.Tx, sb sandbox.Sandbox) error {
	var fee amount.Amount
	if trx.IsSubsidyTx() {
		fee = 0
	} else {
		fee = CalculateFee(trx.Payload().Value(), trx.Payload().Type(), sb.Params())
	}

	if fee == 0 {
		if trx.Fee() != 0 {
			return InvalidFeeError{
				Fee:      trx.Fee(),
				Expected: fee,
			}
		}
	} else {
		// Check if the absolute difference between the calculated fee and the transaction fee
		// is greater than 1 PAC, indicating an invalid fee.
		if math.Abs(float64(fee-trx.Fee())) > 1 {
			return InvalidFeeError{
				Fee:      trx.Fee(),
				Expected: fee,
			}
		}
	}

	return nil
}

func CalculateFee(amt amount.Amount, payloadType payload.Type, params *param.Params) amount.Amount {
	switch payloadType {
	case payload.TypeUnbond,
		payload.TypeSortition:

		return 0

	case payload.TypeTransfer,
		payload.TypeBond,
		payload.TypeWithdraw:
		fee := amt.MulF64(params.FeeFraction)
		fee = util.Max(fee, params.MinimumFee)
		fee = util.Min(fee, params.MaximumFee)

		return fee

	default:
		return 0
	}
}
