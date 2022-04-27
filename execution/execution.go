package execution

import (
	"github.com/zarbchain/zarb-go/execution/executor"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/tx/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/errors"
)

type Executor interface {
	Execute(trx *tx.Tx, sb sandbox.Sandbox) error
	Fee() int64
}
type Execution struct {
	executors      map[payload.Type]Executor
	accumulatedFee int64
}

func newExecution(strict bool) *Execution {
	execs := make(map[payload.Type]Executor)
	execs[payload.PayloadTypeSend] = executor.NewSendExecutor(strict)
	execs[payload.PayloadTypeBond] = executor.NewBondExecutor(strict)
	execs[payload.PayloadTypeSortition] = executor.NewSortitionExecutor(strict)
	execs[payload.PayloadTypeUnbond] = executor.NewUnbondExecutor(strict)
	execs[payload.PayloadTypeWithdraw] = executor.NewWithdrawExecutor(strict)

	return &Execution{
		executors: execs,
	}
}
func NewExecutor() *Execution {
	return newExecution(true)
}

func NewChecker() *Execution {
	return newExecution(false)
}

func (exe *Execution) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	if err := trx.SanityCheck(); err != nil {
		return err
	}
	if err := exe.checkStamp(trx, sb); err != nil {
		return err
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

func (exe *Execution) checkStamp(trx *tx.Tx, sb sandbox.Sandbox) error {
	curHeight := sb.CurrentHeight()
	height := sb.BlockHeightByStamp(trx.Stamp())
	interval := sb.TransactionToLiveInterval()

	if trx.IsSubsidyTx() {
		interval = 1
	} else if trx.IsSortitionTx() {
		interval = 7
	}

	if height == -1 || curHeight-height > interval {
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
		fee := int64(float64(trx.Payload().Value()) * sb.FeeFraction())
		fee = util.Max64(fee, sb.MinFee())
		if trx.Fee() != fee {
			return errors.Errorf(errors.ErrInvalidTx, "fee is wrong, expected: %v, got: %v", fee, trx.Fee())
		}
	}
	return nil
}
