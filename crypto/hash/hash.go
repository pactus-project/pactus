package hash

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160" //nolint:all // used to hash the public key to generate the address.
)

const HashSize = 32

type Hash [HashSize]byte

var UndefHash = Hash{0}

func Hash256(data []byte) []byte {
	h := blake2b.Sum256(data)

	return h[:]
}

func Hash160(data []byte) []byte {
	//nolint:all // used to hash the public key to generate the address.
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

// FromString decodes the input string and returns the Hash
// if the string is a valid hexadecimal encoding of a hash.
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

// FromBytes constructs a Hash from the raw bytes.
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

// String returns the hex-encoded string representation of the hash.
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// Bytes returns the raw byte representation of the hash.
func (h Hash) Bytes() []byte {
	return h[:]
}

// LogString returns a concise string representation intended for use in logs.
func (h Hash) LogString() string {
	return fmt.Sprintf("%X", h[:6])
}

func (h Hash) IsUndef() bool {
	return h == UndefHash
}
