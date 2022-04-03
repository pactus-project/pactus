package wallet

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"hash/crc32"
	"time"

	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type Store struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Network   int       `json:"network"`
	Encrypted bool      `json:"encrypted"`
	VaultCRC  uint32    `json:"crc"`
	Vault     *vault    `json:"vault"`
}

type vault struct {
	Addresses []address `json:"addresses"`
	Seed      seed      `json:"seed"`
	Keystore  keystore  `json:"keystore"`
}

type address struct {
	Method  string `json:"method"`
	Address string `json:"address"`
	Params  params `json:"params"`
}

type seed struct {
	Method     string    `json:"method"`
	ParentSeed encrypted `json:"seed"`
	ParentKey  encrypted `json:"prv"`
}

type keystore struct {
	Prv []encrypted `json:"prv"`
}

func RecoverStore(mnemonic string, net int) *Store {
	return createStoreFromMnemonic("", mnemonic, 1)
}

func NewStore(passphrase string, net int) *Store {
	entropy, err := bip39.NewEntropy(128)
	exitOnErr(err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	exitOnErr(err)
	return createStoreFromMnemonic(passphrase, mnemonic, 1)
}

func createStoreFromMnemonic(passphrase string, mnemonic string, net int) *Store {
	ikm, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	exitOnErr(err)
	parentKey, err := bls.PrivateKeyFromSeed(ikm, nil)
	exitOnErr(err)

	e := newEncrypter(passphrase)

	s := &Store{
		Version:   1,
		CreatedAt: time.Now(),
		Network:   net,
		Encrypted: false,
		Vault: &vault{
			Seed: seed{
				Method:     "BIP-39",
				ParentSeed: e.encrypt(mnemonic),
				ParentKey:  e.encrypt(parentKey.String()),
			},
		},
	}

	s.generateStartKeys(passphrase, 21)
	return s
}

func (s *Store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	exitOnErr(err)
	return crc32.ChecksumIEEE(d)
}

func (s *Store) Addresses() []crypto.Address {
	addrs := make([]crypto.Address, len(s.Vault.Addresses))
	for i, a := range s.Vault.Addresses {
		addr, err := crypto.AddressFromString(a.Address)
		exitOnErr(err)
		addrs[i] = addr
	}

	return addrs
}

func (s *Store) ImportPrivateKey(passphrase string, prv *bls.PrivateKey) error {
	if s.Contains(prv.PublicKey().Address()) {
		return errors.New("address already exists")
	}

	e := newEncrypter(passphrase)
	s.Vault.Keystore.Prv = append(s.Vault.Keystore.Prv, e.encrypt(prv.String()))

	p := newParams()
	p.SetUint32("index", uint32(len(s.Vault.Keystore.Prv)-1))
	s.Vault.Addresses = append(s.Vault.Addresses, address{
		Method:  "IMPORTED",
		Address: prv.PublicKey().Address().String(),
		Params:  p,
	})

	return nil
}

func (s *Store) deriveNewKeySeed(passphrase string) []byte {
	data := []byte{0}
	parentSeed := s.ParentSeed(passphrase)
	hmacKey := sha256.Sum256(parentSeed)

	checkKeySeed := func(seed []byte) bool {
		for _, a := range s.Vault.Addresses {
			if bytes.Equal(a.Params.GetBytes("seed"), seed) {
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

/// Note:
/// 1- Deriving Child key seeds from parent seed
/// 2- Exposing any child key, should not expose parnet key or any other child keys

func (s *Store) derivePrivayeKey(passphrase string, keySeed []byte) *bls.PrivateKey {
	keyInfo := []byte{} // TODO, update for testnet
	parnetKey := s.ParentKey(passphrase)

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

func (s *Store) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	for _, a := range s.Vault.Addresses {
		if a.Address == addr {
			switch a.Method {
			case "IMPORTED":
				{
					e := newEncrypter(passphrase)
					index := a.Params.GetUint32("index")
					prvStr, err := e.decrypt(s.Vault.Keystore.Prv[index])
					exitOnErr(err)
					prv, err := bls.PrivateKeyFromString(prvStr)
					exitOnErr(err)
					return prv, nil
				}
			case "BLS_KDF_CHAIN":
				{
					seed := a.Params.GetBytes("seed")
					prv := s.derivePrivayeKey(passphrase, seed)
					return prv, nil
				}
			}
		}
	}

	return nil, errors.New("address not found")
}

func (s *Store) generateStartKeys(passphrase string, count int) {
	for i := 0; i < count; i++ {
		seed := s.deriveNewKeySeed(passphrase)
		prv := s.derivePrivayeKey(passphrase, seed)

		a := address{}
		a.Address = prv.PublicKey().Address().String()
		a.Params = newParams()
		a.Params.SetBytes("seed", seed)
		a.Method = "BLS_KDF_CHAIN"
		s.Vault.Addresses = append(s.Vault.Addresses, a)
	}
}

func (s *Store) Contains(addr crypto.Address) bool {
	return s.getAddressInfo(addr) != nil
}

func (s *Store) getAddressInfo(addr crypto.Address) *address {
	for _, a := range s.Vault.Addresses {
		if a.Address == addr.String() {
			return &a
		}
	}
	return nil
}

func (s *Store) ParentSeed(passphrase string) []byte {
	h, err := bip39.NewSeedWithErrorChecking(s.Mnemonic(passphrase), "")
	exitOnErr(err)

	return h
}

func (s *Store) Mnemonic(passphrase string) string {
	m, err := newEncrypter(passphrase).decrypt(s.Vault.Seed.ParentSeed)
	exitOnErr(err)

	return m
}

func (s *Store) ParentKey(passphrase string) *bls.PrivateKey {
	m, err := newEncrypter(passphrase).decrypt(s.Vault.Seed.ParentKey)
	exitOnErr(err)

	prv, err := bls.PrivateKeyFromString(m)
	exitOnErr(err)

	return prv
}
