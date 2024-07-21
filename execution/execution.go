package execution

import (
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
)

func Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	exe, err := executor.MakeExecutor(trx, sb)
	if err != nil {
		return err
	}

	exe.Execute()
	sb.CommitTransaction(trx)

	return nil
}

func CheckAndExecute(trx *tx.Tx, sb sandbox.Sandbox, strict bool) error {
	exe, err := executor.MakeExecutor(trx, sb)
	if err != nil {
		return err
	}

	if sb.IsBanned(trx.Payload().Signer()) {
		return SignerBannedError{
			addr: trx.Payload().Signer(),
		}
	}

	if exists := sb.AnyRecentTransaction(trx.ID()); exists {
		return TransactionCommittedError{
			ID: trx.ID(),
		}
	}

	if err := CheckLockTime(trx, sb, strict); err != nil {
		return err
	}

	if err := CheckFee(trx); err != nil {
		return err
	}

	if err := exe.Check(strict); err != nil {
		return err
	}

	exe.Execute()
	sb.CommitTransaction(trx)

	return nil
}

func CheckLockTime(trx *tx.Tx, sb sandbox.Sandbox, strict bool) error {
	interval := sb.Params().TransactionToLiveInterval

	if trx.IsSubsidyTx() {
		interval = 0
	} else if trx.IsSortitionTx() {
		interval = sb.Params().SortitionInterval
	}

	if sb.CurrentHeight() > interval {
		if trx.LockTime() < sb.CurrentHeight()-interval {
			return LockTimeExpiredError{
				LockTime: trx.LockTime(),
			}
		}
	}

	if strict {
		// In strict mode, transactions with future lock times are rejected.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if trx.LockTime() > sb.CurrentHeight() {
			return LockTimeInFutureError{
				LockTime: trx.LockTime(),
			}
		}
	}

	return nil
}

func CheckFee(trx *tx.Tx) error {
	// TODO: This check maybe can be done in BasicCheck?
	if trx.IsFreeTx() {
		if trx.Fee() != 0 {
			return InvalidFeeError{
				Fee:      trx.Fee(),
				Expected: 0,
			}
		}
	}

	return nil
}
