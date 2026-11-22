package certificate

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrTooManyCommitters = errors.New("too many committers in certificate")
	ErrTooManyAbsentees  = errors.New("too many absentees in certificate")

	// ErrInvalidValidatorNumber is returned when a validator number is out of range.
	ErrInvalidValidatorNumber = errors.New("invalid validator number in certificate")
	// ErrTrailingBytes is returned when there are bytes after the encoded certificate.
	ErrTrailingBytes = errors.New("trailing bytes after the certificate")
	// ErrNoSignature is returned when the certificate has no signature.
	ErrNoSignature = errors.New("certificate has no signature")
)

// BasicCheckError is returned when the basic check on the certificate fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}

// UnexpectedCommittersError is returned when the list of committers
// does not match the expectations.
type UnexpectedCommittersError struct {
	Committers []int32
}

func (e UnexpectedCommittersError) Error() string {
	return fmt.Sprintf("certificate has an unexpected committers: %v",
		e.Committers)
}

func (e UnexpectedCommittersError) Is(target error) bool {
	return reflect.DeepEqual(e, target)
}

// InsufficientPowerError is returned when the accumulated power does not meet
// the required threshold.
type InsufficientPowerError struct {
	SignedPower   int64
	RequiredPower int64
}

func (e InsufficientPowerError) Error() string {
	return fmt.Sprintf("accumulated power is %v, should be at least %v",
		e.SignedPower, e.RequiredPower)
}
