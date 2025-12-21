package types

import (
	"github.com/pactus-project/pactus/crypto"
)

// ListAddressOptions contains options for filtering addresses.
type ListAddressOptions struct {
	AddressTypes []crypto.AddressType
}

// ListAddressOption is a functional option for ListAddresses.
type ListAddressOption func(*ListAddressOptions)

// WithAddressType filters addresses by the specified type.
func WithAddressType(addressType crypto.AddressType) ListAddressOption {
	return func(opt *ListAddressOptions) {
		opt.AddressTypes = []crypto.AddressType{addressType}
	}
}

// OnlyValidatorAddresses filters to show only validator addresses.
func OnlyValidatorAddresses() ListAddressOption {
	return func(opt *ListAddressOptions) {
		opt.AddressTypes = []crypto.AddressType{crypto.AddressTypeValidator}
	}
}

// OnlyAccountAddresses filters to show only account addresses (BLS and Ed25519).
func OnlyAccountAddresses() ListAddressOption {
	return func(opt *ListAddressOptions) {
		opt.AddressTypes = []crypto.AddressType{
			crypto.AddressTypeBLSAccount,
			crypto.AddressTypeEd25519Account,
		}
	}
}

// NewAddressOptions contains options for creating new addresses.
type NewAddressOptions struct {
	Password string
}

// NewAddressOption is a functional option for NewAddress.
type NewAddressOption func(*NewAddressOptions)

// WithPassword sets the password for address creation required for Ed25519 accounts.
func WithPassword(password string) NewAddressOption {
	return func(opt *NewAddressOptions) {
		opt.Password = password
	}
}
