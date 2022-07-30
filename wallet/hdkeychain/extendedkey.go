package hdkeychain

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"

	herumi "github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
)

// Public key Generator for BLS12-381 curve used in Zarb
var g2Gen herumi.G2

func init() {
	err := g2Gen.SetString(`
	1
	24aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8
	13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e
	ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801
	606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be`, 16)

	if err != nil {
		panic(err)
	}
}

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

	// maxUint8 is the max positive integer which can be serialized in a uint8
	maxUint8 = 1<<8 - 1
)

var (
	// ErrDeriveHardFromPublic describes an error in which the caller
	// attempted to derive a hardened extended key from a public key.
	ErrDeriveHardFromPublic = errors.New("cannot derive a hardened key " +
		"from a public key")

	// ErrDeriveBeyondMaxDepth describes an error in which the caller
	// has attempted to derive more than 255 keys from a root key.
	ErrDeriveBeyondMaxDepth = errors.New("cannot derive a key with more than " +
		"255 indices in its path")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New("unable to create private keys from a " +
		"public extended key")

	// ErrInvalidChild describes an error in which the child at a specific
	// index is invalid due to the derived key falling outside of the valid
	// range for BLS private keys.  This error indicates the caller
	// should simply ignore the invalid child extended key at this index and
	// increment to the next index.
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// BLS private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)

	//
	ErrInvalidKey = errors.New("key is invalid")
)

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key.
type ExtendedKey struct {
	key       []byte // This will be the pubkey for extended pub keys
	pubKey    []byte // This will only be set for extended priv keys
	chainCode []byte
	depth     uint8
	childNum  uint32
	isPrivate bool
}

// NewExtendedKey returns a new instance of an extended key with the given
// fields. No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func NewExtendedKey(key, chainCode []byte, depth uint8, childNum uint32,
	isPrivate bool) *ExtendedKey {
	// NOTE: The pubKey field is intentionally left nil so it is only
	// computed and memoized as required.
	return &ExtendedKey{
		key:       key,
		pubKey:    nil,
		chainCode: chainCode,
		depth:     depth,
		childNum:  childNum,
		isPrivate: isPrivate,
	}
}

// pubKeyBytes returns bytes for the serialized compressed public key associated
// with this extended key in an efficient manner including memoization as
// necessary.
//
// When the extended key is already a public key, the key is simply returned as
// is since it's already in the correct form.  However, when the extended key is
// a private key, the public key will be calculated and memoized so future
// accesses can simply return the cached result.
func (k *ExtendedKey) pubKeyBytes() []byte {
	// Just return the key if it's already an extended public key.
	if !k.isPrivate {
		return k.key
	}

	// This is a private extended key, so calculate and memoize the public
	// key if needed.
	if len(k.pubKey) == 0 {
		privKey := new(herumi.Fr)
		pubPoint := new(herumi.G2)
		err := privKey.Deserialize(k.key)
		if err != nil {
			panic(err)
		}

		herumi.G2Mul(pubPoint, &g2Gen, privKey)
		k.pubKey = pubPoint.Serialize()
	}

	return k.pubKey
}

// IsPrivate returns whether or not the extended key is a private extended key.
//
// A private extended key can be used to derive both hardened and non-hardened
// child private and public extended keys. A public extended key can only be
// used to derive non-hardened child public extended keys.
func (k *ExtendedKey) IsPrivate() bool {
	return k.isPrivate
}

