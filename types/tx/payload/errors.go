package payload

import "errors"

// ErrInvalidPublicKeySize is returned when the public key size is not valid.
var ErrInvalidPublicKeySize = errors.New("invalid public key size")

// BasicCheckError describes is returned when the basic check on the transaction's payload fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}
