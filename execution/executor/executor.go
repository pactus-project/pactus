package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type Executor interface {
	Check(strict bool) error
	Execute()
}

func MakeExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (Executor, error) {
	var exe Executor
	var err error
	switch typ := trx.Payload().Type(); typ {
	case payload.TypeTransfer:
		exe, err = newTransferExecutor(trx, sbx)
	case payload.TypeBond:
		exe, err = newBondExecutor(trx, sbx)
	case payload.TypeUnbond:
		exe, err = newUnbondExecutor(trx, sbx)
	case payload.TypeWithdraw:
		exe, err = newWithdrawExecutor(trx, sbx)
	case payload.TypeSortition:
		exe, err = newSortitionExecutor(trx, sbx)
	case payload.TypeBatchTransfer:
		exe, err = newBatchTransferExecutor(trx, sbx)
	default:
		return nil, InvalidPayloadTypeError{
			PayloadType: typ,
		}
	}

	return exe, err
}
