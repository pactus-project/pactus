package vault

import (
	"github.com/cosmos/go-bip39"
	"github.com/pactus-project/pactus/wallet/addresspath"
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
	_, err := bip39.MnemonicToByteArray(mnemonic)
	return err
}

// _H hardens the integer value 'i' by adding 0x80000000 (2^31) to it.
// This function does not check if 'i' is already hardened.
func _H[T constraints.Integer](i T) uint32 {
	return uint32(i) + addresspath.HardenedKeyStart
}

// _N de-hardens the integer value 'i' by subtracting 0x80000000 (2^31) from it.
// This function does not check if 'i' is already non-hardened.
func _N[T constraints.Integer](i T) uint32 {
	return uint32(i) - addresspath.HardenedKeyStart
}
