package proposal

import "errors"

// BasicCheckError is returned when the basic check on the proposal fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}

// ErrNoSignature is returned when the proposal has no signature.
var ErrNoSignature = errors.New("no signature")
