package execution

import (
	"fmt"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

// TransactionCommittedError is returned when an attempt is made
// to replay a transaction that has already been committed.
// This is to prevent replay attacks where an attacker tries to
// submit the same transaction more than once.
type TransactionCommittedError struct {
	ID tx.ID
}

func (e TransactionCommittedError) Error() string {
	return fmt.Sprintf("the transaction committed before: %s",
		e.ID.String())
}

// UnknownPayloadTypeError is returned when transaction payload type
// is not valid.
type UnknownPayloadTypeError struct {
	PayloadType payload.Type
}

func (e UnknownPayloadTypeError) Error() string {
	return fmt.Sprintf("unknown payload type: %s",
		e.PayloadType.String())
}

// PastLockTimeError is returned when the lock time of a transaction
// is in the past and has expired,
// indicating the transaction can no longer be executed.
type PastLockTimeError struct {
	LockTime uint32
}

func (e PastLockTimeError) Error() string {
	return fmt.Sprintf("lock time is in the past: %v", e.LockTime)
}

// FutureLockTimeError is returned when the lock time of a transaction
// is in the future,
// indicating the transaction is not yet eligible for processing.
type FutureLockTimeError struct {
	LockTime uint32
}

func (e FutureLockTimeError) Error() string {
	return fmt.Sprintf("lock time is in the future: %v", e.LockTime)
}

// InvalidFeeError is returned when the transaction fee is not valid.
type InvalidFeeError struct {
	Fee      amount.Amount
	Expected amount.Amount
}

func (e InvalidFeeError) Error() string {
	return fmt.Sprintf("fee is invalid, expected: %s, got: %s", e.Expected, e.Fee)
}
