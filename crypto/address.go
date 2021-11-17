package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
	cbor "github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/errors"
)

// Address format
// `zc` + type + data
// type is 1 for schnorr signatures

const (
	addressSize           = 20
	hrpAddress            = "zc"
	treasuryAddressString = "0000000000000000000000000000000000000000"
)

var TreasuryAddress = Address{
	data: addressData{
		Address: [addressSize]byte{0},
	},
}

type Address struct {
	data addressData
}

type addressData struct {
	Address [addressSize]byte
}

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
	return AddressFromRawBytes(data)

}

func AddressFromRawBytes(bs []byte) (Address, error) {
	if len(bs) != addressSize {
		return Address{}, fmt.Errorf("address should be %d bytes, but it is %v bytes", addressSize, len(bs))
	}

	var addr Address
	copy(addr.data.Address[:], bs[:])

	return addr, nil
}

func (addr Address) RawBytes() []byte {
	return addr.data.Address[:]
}

func (addr Address) Fingerprint() string {
	return addr.String()[0:12]
}

func (addr Address) String() string {
	if addr.EqualsTo(TreasuryAddress) {
		return treasuryAddressString
	}
	str, err := bech32.EncodeFromBase256(hrpAddress, addr.data.Address[:])
	if err != nil {
		panic(fmt.Sprintf("Invalid address. %v", err))
	}

	return str
}

func (addr Address) MarshalText() ([]byte, error) {
	return []byte(addr.String()), nil
}

func (addr *Address) UnmarshalText(text []byte) error {
	/// Unmarshal empty value
	if len(text) == 0 {
		return nil
	}

	a, err := AddressFromString(string(text))
	if err != nil {
		return err
	}

	*addr = a
	return nil
}

func (addr Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(addr.String())
}

func (addr *Address) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return addr.UnmarshalText([]byte(text))
}

func (addr Address) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(addr.data.Address)
}

func (addr *Address) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &addr.data.Address)
}

func (addr *Address) SanityCheck() error {
	if addr.EqualsTo(TreasuryAddress) {
		return errors.Errorf(errors.ErrInvalidAddress, "")
	}
	return nil
}

func (addr Address) Verify(pb PublicKey) bool {
	return pb.Address().EqualsTo(addr)
}

func (addr Address) EqualsTo(right Address) bool {
	return bytes.Equal(addr.RawBytes(), right.RawBytes())
}
