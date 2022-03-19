package bls

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
	"golang.org/x/crypto/hkdf"
)

const PrivateKeySize = 32

type PrivateKey struct {
	data privateKeyData
}

type privateKeyData struct {
	SecretKey *bls.SecretKey
}

func PrivateKeyFromString(text string) (*PrivateKey, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromRawBytes(data)
}

// PrivateKeyFromSeed generates a private key deterministically from
// a secret octet string IKM.
// Based on https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-02#section-2.3
func PrivateKeyFromSeed(ikm []byte) (*PrivateKey, error) {
	// L is `ceil((3 * ceil(log2(r))) / 16) = 48`,
	//    where `r` is the order of the BLS 12-381 curve
	//    r: 0x73eda753 299d7d48 3339d808 09a1d805 53bda402 fffe5bfe ffffffff 00000001
	// 	  https://datatracker.ietf.org/doc/html/draft-yonezawa-pairing-friendly-curves-02#section-4.2.2
	//

	if len(ikm) < 32 {
		return nil, fmt.Errorf("ikm is too short")
	}

	salt := []byte("BLS-SIG-KEYGEN-SALT-")
	L := 48
	okm := make([]byte, L)
	_, _ = hkdf.New(sha256.New, append(ikm, 0), salt, []byte{0, byte(L)}).Read(okm)

	// OS2IP: https://datatracker.ietf.org/doc/html/rfc8017#section-4.2
	// OS2IP converts an octet string to a nonnegative integer.

	r, _ := new(big.Int).SetString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001", 16)
	x := new(big.Int).Mod(new(big.Int).SetBytes(okm), r)
	buf := [32]byte{}
	x.FillBytes(buf[:])

	return PrivateKeyFromRawBytes(buf[:])
}

func PrivateKeyFromRawBytes(data []byte) (*PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return nil, fmt.Errorf("invalid private key")
	}
	sc := new(bls.SecretKey)
	if err := sc.Deserialize(data); err != nil {
		return nil, err
	}

	var prv PrivateKey
	prv.data.SecretKey = sc

	return &prv, nil
}

func (prv PrivateKey) RawBytes() []byte {
	if prv.data.SecretKey == nil {
		return nil
	}
	return prv.data.SecretKey.Serialize()
}

func (prv PrivateKey) String() string {
	if prv.data.SecretKey == nil {
		return ""
	}
	return prv.data.SecretKey.SerializeToHexStr()
}

func (prv *PrivateKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(prv.String())
}

func (prv *PrivateKey) MarshalCBOR() ([]byte, error) {
	if prv.data.SecretKey == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return cbor.Marshal(prv.RawBytes())
}

func (prv *PrivateKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	p, err := PrivateKeyFromRawBytes(data)
	if err != nil {
		return err
	}

	*prv = *p
	return nil
}

func (prv *PrivateKey) SanityCheck() error {
	if prv.data.SecretKey.IsZero() {
		return fmt.Errorf("private key is zero")
	}
	return nil
}

func (prv *PrivateKey) Sign(msg []byte) crypto.Signature {
	sig := new(Signature)
	sig.data.Signature = prv.data.SecretKey.SignByte(msg)

	return sig
}

func (prv *PrivateKey) PublicKey() crypto.PublicKey {
	pb := new(PublicKey)
	pb.data.PublicKey = prv.data.SecretKey.GetPublicKey()

	return pb
}

func (prv *PrivateKey) EqualsTo(right crypto.PrivateKey) bool {
	return prv.data.SecretKey.IsEqual(right.(*PrivateKey).data.SecretKey)
}
