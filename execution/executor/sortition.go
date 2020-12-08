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

func NewSortitionExecutor(sandbox sandbox.Sandbox) *SortitionExecutor {
	return &SortitionExecutor{sandbox}
}

func (e *SortitionExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.SortitionPayload)

	val := e.sandbox.Validator(pld.Address)
	if val == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve validator")
	}
	if trx.Fee() != 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: 0, got: %v", trx.Fee())
	}

	e.sandbox.AddToSet(val)

	return nil
}

func (e *SortitionExecutor) Fee() int64 {
	return 0
}
