package vault

import (
	"strings"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util/bip39"
)

// GenerateMnemonic generates a new mnemonic (seed phrase) based on BIP-39
// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
func GenerateMnemonic(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

// CheckMnemonic validates a mnemonic (seed phrase) based on BIP-39.
func CheckMnemonic(mnemonic string) error {
	_, err := bip39.EntropyFromMnemonic(mnemonic)

	return err
}

func PrivateKeyFromString(str string) (crypto.PrivateKey, error) {
	str = strings.TrimSpace(strings.ToLower(str))

	maybeBLSPrivateKey := func(str string) bool {
		// BLS private keys start with "SECRET1P..." or "TSECRET1P...".
		return strings.Contains(str, "secret1p")
	}

	maybeEd25519PrivateKey := func(str string) bool {
		// Ed25519 private keys start with "SECRET1R..." or "TSECRET1R...".
		return strings.Contains(str, "secret1r")
	}

	maybeSecp256k1PrivateKey := func(str string) bool {
		// Secp256k1 private keys start with "SECRET1Y..." or "TSECRET1Y...".
		return strings.Contains(str, "secret1y")
	}

	var prv crypto.PrivateKey
	var err error
	switch {
	case maybeBLSPrivateKey(str):
		prv, err = bls.PrivateKeyFromString(str)

	case maybeEd25519PrivateKey(str):
		prv, err = ed25519.PrivateKeyFromString(str)

	case maybeSecp256k1PrivateKey(str):
		prv, err = secp256k1.PrivateKeyFromString(str)

	default:
		err = ErrInvalidPrivateKey
	}

	return prv, err
}

func PublicKeyFromString(str string) (crypto.PublicKey, error) {
	str = strings.TrimSpace(strings.ToLower(str))

	maybeBLSPublicKey := func(str string) bool {
		// BLS public keys start with "public1p..." or "tpublic1p...".
		return strings.Contains(str, "public1p")
	}

	maybeEd25519PublicKey := func(str string) bool {
		// Ed25519 public keys start with "public1r..." or "tpublic1r...".
		return strings.Contains(str, "public1r")
	}

	maybeSecp256k1PublicKey := func(str string) bool {
		// Secp256k1 public keys start with "public1y..." or "tpublic1y...".
		return strings.Contains(str, "public1y")
	}

	var pub crypto.PublicKey
	var err error
	switch {
	case maybeBLSPublicKey(str):
		pub, err = bls.PublicKeyFromString(str)

	case maybeEd25519PublicKey(str):
		pub, err = ed25519.PublicKeyFromString(str)

	case maybeSecp256k1PublicKey(str):
		pub, err = secp256k1.PublicKeyFromString(str)

	default:
		err = ErrInvalidPublicKey
	}

	return pub, err
}
