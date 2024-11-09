package execution

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/tx"
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

// LockTimeExpiredError is returned when the lock time of a transaction
// is in the past and has expired,
// indicating the transaction can no longer be executed.
type LockTimeExpiredError struct {
	LockTime uint32
}

func (e LockTimeExpiredError) Error() string {
	return fmt.Sprintf("lock time expired: %v", e.LockTime)
}

// LockTimeInFutureError is returned when the lock time of a transaction
// is in the future,
// indicating the transaction is not yet eligible for processing.
type LockTimeInFutureError struct {
	LockTime uint32
}

func (e LockTimeInFutureError) Error() string {
	return fmt.Sprintf("lock time is in the future: %v", e.LockTime)
}

// SignerBannedError is returned when the signer of transaction is banned and its assets is freezed.
type SignerBannedError struct {
	addr crypto.Address
}

func (e SignerBannedError) Error() string {
	return fmt.Sprintf("the signer is banned: %s", e.addr.String())
}
