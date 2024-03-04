package hash

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2b"
	//nolint:staticcheck // use to hash the public key to get the address
	"golang.org/x/crypto/ripemd160"
)

const HashSize = 32

type Hash [HashSize]byte

var UndefHash = Hash{0}

func Hash256(data []byte) []byte {
	h := blake2b.Sum256(data)

	return h[:]
}

func Hash160(data []byte) []byte {
	h := ripemd160.New()
	n, err := h.Write(data)
	if err != nil {
		return nil
	}
	if n != len(data) {
		return nil
	}

	return h.Sum(nil)
}

func FromString(str string) (Hash, error) {
	data, err := hex.DecodeString(str)
	if err != nil {
		return Hash{}, err
	}
	if len(data) != HashSize {
		return Hash{}, fmt.Errorf("Hash should be %d bytes, but it is %v bytes", HashSize, len(data))
	}

	return FromBytes(data)
}

func FromBytes(data []byte) (Hash, error) {
	if len(data) != HashSize {
		return Hash{}, fmt.Errorf("Hash should be %d bytes, but it is %v bytes", HashSize, len(data))
	}

	var h Hash
	copy(h[:], data[:HashSize])

	return h, nil
}

func CalcHash(data []byte) Hash {
	h, _ := FromBytes(Hash256(data))

	return h
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) ShortString() string {
	return fmt.Sprintf("%X", h[:6])
}

func (h Hash) IsUndef() bool {
	return h == UndefHash
}
