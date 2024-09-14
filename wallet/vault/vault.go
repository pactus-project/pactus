package vault

import (
	"cmp"
	"encoding/json"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/crypto/ed25519"
	ed25519hdkeychain "github.com/pactus-project/pactus/crypto/ed25519/hdkeychain"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/exp/slices"
)

//
// Deterministic Hierarchical Derivation Path
//
// Overview:
//
// This specification defines a hierarchical derivation path for generating addresses, based on BIP32.
// The path is structured into four distinct levels:
//
// m / purpose' / coin_type' / address_type' / address_index
//
// Explanation:
//
//   `m` Denotes the master node (or root) of the tree
//   `'` Apostrophe in the path indicates that BIP32 hardened derivation is used.
//   `/` Separates the tree into depths, thus i / j signifies that j is a child of i
//
// Path Components:
//
// * `purpose`: Indicates the specific use case for the derived addresses:
//    - 12381: Used for the BLS12-381 curve, based on PIP-8.
//    - 65535: Used for imported private keys, based on PIP-13.
//    - 44: A comprehensive purpose for standard curves, based on BIP-44.
//
// * `coin_type`: Identifies the coin type:
//    - 21888: Pactus Mainnet
//    - 21777: Pactus Testnet
//
// * `address_type`: Specifies the type of address.
//
// * `address_index`: A sequential number and increase when a new address is derived.
//
// References:
//  - https://pips.pactus.org/PIPs/pip-8
//  - https://pips.pactus.org/PIPs/pip-13
//  - https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
//

const (
	TypeFull     = int(1)
	TypeNeutered = int(2)
)

type AddressInfo struct {
	Address   string `json:"address"`    // Address in the wallet
	PublicKey string `json:"public_key"` // Public key associated with the address
	Label     string `json:"label"`      // Label for the address
	Path      string `json:"path"`       // Path for the address
}

const (
	PurposeBLS12381         = uint32(12381)
	PurposeBIP44            = uint32(44)
	PurposeImportPrivateKey = uint32(65535)
)

type Vault struct {
	Type      int                    `json:"type"`      // Wallet type. 1: Full keys, 2: Neutered
	CoinType  uint32                 `json:"coin_type"` // Coin type: 21888 for Mainnet, 21777 for Testnet
	Addresses map[string]AddressInfo `json:"addresses"` // All addresses that are stored in the wallet
	Encrypter encrypter.Encrypter    `json:"encrypter"` // Encryption algorithm
	KeyStore  string                 `json:"key_store"` // KeyStore that stores the secrets and encrypts using Encrypter
	Purposes  purposes               `json:"purposes"`  // Contains Purpose 12381 for BLS signature
}

type keyStore struct {
	MasterNode   masterNode `json:"master_node"`   // HD Root Tree (Master node)
	ImportedKeys []string   `json:"imported_keys"` // Imported private keys
}

type masterNode struct {
	Mnemonic string `json:"seed,omitempty"` // Seed phrase or mnemonic (encrypted)
}

type purposes struct {
	PurposeBLS   purposeBLS   `json:"purpose_bls"`   // BLS Purpose: m/12381'/21888'/1' or 2'/0
	PurposeBIP44 purposeBIP44 `json:"purpose_bip44"` // BIP44 Purpose: m/44'/21888'/3'/0'
}

type purposeBLS struct {
	XPubValidator      string `json:"xpub_account"`         // Extended public key for account: m/12381'/21888'/1'/0
	XPubAccount        string `json:"xpub_validator"`       // Extended public key for validator: m/12381'/21888'/2'/0
	NextAccountIndex   uint32 `json:"next_account_index"`   // Index of next derived account
	NextValidatorIndex uint32 `json:"next_validator_index"` // Index of next derived validator
}

type purposeBIP44 struct {
	NextEd25519Index uint32 `json:"next_ed25519_index"` // Index of next Ed25519 derived account: m/44'/21888/3'/0'
}

func CreateVaultFromMnemonic(mnemonic string, coinType uint32) (*Vault, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := blshdkeychain.NewMaster(seed, false)
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
		Addresses: make(map[string]AddressInfo),
		KeyStore:  string(keyStoreDate),
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
		Addresses: make(map[string]AddressInfo),
		KeyStore:  "",
		Purposes:  v.Purposes,
	}

	for addr, info := range v.Addresses {
		neutered.Addresses[addr] = info
	}

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
	info, ok := v.Addresses[addr]
	if !ok {
		return ""
	}

	return info.Label
}

func (v *Vault) SetLabel(addr, label string) error {
	info, ok := v.Addresses[addr]
	if !ok {
		return NewErrAddressNotFound(addr)
	}

	info.Label = label
	v.Addresses[addr] = info

	return nil
}

