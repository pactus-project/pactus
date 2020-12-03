package execution

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution/executor"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

type Executor interface {
	Execute(trx *tx.Tx) error
	Fee() int64
}
type Execution struct {
	executors      map[tx.PayloadType]Executor
	accumulatedFee int64
}

func NewExecution(sb sandbox.Sandbox) *Execution {
	execs := make(map[tx.PayloadType]Executor)
	execs[tx.PayloadTypeSend] = executor.NewSendExecutor(sb)
	execs[tx.PayloadTypeBond] = executor.NewBondExecutor(sb)
	execs[tx.PayloadTypeSortition] = executor.NewSendExecutor(sb)

	return &Execution{
		executors: execs,
	}
}

func (exe *Execution) Execute(trx *tx.Tx, isMintbaseTx bool) error {
	if err := trx.SanityCheck(); err != nil {
		return err
	}

	if isMintbaseTx {
		if !trx.IsMintbaseTx() {
			return errors.Errorf(errors.ErrInvalidTx, "Not a mintbase transaction")
		}
	} else {
		if trx.IsMintbaseTx() {
			return errors.Errorf(errors.ErrInvalidTx, "Duplicated mintbase transaction")
		}
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
