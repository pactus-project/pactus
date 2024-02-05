package hdkeychain

// References:
//  [PIP-11]: Deterministic key hierarchy for BLS12-381 curve
//  https://pips.pactus.org/PIPs/pip-11

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"math/big"
	"strings"

	bls12381 "github.com/kilic/bls12-381"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	// HardenedKeyStart is the index at which a hardened key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardened child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = uint32(0x80000000) // 2^31

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
	isPrivate bool
	pubOnG1   bool
}

// newExtendedKey returns a new instance of an extended key with the given
// fields. No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func newExtendedKey(key, chainCode []byte, path []uint32, isPrivate, pubOnG1 bool) *ExtendedKey {
	return &ExtendedKey{
		key:       key,
		chainCode: chainCode,
		path:      path,
		isPrivate: isPrivate,
		pubOnG1:   pubOnG1,
	}
}

// pubKeyBytes returns bytes for the serialized public key associated with this
// extended key.
//
// When the extended key is already a public key, the key is simply returned as
// is since it's already in the correct form.  However, when the extended key is
// a private key, the public key will be calculated.
func (k *ExtendedKey) pubKeyBytes() []byte {
	// Just return the key if it's already an extended public key.
	if !k.isPrivate {
		return k.key
	}
	g1 := bls12381.NewG1()

	privKey := bls12381.NewFr()
	privKey.FromBytes(k.key)
	if k.pubOnG1 {
		pub := new(bls12381.PointG1)
		g1.MulScalar(pub, g1.One(), privKey)

		return g1.ToCompressed(pub)
	}

	g2 := bls12381.NewG2()

	pub := new(bls12381.PointG2)
	g2.MulScalar(pub, g2.One(), privKey)

	return g2.ToCompressed(pub)
}

