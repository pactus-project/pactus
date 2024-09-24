package errors

import (
	"fmt"
)

const (
	ErrNone = iota
	ErrGeneric
	ErrInvalidAddress
	ErrInvalidPublicKey
	ErrInvalidPrivateKey
	ErrInvalidSignature
	ErrInvalidHeight
	ErrInvalidRound
	ErrInvalidVote
	ErrInvalidMessage
	ErrDuplicateVote

	ErrCount
)

var messages = map[int]string{
	ErrNone:              "no error",
	ErrGeneric:           "generic error",
	ErrInvalidAddress:    "invalid address",
	ErrInvalidPublicKey:  "invalid public key",
	ErrInvalidPrivateKey: "invalid private key",
	ErrInvalidSignature:  "invalid signature",
	ErrInvalidHeight:     "invalid height",
	ErrInvalidRound:      "invalid round",
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

func Errorf(code int, format string, a ...any) error {
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
	e, ok := err.(i)
	if !ok {
		return ErrGeneric
	}

	return e.Code()
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
