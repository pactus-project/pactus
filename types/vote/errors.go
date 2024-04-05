package vote

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
)

// BasicCheckError is returned when the basic check on the transaction fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}

// InvalidSignerError is returned when the vote signer does not match with the
// public key.
type InvalidSignerError struct {
	Expected crypto.Address
	Got      crypto.Address
}

func (e InvalidSignerError) Error() string {
	return fmt.Sprintf("invalid signer, expected: %s, got: %s",
		e.Expected.String(), e.Got.String())
}
