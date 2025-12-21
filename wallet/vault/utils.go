package vault

import (
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
