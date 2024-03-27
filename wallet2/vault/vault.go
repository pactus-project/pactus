package vault

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/wallet2/addresspath"
	"github.com/pactus-project/pactus/wallet2/db"
	"github.com/pactus-project/pactus/wallet2/encrypter"
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
	db        db.DB
	Type      int                 `json:"type"`      // Wallet type. 1: Full keys, 2: Neutered
	CoinType  uint32              `json:"coin_type"` // Coin type: 21888 for Mainnet, 21777 for Testnet
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
	XPubValidator      string `json:"xpub_account"`         // Extended public key for validator: m/12381'/218
	XPubAccount        string `json:"xpub_validator"`       // Extended public key for account: m/12381'/21888/1'/0
	NextAccountIndex   uint32 `json:"next_account_index"`   // Index of next derived account
	NextValidatorIndex uint32 `json:"next_validator_index"` // Index of next derived validator
}

func CreateVaultFromMnemonic(mnemonic string, coinType uint32, database db.DB) (*Vault, error) {
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
		db:        database,
		Type:      TypeFull,
		CoinType:  coinType,
		Encrypter: enc,
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
		db:        v.db,
		Type:      TypeNeutered,
		CoinType:  v.CoinType,
		Encrypter: encrypter.NopeEncrypter(),
		KeyStore:  "",
		Purposes:  v.Purposes,
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
	info, err := v.db.GetAddressInfoByAddress(addr)
	if err != nil {
		return ""
	}

	return info.Label
}

func (v *Vault) SetLabel(address, label string) error {
	addr, err := v.db.GetAddressInfoByAddress(address)
	if err != nil {
		return NewErrAddressNotFound(address)
	}

	addr.Label = label
	err = v.db.UpdateAddressLabel(addr.Label, addr.Address)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Addresses() []db.AddressInfo {
	addrs, err := v.db.GetAllAddressInfos()
	if err != nil {
		return nil
	}

	v.sortAddressesByAddressIndex(addrs...)
	v.sortAddressesByAddressType(addrs...)
	v.sortAddressesByPurpose(addrs...)

	return addrs
}

func (v *Vault) AllValidatorAddresses() []db.AddressInfo {
	addrs, err := v.db.GetAllAddressInfos()
	if err != nil {
		return nil
	}

	result := make([]db.AddressInfo, 0, len(addrs)/2)
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

func (v *Vault) AllImportedPrivateKeyAddresses() []db.AddressInfo {
	addrs, err := v.db.GetAllAddressInfos()
	if err != nil {
		return nil
	}

	result := make([]db.AddressInfo, 0, len(addrs)/2)
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

func (v *Vault) IsEncrypted() bool {
	return v.Encrypter.IsEncrypted()
}

func (v *Vault) AddressFromPath(p string) *db.AddressInfo {
	addr, err := v.db.GetAddressByPath(p)
	if err != nil {
		return nil
	}

	return addr
}

func (v *Vault) ImportPrivateKey(password string, prv *bls.PrivateKey) error {
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

	importedPrvLabelCounter := (len(v.AllImportedPrivateKeyAddresses()) / 2) + 1
	rewardAddr := &db.AddressInfo{
		Address:   accAddr.String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Reward Address %d", importedPrvLabelCounter),
		Path:      blsAccPathStr,
	}

	if _, err = v.db.InsertAddressInfo(rewardAddr); err != nil {
		return err
	}

	validatorAddr := &db.AddressInfo{
		Address:   valAddr.String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Validator Address %d", importedPrvLabelCounter),
		Path:      blsValidatorPathStr,
	}

	if _, err = v.db.InsertAddressInfo(validatorAddr); err != nil {
		return err
	}

	keyStore.ImportedKeys = append(keyStore.ImportedKeys, prv.String())

	err = v.encryptKeyStore(keyStore, password)
	if err != nil {
		return err
	}

	return nil
}

//nolint:gocognit // refactor me and reduce the code complexity
func (v *Vault) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	if v.IsNeutered() {
		return nil, ErrNeutered
	}

	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return nil, err
	}

	keys := make([]crypto.PrivateKey, len(addrs))
	for i, addr := range addrs {
		info := v.Address(addr)
		if info == nil {
			return nil, NewErrAddressNotFound(addr)
		}

		path, err := addresspath.NewPathFromString(info.Path)
		if err != nil {
			return nil, err
		}

		if path.CoinType() != H(v.CoinType) {
			return nil, ErrInvalidCoinType
		}

		switch path.Purpose() {
		case H(PurposeBLS12381):
			seed, err := bip39.NewSeedWithErrorChecking(keyStore.MasterNode.Mnemonic, "")
			if err != nil {
				return nil, err
			}
			masterKey, err := hdkeychain.NewMaster(seed, false)
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

			prvKey, err := bls.PrivateKeyFromBytes(prvBytes)
			if err != nil {
				return nil, err
			}

			keys[i] = prvKey
		case H(PurposeImportPrivateKey):
			index := path.AddressIndex() - hdkeychain.HardenedKeyStart
			if index > uint32(len(keyStore.ImportedKeys)-1) {
				return nil, ErrIndexOutOfRange
			}
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

func (v *Vault) NewBLSAccountAddress(label string) (string, error) {
	ext, err := hdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubAccount)
	if err != nil {
		return "", err
	}
	index := v.Purposes.PurposeBLS.NextAccountIndex
	ext, err = ext.DerivePath([]uint32{index})
	if err != nil {
		return "", err
	}

	blsPubKey, err := bls.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return "", err
	}

	addr := blsPubKey.AccountAddress().String()
	address := &db.AddressInfo{
		Address: addr,
		Label:   label,
		Path:    addresspath.NewPath(ext.Path()...).String(),
	}
	if _, err := v.db.InsertAddressInfo(address); err != nil {
		return "", err
	}

	v.Purposes.PurposeBLS.NextAccountIndex++

	return addr, nil
}

func (v *Vault) NewValidatorAddress(label string) (string, error) {
	ext, err := hdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubValidator)
	if err != nil {
		return "", err
	}
	index := v.Purposes.PurposeBLS.NextValidatorIndex
	ext, err = ext.DerivePath([]uint32{index})
	if err != nil {
		return "", err
	}

	blsPubKey, err := bls.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return "", err
	}

	addr := blsPubKey.ValidatorAddress().String()
	address := &db.AddressInfo{
		Address: addr,
		Label:   label,
		Path:    addresspath.NewPath(ext.Path()...).String(),
	}
	if _, err := v.db.InsertAddressInfo(address); err != nil {
		return "", err
	}

	v.Purposes.PurposeBLS.NextValidatorIndex++

	return addr, nil
}

