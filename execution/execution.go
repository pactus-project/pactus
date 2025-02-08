package execution

import (
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
)

func Execute(trx *tx.Tx, sbx sandbox.Sandbox) error {
	exe, err := executor.MakeExecutor(trx, sbx)
	if err != nil {
		return err
	}

	exe.Execute()
	sbx.CommitTransaction(trx)

	return nil
}

func CheckAndExecute(trx *tx.Tx, sbx sandbox.Sandbox, strict bool) error {
	exe, err := executor.MakeExecutor(trx, sbx)
	if err != nil {
		return err
	}

	if sbx.IsBanned(trx) {
		return SignerBannedError{
			addr: trx.Payload().Signer(),
		}
	}

	if exists := sbx.RecentTransaction(trx.ID()); exists {
		return TransactionCommittedError{
			ID: trx.ID(),
		}
	}

	if err := CheckLockTime(trx, sbx, strict); err != nil {
		return err
	}

	if err := exe.Check(strict); err != nil {
		return err
	}

	exe.Execute()
	sbx.CommitTransaction(trx)

	return nil
}

func CheckLockTime(trx *tx.Tx, sbx sandbox.Sandbox, strict bool) error {
	interval := sbx.Params().TransactionToLiveInterval

	if trx.IsSubsidyTx() {
		interval = 0
	} else if trx.IsSortitionTx() {
		interval = sbx.Params().SortitionInterval
	}

	if sbx.CurrentHeight() > interval {
		if trx.LockTime() < sbx.CurrentHeight()-interval {
			return LockTimeExpiredError{
				LockTime: trx.LockTime(),
			}
		}
	}

	if strict {
		// In strict mode, transactions with future lock times are rejected.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if trx.LockTime() > sbx.CurrentHeight() {
			return LockTimeInFutureError{
				LockTime: trx.LockTime(),
			}
		}
	}

	return nil
}
