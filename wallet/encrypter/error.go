package encrypter

import (
	"errors"
)

// ErrInvalidPassword describes an error in which the password is invalid.
var ErrInvalidPassword = errors.New("invalid password")

// ErrInvalidParam describes an error in which the encryption parameter is invalid.
var ErrInvalidParam = errors.New("invalid param")

// ErrInvalidCipher describes an error in which the cipher message is invalid.
var ErrInvalidCipher = errors.New("invalid cipher message")

// ErrMethodNotSupported describes an error in which the cipher method is not known.
var ErrMethodNotSupported = errors.New("cipher method is not supported")