func (v *Vault) Address(address string) *db.AddressInfo {
	info, err := v.db.GetAddressInfoByAddress(address)
	if err != nil {
		return nil
	}

	path, err := addresspath.NewPathFromString(info.Path)
	if err != nil {
		return nil
	}

	if path.CoinType() != H(v.CoinType) {
		return nil
	}

	switch path.Purpose() {
	case H(PurposeBLS12381):
		addr, err := crypto.AddressFromString(info.Address)
		if err != nil {
			return nil
		}

		var xPub string
		if addr.IsAccountAddress() {
			xPub = v.Purposes.PurposeBLS.XPubAccount
		} else if addr.IsValidatorAddress() {
			xPub = v.Purposes.PurposeBLS.XPubValidator
		}

		ext, err := hdkeychain.NewKeyFromString(xPub)
		if err != nil {
			return nil
		}

		p, err := addresspath.NewPathFromString(info.Path)
		if err != nil {
			return nil
		}

		extendedKey, err := ext.Derive(p.AddressIndex())
		if err != nil {
			return nil
		}

		blsPubKey, err := bls.PublicKeyFromBytes(extendedKey.RawPublicKey())
		if err != nil {
			return nil
		}

		info.PublicKey = blsPubKey.String()
	case H(PurposeImportPrivateKey):
	default:
		return nil
	}

	return info
}

func (v *Vault) Contains(addr string) bool {
	return v.Address(addr) != nil
}

func (v *Vault) Mnemonic(password string) (string, error) {
	if v.IsNeutered() {
		return "", ErrNeutered
	}
	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return "", err
	}

	return keyStore.MasterNode.Mnemonic, nil
}

func (v *Vault) AddressCount() (int, error) {
	totalRecords, err := v.db.GetTotalRecords(db.AddressTable, db.EmptyQuery)

	return int(totalRecords), err
}

func (v *Vault) sortAddressesByAddressIndex(addrs ...db.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b db.AddressInfo) int {
		pathA, _ := addresspath.NewPathFromString(a.Path)
		pathB, _ := addresspath.NewPathFromString(b.Path)

		return cmp.Compare(pathA.AddressIndex(), pathB.AddressIndex())
	})
}

func (v *Vault) sortAddressesByAddressType(addrs ...db.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b db.AddressInfo) int {
		pathA, _ := addresspath.NewPathFromString(a.Path)
		pathB, _ := addresspath.NewPathFromString(b.Path)

		return cmp.Compare(pathA.AddressType(), pathB.AddressType())
	})
}

func (v *Vault) sortAddressesByPurpose(addrs ...db.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b db.AddressInfo) int {
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
