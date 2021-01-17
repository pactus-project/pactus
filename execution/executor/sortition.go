package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type SortitionExecutor struct {
	sandbox sandbox.Sandbox
}

func NewSortitionExecutor(sb sandbox.Sandbox) *SortitionExecutor {
	return &SortitionExecutor{sandbox: sb}
}

func (e *SortitionExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.SortitionPayload)

	val := e.sandbox.Validator(pld.Address)
	if val == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve validator")
	}
	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}
	if e.sandbox.CurrentHeight()-val.BondingHeight() < e.sandbox.TransactionToLiveInterval()+10 {
		return errors.Errorf(errors.ErrInvalidTx, "Too early to send sortition")
	}
	if !e.sandbox.VerifySortition(trx.Stamp(), pld.Proof, val) {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid proof or index")
	}
	if err := e.sandbox.AddToSet(trx.Stamp(), val.Address()); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, err.Error())
	}
	val.IncSequence()

	e.sandbox.UpdateValidator(val)

	return nil
}

func (e *SortitionExecutor) Fee() int64 {
	return 0
}
