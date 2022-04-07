package crypto

import (
	"bytes"
	"crypto/rand"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/zarbchain/zarb-go/errors"
)

// Address format:
// `zc1` + type + data + checksum

const (
	SignatureTypeBLS byte = 1
)

const (
	AddressSize           = 21
	hrpAddress            = "zc"
	treasuryAddressString = "000000000000000000000000000000000000000000"
)

var TreasuryAddress = Address{0}

type Address [AddressSize]byte

func AddressFromString(text string) (Address, error) {
	if text == treasuryAddressString {
		return TreasuryAddress, nil
	}

	hrp, data, err := bech32.DecodeToBase256(text)
	if err != nil {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, err.Error())
	}
	if hrp != hrpAddress {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "invalid hrp: %v", hrp)
	}
	if len(data) != AddressSize {
		return Address{}, errors.Errorf(errors.ErrInvalidAddress, "address should be %d bytes, but it is %v bytes", AddressSize, len(data))
	}
	var addr Address
	copy(addr[:], data[:])

	return addr, nil
}

func (addr Address) Bytes() []byte {
	return addr[:]
}

func (addr Address) Fingerprint() string {
	return addr.String()[0:12]
}

func (addr Address) String() string {
	if addr.EqualsTo(TreasuryAddress) {
		return treasuryAddressString
	}
	str, err := bech32.EncodeFromBase256(hrpAddress, addr.Bytes())
	if err != nil {
		panic(err.Error())
	}

	return str
}

func (addr *Address) SanityCheck() error {
	if addr[0] != 0 && addr[0] != 1 {
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
