package execution

import (
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
)

type Executor interface {
	Execute(trx *tx.Tx, sb sandbox.Sandbox) error
	Fee() int64
}
type Execution struct {
	executors      map[payload.Type]Executor
	accumulatedFee int64
	strict         bool
}

func newExecution(strict bool) *Execution {
	execs := make(map[payload.Type]Executor)
	execs[payload.PayloadTypeTransfer] = executor.NewTransferExecutor(strict)
	execs[payload.PayloadTypeBond] = executor.NewBondExecutor(strict)
	execs[payload.PayloadTypeSortition] = executor.NewSortitionExecutor(strict)
	execs[payload.PayloadTypeUnbond] = executor.NewUnbondExecutor(strict)
	execs[payload.PayloadTypeWithdraw] = executor.NewWithdrawExecutor(strict)

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
	if err := trx.BasicCheck(); err != nil {
		return err
	}
	if trx.IsLockTime() {
		if err := exe.checkLockTime(trx, sb); err != nil {
			return err
		}
	} else {
		if err := exe.checkStamp(trx, sb); err != nil {
			return err
		}
	}

	if err := exe.checkFee(trx, sb); err != nil {
		return err
	}

	e, ok := exe.executors[trx.Payload().Type()]
	if !ok {
		return errors.Errorf(errors.ErrInvalidTx, "unknown transaction type: %v", trx.Payload().Type())
	}

	if err := e.Execute(trx, sb); err != nil {
		return err
	}

	exe.accumulatedFee += e.Fee()

	return nil
}

func (exe *Execution) AccumulatedFee() int64 {
	return exe.accumulatedFee
}

func (exe *Execution) checkLockTime(trx *tx.Tx, sb sandbox.Sandbox) error {
	curHeight := sb.CurrentHeight()
	lockTimeHeight := trx.LockTime()
	interval := sb.Params().TransactionToLiveInterval

	if trx.IsSubsidyTx() || trx.IsSortitionTx() {
		return errors.Errorf(errors.ErrInvalidTx, "invalid lock time")
	}

	if curHeight > lockTimeHeight+interval {
		return errors.Errorf(errors.ErrInvalidTx, "expired lock time")
	}

	if curHeight < lockTimeHeight {
		if exe.strict {
			return errors.Errorf(errors.ErrInvalidTx, "unfinalized transaction")
		}
	}

	return nil
}

func (exe *Execution) checkStamp(trx *tx.Tx, sb sandbox.Sandbox) error {
	curHeight := sb.CurrentHeight()
	height, _ := sb.RecentBlockByStamp(trx.Stamp())
	interval := sb.Params().TransactionToLiveInterval

	if trx.IsSubsidyTx() {
		interval = 1
	} else if trx.IsSortitionTx() {
		interval = sb.Params().SortitionInterval
	}

	if curHeight-height > interval {
		return errors.Errorf(errors.ErrInvalidTx, "invalid stamp")
	}

	return nil
}

func (exe *Execution) checkFee(trx *tx.Tx, sb sandbox.Sandbox) error {
	if trx.IsFreeTx() {
		if trx.Fee() != 0 {
			return errors.Errorf(errors.ErrInvalidTx, "fee is wrong, expected: 0, got: %v", trx.Fee())
		}
	} else {
		fee := CalculateFee(trx.Payload().Value(), sb.Params())
		if trx.Fee() != fee {
			return errors.Errorf(errors.ErrInvalidFee, "fee is wrong, expected: %v, got: %v", fee, trx.Fee())
		}
	}
	return nil
}

func CalculateFee(amt int64, params param.Params) int64 {
	fee := int64(float64(amt) * params.FeeFraction)
	fee = util.Max(fee, params.MinimumFee)
	fee = util.Min(fee, params.MaximumFee)
	return fee
}