// Derive returns a derived child extended key at the given index.
//
// When this extended key is a private extended key (as determined by the IsPrivate
// function), a private extended key will be derived. Otherwise, the derived
// extended key will be also be a public extended key.
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
func (k *ExtendedKey) Derive(i uint32) (*ExtendedKey, error) {
	// Prevent derivation of children beyond the max allowed depth.
	if k.depth == maxUint8 {
		return nil, ErrDeriveBeyondMaxDepth
	}

	// There are four scenarios that could happen here:
	// 1) Private extended key -> Hardened child private extended key
	// 2) Private extended key -> Non-hardened child private extended key
	// 3) Public extended key -> Non-hardened child public extended key
	// 4) Public extended key -> Hardened child public extended key (INVALID!)

	isChildHardened := i >= HardenedKeyStart

	if k.depth > 0 {
		isParentHardened := k.childNum >= HardenedKeyStart
		// In case we try to extend a non-hardened child key from hardened or vice versa.
		if isChildHardened && !isParentHardened ||
			!isChildHardened && isParentHardened {
			return nil, ErrInvalidChild
		}
	}

	// The data used to derive the child key depends on whether or not the
	// child is hardened.
	//
	// For hardened children:
	//   data (36 bytes) = parent_private_key (32 bytes)  || index (4 bytes)
	//
	// For normal children:
	//   data (100 bytes) = parent_public_key (96 bytes)  || index (4 bytes)
	data := make([]byte, 0, 100)
	if isChildHardened {
		// Case #1.
		// When the child is a hardened child, the key is known to be a
		// private key.
		data = append(data, k.key...)
		if len(data) != 32 {
			return nil, ErrInvalidKey
		}
	} else {
		// Case #2 or #3.
		// This is either a public or private extended key, but in
		// either case, the data which is used to derive the child key
		// starts with the BLS public key bytes.
		data = append(data, k.pubKeyBytes()...)
		if len(data) != 96 {
			return nil, ErrInvalidKey
		}
	}
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	data = append(data, bs...)

	// Take the HMAC-SHA512 of the current key's chain code and the derived
	// data:
	//   I = HMAC-SHA512(Key = chainCode, Data = data)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	_, _ = hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = intermediate key used to derive the child private key
	//   Ir = child chain code
	ikm := ilr[:len(ilr)/2]
	childChainCode := ilr[len(ilr)/2:]

	derivedPrivKey, err := bls.PrivateKeyFromSeed(ikm, nil)
	if err != nil {
		return nil, err
	}

	var childKey []byte
	if isChildHardened {
		if k.isPrivate {
			// Case #1
			// Corresponding private key is same as intermediate private key

			childKey = derivedPrivKey.Bytes()
		} else {
			// Case #4
			// A hardened child extended key may not be created from a public
			// extended key.
			return nil, ErrDeriveHardFromPublic
		}
	} else {
		if k.isPrivate {
			// Case #2
			// Calculate the corresponding private key for the
			// intermediate private key

			scalar1 := new(herumi.Fr)
			scalar2 := new(herumi.Fr)
			scalarAdd := new(herumi.Fr)

			if err := scalar1.Deserialize(k.key); err != nil {
				return nil, ErrInvalidKey
			}
			if err := scalar2.Deserialize(derivedPrivKey.Bytes()); err != nil {
				// impossible
				return nil, ErrInvalidKey
			}

			herumi.FrAdd(scalarAdd, scalar1, scalar2)

			childKey = scalarAdd.Serialize()
		} else {
			// Case #3.
			// Calculate the corresponding public key for the
			// intermediate private key

			ilScalar := new(herumi.Fr)
			err := ilScalar.Deserialize(derivedPrivKey.Bytes())
			if err != nil {
				// impossible
				return nil, err
			}

			ilPoint := new(herumi.G2)
			herumi.G2Mul(ilPoint, &g2Gen, ilScalar)

			pubKey := new(herumi.G2)
			err = pubKey.Deserialize(k.key)
			if err != nil {
				return nil, err
			}
			childPubKey := new(herumi.G2)
			herumi.G2Add(childPubKey, pubKey, ilPoint)

			childKey = childPubKey.Serialize()
		}
	}

	return NewExtendedKey(childKey, childChainCode,
		k.depth+1, i, k.isPrivate), nil
}

// ChildNum returns the index at which the child extended key was derived.
//
// Extended keys with ChildNum value between 0 and 2^31-1 are normal child
// keys, and those with a value between 2^31 and 2^32-1 are hardened keys.
func (k *ExtendedKey) ChildIndex() uint32 {
	return k.childNum
}

// Neuter returns a new extended public key from this extended private key.  The
// same extended key will be returned unaltered if it is already an extended
// public key.
//
// As the name implies, an extended public key does not have access to the
// private key, so it is not capable of signing transactions or deriving
// child extended private keys.  However, it is capable of deriving further
// child extended public keys.
func (k *ExtendedKey) Neuter() (*ExtendedKey, error) {
	// Already an extended public key.
	if !k.isPrivate {
		return k, nil
	}

	// Convert it to an extended public key.  The key for the new extended
	// key will simply be the pubkey of the current extended private key.
	//
	// This is the function N((k,c)) -> (K, c) from [BIP32].
	return NewExtendedKey(k.pubKeyBytes(), k.chainCode,
		k.depth, k.childNum, false), nil
}

// BLSPublicKey converts the extended key to a BLS public key and returns it.
func (k *ExtendedKey) BLSPublicKey() (*bls.PublicKey, error) {
	return bls.PublicKeyFromBytes(k.pubKeyBytes())
}

// BLSPrivateKey converts the extended key to a BLS private key and returns it.
// As you might imagine this is only possible if the extended key is a private
// extended key (as determined by the IsPrivate function).  The ErrNotPrivExtKey
// error will be returned if this function is called on a public extended key.
func (k *ExtendedKey) BLSPrivateKey() (*bls.PrivateKey, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	return bls.PrivateKeyFromBytes(k.key)
}

// Address converts the extended key to address
func (k *ExtendedKey) Address() crypto.Address {
	pub, _ := k.BLSPublicKey()
	return pub.Address()
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
//
func NewMaster(seed []byte) (*ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// masterKey is the master key used along with a random seed used to generate
	// the master node in the hierarchical tree.
	var masterKey = []byte("Zarb seed")

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Zarb seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	_, _ = hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master ikm
	//   Ir = master chain code
	ikm := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	privKey, err := bls.PrivateKeyFromSeed(ikm, nil)
	if err != nil {
		return nil, err
	}

	return NewExtendedKey(privKey.Bytes(), chainCode, 0, 0, true), nil
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
