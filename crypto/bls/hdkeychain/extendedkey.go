package hdkeychain

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"strings"

	bls12381 "github.com/kilic/bls12-381"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
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
	path      Path
	isPrivate bool
	pubOnG1   bool
}

// newExtendedKey returns a new instance of an extended key with the given
// fields. No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func newExtendedKey(key, chainCode []byte, path Path, isPrivate bool, pubOnG1 bool) *ExtendedKey {
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

// Derive returns a derived child extended key from this master key at the
// given path.
func (k *ExtendedKey) DerivePath(path Path) (*ExtendedKey, error) {
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
	//   data (36 bytes) = parent_private_key (32 bytes)  || index (4 bytes)
	//
	// For normal children:
	//   data (52 bytes)  = parent_public_key_g1 (48 bytes)  || index (4 bytes)
	//   data (100 bytes) = parent_public_key_g2 (96 bytes)  || index (4 bytes)
	data := make([]byte, 0, 100)
	if isChildHardened {
		// Case #1 and #4.
		if k.isPrivate {
			// Case #1
			//
			// When the child is a hardened child, the key is known to be a
			// private key.
			data = append(data, k.key...)
			if len(data) != 32 {
				return nil, ErrInvalidKeyData
			}
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
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, index)
	data = append(data, bs...)

	// Take the HMAC-SHA512 of the current key's chain code and the derived
	// data:
	//   I = HMAC-SHA512(Key = chainCode, Data = data)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	_, _ = hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = intermediate key used to derive the child
	//   Ir = child chain code
	il := ilr[:len(ilr)/2]
	childChainCode := ilr[len(ilr)/2:]

	ilNum := new(bls12381.Fr)
	ilNum.FromBytes(il)

	var childKey []byte
	if k.isPrivate {
		// Case #1 or #2.
		// Add the parent private key to the intermediate private key to
		// derive the final child key.
		//
		// childKey = parse256(Il) + parenKey

		keyNum := new(bls12381.Fr)
		keyNum.FromBytes(k.key)

		childKeyNum := bls12381.NewFr()
		childKeyNum.Add(keyNum, ilNum)

		if childKeyNum.IsZero() {
			return nil, ErrInvalidKeyData
		}
		childKey = childKeyNum.ToBytes()
	} else {
		// Case #3.
		// Calculate the corresponding intermediate public key for the
		// intermediate private key

		if k.pubOnG1 {
			g1 := bls12381.NewG1()

			// Public key is in G1 subgroup
			ilPoint := new(bls12381.PointG1)
			g1.MulScalar(ilPoint, g1.One(), ilNum)

			pubKey, err := g1.FromCompressed(k.key)
			if err != nil {
				return nil, err
			}
			childPubKey := new(bls12381.PointG1)
			g1.Add(childPubKey, pubKey, ilPoint)

			if g1.IsZero(childPubKey) {
				return nil, ErrInvalidKeyData
			}
			childKey = g1.ToCompressed(childPubKey)
		} else {
			g2 := bls12381.NewG2()

			// Public key is in G2 subgroup
			ilPoint := new(bls12381.PointG2)
			g2.MulScalar(ilPoint, g2.One(), ilNum)

			pubKey, err := g2.FromCompressed(k.key)
			if err != nil {
				return nil, err
			}
			childPubKey := new(bls12381.PointG2)
			g2.Add(childPubKey, pubKey, ilPoint)

			if g2.IsZero(childPubKey) {
				return nil, ErrInvalidKeyData
			}
			childKey = g2.ToCompressed(childPubKey)
		}
	}

	newPath := make(Path, 0, len(k.path)+1)
	copy(newPath, k.path)
	newPath = append(k.path, index)
	return newExtendedKey(childKey, childChainCode,
		newPath, k.isPrivate, k.pubOnG1), nil
}

// Path returns the path of derived key.
//
// Path with values between 0 and 2^31-1 are normal child keys,
// and those values between 2^31 and 2^32-1 are hardened keys.
func (k *ExtendedKey) Path() Path {
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
	// The serialized format is:
	// path (variant) || chain code (32) || pubkey length (1 byte) || key data (32, 48 or 96)
	w := bytes.NewBuffer(make([]byte, 0))
	err := encoding.WriteElement(w, byte(len(k.path)))
	util.ExitOnErr(err)

	for _, p := range k.path {
		err := encoding.WriteVarInt(w, uint64(p))
		util.ExitOnErr(err)
	}
	err = encoding.WriteElement(w, k.chainCode)
	util.ExitOnErr(err)

	hrp := crypto.XPublicKeyHRP
	if k.isPrivate {
		hrp = crypto.XPrivateKeyHRP

		pubKeyLen := byte(96)
		if k.pubOnG1 {
			pubKeyLen = 48
		}

		err := encoding.WriteElement(w, pubKeyLen)
		util.ExitOnErr(err)

		err = encoding.WriteElement(w, k.key)
		util.ExitOnErr(err)
	} else {
		err = encoding.WriteVarBytes(w, k.key)
		util.ExitOnErr(err)
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
func NewKeyFromString(key string) (*ExtendedKey, error) {
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(strings.ToLower(key))
	if err != nil {
		return nil, err
	}

	if typ != crypto.SignatureTypeBLS {
		return nil, ErrInvalidKeyData
	}

	r := bytes.NewReader(data)
	path := Path{}
	pathLen := byte(0)
	err = encoding.ReadElement(r, &pathLen)
	if err != nil {
		return nil, err
	}
	for i := byte(0); i < pathLen; i++ {
		p, err := encoding.ReadVarInt(r)
		if err != nil {
			return nil, err
		}
		path = append(path, uint32(p))
	}

	switch hrp {
	case crypto.XPrivateKeyHRP:
		if r.Len() != 65 {
			return nil, ErrInvalidKeyData
		}
		chainCode := make([]byte, 32)
		key := make([]byte, 32)

		err := encoding.ReadElement(r, chainCode)
		util.ExitOnErr(err)

		pubKeyLen, _ := encoding.ReadVarInt(r)

		err = encoding.ReadElement(r, key)
		util.ExitOnErr(err)

		pubOnG1 := pubKeyLen == 48
		return newExtendedKey(key[:], chainCode[:], path, true, pubOnG1), nil

	case crypto.XPublicKeyHRP:
		if r.Len() != 64 && r.Len() != 81 && r.Len() != 129 {
			return nil, ErrInvalidKeyData
		}
		chainCode := make([]byte, 32)

		err := encoding.ReadElement(r, chainCode)
		util.ExitOnErr(err)

		key, err := encoding.ReadVarBytes(r)
		util.ExitOnErr(err)

		pubOnG1 := len(key) == 48
		return newExtendedKey(key[:], chainCode[:], path, false, pubOnG1), nil

	default:
		return nil, ErrInvalidKeyData
	}
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
	masterKey := []byte("BLS12381-HD-KEYCHAIN")

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "BLS12381-HD-KEYCHAIN", Data = S)
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

	return newExtendedKey(privKey.Bytes(), chainCode, Path{}, true, pubOnG1), nil
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