// IsPrivate returns whether or not the extended key is a private extended key.
//
// A private extended key can be used to derive both hardened and non-hardened
// child private and public extended keys. A public extended key can only be
// used to derive non-hardened child public extended keys.
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
//
// When this extended key is a private extended key (as determined by the IsPrivate
// function), a private extended key will be derived. Otherwise, the derived
// extended key will be a public extended key.
//
// When the index is greater to or equal than the HardenedKeyStart constant, the
// derived extended key will be a hardened extended key.  It is only possible to
// derive a hardened extended key from a private extended key. Consequently,
// this function will return ErrDeriveHardFromPublic if a hardened child
// extended key is requested from a public extended key.
//
// A hardened extended key is useful since, as previously mentioned, it requires
// a parent private extended key to derive. In other words, normal child
// extended public keys can be derived from a parent public extended key (no
// knowledge of the parent private key) whereas hardened extended keys may not
// be.
//
//nolint:nestif // complexity can't be reduced more.
func (k *ExtendedKey) Derive(index uint32) (*ExtendedKey, error) {
	// There are four scenarios that could happen here:
	// 1) Private extended key -> Hardened child private extended key
	// 2) Private extended key -> Non-hardened child private extended key
	// 3) Public extended key -> Non-hardened child public extended key
	// 4) Public extended key -> Hardened child public extended key (INVALID!)

	isChildHardened := index >= HardenedKeyStart

	// The data used to derive the child key depends on whether or not the
	// child is hardened.
	//
	// For hardened children:
	//   G1: 0x01 || ser256(parentKey) || ser32(i)
	//   G2: 0x00 || ser256(parentKey) || ser32(i)
	//
	// For normal children:
	//   G1: serG1(parentPubKey) || ser32(i)
	//   G2: serG2(parentPubKey) || ser32(i)
	//
	data := make([]byte, 0, 100)
	if isChildHardened {
		// Case #1 and #4.
		if k.isPrivate {
			// Case #1
			//
			// When the child is a hardened child, the key is known to be a
			// private key.
			// Pad it with a leading zero as required by [BIP32] for deriving the child.
			if len(k.key) != 32 {
				return nil, ErrInvalidKeyData
			}
			if k.pubOnG1 {
				data = append(data, 0x01)
			} else {
				data = append(data, 0x00)
			}
			data = append(data, k.key...)
		} else {
			// Case #4
			//
			// A hardened child extended key may not be created from a public
			// extended key.
			return nil, ErrDeriveHardFromPublic
		}
	} else {
		// Case #2 or #3.
		//
		// This is either a public or private extended key, but in
		// either case, the data which is used to derive the child key
		// starts with the BLS public key bytes.
		data = append(data, k.pubKeyBytes()...)
		if k.pubOnG1 && len(data) != 48 {
			return nil, ErrInvalidKeyData
		}
		if !k.pubOnG1 && len(data) != 96 {
			return nil, ErrInvalidKeyData
		}
	}
	indexData := make([]byte, 4)
	binary.BigEndian.PutUint32(indexData, index)
	data = append(data, indexData...)

	// The order is same for all three groups (g1, g2, and gt).
	gt := bls12381.NewGT()
	order := gt.Q()

	var childChainCode, il []byte

	for {
		// Take the HMAC-SHA512 of the current key's chain code and the derived
		// data:
		//   I = HMAC-SHA512(Key = chainCode, Data = data)
		hmac512 := hmac.New(sha512.New, k.chainCode)
		_, _ = hmac512.Write(data)
		ilr := hmac512.Sum(nil)

		// Split "I" into two 32-byte sequences Il and Ir where:
		//   Il = intermediate key used to derive the child
		//   Ir = child chain code
		il = ilr[:len(ilr)/2]
		childChainCode = ilr[len(ilr)/2:]

		// If Il greater or equal to the order of the group, or it is zero,
		// generate a new "I" with data equals to 0x01 || Ir || ser32(i)
		ilNum := big.Int{}
		ilNum.SetBytes(il)
		if ilNum.Cmp(order) == -1 && ilNum.Cmp(big.NewInt(0)) != 0 {
			break
		}

		data = []byte{0x01}
		data = append(data, childChainCode...)
		data = append(data, indexData...)
	}

	ilFr := new(bls12381.Fr)
	ilFr.FromBytes(il)

	var childKey []byte
	if k.isPrivate {
		// Case #1 or #2.
		// Add the parent private key to the intermediate private key to
		// derive the final child key.
		//
		// childKey = parse256(Il) + parentKey

		keyNum := new(bls12381.Fr)
		keyNum.FromBytes(k.key)

		childKeyNum := bls12381.NewFr()
		childKeyNum.Add(keyNum, ilFr)

		childKey = childKeyNum.ToBytes()
	} else {
		// Case #3.
		// Calculate the corresponding intermediate public key for the
		// intermediate private key.
		//
		if k.pubOnG1 {
			// Public key is in G1 subgroup
			//
			// childKey = pointG1(parse256(Il)) + parentKey
			g1 := bls12381.NewG1()

			ilPoint := new(bls12381.PointG1)
			g1.MulScalar(ilPoint, g1.One(), ilFr)

			pubKey, err := g1.FromCompressed(k.key)
			if err != nil {
				return nil, err
			}
			childPubKey := new(bls12381.PointG1)
			g1.Add(childPubKey, pubKey, ilPoint)

			childKey = g1.ToCompressed(childPubKey)
		} else {
			// Public key is in G2 subgroup
			//
			// childKey = pointG2(parse256(Il)) + parentKey
			g2 := bls12381.NewG2()

			ilPoint := new(bls12381.PointG2)
			g2.MulScalar(ilPoint, g2.One(), ilFr)

			pubKey, err := g2.FromCompressed(k.key)
			if err != nil {
				return nil, err
			}
			childPubKey := new(bls12381.PointG2)
			g2.Add(childPubKey, pubKey, ilPoint)

			childKey = g2.ToCompressed(childPubKey)
		}
	}

	newPath := make([]uint32, 0, len(k.path)+1)
	newPath = append(newPath, k.path...)
	newPath = append(newPath, index)

	return newExtendedKey(childKey, childChainCode,
		newPath, k.isPrivate, k.pubOnG1), nil
}

// Path returns the path of derived key.
//
// Path with values between 0 and 2^31-1 are normal child keys,
// and those values between 2^31 and 2^32-1 are hardened keys.
func (k *ExtendedKey) Path() []uint32 {
	return k.path
}

// RawPrivateKey returns the raw bytes of the private key.
// As you might imagine this is only possible if the extended key is a private
// extended key (as determined by the IsPrivate function).  The ErrNotPrivExtKey
// error will be returned if this function is called on a public extended key.
func (k *ExtendedKey) RawPrivateKey() ([]byte, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	return k.key, nil
}

