package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
)

// Address format:
// `zc` + type + data + checksum
// type is 1 for BLS signatures

const (
	AddressTypeBLS byte = 1
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
		return Address{}, err
	}
	if hrp != hrpAddress {
		return Address{}, fmt.Errorf("invalid hrp: %v", hrp)
	}
	// TODO: fix me, Get type from decode function DecodeToBase256
	data = append([]byte{AddressTypeBLS}, data...)
	return AddressFromBytes(data)

}

func AddressFromBytes(bs []byte) (Address, error) {
	if len(bs) != AddressSize {
		return Address{}, fmt.Errorf("address should be %d bytes, but it is %v bytes", AddressSize, len(bs))
	}

	var addr Address
	copy(addr[:], bs[:])

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
	str, err := bech32.EncodeFromBase256(hrpAddress, addr[1:])
	if err != nil {
		panic(fmt.Sprintf("Invalid address. %v", err))
	}

	return str
}

func (addr *Address) SanityCheck() error {
	if addr[0] != 0 && addr[0] != 1 {
		return fmt.Errorf("invalid type")
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
	addr, _ := AddressFromBytes(data)
	return addr
}
