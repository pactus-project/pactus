package state

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/vote"
)

// ErrInvalidBlockVersion indicates that the block version is not valid.
var ErrInvalidBlockVersion = errors.New("invalid block version")

// ErrInvalidSubsidyTransaction indicates that the subsidy transaction is not valid.
var ErrInvalidSubsidyTransaction = errors.New("invalid subsidy transaction")

// ErrDuplicatedSubsidyTransaction indicates that there is more than one subsidy transaction
// inside the block.
var ErrDuplicatedSubsidyTransaction = errors.New("duplicated subsidy transaction")

// ErrInvalidSortitionSeed indicates that the block's sortition seed is either invalid or unverifiable.
var ErrInvalidSortitionSeed = errors.New("invalid sortition seed")

// ErrInvalidCertificate indicates that the block certificate is invalid.
var ErrInvalidCertificate = errors.New("invalid certificate")

// InvalidSubsidyAmountError is returned when the amount of the subsidy transaction is not as expected.
type InvalidSubsidyAmountError struct {
	Expected amount.Amount
	Got      amount.Amount
}

func (e InvalidSubsidyAmountError) Error() string {
	return fmt.Sprintf("invalid subsidy amount, expected: %v, got: %v", e.Expected, e.Got)
}

// InvalidVoteForCertificateError is returned when an attempt to update
// the last certificate with an invalid vote is made.
type InvalidVoteForCertificateError struct {
	Vote *vote.Vote
}

func (e InvalidVoteForCertificateError) Error() string {
	return fmt.Sprintf("invalid vote to update the last certificate: %s",
		e.Vote.Type().String())
}

// InvalidStateRootHashError is returned when the state root hash of the block
// does not match the current state root hash.
type InvalidStateRootHashError struct {
	Expected hash.Hash
	Got      hash.Hash
}

func (e InvalidStateRootHashError) Error() string {
	return fmt.Sprintf("invalid state root hash, expected: %s, got: %s",
		e.Expected, e.Got)
}

// InvalidProposerError is returned when the block proposer is not as expected.
type InvalidProposerError struct {
	Expected crypto.Address
	Got      crypto.Address
}

func (e InvalidProposerError) Error() string {
	return fmt.Sprintf("invalid block proposer, expected: %s, got: %s",
		e.Expected, e.Got)
}

// InvalidBlockTimeError is returned when the block time is not valid.
type InvalidBlockTimeError struct {
	Reason string
}

func (e InvalidBlockTimeError) Error() string {
	return fmt.Sprintf("invalid block time: %s", e.Reason)
}
