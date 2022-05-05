package vault

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidPassword describes an error in which the password is invalid.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrAddressExists describes an error in which the address already exist
	// in wallet.
	ErrAddressExists = errors.New("address already exists")
)

// ErrAddressNotFound describes an error in which the address doesn't exist
// in wallet.
type ErrAddressNotFound struct {
	addr string
}

func NewErrAddressNotFound(addr string) error {
	return ErrAddressNotFound{addr: addr}
}

func (e ErrAddressNotFound) Error() string {
	return fmt.Sprintf("address not found: %s", e.addr)
}

// ErrUnknownMethod describes an error in which the method is not known.
type ErrUnknownMethod struct {
	name string
}

func NewErrUnknownMethod(name string) error {
	return ErrUnknownMethod{name: name}
}

func (e ErrUnknownMethod) Error() string {
	return fmt.Sprintf("unknown method: %s", e.name)
}
