package hash

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/zarbchain/zarb-go/util"
	"golang.org/x/crypto/blake2b"
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

func (h Hash) Stamp() Stamp {
	var stamp Stamp
	copy(stamp[:], h[0:4])
	return stamp
}

func (h Hash) Fingerprint() string {
	return fmt.Sprintf("%X", h[:6])
}

func (h Hash) IsUndef() bool {
	return h.EqualsTo(UndefHash)
}

func (h *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h Hash) SanityCheck() error {
	if h.IsUndef() {
		return fmt.Errorf("hash is zero")
	}

	return nil
}

func (h Hash) EqualsTo(r Hash) bool {
	return h == r
}

// ---------
// For tests
func GenerateTestHash() Hash {
	return CalcHash(util.Uint64ToSlice(util.RandUint64(0)))
}

func GenerateTestStamp() Stamp {
	return GenerateTestHash().Stamp()
}
