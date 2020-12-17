package execution

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution/executor"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type Executor interface {
	Execute(trx *tx.Tx) error
	Fee() int64
}
type Execution struct {
	executors      map[payload.PayloadType]Executor
	sandbox        sandbox.Sandbox
	accumulatedFee int64
}

func NewExecution(sb sandbox.Sandbox) *Execution {
	execs := make(map[payload.PayloadType]Executor)
	execs[payload.PayloadTypeSend] = executor.NewSendExecutor(sb)
	execs[payload.PayloadTypeBond] = executor.NewBondExecutor(sb)
	execs[payload.PayloadTypeSortition] = executor.NewSortitionExecutor(sb)

	return &Execution{
		executors: execs,
		sandbox:   sb,
	}
}

func (exe *Execution) Execute(trx *tx.Tx) error {
	if err := trx.SanityCheck(); err != nil {
		return err
	}

	curHeight := exe.sandbox.CurrentHeight()
	height := exe.sandbox.RecentBlockHeight(trx.Stamp())
	interval := exe.sandbox.TransactionToLiveInterval()

	if height == -1 || curHeight-height > interval {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid stamp")
	}
	if len(trx.Memo()) > exe.sandbox.MaxMemoLength() {
		return errors.Errorf(errors.ErrInvalidTx, "Memo length exceeded")
	}

	e, ok := exe.executors[trx.PayloadType()]
	if !ok {
		return errors.Errorf(errors.ErrInvalidTx, "unknown transaction type: %v", trx.PayloadType())
	}

	if err := e.Execute(trx); err != nil {
		return err
	}

	exe.accumulatedFee += e.Fee()

	return nil
}

func (exe *Execution) AccumulatedFee() int64 {
	return exe.accumulatedFee
}
