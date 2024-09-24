package voteset

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
)

// ErrDuplicatedVote is returned when a duplicated vote from a validator is detected.
var ErrDuplicatedVote = errors.New("duplicated vote")

// IneligibleVoterError is returned when the voter is not a member of the committee.
type IneligibleVoterError struct {
	Address crypto.Address
}

func (e IneligibleVoterError) Error() string {
	return fmt.Sprintf("validator %s is not part of the committee", e.Address)
}
