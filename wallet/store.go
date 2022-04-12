package wallet

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type Store struct {
	Version   int       `json:"version"`
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Network   int       `json:"network"`
	Encrypted bool      `json:"encrypted"`
	VaultCRC  uint32    `json:"crc"`
	Vault     *vault    `json:"vault"`
}

type vault struct {
	Addresses []address   `json:"addresses"`
	Seed      seed        `json:"seed"`
	Keystore  []encrypted `json:"keystore"`
}

type address struct {
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

func CreateStoreFromMnemonic(mnemonic, passphrase string, net int) (*Store, error) {
	keyInfo := []byte{} // TODO, update for testnet
	parentSeed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	exitOnErr(err)
	parentKey, err := bls.PrivateKeyFromSeed(parentSeed, keyInfo)
	exitOnErr(err)

	e := newEncrypter(passphrase, net)
	s := &Store{
		Version:   1,
		UUID:      uuid.New(),
		CreatedAt: time.Now().Round(time.Second).UTC(),
		Network:   net,
		Encrypted: len(passphrase) != 0,
		Vault: &vault{
			Seed: seed{
				Method:     "BIP-39",
				ParentSeed: e.encrypt(mnemonic),
				ParentKey:  e.encrypt(base64.StdEncoding.EncodeToString(parentKey.Bytes())),
			},
		},
	}
	return s, nil
}

func (s *Store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	exitOnErr(err)
	return crc32.ChecksumIEEE(d)
}

func (s *Store) UpdatePassword(oldPassphrase, newPassphrase string) error {
	oldEncrypter := newEncrypter(oldPassphrase, s.Network)
	newEncrypter := newEncrypter(newPassphrase, s.Network)

	// Updating parent seed
	parentSeed, err := oldEncrypter.decrypt(s.Vault.Seed.ParentSeed)
	if err != nil {
		return err
	}
	s.Vault.Seed.ParentSeed = newEncrypter.encrypt(parentSeed)

	// Updating parent key
	parentKey, err := oldEncrypter.decrypt(s.Vault.Seed.ParentKey)
	exitOnErr(err)
	s.Vault.Seed.ParentKey = newEncrypter.encrypt(parentKey)

	// Updating private keys
	for i, prv := range s.Vault.Keystore {
		key, err := oldEncrypter.decrypt(prv)
		exitOnErr(err)

		s.Vault.Keystore[i] = newEncrypter.encrypt(key)
	}
	s.Encrypted = len(newPassphrase) > 0
	return nil
}

func (s *Store) Addresses() []AddressInfo {
	addrs := make([]AddressInfo, len(s.Vault.Addresses))
	for i, a := range s.Vault.Addresses {
		addrs[i].Address = a.Address
		addrs[i].Label = a.Label
		addrs[i].Imported = (a.Method == "IMPORTED")
	}

	return addrs
}

func (s *Store) AddressCount() int {
	return len(s.Vault.Addresses)
}

func (s *Store) ImportPrivateKey(passphrase string, prv *bls.PrivateKey) error {
	/// Decrypt parent key to make sure the passphrase is correct
	_, err := s.parentKey(passphrase)
	if err != nil {
		return err
	}
	if s.Contains(prv.PublicKey().Address().String()) {
		return ErrAddressExists
	}

	e := newEncrypter(passphrase, s.Network)
	s.Vault.Keystore = append(s.Vault.Keystore, e.encrypt(prv.String()))

	p := newParams()
	p.SetUint32("index", uint32(len(s.Vault.Keystore)-1))
	s.Vault.Addresses = append(s.Vault.Addresses, address{
		Method:  "IMPORTED",
		Address: prv.PublicKey().Address().String(),
		Params:  p,
	})

	return nil
}

func (s *Store) newKeySeed(passphrase string) ([]byte, error) {
	mnemonic, err := s.Mnemonic(passphrase)
	if err != nil {
		return nil, err
	}
	parentSeed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	exitOnErr(err)
	data := []byte{0}
	hmacKey := sha256.Sum256(parentSeed)

	checkKeySeed := func(seed []byte) bool {
		for _, a := range s.Vault.Addresses {
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

func (s *Store) derivePrivateKey(parentKey []byte, keySeed []byte) *bls.PrivateKey {
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

func (s *Store) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	for _, a := range s.Vault.Addresses {
		if a.Address == addr {
			var prv *bls.PrivateKey
			switch a.Method {
			case "IMPORTED":
				{
					e, err := s.makeEncrypter(passphrase)
					if err != nil {
						return nil, err
					}
					index := a.Params.GetUint32("index")
					prvStr, err := e.decrypt(s.Vault.Keystore[index])
					if err != nil {
						return nil, err
					}
					prv, err = bls.PrivateKeyFromString(prvStr)
					exitOnErr(err)
				}
			case "KDF-CHAIN":
				{
					seed := a.Params.GetBytes("seed")
					parentKey, err := s.parentKey(passphrase)
					if err != nil {
						return nil, err
					}
					prv = s.derivePrivateKey(parentKey, seed)
				}
			}

			if prv.PublicKey().Address().String() != addr {
				// If you see tis error, please report it
				exitOnErr(fmt.Errorf("invalid private key for given address"))
			}
			return prv, nil
		}
	}

	return nil, NewErrAddressNotFound(addr)
}

func (s *Store) NewAddress(passphrase, label string) (string, error) {
	keySeed, err := s.newKeySeed(passphrase)
	if err != nil {
		return "", err
	}
	parentKey, err := s.parentKey(passphrase)
	exitOnErr(err) // Password has been checked in previous line

	prv := s.derivePrivateKey(parentKey, keySeed)

	params := newParams()
	params.SetBytes("seed", keySeed)
	a := address{
		Method:  "KDF-CHAIN",
		Address: prv.PublicKey().Address().String(),
		Label:   label,
		Params:  params,
	}

	s.Vault.Addresses = append(s.Vault.Addresses, a)

	return a.Address, nil
}

func (s *Store) Contains(addr string) bool {
	return s.getAddressInfo(addr) != nil
}

func (s *Store) getAddressInfo(addr string) *address {
	for _, a := range s.Vault.Addresses {
		if a.Address == addr {
			return &a
		}
	}
	return nil
}

func (s *Store) Mnemonic(passphrase string) (string, error) {
	e, err := s.makeEncrypter(passphrase)
	if err != nil {
		return "", err
	}
	return e.decrypt(s.Vault.Seed.ParentSeed)
}

func (s *Store) parentKey(passphrase string) ([]byte, error) {
	e, err := s.makeEncrypter(passphrase)
	if err != nil {
		return nil, err
	}
	m, err := e.decrypt(s.Vault.Seed.ParentKey)
	if err != nil {
		return nil, err
	}
	parentKey, err := base64.StdEncoding.DecodeString(m)
	exitOnErr(err)

	return parentKey, nil
}

func (s *Store) makeEncrypter(passphrase string) (encrypter, error) {
	if s.Encrypted && len(passphrase) != 0 {
		return newArgon2Encrypter(passphrase), nil
	}

	if !s.Encrypted && len(passphrase) == 0 {
		return newNopeEncrypter(), nil
	}

	return nil, ErrInvalidPassword
}
