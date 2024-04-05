package errors

import (
	"fmt"
)

const (
	ErrNone = iota
	ErrGeneric
	ErrInvalidBlock
	ErrInvalidAmount
	ErrInvalidAddress
	ErrInvalidPublicKey
	ErrInvalidPrivateKey
	ErrInvalidSignature
	ErrInvalidTx
	ErrInvalidProof
	ErrInvalidHeight
	ErrInvalidRound
	ErrInvalidProposal
	ErrInvalidVote
	ErrInvalidMessage
	ErrDuplicateVote

	ErrCount
)

var messages = map[int]string{
	ErrNone:              "no error",
	ErrGeneric:           "generic error",
	ErrInvalidBlock:      "invalid block",
	ErrInvalidAmount:     "invalid amount",
	ErrInvalidAddress:    "invalid address",
	ErrInvalidPublicKey:  "invalid public key",
	ErrInvalidPrivateKey: "invalid private key",
	ErrInvalidSignature:  "invalid signature",
	ErrInvalidTx:         "invalid transaction",
	ErrInvalidProof:      "invalid proof",
	ErrInvalidHeight:     "invalid height",
	ErrInvalidRound:      "invalid round",
	ErrInvalidProposal:   "invalid proposal",
	ErrInvalidVote:       "invalid vote",
	ErrInvalidMessage:    "invalid message",
	ErrDuplicateVote:     "duplicate vote",
}

type withCodeError struct {
	code    int
	message string
}

func Error(code int) error {
	message, ok := messages[code]
	if !ok {
		message = "Unknown error code"
	}

	return &withCodeError{
		code:    code,
		message: message,
	}
}

func Errorf(code int, format string, a ...interface{}) error {
	message, ok := messages[code]
	if !ok {
		message = "Unknown error code"
	}

	return &withCodeError{
		code:    code,
		message: message + ": " + fmt.Sprintf(format, a...),
	}
}

func Code(err error) int {
	type i interface {
		Code() int
	}

	if err == nil {
		return ErrNone
	}
	_e, ok := err.(i) //nolint
	if !ok {
		return ErrGeneric
	}

	return _e.Code()
}

func (e *withCodeError) Error() string {
	return e.message
}

func (e *withCodeError) Code() int {
	return e.code
}

func (e *withCodeError) Is(target error) bool {
	return e.code == Code(target)
}
