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

	validator := e.sandbox.Validator(pld.Address)
	if validator == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve validator")
	}

	return nil
}
