package voteset

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
)

// ErrDoubleVote is returned when a validator casts multiple different votes in the same round.
var ErrDoubleVote = errors.New("double vote")

// IneligibleVoterError is returned when the voter is not a member of the committee.
type IneligibleVoterError struct {
	Address crypto.Address
}

func (e IneligibleVoterError) Error() string {
	return fmt.Sprintf("validator %s is not part of the committee", e.Address)
}
