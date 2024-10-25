package bls

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/bech32m"
	"golang.org/x/crypto/hkdf"
)

var _ crypto.PrivateKey = &PrivateKey{}

const PrivateKeySize = 32

type PrivateKey struct {
	scalar *big.Int
}

// PrivateKeyFromString decodes the input string and returns the PrivateKey
// if the string is a valid bech32m encoding of a BLS public key.
func PrivateKeyFromString(text string) (*PrivateKey, error) {
	// Decode the bech32m encoded private key.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return nil, err
	}

	// Check if hrp is valid
	if hrp != crypto.PrivateKeyHRP {
		return nil, crypto.InvalidHRPError(hrp)
	}

	if typ != crypto.SignatureTypeBLS {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return PrivateKeyFromBytes(data)
}

// KeyGen generates a private key deterministically from a secret octet string
// IKM and an optional octet string keyInfo.
// Based on https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-2.3
func KeyGen(ikm, keyInfo []byte) (*PrivateKey, error) {
	// L is `ceil((3 * ceil(log2(r))) / 16) = 48`,
	//    where `r` is the order of the BLS 12-381 curve
	//    r:  0x73eda753 299d7d48 3339d808 09a1d805 53bda402 fffe5bfe ffffffff 00000001
	// 	  https://datatracker.ietf.org/doc/html/draft-yonezawa-pairing-friendly-curves-02#section-4.2.2
	//

	if len(ikm) < 32 {
		return nil, fmt.Errorf("ikm is too short")
	}

	secret := make([]byte, 0, len(ikm)+1)
	secret = append(secret, ikm...)
	secret = append(secret, util.I2OSP(big.NewInt(0), 1)...)

	l := int64(48)
	pseudoRandomKey := make([]byte, 0, len(keyInfo)+2)
	pseudoRandomKey = append(pseudoRandomKey, keyInfo...)
	pseudoRandomKey = append(pseudoRandomKey, util.I2OSP(big.NewInt(l), 2)...)

	salt := []byte("BLS-SIG-KEYGEN-SALT-")
	num := big.NewInt(0)
	for num.Sign() == 0 {
		hash := sha256.Sum256(salt)
		salt = hash[:]

		okm := make([]byte, l)
		prk := hkdf.Extract(sha256.New, secret, salt)
		reader := hkdf.Expand(sha256.New, prk, pseudoRandomKey)
		_, _ = reader.Read(okm)

		order := fr.Modulus()
		num = new(big.Int).Mod(util.OS2IP(okm), order)
	}

	sk := make([]byte, 32)
	num.FillBytes(sk)

	return PrivateKeyFromBytes(sk)
}

// PrivateKeyFromBytes constructs a BLS private key from the raw bytes.
func PrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return nil, crypto.InvalidLengthError(len(data))
	}

	scalar := new(big.Int)
	scalar = scalar.SetBytes(data)

	return &PrivateKey{scalar: scalar}, nil
}

// String returns a human-readable string for the BLS private key.
func (prv *PrivateKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PrivateKeyHRP,
		crypto.SignatureTypeBLS,
		prv.Bytes())

	return strings.ToUpper(str)
}

// Bytes return the raw bytes of the private key.
func (prv *PrivateKey) Bytes() []byte {
	data := prv.scalar.Bytes()
	data = util.Extend(data, PrivateKeySize)

	return data
}

// Sign calculates the signature from the private key and given message.
// It's defined in section 2.6 of the spec: CoreSign.
func (prv *PrivateKey) Sign(msg []byte) crypto.Signature {
	return prv.SignNative(msg)
}

func (prv *PrivateKey) SignNative(msg []byte) *Signature {
	qAffine, err := bls12381.HashToG1(msg, dst)
	if err != nil {
		panic(err)
	}
	qJac := new(bls12381.G1Jac).FromAffine(&qAffine)
	sigJac := qJac.ScalarMultiplication(qJac, prv.scalar)
	sigAffine := new(bls12381.G1Affine).FromJacobian(sigJac)
	data := sigAffine.Bytes()

	return &Signature{
		data:    data[:],
		pointG1: sigAffine,
	}
}

func (prv *PrivateKey) PublicKeyNative() *PublicKey {
	pkJac := new(bls12381.G2Jac).ScalarMultiplication(&gen2Jac, prv.scalar)
	pkAffine := new(bls12381.G2Affine).FromJacobian(pkJac)
	data := pkAffine.Bytes()

	return &PublicKey{
		data:    data[:],
		pointG2: pkAffine,
	}
}

func (prv *PrivateKey) PublicKey() crypto.PublicKey {
	return prv.PublicKeyNative()
}

func (prv *PrivateKey) EqualsTo(x crypto.PrivateKey) bool {
	xBLS, ok := x.(*PrivateKey)
	if !ok {
		return false
	}

	return prv.scalar.Cmp(xBLS.scalar) == 0
}
