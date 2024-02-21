package vault

import (
	"cmp"
	"encoding/json"
	"slices"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet2/addresspath"
	"github.com/pactus-project/pactus/wallet2/db"
	"github.com/tyler-smith/go-bip39"
)

//
// Deterministic Hierarchy derivation path
//
// Specification
//
// We define the following 4 levels in BIP32 path:
//
// m / purpose' / coin_type' / address_type' / address_index
//
// Where:
//   `'` Apostrophe in the path indicates that BIP32 hardened derivation is used.
//   `m` Denotes the master node (or root) of the tree
//   `/` Separates the tree into depths, thus i / j signifies that j is a child of i
//   `purpose` is set to 12381 which is the name of the new curve (BLS12-381).
//   `coin_type` is set 21888 for Mainnet, 21777 for Testnet
//   `address_type` determine the type of address
//   `address_index` is a sequential number and increase when a new address is derived.
//
// References:
// PIP-8: https://pips.pactus.org/PIPs/pip-8

const (
	TypeFull     = int(1)
	TypeNeutered = int(2)
)

const (
	PurposeBLS12381         = uint32(12381)
	PurposeImportPrivateKey = uint32(65535)
)

type Vault struct {
	db       db.DB
	Type     int    `json:"type"`      // Wallet type. 1: Full keys, 2: Neutered
	CoinType uint32 `json:"coin_type"` // Coin type: 21888 for Mainnet, 21777 for Testnet
	// Addresses map[string]AddressInfo `json:"addresses"` // All addresses that are stored in the wallet
	Encrypter encrypter.Encrypter `json:"encrypter"` // Encryption algorithm
	KeyStore  string              `json:"key_store"` // KeyStore that stores the secrets and encrypts using Encrypter
	Purposes  purposes            `json:"purposes"`  // Contains Purpose 12381 for BLS signature
}

type keyStore struct {
	MasterNode   masterNode `json:"master_node"`   // HD Root Tree (Master node)
	ImportedKeys []string   `json:"imported_keys"` // Imported private keys
}

type masterNode struct {
	Mnemonic string `json:"seed,omitempty"` // Seed phrase or mnemonic (encrypted)
}

type purposes struct {
	PurposeBLS purposeBLS `json:"purpose_bls"` // BLS Purpose: m/12381'/21888/0'/0'
}

type purposeBLS struct {
	XPubValidator      string `json:"xpub_account"`         // Extended public key for account: m/12381'/21888/1'/0
	XPubAccount        string `json:"xpub_validator"`       // Extended public key for validator: m/12381'/218
	NextAccountIndex   uint32 `json:"next_account_index"`   // Index of next derived account
	NextValidatorIndex uint32 `json:"next_validator_index"` // Index of next derived validator
}

func CreateVaultFromMnemonic(mnemonic string, coinType uint32) (*Vault, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, false)
	if err != nil {
		return nil, err
	}
	enc := encrypter.NopeEncrypter()

	xPubValidator, err := masterKey.DerivePath([]uint32{
		H(PurposeBLS12381),
		H(coinType),
		H(crypto.AddressTypeValidator),
	})
	if err != nil {
		return nil, err
	}

	xPubAccount, err := masterKey.DerivePath([]uint32{
		H(PurposeBLS12381),
		H(coinType),
		H(crypto.AddressTypeBLSAccount),
	})
	if err != nil {
		return nil, err
	}

	ks := keyStore{
		MasterNode: masterNode{
			Mnemonic: mnemonic,
		},
		ImportedKeys: make([]string, 0),
	}

	keyStoreDate, err := json.Marshal(ks)
	if err != nil {
		return nil, err
	}

	return &Vault{
		Type:      TypeFull,
		CoinType:  coinType,
		Encrypter: enc,
		// Addresses: make(map[string]AddressInfo),
		KeyStore: string(keyStoreDate),
		Purposes: purposes{
			PurposeBLS: purposeBLS{
				XPubValidator: xPubValidator.Neuter().String(),
				XPubAccount:   xPubAccount.Neuter().String(),
			},
		},
	}, nil
}

func (v *Vault) Neuter() *Vault {
	neutered := &Vault{
		Type:      TypeNeutered,
		CoinType:  v.CoinType,
		Encrypter: encrypter.NopeEncrypter(),
		// Addresses: make(map[string]AddressInfo),
		KeyStore: "",
		Purposes: v.Purposes,
	}

	// for addr, info := range v.Addresses {
	// 	neutered.Addresses[addr] = info
	// }

	return neutered
}

