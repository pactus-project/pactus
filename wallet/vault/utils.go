package vault

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/util"
	"github.com/tyler-smith/go-bip39"
)

func derivePathToString(path []uint32) string {
	str := "m"
	for _, i := range path {
		if i >= hdkeychain.HardenedKeyStart {
			str += fmt.Sprintf("/%d'", i-hdkeychain.HardenedKeyStart)
		} else {
			str += fmt.Sprintf("/%d", i)
		}
	}
	return str
}

func stringToDerivePath(str string) ([]uint32, error) {
	sub := strings.Split(str, "/")
	if sub[0] != "m" {
		return nil, ErrInvalidPath
	}
	path := []uint32{}
	for i := 1; i < len(sub); i++ {
		indexStr := sub[i]
		added := uint32(0)
		if indexStr[len(indexStr)-1] == '\'' {
			added = hdkeychain.HardenedKeyStart
			indexStr = indexStr[:len(indexStr)-1]
		}
		val, err := strconv.ParseInt(indexStr, 10, 32)
		if err != nil {
			return nil, err
		}
		path = append(path, uint32(val)+added)
	}

	return path, nil
}

// GenerateMnemonic generates a new mnemonic (seed phrase) based on BIP-39
// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
func GenerateMnemonic(bitSize int) string {
	entropy, err := bip39.NewEntropy(bitSize)
	util.ExitOnErr(err)

	mnemonic, err := bip39.NewMnemonic(entropy)
	util.ExitOnErr(err)

	return mnemonic
}

// CheckMnemonic validates a mnemonic (seed phrase) based on BIP-39
func CheckMnemonic(mnemonic string) error {
	_, err := bip39.EntropyFromMnemonic(mnemonic)
	return err
}
