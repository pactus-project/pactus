package execution

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution/executor"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/util"
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

	if err := exe.checkStamp(trx); err != nil {
		return err
	}
	if err := exe.checkMemo(trx); err != nil {
		return err
	}
	if err := exe.checkFee(trx); err != nil {
		return err
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

func (exe *Execution) ResetFee() {
	exe.accumulatedFee = 0
}

func (exe *Execution) AccumulatedFee() int64 {
	return exe.accumulatedFee
}

func (exe *Execution) ClaimAccumulatedFee() {
	acc := exe.sandbox.Account(crypto.TreasuryAddress)
	acc.AddToBalance(exe.accumulatedFee)
	exe.sandbox.UpdateAccount(acc)
}

func (exe *Execution) checkMemo(trx *tx.Tx) error {
	if len(trx.Memo()) > exe.sandbox.MaxMemoLength() {
		return errors.Errorf(errors.ErrInvalidTx, "Memo length exceeded")
	}
	return nil
}

func (exe *Execution) checkStamp(trx *tx.Tx) error {
	curHeight := exe.sandbox.CurrentHeight()
	height := exe.sandbox.RecentBlockHeight(trx.Stamp())
	interval := exe.sandbox.TransactionToLiveInterval()

	if trx.IsMintbaseTx() {
		interval = 1
	} else if trx.IsSortitionTx() {
		interval = exe.sandbox.MaximumPower()
	}

	if height == -1 || curHeight-height > interval {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid stamp")
	}

	return nil
}

func (exe *Execution) checkFee(trx *tx.Tx) error {
	if trx.IsMintbaseTx() || trx.IsSortitionTx() {
		if trx.Fee() != 0 {
			return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: 0, got: %v", trx.Fee())
		}
	} else {
		fee := int64(float64(trx.Payload().Value()) * exe.sandbox.FeeFraction())
		fee = util.Max64(fee, exe.sandbox.MinFee())
		if trx.Fee() != fee {
			return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: %v, got: %v", fee, trx.Fee())
		}
	}
	return nil
}
