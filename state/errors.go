package state

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/vote"
)

// ErrInvalidSubsidyTransaction indicates that the subsidy transaction is not valid.
var ErrInvalidSubsidyTransaction = errors.New("invalid subsidy transaction")

// ErrDuplicatedSubsidyTransaction indicates that there is more than one subsidy transaction
// inside the block.
var ErrDuplicatedSubsidyTransaction = errors.New("duplicated subsidy transaction")

// InvalidSubsidyAmountError is returned when the amount of the subsidy transaction is not as expected.
type InvalidSubsidyAmountError struct {
	Expected amount.Amount
	Got      amount.Amount
}

func (e InvalidSubsidyAmountError) Error() string {
	return fmt.Sprintf("invalid subsidy amount, expected %v, got %v", e.Expected, e.Got)
}

// InvalidVoteForCertificateError is returned when an attempt to update
// the last certificate with an invalid vote is made.
type InvalidVoteForCertificateError struct {
	Vote *vote.Vote
}

func (e InvalidVoteForCertificateError) Error() string {
	return fmt.Sprintf("invalid vote to update the last certificate: %s",
		e.Vote.String())
}

// InvalidBlockCertificateError is returned when the given certificate is invalid.
type InvalidBlockCertificateError struct {
	Cert *certificate.BlockCertificate
}

func (e InvalidBlockCertificateError) Error() string {
	return fmt.Sprintf("invalid certificate for block %d",
		e.Cert.Height())
}
