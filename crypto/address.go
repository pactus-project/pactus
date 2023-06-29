package crypto

import (
	"bytes"

	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/errors"
)

// Address format:
// hrp + `1` + type + data + checksum

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
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, err.Error())
	}

	// Check if hrp is valid
	if hrp != AddressHRP {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "invalid hrp: %v", hrp)
	}

	if typ != SignatureTypeBLS {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "invalid address key type: %v", typ)
	}

	// The regrouped data must be 20 bytes.
	if len(data) != 20 {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress,
			"address should be %d bytes, but it is %v bytes", AddressSize, len(data)+1)
	}

	var addr Address
	addr[0] = typ
	copy(addr[1:], data[:])

	return addr, nil
}

// Bytes returns the 21 bytes of the address data.
func (addr Address) Bytes() []byte {
	return addr[:]
}

// Fingerprint returns a short string for the address useful for logger.
func (addr Address) Fingerprint() string {
	return addr.String()[0:12]
}

// String returns a human-readable string for the address.
func (addr Address) String() string {
	if addr.EqualsTo(TreasuryAddress) {
		return treasuryAddressString
	}

	str, err := bech32m.EncodeFromBase256WithType(
		AddressHRP,
		SignatureTypeBLS,
		addr[1:])
	if err != nil {
		panic(err.Error())
	}

	return str
}

func (addr *Address) SanityCheck() error {
	if addr[0] == 0 {
		if !addr.EqualsTo(TreasuryAddress) {
			return errors.Errorf(errors.ErrInvalidAddress, "invalid address data")
		}
	} else if addr[0] != SignatureTypeBLS {
		return errors.Errorf(errors.ErrInvalidAddress, "invalid address type")
	}
	return nil
}

func (addr Address) EqualsTo(right Address) bool {
	return bytes.Equal(addr.Bytes(), right.Bytes())
}