// RawPublicKey returns the raw bytes of the public key.
func (k *ExtendedKey) RawPublicKey() []byte {
	return k.pubKeyBytes()
}

// Neuter returns a new extended public key from this extended private key.  The
// same extended key will be returned unaltered if it is already an extended
// public key.
//
// As the name implies, an extended public key does not have access to the
// private key, so it is not capable of signing transactions or deriving
// child extended private keys.  However, it is capable of deriving further
// child extended public keys.
func (k *ExtendedKey) Neuter() *ExtendedKey {
	// Already an extended public key.
	if !k.isPrivate {
		return k
	}

	// Convert it to an extended public key.  The key for the new extended
	// key will simply be the pubkey of the current extended private key.
	//
	// This is the function N((k,c)) -> (K, c) from [BIP32].
	return newExtendedKey(k.pubKeyBytes(), k.chainCode,
		k.path, false, k.pubOnG1)
}

// String returns the extended key as a bech32-encoded string.
func (k *ExtendedKey) String() string {
	//
	// The serialized format is structured as follows:
	// +-------+---------+------------+-------+------------+----------+
	// | Depth | Path    | Chain code | G1/G2 | Key length | Key data |
	// +-------+---------+------------+-------+------------+----------+
	// | 1     | depth*4 | 32         | 1     | 1          | 32/48/96 |
	// +-------+---------+------------+-------+------------+----------+
	//
	// Description:
	// - Depth: 1 byte representing the depth of derivation path.
	// - Path: serialized BIP-32 path; each entry is encoded as 32-bit unsigned integer, least significant byte first
	// - Chain code: 32 bytes chain code
	// - G1 or G2: 1 byte to specify the group.
	// - Key length: 1 byte representing the length of the key data.
	// - Key data: Can be 32, 48, or 96 bytes.
	//

	w := bytes.NewBuffer(make([]byte, 0))
	err := encoding.WriteElement(w, byte(len(k.path)))
	if err != nil {
		return err.Error()
	}

	for _, p := range k.path {
		err := encoding.WriteElement(w, p)
		if err != nil {
			return err.Error()
		}
	}
	err = encoding.WriteVarBytes(w, k.chainCode)
	if err != nil {
		return err.Error()
	}

	err = encoding.WriteElement(w, k.pubOnG1)
	if err != nil {
		return err.Error()
	}

	err = encoding.WriteVarBytes(w, k.key)
	if err != nil {
		return err.Error()
	}

	hrp := crypto.XPublicKeyHRP
	if k.isPrivate {
		hrp = crypto.XPrivateKeyHRP
	}

	str, err := bech32m.EncodeFromBase256WithType(hrp, crypto.SignatureTypeBLS, w.Bytes())
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

	if typ != crypto.SignatureTypeBLS {
		return nil, ErrInvalidKeyData
	}

	r := bytes.NewReader(data)
	depth := uint8(0)
	err = encoding.ReadElement(r, &depth)
	if err != nil {
		return nil, err
	}

	path := make([]uint32, depth)
	for i := byte(0); i < depth; i++ {
		err := encoding.ReadElement(r, &path[i])
		if err != nil {
			return nil, err
		}
	}

	chainCode, err := encoding.ReadVarBytes(r)
	if err != nil {
		return nil, err
	}

	var pubOnG1 bool
	err = encoding.ReadElement(r, &pubOnG1)
	if err != nil {
		return nil, err
	}

	key, err := encoding.ReadVarBytes(r)
	if err != nil {
		return nil, err
	}

	isPrivate := true
	if hrp == crypto.XPublicKeyHRP {
		isPrivate = false
	}

	return newExtendedKey(key, chainCode, path, isPrivate, pubOnG1), nil
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
func NewMaster(seed []byte, pubOnG1 bool) (*ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// masterKey is the master key used along with a random seed used to generate
	// the master node in the hierarchical tree.
	masterKey := []byte("BLS12381 seed")

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "BLS12381-HD seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	_, _ = hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master IKM
	//   Ir = master chain code
	ikm := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	// Using BLS KeyGen to generate the master private key from the IKM.
	privKey, err := bls.KeyGen(ikm, nil)
	if err != nil {
		return nil, err
	}

	return newExtendedKey(privKey.Bytes(), chainCode, []uint32{}, true, pubOnG1), nil
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