func (v *Vault) AddressInfos() []AddressInfo {
	addrs := make([]AddressInfo, 0, 1)
	for _, addrInfo := range v.Addresses {
		addrs = append(addrs, addrInfo)
	}

	v.SortAddressesByAddressIndex(addrs...)
	v.SortAddressesByAddressType(addrs...)
	v.SortAddressesByPurpose(addrs...)

	return addrs
}

func (v *Vault) AllValidatorAddresses() []AddressInfo {
	addrs := make([]AddressInfo, 0, v.AddressCount()/2)
	for _, addrInfo := range v.Addresses {
		addrPath, _ := addresspath.FromString(addrInfo.Path)
		if addrPath.AddressType() == H(crypto.AddressTypeValidator) {
			addrs = append(addrs, addrInfo)
		}
	}

	v.SortAddressesByAddressIndex(addrs...)
	v.SortAddressesByPurpose(addrs...)

	return addrs
}

func (v *Vault) AllAccountAddresses() []AddressInfo {
	addrs := make([]AddressInfo, 0, v.AddressCount()/2)
	for _, addrInfo := range v.Addresses {
		addrPath, _ := addresspath.FromString(addrInfo.Path)
		if addrPath.AddressType() != H(crypto.AddressTypeValidator) {
			addrs = append(addrs, addrInfo)
		}
	}

	v.SortAddressesByAddressIndex(addrs...)
	v.SortAddressesByPurpose(addrs...)

	return addrs
}

func (v *Vault) AllImportedPrivateKeysAddresses() []AddressInfo {
	addrs := make([]AddressInfo, 0, v.AddressCount()/2)
	for _, addrInfo := range v.Addresses {
		addrPath, _ := addresspath.FromString(addrInfo.Path)
		if addrPath.Purpose() == H(PurposeImportPrivateKey) {
			addrs = append(addrs, addrInfo)
		}
	}

	v.SortAddressesByAddressIndex(addrs...)
	v.SortAddressesByAddressType(addrs...)

	return addrs
}

func (*Vault) SortAddressesByPurpose(addrs ...AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.Purpose(), pathB.Purpose())
	})
}

func (*Vault) SortAddressesByAddressType(addrs ...AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressType(), pathB.AddressType())
	})
}

func (*Vault) SortAddressesByAddressIndex(addrs ...AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressIndex(), pathB.AddressIndex())
	})
}

func (v *Vault) IsEncrypted() bool {
	return v.Encrypter.IsEncrypted()
}

func (v *Vault) AddressCount() int {
	return len(v.Addresses)
}

func (v *Vault) AddressFromPath(p string) *AddressInfo {
	for _, addressInfo := range v.Addresses {
		if addressInfo.Path == p {
			return &addressInfo
		}
	}

	return nil
}

func (v *Vault) ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error {
	if v.IsNeutered() {
		return ErrNeutered
	}

	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return err
	}

	addressIndex := len(keyStore.ImportedKeys)
	pub := prv.PublicKeyNative()

	accAddr := pub.AccountAddress()
	if v.Contains(accAddr.String()) {
		return ErrAddressExists
	}

	valAddr := pub.ValidatorAddress()
	if v.Contains(valAddr.String()) {
		return ErrAddressExists
	}

	blsAccPathStr := addresspath.NewPath(
		H(PurposeImportPrivateKey),
		H(v.CoinType),
		H(crypto.AddressTypeBLSAccount),
		H(addressIndex)).String()

	blsValidatorPathStr := addresspath.NewPath(
		H(PurposeImportPrivateKey),
		H(v.CoinType),
		H(crypto.AddressTypeValidator),
		H(addressIndex)).String()

	importedPrvLabelCounter := (len(v.AllImportedPrivateKeysAddresses()) / 2) + 1
	v.Addresses[accAddr.String()] = AddressInfo{
		Address:   accAddr.String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported BLS Account Address %d", importedPrvLabelCounter),
		Path:      blsAccPathStr,
	}

	v.Addresses[valAddr.String()] = AddressInfo{
		Address:   valAddr.String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Validator Address %d", importedPrvLabelCounter),
		Path:      blsValidatorPathStr,
	}

	keyStore.ImportedKeys = append(keyStore.ImportedKeys, prv.String())

	err = v.encryptKeyStore(keyStore, password)
	if err != nil {
		return err
	}

	return nil
}

