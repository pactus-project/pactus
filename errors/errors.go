package errors

import (
	"fmt"
)

const (
	ErrNone = iota
	ErrGeneric
	ErrNetwork
	ErrInvalidBlock
	ErrInvalidAddress
	ErrInvalidPublicKey
	ErrInvalidPrivateKey
	ErrInvalidSignature
	ErrInvalidAmount
	ErrInvalidSequence
	ErrInvalidTx
	ErrInvalidProposal
	ErrInvalidVote
	ErrInvalidMessage
	ErrInvalidConfig
	ErrDuplicateVote
	ErrInsufficientFunds

	ErrCount
)

var messages = map[int]string{
	ErrNone:              "No error",
	ErrGeneric:           "Generic error",
	ErrNetwork:           "Network error",
	ErrInvalidBlock:      "Invalid block",
	ErrInvalidAddress:    "Invalid address",
	ErrInvalidPublicKey:  "Invalid public key",
	ErrInvalidPrivateKey: "Invalid private key",
	ErrInvalidSignature:  "Invalid signature",
	ErrInvalidAmount:     "Invalid amount",
	ErrInvalidSequence:   "Invalid sequence",
	ErrInvalidTx:         "Invalid transaction",
	ErrInvalidProposal:   "Invalid proposal",
	ErrInvalidVote:       "Invalid vote",
	ErrInvalidMessage:    "Invalid message",
	ErrInvalidConfig:     "Invalid config",
	ErrDuplicateVote:     "Duplicate vote",
	ErrInsufficientFunds: "Insufficient funds",
}

type withCode struct {
	code    int
	message string
}

func Error(code int) error {
	message, ok := messages[code]
	if !ok {
		message = "Unknown error code"
	}

	return &withCode{
		code:    code,
		message: message,
	}
}

func Errorf(code int, format string, a ...interface{}) error {
	message, ok := messages[code]
	if !ok {
		message = "Unknown error code"
	}

	return &withCode{
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
	_e, ok := err.(i)
	if !ok {
		return ErrGeneric
	}
	return _e.Code()
}

func (e *withCode) Error() string {
	return e.message
}

func (e *withCode) Code() int {
	return e.code
}
