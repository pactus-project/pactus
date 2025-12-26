package crypto

import (
	"database/sql/driver"
	"io"
	"slices"

	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

// Address format:
// hrp + `1` + type + data + checksum

const (
	SignatureTypeBLS     byte = 1
	SignatureTypeEd25519 byte = 3
)

const (
	AddressSize           = 21
	treasuryAddressString = "000000000000000000000000000000000000000000"
)

var TreasuryAddress = Address{0}

type Address [AddressSize]byte

// AddressFromString decodes the input string and returns the Address
// if the string is a valid bech32m encoding of an address.
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
	validTypes := []AddressType{
		AddressTypeValidator,
		AddressTypeBLSAccount,
		AddressTypeEd25519Account,
	}
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

// LogString returns a concise string representation intended for use in logs.
func (addr Address) LogString() string {
	return addr.String()[0:12]
}

// String returns a human-readable string for the address.
func (addr Address) String() string {
	if addr == TreasuryAddress {
		return treasuryAddressString
	}

	str, _ := bech32m.EncodeFromBase256WithType(
		AddressHRP,
		addr[0],
		addr[1:])

	return str
}

func (addr Address) Type() AddressType {
	return AddressType(addr[0])
}

func (addr Address) Encode(w io.Writer) error {
	switch typ := addr.Type(); typ {
	case AddressTypeTreasury:
		return encoding.WriteElement(w, uint8(0))
	case AddressTypeValidator,
		AddressTypeBLSAccount,
		AddressTypeEd25519Account:
		return encoding.WriteElement(w, addr)
	default:
		return InvalidAddressTypeError(typ)
	}
}

func (addr *Address) Decode(r io.Reader) error {
	err := encoding.ReadElement(r, &addr[0])
	if err != nil {
		return err
	}
	switch typ := addr.Type(); typ {
	case AddressTypeTreasury:
		return nil
	case AddressTypeValidator,
		AddressTypeBLSAccount,
		AddressTypeEd25519Account:
		return encoding.ReadElement(r, addr[1:])
	default:
		return InvalidAddressTypeError(typ)
	}
}

// SerializeSize returns the number of bytes it would take to serialize the address.
func (addr Address) SerializeSize() int {
	switch typ := addr.Type(); typ {
	case AddressTypeTreasury:
		return 1
	case AddressTypeValidator,
		AddressTypeBLSAccount,
		AddressTypeEd25519Account:
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
		addr.Type() == AddressTypeBLSAccount ||
		addr.Type() == AddressTypeEd25519Account
}

func (addr Address) IsValidatorAddress() bool {
	return addr.Type() == AddressTypeValidator
}

// Value implements driver.Valuer to store the address as raw bytes.
func (addr Address) Value() (driver.Value, error) {
	return addr[:], nil
}

// Scan implements sql.Scanner to read the address from raw bytes.
func (addr *Address) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		copy(addr[:], v)

		return nil
	default:
		return ErrInvalidSQLType
	}
}
