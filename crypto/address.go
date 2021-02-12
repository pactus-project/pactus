package crypto

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/errors"
)

const AddressSize = 20

type Address struct {
	data addressData
}

type addressData struct {
	Address [AddressSize]byte
}

func AddressFromString(text string) (Address, error) {
	bs, err := hex.DecodeString(text)
	if err != nil {
		return Address{}, err
	}

	return AddressFromRawBytes(bs)
}

func AddressFromRawBytes(bs []byte) (Address, error) {
	if len(bs) != AddressSize {
		return Address{}, fmt.Errorf("Address should be %d bytes, but it is %v bytes", AddressSize, len(bs))
	}

	var addr Address
	copy(addr.data.Address[:], bs[:])

	return addr, nil
}

/// -------
/// CASTING

func (addr Address) RawBytes() []byte {
	return addr.data.Address[:]
}

func (addr Address) Fingerprint() string {
	return hex.EncodeToString(addr.data.Address[:6])
}

func (addr Address) String() string {
	return hex.EncodeToString(addr.data.Address[:])
}

/// ----------
/// MARSHALING

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
	bz, err := addr.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(bz))
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
