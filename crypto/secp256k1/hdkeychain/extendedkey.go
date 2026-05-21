package hdkeychain

// References:
// - BIP-32: Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
// - SLIP-0010 : Universal private key derivation from master private key
//   https://github.com/satoshilabs/slips/blob/master/slip-0010.md

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"strings"

	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	// hardenedKeyStart is the index at which a hardened key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardened child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	hardenedKeyStart = uint32(0x80000000) // 2^31

	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits

	privateKeyLen = 32
)

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key per BIP-32.
type ExtendedKey struct {
	key       []byte // 32-byte private key or 33-byte compressed public key
	chainCode []byte
	path      []uint32
	isPrivate bool
}

// newExtendedKey returns a new instance of an extended key with the given
// fields. No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func newExtendedKey(key, chainCode []byte, path []uint32, isPrivate bool) *ExtendedKey {
	return &ExtendedKey{
		key:       key,
		chainCode: chainCode,
		path:      path,
		isPrivate: isPrivate,
	}
}

func isValidPrivateKey(key []byte) bool {
	var scalar secp.ModNScalar
	overflow := scalar.SetByteSlice(key)

	return !overflow && !scalar.IsZero()
}

func (k *ExtendedKey) pubKeyBytes() []byte {
	if !k.isPrivate {
		return k.key
	}

	var scalar secp.ModNScalar
	_ = scalar.SetByteSlice(k.key)

	return secp.NewPrivateKey(&scalar).PubKey().SerializeCompressed()
}

