package vault

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/types/crypto/bls/hdkeychain"
	"github.com/zarbchain/zarb-go/util"
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
func GenerateMnemonic() string {
	entropy, err := bip39.NewEntropy(128)
	util.ExitOnErr(err)

	mnemonic, err := bip39.NewMnemonic(entropy)
	util.ExitOnErr(err)

	return mnemonic
}