// PrivateKeys retrieves the private keys for the given addresses using the provided password.
func (v *Vault) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	if v.IsNeutered() {
		return nil, ErrNeutered
	}

	// Decrypt the key store once to avoid decrypting for each key.
	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return nil, err
	}
	mnemonicSeed := bip39.NewSeed(keyStore.MasterNode.Mnemonic, "")

	keys := make([]crypto.PrivateKey, len(addrs))
	for i, addr := range addrs {
		info := v.AddressInfo(addr)
		if info == nil {
			return nil, NewErrAddressNotFound(addr)
		}

		hdPath, err := addresspath.FromString(info.Path)
		if err != nil {
			return nil, err
		}

		if hdPath.CoinType() != H(v.CoinType) {
			return nil, ErrInvalidCoinType
		}

		switch hdPath.Purpose() {
		case H(PurposeBLS12381):
			prvKey, err := v.deriveBLSPrivateKey(mnemonicSeed, hdPath)
			if err != nil {
				return nil, err
			}
			keys[i] = prvKey
		case H(PurposeBIP44):
			prvKey, err := v.deriveEd25519PrivateKey(mnemonicSeed, hdPath)
			if err != nil {
				return nil, err
			}
			keys[i] = prvKey
		case H(PurposeImportPrivateKey):
			index := hdPath.AddressIndex() - blshdkeychain.HardenedKeyStart
			// TODO: Check if index is within the valid range.
			str := keyStore.ImportedKeys[index]
			prv, err := bls.PrivateKeyFromString(str)
			if err != nil {
				return nil, err
			}
			keys[i] = prv
		default:
			return nil, ErrUnsupportedPurpose
		}
	}

	return keys, nil
}

func (v *Vault) NewValidatorAddress(label string) (*AddressInfo, error) {
	ext, err := blshdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubValidator)
	if err != nil {
		return nil, err
	}
	index := v.Purposes.PurposeBLS.NextValidatorIndex
	ext, err = ext.DerivePath([]uint32{index})
	if err != nil {
		return nil, err
	}

	blsPubKey, err := bls.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := blsPubKey.ValidatorAddress().String()
	info := AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: blsPubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}
	v.Addresses[addr] = info
	v.Purposes.PurposeBLS.NextValidatorIndex++

	return &info, nil
}

func (v *Vault) NewBLSAccountAddress(label string) (*AddressInfo, error) {
	ext, err := blshdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubAccount)
	if err != nil {
		return nil, err
	}
	index := v.Purposes.PurposeBLS.NextAccountIndex
	ext, err = ext.DerivePath([]uint32{index})
	if err != nil {
		return nil, err
	}

	blsPubKey, err := bls.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := blsPubKey.AccountAddress().String()
	info := AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: blsPubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}
	v.Addresses[addr] = info
	v.Purposes.PurposeBLS.NextAccountIndex++

	return &info, nil
}

func (v *Vault) NewEd25519AccountAddress(label, password string) (*AddressInfo, error) {
	seed, err := v.MnemonicSeed(password)
	if err != nil {
		return nil, err
	}

	masterKey, err := ed25519hdkeychain.NewMaster(seed)
	if err != nil {
		return nil, err
	}

	index := v.Purposes.PurposeBIP44.NextEd25519Index
	ext, err := masterKey.DerivePath([]uint32{
		H(PurposeBIP44),
		H(v.CoinType),
		H(crypto.AddressTypeEd25519Account),
		H(index),
	})
	if err != nil {
		return nil, err
	}

	ed25519PubKey, err := ed25519.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := ed25519PubKey.AccountAddress().String()
	info := AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: ed25519PubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}
	v.Addresses[addr] = info
	v.Purposes.PurposeBIP44.NextEd25519Index++

	return &info, nil
}

// AddressInfo like it can return bls.PublicKey instead of string.
func (v *Vault) AddressInfo(addr string) *AddressInfo {
	info, ok := v.Addresses[addr]
	if !ok {
		return nil
	}

	return &info
}

func (v *Vault) Contains(addr string) bool {
	return v.AddressInfo(addr) != nil
}

func (v *Vault) Mnemonic(password string) (string, error) {
	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return "", err
	}

	return keyStore.MasterNode.Mnemonic, nil
}

func (v *Vault) MnemonicSeed(password string) ([]byte, error) {
	mnemonic, err := v.Mnemonic(password)
	if err != nil {
		return nil, err
	}
	seed := bip39.NewSeed(mnemonic, "")

	return seed, nil
}

func (v *Vault) decryptKeyStore(password string) (*keyStore, error) {
	if v.IsNeutered() {
		return nil, ErrNeutered
	}

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

func (*Vault) deriveBLSPrivateKey(mnemonicSeed []byte, path []uint32) (*bls.PrivateKey, error) {
	masterKey, err := blshdkeychain.NewMaster(mnemonicSeed, false)
	if err != nil {
		return nil, err
	}
	ext, err := masterKey.DerivePath(path)
	if err != nil {
		return nil, err
	}
	prvBytes, err := ext.RawPrivateKey()
	if err != nil {
		return nil, err
	}

	return bls.PrivateKeyFromBytes(prvBytes)
}

func (*Vault) deriveEd25519PrivateKey(mnemonicSeed []byte, path []uint32) (*ed25519.PrivateKey, error) {
	masterKey, err := ed25519hdkeychain.NewMaster(mnemonicSeed)
	if err != nil {
		return nil, err
	}
	ext, err := masterKey.DerivePath(path)
	if err != nil {
		return nil, err
	}
	prvBytes := ext.RawPrivateKey()

	return ed25519.PrivateKeyFromBytes(prvBytes)
}