// IsPrivate returns whether or not the extended key is a private extended key.
func (k *ExtendedKey) IsPrivate() bool {
	return k.isPrivate
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
func (k *ExtendedKey) Derive(index uint32) (*ExtendedKey, error) {
	isChildHardened := index >= hardenedKeyStart

	data := make([]byte, 0, 37)
	if isChildHardened {
		if !k.isPrivate {
			return nil, ErrDeriveHardFromPublic
		}
		if len(k.key) != privateKeyLen {
			return nil, ErrInvalidKeyData
		}

		data = append(data, 0x00)
		data = append(data, k.key...)
	} else {
		pubKey := k.pubKeyBytes()
		if len(pubKey) != secp.PubKeyBytesLenCompressed {
			return nil, ErrInvalidKeyData
		}
		data = append(data, pubKey...)
	}

	indexData := make([]byte, 4)
	binary.BigEndian.PutUint32(indexData, index)
	data = append(data, indexData...)

	var childChainCode, il []byte

	hmac512 := hmac.New(sha512.New, k.chainCode)
	_, _ = hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	il = ilr[:len(ilr)/2]
	childChainCode = ilr[len(ilr)/2:]

	// Both derived public or private keys rely on treating the left 32-byte
	// sequence calculated above (Il) as a 256-bit integer that must be
	// within the valid range for a secp256k1 private key.  There is a small
	// chance (< 1 in 2^127) this condition will not hold, and in that case,
	// a child extended key can't be created for this index.
	if !isValidPrivateKey(il) {
		return nil, ErrInvalidKeyData
	}

	var childKey []byte
	if k.isPrivate {
		var parentKey, ilScalar secp.ModNScalar
		_ = parentKey.SetByteSlice(k.key)
		_ = ilScalar.SetByteSlice(il)

		childScalar := new(secp.ModNScalar).Add2(&parentKey, &ilScalar)
		childKeyBytes := childScalar.Bytes()
		childKey = childKeyBytes[:]
	} else {
		var ilScalar secp.ModNScalar
		_ = ilScalar.SetByteSlice(il)

		var ilPoint secp.JacobianPoint
		secp.ScalarBaseMultNonConst(&ilScalar, &ilPoint)

		parentPub, err := secp.ParsePubKey(k.key)
		if err != nil {
			return nil, err
		}

		var parentPoint secp.JacobianPoint
		parentPub.AsJacobian(&parentPoint)

		var childPoint secp.JacobianPoint
		secp.AddNonConst(&ilPoint, &parentPoint, &childPoint)
		childPoint.ToAffine()

		childPub := secp.NewPublicKey(&childPoint.X, &childPoint.Y)
		childKey = childPub.SerializeCompressed()
	}

	newPath := make([]uint32, 0, len(k.path)+1)
	newPath = append(newPath, k.path...)
	newPath = append(newPath, index)

	return newExtendedKey(childKey, childChainCode, newPath, k.isPrivate), nil
}

// Path returns the path of derived key.
func (k *ExtendedKey) Path() []uint32 {
	return k.path
}

// ChainCode returns a copy of the chain code.
func (k *ExtendedKey) ChainCode() []byte {
	return bytes.Clone(k.chainCode)
}

// RawPrivateKey returns the raw bytes of the private key.
func (k *ExtendedKey) RawPrivateKey() ([]byte, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	return bytes.Clone(k.key), nil
}

// RawPublicKey returns the raw bytes of the compressed public key.
func (k *ExtendedKey) RawPublicKey() []byte {
	return bytes.Clone(k.pubKeyBytes())
}

// Neuter returns a new extended public key from this extended private key.
func (k *ExtendedKey) Neuter() *ExtendedKey {
	if !k.isPrivate {
		return k
	}

	return newExtendedKey(k.pubKeyBytes(), k.chainCode, k.path, false)
}

// String returns the extended key as a bech32-encoded string.
func (k *ExtendedKey) String() string {
	//
	// The serialized format is structured as follows:
	// +-------+---------+------------+----------+------------+----------+
	// | Depth | Path    | Chain code | Reserved | Key length | Key data |
	// +-------+---------+------------+----------+------------+----------+
	// | 1     | depth*4 | 32         | 1        | 1          | 32/33    |
	// +-------+---------+------------+----------+------------+----------+
	//
	// Description:
	// - Depth: 1 byte representing the depth of derivation path.
	// - Path: serialized BIP-32 path; each entry is encoded as 32-bit unsigned integer, least significant byte first
	// - Chain code: 32 bytes chain code
	// - Reserved: 1 byte reserved and should set to 0.
	// - Key length: 1 byte representing the length of the key data.
	// - Key data: 32 bytes for a private key or 33 bytes for a compressed public key.
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

	hrp := crypto.XPublicKeyHRP
	if k.isPrivate {
		hrp = crypto.XPrivateKeyHRP
	}

	str, err := bech32m.EncodeFromBase256WithType(hrp, crypto.SignatureTypeSecp256k1, buf.Bytes())
	if err != nil {
		return err.Error()
	}

	if k.isPrivate {
		str = strings.ToUpper(str)
	}

	return str
}

// NewKeyFromString returns a new extended key instance from a bech32-encoded string.
func NewKeyFromString(str string) (*ExtendedKey, error) {
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(strings.ToLower(str))
	if err != nil {
		return nil, err
	}

	if typ != crypto.SignatureTypeSecp256k1 {
		return nil, ErrInvalidKeyData
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

	var isPrivate bool
	switch hrp {
	case crypto.XPrivateKeyHRP:
		isPrivate = true

	case crypto.XPublicKeyHRP:
		isPrivate = false

	default:
		return nil, crypto.InvalidHRPError(hrp)
	}

	return newExtendedKey(key, chainCode, path, isPrivate), nil
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
func NewMaster(seed []byte) (*ExtendedKey, error) {
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	masterKey := []byte("Bitcoin seed") // this is compatible with BIP-0032

	hmac512 := hmac.New(sha512.New, masterKey)
	_, _ = hmac512.Write(seed)
	ilr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	key := ilr[:32]
	chainCode := ilr[32:]

	if !isValidPrivateKey(key) {
		return nil, ErrUnusableSeed
	}

	return newExtendedKey(key, chainCode, []uint32{}, true), nil
}
