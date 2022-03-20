package hash

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/util"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
)

const HashSize = 32

type Hash struct {
	data hashData
}

type hashData struct {
	Hash [HashSize]byte
}

var UndefHash = Hash{data: hashData{Hash: [HashSize]byte{0}}}

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
	b, err := hex.DecodeString(str)
	if err != nil {
		return Hash{}, err
	}
	return FromRawBytes(b)
}

func FromRawBytes(data []byte) (Hash, error) {
	if len(data) != HashSize {
		return Hash{}, fmt.Errorf("Hash should be %d bytes, but it is %v bytes", HashSize, len(data))
	}
	var h Hash
	copy(h.data.Hash[:], data[:HashSize])
	return h, nil
}

func CalcHash(data []byte) Hash {
	h, _ := FromRawBytes(Hash256(data))
	return h
}

func (h Hash) String() string {
	return hex.EncodeToString(h.data.Hash[:])
}

func (h Hash) RawBytes() []byte {
	return h.data.Hash[:]
}

func (h Hash) Stamp() Stamp {
	var stamp Stamp
	copy(stamp[:], h.data.Hash[0:4])
	return stamp
}

func (h Hash) Fingerprint() string {
	return fmt.Sprintf("%X", h.data.Hash[:6])
}

func (h Hash) IsUndef() bool {
	return h.EqualsTo(UndefHash)
}

func (h *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hash) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(h.data.Hash)
}

func (h *Hash) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &h.data.Hash)
}

func (h Hash) SanityCheck() error {
	if h.IsUndef() {
		return fmt.Errorf("hash is zero")
	}

	return nil
}

func (h Hash) EqualsTo(r Hash) bool {
	return h.data.Hash == r.data.Hash
}

// ---------
// For tests
func GenerateTestHash() Hash {
	return CalcHash(util.IntToSlice(util.RandInt(999999999)))
}

func GenerateTestStamp() Stamp {
	return GenerateTestHash().Stamp()
}
