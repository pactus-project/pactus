package vault

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"

	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type AddressInfo struct {
	Address  string
	Label    string
	Imported bool
}

type Vault struct {
	Encrypted bool          `json:"encrypted"`
	Seed      seed          `json:"seed"`
	Keystore  []encrypted   `json:"keystore"`
	Addresses []addressInfo `json:"addresses"`
}

type addressInfo struct {
	Method  string `json:"method"`
	Address string `json:"address"`
	Label   string `json:"label"`
	Params  params `json:"params"`
}

type seed struct {
	Method     string    `json:"method"`
	ParentSeed encrypted `json:"seed"`
	ParentKey  encrypted `json:"key"`
}

// GenerateMnemonic generates a new mnemonic (seed phrase based on BIP-32)
func GenerateMnemonic() string {
	entropy, err := bip39.NewEntropy(128)
	exitOnErr(err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	exitOnErr(err)
	return mnemonic
}

func CreateVaultFromMnemonic(mnemonic, password string, keyInfo []byte) (*Vault, error) {
	parentSeed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	parentKey, err := bls.PrivateKeyFromSeed(parentSeed, keyInfo)
	exitOnErr(err)

	e := newEncrypter(password)
	return &Vault{
		Encrypted: len(password) > 0,
		Seed: seed{
			Method:     "BIP-39",
			ParentSeed: e.encrypt(mnemonic),
			ParentKey:  e.encrypt(base64.StdEncoding.EncodeToString(parentKey.Bytes())),
		},
		Keystore:  make([]encrypted, 0),
		Addresses: make([]addressInfo, 0),
	}, nil

}

func (v *Vault) UpdatePassword(oldPassword, newPassword string) error {
	oldEncrypter := newEncrypter(oldPassword)
	newEncrypter := newEncrypter(newPassword)

	// Updating parent seed
	parentSeed, err := oldEncrypter.decrypt(v.Seed.ParentSeed)
	if err != nil {
		return err
	}
	v.Seed.ParentSeed = newEncrypter.encrypt(parentSeed)

	// Updating parent key
	parentKey, err := oldEncrypter.decrypt(v.Seed.ParentKey)
	exitOnErr(err)
	v.Seed.ParentKey = newEncrypter.encrypt(parentKey)

	// Updating private keys
	for i, prv := range v.Keystore {
		key, err := oldEncrypter.decrypt(prv)
		exitOnErr(err)

		v.Keystore[i] = newEncrypter.encrypt(key)
	}
	v.Encrypted = len(newPassword) > 0
	return nil
}

func (v *Vault) AddressInfos() []AddressInfo {
	addrs := make([]AddressInfo, 0, len(v.Addresses))

	for _, info := range v.Addresses {
		addrs = append(addrs, AddressInfo{
			Address:  info.Address,
			Label:    info.Label,
			Imported: (info.Method == "IMPORTED"),
		})
	}

	return addrs
}

func (v *Vault) IsEncrypted() bool {
	return v.Encrypted
}

func (v *Vault) AddressCount() int {
	return len(v.Addresses)
}

func (v *Vault) ImportPrivateKey(password string, prvStr string) error {
	prv, err := bls.PrivateKeyFromString(prvStr)
	if err != nil {
		return err
	}
	addr := prv.PublicKey().Address().String()
	if v.Contains(addr) {
		return ErrAddressExists
	}
	/// Decrypt parent key to make sure the password is correct
	_, err = v.parentKey(password)
	if err != nil {
		return err
	}

	e := newEncrypter(password)
	v.Keystore = append(v.Keystore, e.encrypt(prv.String()))

	p := newParams()
	p.SetUint32("index", uint32(len(v.Keystore)-1))
	v.Addresses = append(v.Addresses, addressInfo{
		Method:  "IMPORTED",
		Address: addr,
		Params:  p,
	})
	return nil
}

func (v *Vault) newKeySeed(password string) ([]byte, error) {
	mnemonic, err := v.Mnemonic(password)
	if err != nil {
		return nil, err
	}
	parentSeed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	exitOnErr(err)
	data := []byte{0}
	hmacKey := sha256.Sum256(parentSeed)

	checkKeySeed := func(seed []byte) bool {
		for _, a := range v.Addresses {
			if safeCmp(seed, a.Params.GetBytes("seed")) {
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
			return keySeed, nil
		}

		data = nextData
	}
}

// Note:
// 1- Deriving Child key seeds from parent seed
// 2- Exposing any child key, should not expose parent key or any other child keys

func (v *Vault) derivePrivateKey(parentKey []byte, keySeed []byte) *bls.PrivateKey {
	keyInfo := []byte{} // TODO, update for testnet

	// To derive a new key, we need:
	//    1- Parent Key
	//    2- Key seed.
	//

	hmac512 := hmac.New(sha512.New, parentKey)
	_, err := hmac512.Write(keySeed)
	exitOnErr(err)
	ikm := hmac512.Sum(nil)

	prv, err := bls.PrivateKeyFromSeed(ikm, keyInfo)
	exitOnErr(err)

	return prv
}

func (v *Vault) getPrivateKey(password, addr string) (crypto.PrivateKey, error) {
	info := v.getAddressInfo(addr)
	if info == nil {
		return nil, NewErrAddressNotFound(addr)
	}

	var prv crypto.PrivateKey
	switch info.Method {
	case "IMPORTED":
		{
			e, err := v.makeEncrypter(password)
			if err != nil {
				return nil, err
			}
			index := info.Params.GetUint32("index")
			prvStr, err := e.decrypt(v.Keystore[index])
			if err != nil {
				return nil, err
			}
			prv, err = bls.PrivateKeyFromString(prvStr)
			exitOnErr(err)
		}
	case "KDF-CHAIN":
		{
			seed := info.Params.GetBytes("seed")
			parentKey, err := v.parentKey(password)
			if err != nil {
				return nil, err
			}
			prv = v.derivePrivateKey(parentKey, seed)
		}
	}

	if prv.PublicKey().Address().String() != addr {
		// If you see this error, please report it
		exitOnErr(errors.New("invalid private key for given address"))
	}
	return prv, nil
}

func (v *Vault) PrivateKey(password, addr string) (string, error) {
	prv, err := v.getPrivateKey(password, addr)
	if err != nil {
		return "", err
	}
	return prv.String(), nil
}

func (v *Vault) PublicKey(password, addr string) (string, error) {
	prv, err := v.getPrivateKey(password, addr)
	if err != nil {
		return "", err
	}
	return prv.PublicKey().String(), nil
}

func (v *Vault) MakeNewAddress(password, label string) (string, error) {
	keySeed, err := v.newKeySeed(password)
	if err != nil {
		return "", err
	}
	parentKey, err := v.parentKey(password)
	exitOnErr(err) // Password has been checked in previous line

	prv := v.derivePrivateKey(parentKey, keySeed)
	addr := prv.PublicKey().Address().String()
	params := newParams()
	params.SetBytes("seed", keySeed)

	v.Addresses = append(v.Addresses, addressInfo{
		Method:  "KDF-CHAIN",
		Address: addr,
		Label:   label,
		Params:  params,
	})

	return addr, nil
}

func (v *Vault) getAddressInfo(addr string) *addressInfo {
	for _, info := range v.Addresses {
		if info.Address == addr {
			return &info
		}
	}

	return nil
}

func (v *Vault) Contains(addr string) bool {
	return v.getAddressInfo(addr) != nil
}

func (v *Vault) Mnemonic(password string) (string, error) {
	e, err := v.makeEncrypter(password)
	if err != nil {
		return "", err
	}
	return e.decrypt(v.Seed.ParentSeed)
}

func (v *Vault) parentKey(password string) ([]byte, error) {
	e, err := v.makeEncrypter(password)
	if err != nil {
		return nil, err
	}
	m, err := e.decrypt(v.Seed.ParentKey)
	if err != nil {
		return nil, err
	}
	parentKey, err := base64.StdEncoding.DecodeString(m)
	exitOnErr(err)

	return parentKey, nil
}

func (v *Vault) makeEncrypter(password string) (encrypter, error) {
	if v.Encrypted && len(password) != 0 {
		return newArgon2Encrypter(password), nil
	}

	if !v.Encrypted && len(password) == 0 {
		return newNopeEncrypter(), nil
	}

	return nil, ErrInvalidPassword
}
