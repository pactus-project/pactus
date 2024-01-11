package vault

import (
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/exp/constraints"
)

// GenerateMnemonic generates a new mnemonic (seed phrase) based on BIP-39
// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
func GenerateMnemonic(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)

	return mnemonic, nil
}

// CheckMnemonic validates a mnemonic (seed phrase) based on BIP-39.
func CheckMnemonic(mnemonic string) error {
	_, err := bip39.EntropyFromMnemonic(mnemonic)

	return err
}

// H hardens the value 'i' by adding it to 0x80000000 (2^31).
func H[T constraints.Integer](i T) uint32 {
	return uint32(i) + hdkeychain.HardenedKeyStart
}
