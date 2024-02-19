package crypto

import (
	"io"
	"slices"

	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

// Address format:
// hrp + `1` + type + data + checksum

type AddressType byte

const (
	AddressTypeTreasury   AddressType = 0
	AddressTypeValidator  AddressType = 1
	AddressTypeBLSAccount AddressType = 2
)

const (
	SignatureTypeTreasury byte = 0
	SignatureTypeBLS      byte = 1
)

const (
	AddressSize           = 21
	treasuryAddressString = "000000000000000000000000000000000000000000"
)

var TreasuryAddress = Address{0}

type Address [AddressSize]byte

// AddressFromString decodes the string encoding of an address and returns
// the Address if text is a valid encoding for a known address type.
func AddressFromString(text string) (Address, error) {
	if text == treasuryAddressString {
		return TreasuryAddress, nil
	}

	// Decode the bech32m encoded address.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return Address{}, err
	}

	// Check if hrp is valid
	if hrp != AddressHRP {
		return Address{}, InvalidHRPError(hrp)
	}

	// check type is valid
	validTypes := []AddressType{AddressTypeValidator, AddressTypeBLSAccount}
	if !slices.Contains(validTypes, AddressType(typ)) {
		return Address{}, InvalidAddressTypeError(typ)
	}

	// check length is valid
	if len(data) != 20 {
		return Address{}, InvalidLengthError(len(data) + 1)
	}

	var addr Address
	addr[0] = typ
	copy(addr[1:], data)

	return addr, nil
}

// NewAddress create a new address based.
func NewAddress(typ AddressType, data []byte) Address {
	var addr Address
	addr[0] = byte(typ)
	copy(addr[1:], data)

	return addr
}

// Bytes returns the 21 bytes of the address data.
func (addr Address) Bytes() []byte {
	return addr[:]
}

// ShortString returns a short string for the address useful for logger.
func (addr Address) ShortString() string {
	return addr.String()[0:12]
}

// String returns a human-readable string for the address.
func (addr Address) String() string {
	if addr == TreasuryAddress {
		return treasuryAddressString
	}

	str, err := bech32m.EncodeFromBase256WithType(
		AddressHRP,
		addr[0],
		addr[1:])
	if err != nil {
		panic(err.Error())
	}

	return str
}

func (addr Address) Type() AddressType {
	return AddressType(addr[0])
}

func (addr Address) Encode(w io.Writer) error {
	switch t := addr.Type(); t {
	case AddressTypeTreasury:
		return encoding.WriteElement(w, uint8(0))
	case AddressTypeValidator,
		AddressTypeBLSAccount:
		return encoding.WriteElement(w, addr)
	default:
		return InvalidAddressTypeError(t)
	}
}

func (addr *Address) Decode(r io.Reader) error {
	err := encoding.ReadElement(r, &addr[0])
	if err != nil {
		return err
	}
	switch t := addr.Type(); t {
	case AddressTypeTreasury:
		return nil
	case AddressTypeValidator,
		AddressTypeBLSAccount:
		return encoding.ReadElement(r, addr[1:])
	default:
		return InvalidAddressTypeError(t)
	}
}

// SerializeSize returns the number of bytes it would take to serialize the address.
func (addr Address) SerializeSize() int {
	switch t := addr.Type(); t {
	case AddressTypeTreasury:
		return 1
	case AddressTypeValidator,
		AddressTypeBLSAccount:
		return AddressSize
	default:
		return 0
	}
}

func (addr Address) IsTreasuryAddress() bool {
	return addr.Type() == AddressTypeTreasury
}

func (addr Address) IsAccountAddress() bool {
	return addr.Type() == AddressTypeTreasury ||
		addr.Type() == AddressTypeBLSAccount
}

func (addr Address) IsValidatorAddress() bool {
	return addr.Type() == AddressTypeValidator
}
