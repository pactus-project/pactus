package block

import "errors"

var ErrTooManyTransactions = errors.New("too many transactions in block")

// ErrTrailingBytes is returned when there are bytes after the encoded block.
var ErrTrailingBytes = errors.New("trailing bytes after the block")

// BasicCheckError is returned when the basic check on the certificate fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}
