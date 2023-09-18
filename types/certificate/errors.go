package certificate

import (
	"fmt"
	"reflect"

	"github.com/pactus-project/pactus/crypto"
)

// BasicCheckError is returned when the basic check on the certificate fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return "certificate basic check failed: " + e.Reason
}

// UnexpectedHeightError is returned when the height of the certificate
// is invalid.
type UnexpectedHeightError struct {
	Expected uint32
	Got      uint32
}

func (e UnexpectedHeightError) Error() string {
	return fmt.Sprintf("certificate height is invalid (expected %v got %v)",
		e.Expected, e.Got)
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

// InvalidSignatureError is returned when the signature in the certificate is invalid.
type InvalidSignatureError struct {
	Signature crypto.Signature
}

func (e InvalidSignatureError) Error() string {
	return fmt.Sprintf("certificate has an invalid signature: %s",
		e.Signature.String())
}
