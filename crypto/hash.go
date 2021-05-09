package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/dchest/blake2b"
	cbor "github.com/fxamacker/cbor/v2"
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

func HashFromString(str string) (Hash, error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return Hash{}, err
	}
	return HashFromRawBytes(b)
}

func HashFromRawBytes(data []byte) (Hash, error) {
	if len(data) != HashSize {
		return Hash{}, fmt.Errorf("Hash should be %d bytes, but it is %v bytes", HashSize, len(data))
	}
	var h Hash
	copy(h.data.Hash[:], data[:HashSize])
	return h, nil
}

func HashH(data []byte) Hash {
	h, _ := HashFromRawBytes(Hash256(data))
	return h
}

func (h Hash) String() string {
	return hex.EncodeToString(h.data.Hash[:])
}

func (h Hash) RawBytes() []byte {
	return h.data.Hash[:]
}

func (h Hash) Fingerprint() string {
	return fmt.Sprintf("%X", h.data.Hash[:6])
}

func (h Hash) IsUndef() bool {
	return h.EqualsTo(UndefHash)
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h *Hash) UnmarshalText(text []byte) error {
	/// Unmarshal empty value
	if len(text) == 0 {
		return nil
	}

	hash, err := HashFromString(string(text))
	if err != nil {
		return err
	}

	*h = hash
	return nil
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hash) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return h.UnmarshalText([]byte(text))
}

func (h Hash) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(h.data.Hash)
}

func (h *Hash) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &h.data.Hash)
}

func (h Hash) SanityCheck() error {
	if h.EqualsTo(UndefHash) {
		return fmt.Errorf("Hash is not defined")
	}

	return nil
}

func (h Hash) EqualsTo(r Hash) bool {
	return h.data.Hash == r.data.Hash
}

// ---------
// For tests
func GenerateTestHash() Hash {
	p := make([]byte, 10)
	random := rand.Reader
	_, err := random.Read(p)
	if err != nil {
		return UndefHash
	}
	return HashH(p)
}
