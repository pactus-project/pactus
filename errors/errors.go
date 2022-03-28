package errors

import (
	"fmt"
)

const (
	ErrNone = iota
	ErrGeneric
	ErrNetwork
	ErrInvalidBlock
	ErrInvalidAmount
	ErrInvalidFee
	ErrInvalidAddress
	ErrInvalidPublicKey
	ErrInvalidPrivateKey
	ErrInvalidSignature
	ErrInvalidSequence
	ErrInvalidTx
	ErrInvalidMemo
	ErrInvalidProof
	ErrInvalidHeight
	ErrInvalidProposal
	ErrInvalidVote
	ErrInvalidMessage
	ErrInvalidConfig
	ErrDuplicateVote
	ErrInsufficientFunds

	ErrCount
)

var messages = map[int]string{
	ErrNone:              "no error",
	ErrGeneric:           "generic error",
	ErrNetwork:           "network error",
	ErrInvalidBlock:      "invalid block",
	ErrInvalidAmount:     "invalid amount",
	ErrInvalidFee:        "invalid fee",
	ErrInvalidAddress:    "invalid address",
	ErrInvalidPublicKey:  "invalid public key",
	ErrInvalidPrivateKey: "invalid private key",
	ErrInvalidSignature:  "invalid signature",
	ErrInvalidSequence:   "invalid sequence",
	ErrInvalidTx:         "invalid transaction",
	ErrInvalidMemo:       "invalid memo",
	ErrInvalidProof:      "invalid proof",
	ErrInvalidHeight:     "invalid height",
	ErrInvalidProposal:   "invalid proposal",
	ErrInvalidVote:       "invalid vote",
	ErrInvalidMessage:    "invalid message",
	ErrInvalidConfig:     "invalid config",
	ErrDuplicateVote:     "duplicate vote",
	ErrInsufficientFunds: "insufficient funds",
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
