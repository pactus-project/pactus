package execution

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution/executor"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/tx"
)

type Executor struct {
	sendExecutor   *executor.SendExecutor
	accumulatedFee int64
	logger         *logger.Logger
}

func NewExecutor(sandbox executor.Sandbox) (*Executor, error) {
	exe := &Executor{
		sendExecutor: executor.NewSendExecutor(sandbox),
	}
	exe.logger = logger.NewLogger("executor", exe)
	return exe, nil
}

func (exe *Executor) Execute(trx *tx.Tx, isMintbaseTx bool) error {
	if !isMintbaseTx {
		if trx.IsMintbaseTx() {
			return errors.Errorf(errors.ErrInvalidTx, "Duplicated mintbase transaction")
		}
	}

	exe.accumulatedFee += trx.Fee()


	return exe.sendExecutor.Execute(trx)
}