func (v *Vault) IsNeutered() bool {
	return v.Type == TypeNeutered
}

func (v *Vault) UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error {
	if v.IsNeutered() {
		return ErrNeutered
	}

	keyStore, err := v.decryptKeyStore(oldPassword)
	if err != nil {
		return err
	}

	newEncrypter := encrypter.NopeEncrypter()
	if newPassword != "" {
		newEncrypter = encrypter.DefaultEncrypter(opts...)
	}
	v.Encrypter = newEncrypter
	err = v.encryptKeyStore(keyStore, newPassword)
	if err != nil {
		return err
	}

	v.Encrypter = newEncrypter

	return nil
}

func (v *Vault) Label(addr string) string {
	info, err := v.db.GetAddressByAddress(addr)
	if err != nil {
		return ""
	}

	return info.Label
}

func (v *Vault) SetLabel(address, label string) error {
	addr, err := v.db.GetAddressByAddress(address)
	if err != nil {
		return NewErrAddressNotFound(address)
	}

	addr.Label = label
	_, err = v.db.UpdateAddressLabel(addr)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Addresses() []db.Address {
	addrs, err := v.db.GetAllAddresses()
	if err != nil {
		return nil
	}

	v.sortAddressesByAddressIndex(addrs...)
	v.sortAddressesByAddressType(addrs...)
	v.sortAddressesByPurpose(addrs...)

	return addrs
}

func (v *Vault) AllValidatorAddresses() []db.Address {
	addrs, err := v.db.GetAllAddresses()
	if err != nil {
		return nil
	}

	// TODO refactor me with account count method
	result := make([]db.Address, 0, len(addrs)/2)
	for i := range addrs {
		addrPath, _ := addresspath.NewPathFromString(addrs[i].Path)
		if addrPath.AddressType() == H(crypto.AddressTypeValidator) {
			result = append(result, addrs[i])
		}
	}

	v.sortAddressesByAddressIndex(result...)
	v.sortAddressesByPurpose(result...)

	return result
}

func (v *Vault) AllImportedPrivateKeyAddresses() []db.Address {
	addrs, err := v.db.GetAllAddresses()
	if err != nil {
		return nil
	}

	// TODO refactor me with account count method
	result := make([]db.Address, 0, len(addrs)/2)
	for i := range addrs {
		addrPath, _ := addresspath.NewPathFromString(addrs[i].Path)
		if addrPath.Purpose() == H(PurposeImportPrivateKey) {
			result = append(result, addrs[i])
		}
	}

	v.sortAddressesByAddressIndex(result...)
	v.sortAddressesByAddressType(result...)

	return result
}

func (v *Vault) sortAddressesByAddressIndex(addrs ...db.Address) {
	slices.SortStableFunc(addrs, func(a, b db.Address) int {
		pathA, _ := addresspath.NewPathFromString(a.Path)
		pathB, _ := addresspath.NewPathFromString(b.Path)

		return cmp.Compare(pathA.AddressIndex(), pathB.AddressIndex())
	})
}

func (v *Vault) sortAddressesByAddressType(addrs ...db.Address) {
	slices.SortStableFunc(addrs, func(a, b db.Address) int {
		pathA, _ := addresspath.NewPathFromString(a.Path)
		pathB, _ := addresspath.NewPathFromString(b.Path)

		return cmp.Compare(pathA.AddressType(), pathB.AddressType())
	})
}

func (v *Vault) sortAddressesByPurpose(addrs ...db.Address) {
	slices.SortStableFunc(addrs, func(a, b db.Address) int {
		pathA, _ := addresspath.NewPathFromString(a.Path)
		pathB, _ := addresspath.NewPathFromString(b.Path)

		return cmp.Compare(pathA.Purpose(), pathB.Purpose())
	})
}

func (v *Vault) decryptKeyStore(password string) (*keyStore, error) {
	keyStoreData, err := v.Encrypter.Decrypt(v.KeyStore, password)
	if err != nil {
		return nil, err
	}

	keyStore := new(keyStore)
	err = json.Unmarshal([]byte(keyStoreData), keyStore)
	if err != nil {
		return nil, err
	}

	return keyStore, nil
}

func (v *Vault) encryptKeyStore(keyStore *keyStore, password string) error {
	keyStoreData, err := json.Marshal(keyStore)
	if err != nil {
		return err
	}

	keyStoreEnc, err := v.Encrypter.Encrypt(string(keyStoreData), password)
	if err != nil {
		return err
	}
	v.KeyStore = keyStoreEnc

	return nil
}
