package wallet

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"errors"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type vault struct {
	Addresses []address `json:"addresses"`
	Seed      seed      `json:"seed"`
	Keystore  keystore  `json:"keystore"`
}

type address struct {
	Method  string `json:"Method"`
	Address string `json:"address"`
	Params  params `json:"params"`
}

type keystore struct {
	Prv []encrypted `json:"prv"`
}

func recoverVault(mnemonic string) *vault {
	s := recoverSeed(mnemonic)

	v := &vault{
		Seed: s,
	}

	v.generateStartKeys("", 21)

	return v
}

func newVault(passphrase string) *vault {
	s := newSeed(passphrase)

	v := &vault{
		Seed: s,
	}

	v.generateStartKeys(passphrase, 21)

	return v
}

func (v *vault) deriveNewKeySeed(passphrase string) []byte {
	data := []byte{0}
	parentSeed := v.Seed.parentSeed(passphrase)
	hmacKey := sha256.Sum256(parentSeed)

	checkKeySeed := func(seed []byte) bool {
		for _, a := range v.Addresses {
			if bytes.Equal(a.Params.getBytes("seed"), seed) {
				return true
			}
		}
		return false
	}

	for {
		hmac512 := hmac.New(sha512.New, hmacKey[:])
		_, err := hmac512.Write(data[:])
		exitOnErr(err)
		hash512 := hmac512.Sum(nil)
		keySeed := hash512[:32]
		nextData := hash512[32:]

		if !checkKeySeed(keySeed) {
			return keySeed
		}

		data = nextData
	}
}

/// Notes:
/// 1- Derive a parnet key based entropy seed. Entropy seed will be used for wallet recovery
/// 2- Deriving Child keys should be deterministic
/// 3- Exposing any child key, should not expose parnet key or any other child keys
/// 4- Not saving child keys inside wallet.
/// 5- Child key should be recovered from parnet key, hash seed and a derive seed
/// 6- If mater key is exposed, non of the child keys can be derived without knowing the seed hash.

func (v *vault) derivePrivayeKey(passphrase string, keySeed []byte) *bls.PrivateKey {
	keyInfo := []byte{} // TODO, update for testnet
	parnetKey := v.Seed.parentKey(passphrase)

	// To derive a new key, we need:
	//    1- Parent Key
	//    2- Key seed.
	//

	hmac512 := hmac.New(sha512.New, parnetKey.Bytes())
	_, err := hmac512.Write(keySeed) /// Note #6
	exitOnErr(err)
	ikm := hmac512.Sum(nil)

	prv, err := bls.PrivateKeyFromSeed(ikm, keyInfo)
	exitOnErr(err)

	return prv
}

func (v *vault) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	for _, a := range v.Addresses {
		if a.Address == addr {
			seed := a.Params.getBytes("seed")
			prv := v.derivePrivayeKey(passphrase, seed)
			return prv, nil
		}
	}

	return nil, errors.New("address not found")
}
func (v *vault) generateStartKeys(passphrase string, count int) {
	for i := 0; i < count; i++ {
		seed := v.deriveNewKeySeed(passphrase)
		prv := v.derivePrivayeKey(passphrase, seed)

		a := address{}
		a.Address = prv.PublicKey().Address().String()
		a.Params = newParams()
		a.Params.setBytes("seed", seed)
		a.Method = "BLS_KDF_CHAIN"
		v.Addresses = append(v.Addresses, a)
	}
}

func (v *vault) Contains(addr crypto.Address) bool {
	return v.GetAddressInfo(addr) != nil
}

func (v *vault) GetAddressInfo(addr crypto.Address) *address {
	for _, a := range v.Addresses {
		if a.Address == addr.String() {
			return &a
		}
	}
	return nil
}
