package hdkeychain

// References:
//  SLIP-0010: Universal private key derivation from master private key
//  https://github.com/satoshilabs/slips/blob/master/slip-0010.md

import (
	"bytes"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"strings"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	// hardenedKeyStart is the index at which a hardened key starts.
	hardenedKeyStart = uint32(0x80000000) // 2^31

	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits
)

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key.
type ExtendedKey struct {
	key       []byte // This will be the bytes of extended public or private key
	chainCode []byte
	path      []uint32
}

// newExtendedKey returns a new instance of an extended key with the given
// fields. No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func newExtendedKey(key, chainCode []byte, path []uint32) *ExtendedKey {
	return &ExtendedKey{
		key:       key,
		chainCode: chainCode,
		path:      path,
	}
}

// DerivePath returns a derived child extended key from this master key at the
// given path.
func (k *ExtendedKey) DerivePath(path []uint32) (*ExtendedKey, error) {
	ext := k
	var err error
	for _, index := range path {
		ext, err = ext.Derive(index)
		if err != nil {
			return nil, err
		}
	}

	return ext, nil
}

// Derive returns a derived child extended key at the given index.
//
// For ed25519 and curve25519 the private keys are no longer multipliers for the group generator;
// instead the hash of the private key is the multiplier.
// For this reason, our scheme for ed25519 and curve25519 does not support public key derivation and
// uses the produced hashes directly as private keys.
func (k *ExtendedKey) Derive(index uint32) (*ExtendedKey, error) {
	isChildHardened := index >= hardenedKeyStart

	if !isChildHardened {
		return nil, ErrNonHardenedPath
	}

	// Calculate derive Data:
	//   Data = 0x00 || ser_256(k_par) || ser_32(i)
	indexData := make([]byte, 4)
	binary.BigEndian.PutUint32(indexData, index)

	data := make([]byte, 0, 37)
	data = append(data, 0x00)
	data = append(data, k.key...)
	data = append(data, indexData...)

	// Take the HMAC-SHA512 of the current key's chain code and the derived
	// data:
	//   I = HMAC-SHA512(Key = chainCode, Data = data)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	_, _ = hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	// Split I into two 32-byte sequences, IL and IR.
	// The returned chain code ci is IR.
	// The returned child key ki is IL.
	childChainCode := ilr[32:]
	childKey := ilr[:32]

	newPath := make([]uint32, 0, len(k.path)+1)
	newPath = append(newPath, k.path...)
	newPath = append(newPath, index)

	return newExtendedKey(childKey, childChainCode, newPath), nil
}

// Path returns the path of derived key.
//
// Path values are always between 2^31 and 2^32-1 as they are hardened keys.
func (k *ExtendedKey) Path() []uint32 {
	return k.path
}

// RawPrivateKey returns the raw bytes of the private key.
func (k *ExtendedKey) RawPrivateKey() []byte {
	return k.key
}

// RawPublicKey returns the raw bytes of the public key.
func (k *ExtendedKey) RawPublicKey() []byte {
	pub := ed25519.NewKeyFromSeed(k.key).Public()

	return pub.(ed25519.PublicKey)[:]
}

// String returns the extended key as a bech32-encoded string.
func (k *ExtendedKey) String() string {
	//
	// The serialized format is structured as follows:
	// +-------+---------+------------+----------+------------+----------+
	// | Depth | Path    | Chain code | Reserved | Key length | Key data |
	// +-------+---------+------------+----------+------------+----------+
	// | 1     | depth*4 | 32         | 1        | 1          | 32       |
	// +-------+---------+------------+----------+------------+----------+
	//
	// Description:
	// - Depth: 1 byte representing the depth of derivation path.
	// - Path: serialized BIP-32 path; each entry is encoded as 32-bit unsigned integer, least significant byte first
	// - Chain code: 32 bytes chain code
	// - Reserved: 1 byte reserved and should set to 0.
	// - Key length: 1 byte representing the length of the key data that is 32.
	// - Key data: The key data that is 32 bytes.
	//

	buf := bytes.NewBuffer(make([]byte, 0))
	err := encoding.WriteElement(buf, byte(len(k.path)))
	if err != nil {
		return err.Error()
	}

	for _, p := range k.path {
		err := encoding.WriteElement(buf, p)
		if err != nil {
			return err.Error()
		}
	}
	err = encoding.WriteVarBytes(buf, k.chainCode)
	if err != nil {
		return err.Error()
	}

	err = encoding.WriteElement(buf, uint8(0))
	if err != nil {
		return err.Error()
	}

	err = encoding.WriteVarBytes(buf, k.key)
	if err != nil {
		return err.Error()
	}

	str, err := bech32m.EncodeFromBase256WithType(crypto.XPrivateKeyHRP, crypto.SignatureTypeEd25519, buf.Bytes())
	if err != nil {
		return err.Error()
	}

	str = strings.ToUpper(str)

	return str
}

// NewKeyFromString returns a new extended key instance from a bech32-encoded string.
func NewKeyFromString(str string) (*ExtendedKey, error) {
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(strings.ToLower(str))
	if err != nil {
		return nil, err
	}

	if typ != crypto.SignatureTypeEd25519 {
		return nil, ErrInvalidKeyData
	}

	if hrp != crypto.XPrivateKeyHRP {
		return nil, ErrInvalidHRP
	}

	reader := bytes.NewReader(data)
	depth := uint8(0)
	err = encoding.ReadElement(reader, &depth)
	if err != nil {
		return nil, err
	}

	path := make([]uint32, depth)
	for i := byte(0); i < depth; i++ {
		err := encoding.ReadElement(reader, &path[i])
		if err != nil {
			return nil, err
		}
	}

	chainCode, err := encoding.ReadVarBytes(reader)
	if err != nil {
		return nil, err
	}

	var res uint8
	err = encoding.ReadElement(reader, &res)
	if err != nil {
		return nil, err
	}

	key, err := encoding.ReadVarBytes(reader)
	if err != nil {
		return nil, err
	}

	return newExtendedKey(key, chainCode, path), nil
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
func NewMaster(seed []byte) (*ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = Curve, Data = Seed)
	curve := []byte("ed25519 seed")
	hmac512 := hmac.New(sha512.New, curve)
	_, _ = hmac512.Write(seed)
	ilr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master key
	//   Ir = master chain code
	masterChainCode := ilr[32:]
	masterKey := ilr[:32]

	return newExtendedKey(masterKey, masterChainCode, []uint32{}), nil
}

// GenerateSeed returns a cryptographically secure random seed that can be used
// as the input for the NewMaster function to generate a new master node.
//
// The length is in bytes and it must be between 16 and 64 (128 to 512 bits).
// The recommended length is 32 (256 bits) as defined by the RecommendedSeedLen
// constant.
func GenerateSeed(length uint8) ([]byte, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if length < MinSeedBytes || length > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
