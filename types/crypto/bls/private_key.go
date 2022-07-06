package bls

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/bech32m"
	"github.com/zarbchain/zarb-go/util/errors"
	"golang.org/x/crypto/hkdf"
)

const (
	PrivateKeySize = 32
	hrpPrivateKey  = "secret"
)

type PrivateKey struct {
	secretKey bls.SecretKey
}

// PrivateKeyFromString decodes the string encoding of a BLS private key
// and returns the private key if text is a valid encoding for BLS private key.
func PrivateKeyFromString(text string) (*PrivateKey, error) {
	// Decode the bech32m encoded private key.
	hrp, data, err := bech32m.Decode(text)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey, err.Error())
	}

	// Check if hrp is valid
	if hrp != hrpPrivateKey {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey, "invalid hrp: %v", hrp)
	}

	// The first byte of the decoded private key is the signature type, it must
	// exist.
	if len(data) < 1 {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey, "no private key type")
	}

	// ...and should be 1 for BLS signature.
	sigType := data[0]
	if sigType != crypto.SignatureTypeBLS {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey, "invalid private key type: %v", sigType)
	}

	// The remaining characters of the private key returned are grouped into
	// words of 5 bits. In order to restore the original program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32m.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey, err.Error())
	}

	return PrivateKeyFromBytes(regrouped)
}

// PrivateKeyFromSeed generates a private key deterministically from
// a secret octet string IKM.
// Based on https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-2.3
func PrivateKeyFromSeed(ikm []byte, keyInfo []byte) (*PrivateKey, error) {
	// L is `ceil((3 * ceil(log2(r))) / 16) = 48`,
	//    where `r` is the order of the BLS 12-381 curve
	//    r: 0x73eda753 299d7d48 3339d808 09a1d805 53bda402 fffe5bfe ffffffff 00000001
	// 	  https://datatracker.ietf.org/doc/html/draft-yonezawa-pairing-friendly-curves-02#section-4.2.2
	//

	if len(ikm) < 32 {
		return nil, fmt.Errorf("ikm is too short")
	}

	salt := []byte("BLS-SIG-KEYGEN-SALT-")
	x := big.NewInt(0)
	for x.Sign() == 0 {
		h := sha256.Sum256(salt)
		salt = h[:]
		L := int64(48)
		okm := make([]byte, L)
		prk := hkdf.Extract(sha256.New, append(ikm, util.IS2OP(big.NewInt(0), 1)...), salt[:])
		reader := hkdf.Expand(sha256.New, prk, append(keyInfo, util.IS2OP(big.NewInt(L), 2)...))
		_, _ = reader.Read(okm)

		r, _ := new(big.Int).SetString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001", 16)
		x = new(big.Int).Mod(util.OS2IP(okm), r)
	}

	sk := make([]byte, 32)
	x.FillBytes(sk)
	return PrivateKeyFromBytes(sk)
}

// PrivateKeyFromBytes constructs a BLS private key from the raw bytes.
// This method in unexported and should not be called from the outside.
func PrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey,
			"private key should be %d bytes, but it is %v bytes", PrivateKeySize, len(data))
	}
	sc := new(bls.SecretKey)
	if err := sc.Deserialize(data); err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPrivateKey, err.Error())
	}

	var prv PrivateKey
	prv.secretKey = *sc

	return &prv, nil
}

// String returns a human-readable string for the BLS private key.
func (prv PrivateKey) String() string {
	data := prv.secretKey.Serialize()

	// Group the private key bytes into 5 bit groups, as this is what is used to
	// encode each character in the private key string.
	converted, err := bech32m.ConvertBits(data, 8, 5, true)
	if err != nil {
		panic(err.Error())
	}

	// Concatenate the type of the private key which is 1 for BLS and program,
	// and encode the resulting bytes using bech32m encoding.
	combined := make([]byte, len(converted)+1)
	combined[0] = crypto.SignatureTypeBLS
	copy(combined[1:], converted)
	str, err := bech32m.Encode(hrpPrivateKey, combined)
	if err != nil {
		panic(err.Error())
	}

	return strings.ToUpper(str)
}

// Bytes return the raw bytes of the private key.
func (prv PrivateKey) Bytes() []byte {
	return prv.secretKey.Serialize()
}

func (prv *PrivateKey) SanityCheck() error {
	if prv.secretKey.IsZero() {
		return fmt.Errorf("private key is zero")
	}
	return nil
}

func (prv *PrivateKey) Sign(msg []byte) crypto.Signature {
	sig := new(Signature)
	sig.signature = *prv.secretKey.SignByte(msg)

	return sig
}

func (prv *PrivateKey) PublicKey() crypto.PublicKey {
	pub := prv.secretKey.GetPublicKey()
	return &PublicKey{
		publicKey: *pub,
	}
}

func (prv *PrivateKey) EqualsTo(right crypto.PrivateKey) bool {
	return prv.secretKey.IsEqual(&right.(*PrivateKey).secretKey)
}
