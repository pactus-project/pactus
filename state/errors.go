package state

import (
	"fmt"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/vote"
)

// InvalidVoteForCertificateError is returned when an attempt to update
// the last certificate with an invalid vote is made.
type InvalidVoteForCertificateError struct {
	Vote *vote.Vote
}

func (e InvalidVoteForCertificateError) Error() string {
	return fmt.Sprintf("invalid vote to update the last certificate: %s",
		e.Vote.String())
}

// InvalidCertificateError is returned when the given certificate is invalid.
type InvalidCertificateError struct {
	Cert *certificate.Certificate
}

func (e InvalidCertificateError) Error() string {
	return fmt.Sprintf("invalid certificate for block %d",
		e.Cert.Height())
}
