package crypto

import (
	"bytes"
	"crypto/rand"

	"github.com/zarbchain/zarb-go/util/bech32m"
	"github.com/zarbchain/zarb-go/util/errors"
)

// Address format:
// hrp + `1` + type + data + checksum

const (
	SignatureTypeBLS byte = 1
)

const (
	AddressSize           = 21
	treasuryAddressString = "000000000000000000000000000000000000000000"
)

var TreasuryAddress = Address{0}
var DefaultHRP = "zc"

type Address [AddressSize]byte

/// AddressFromString decodes the string encoding of an address and returns
/// the Address if text is a valid encoding for a known address type.
func AddressFromString(text string) (Address, error) {
	if text == treasuryAddressString {
		return TreasuryAddress, nil
	}

	// Decode the bech32m encoded address.
	hrp, data, err := bech32m.Decode(text)
	if err != nil {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, err.Error())
	}

	// Check if hrp is valid
	if hrp != DefaultHRP {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "invalid hrp: %v", hrp)
	}

	// The first byte of the decoded address is the signature type, it must
	// exist.
	if len(data) < 1 {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "no address type")
	}

	// ...and should be 1 for BLS signature.
	sigType := data[0]
	if sigType != SignatureTypeBLS {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "invalid address type: %v", sigType)
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32m.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, err.Error())
	}

	// The regrouped data must be 20 bytes.
	if len(regrouped) != 20 {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "address should be %d bytes, but it is %v bytes", AddressSize, len(data)+1)
	}

	var addr Address
	addr[0] = sigType
	copy(addr[1:], regrouped[:])

	return addr, nil
}

// Bytes returns the 21 bytes of the address data.
func (addr Address) Bytes() []byte {
	return addr[:]
}

/// Fingerprint returns a short string for the address useful for logger.
func (addr Address) Fingerprint() string {
	return addr.String()[0:12]
}

/// String returns a human-readable string for the address.
func (addr Address) String() string {
	if addr.EqualsTo(TreasuryAddress) {
		return treasuryAddressString
	}

	// Group the address bytes into 5 bit groups, as this is what is used to
	// encode each character in the address string.
	converted, err := bech32m.ConvertBits(addr[1:], 8, 5, true)
	if err != nil {
		panic(err.Error())
	}

	// Concatenate the address type and program, and encode the resulting
	// bytes using bech32m encoding.
	combined := make([]byte, len(converted)+1)
	combined[0] = addr[0]
	copy(combined[1:], converted)
	str, err := bech32m.Encode(DefaultHRP, combined)
	if err != nil {
		panic(err.Error())
	}

	return str
}

func (addr *Address) SanityCheck() error {
	if addr[0] == 0 {
		if !addr.EqualsTo(TreasuryAddress) {
			return errors.Errorf(errors.ErrInvalidAddress, "invalid data")
		}
	} else if addr[0] != SignatureTypeBLS {
		return errors.Errorf(errors.ErrInvalidAddress, "invalid type")
	}
	return nil
}

func (addr Address) EqualsTo(right Address) bool {
	return bytes.Equal(addr.Bytes(), right.Bytes())
}

/// For tests
func GenerateTestAddress() Address {
	data := make([]byte, 20)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	data = append([]byte{1}, data...)
	var addr Address
	copy(addr[:], data[:])
	return addr
}
