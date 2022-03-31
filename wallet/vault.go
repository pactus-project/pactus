package wallet

import (
	"crypto/hmac"
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

func newVault(passphrase string) *vault {
	s := newSeed(passphrase)

	v := &vault{
		Seed: s,
	}

	for i := 0; i < 21; i++ {
		v.deriveNewKey(passphrase)
	}

	return v
}
func (v *vault) deriveNewKey(passphrase string) address {
	deriveSeed := []byte{}
	for {
		prv, nextDeriveSeed := v.deriveKey(passphrase, deriveSeed)

		if !v.Contains(prv.PublicKey().Address()) {
			a := address{}
			a.Address = prv.PublicKey().Address().String()
			a.Params = newParams()
			a.Params.setBytes("seed", deriveSeed)
			a.Method = "BLS_HMAC_HKDF_SEED"
			v.Addresses = append(v.Addresses, a)
			return a
		}
		deriveSeed = nextDeriveSeed
	}
}

/// Notes:
/// 1- Derive a parnet key based entropy seed. Entropy seed will be used for wallet recovery
/// 2- Deriving Child keys should be deterministic
/// 3- Exposing any child key, should not expose parnet key or any other child keys
/// 4- Not saving child keys inside wallet.
/// 5- Child key should be recovered from parnet key, hash seed and a derive seed
/// 6- If mater key is exposed, non of the child keys can be derived without knowing the seed hash.

func (v *vault) deriveKey(passphrase string, deriveSeed []byte) (*bls.PrivateKey, []byte) {
	keyInfo := []byte{} // TODO, update for testnet
	hashedSeed := v.Seed.hashedSeed(passphrase)
	parnetKey := v.Seed.parentKey()

	/// To derive a new key, we need these variables:
	///    1- Parent Key
	///    2- mnemonic's seed hash
	///    2- Derive seed.
	///

	hmac512 := hmac.New(sha512.New, parnetKey.Bytes())
	_, err := hmac512.Write(hashedSeed[:]) /// Note #6
	exitOnErr(err)
	_, err = hmac512.Write(deriveSeed[:])
	exitOnErr(err)
	hash512 := hmac512.Sum(nil)
	ikm := hash512[:32]
	nextDeriveSeed := hash512[32:]

	prv, err := bls.PrivateKeyFromSeed(ikm, keyInfo)
	exitOnErr(err)

	return prv, nextDeriveSeed
}

func (v *vault) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	for _, a := range v.Addresses {
		if a.Address == addr {
			seed := a.Params.getBytes("seed")
			prv, _ := v.deriveKey(passphrase, seed)
			return prv, nil
		}
	}

	return nil, errors.New("address not found")
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
