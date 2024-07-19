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

func MakeExecutor(trx *tx.Tx, sb sandbox.Sandbox) (Executor, error) {
	var exe Executor
	var err error
	switch t := trx.Payload().Type(); t {
	case payload.TypeTransfer:
		exe, err = newTransferExecutor(trx, sb)
	case payload.TypeBond:
		exe, err = newBondExecutor(trx, sb)
	case payload.TypeUnbond:

	case payload.TypeWithdraw:

	case payload.TypeSortition:
		exe, err = newSortitionExecutor(trx, sb)
	default:

		return nil, InvalidPayloadTypeError{
			PayloadType: t,
		}
	}

	return exe, err
}
