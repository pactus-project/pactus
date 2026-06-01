package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/crypto/ed25519"
	ed25519hdkeychain "github.com/pactus-project/pactus/crypto/ed25519/hdkeychain"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	secp256k1hdkeychain "github.com/pactus-project/pactus/crypto/secp256k1/hdkeychain"
	"github.com/pactus-project/pactus/util/bip39"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/types"
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

// VaultType represents the type of vault.
type VaultType int

const (
	TypeFull     VaultType = iota + 1 // Full vault with private keys.
	TypeNeutered                      // Neutered vault without private keys.
)

// String returns the string representation of the VaultType.
func (vt VaultType) String() string {
	switch vt {
	case TypeFull:
		return "Full"
	case TypeNeutered:
		return "Neutered"
	default:
		return "Unknown"
	}
}

// AddressGapLimit is the maximum number of consecutive inactive addresses before stopping recovery.
const AddressGapLimit = 8

type Vault struct {
	Type      VaultType            `json:"type"`      // Vault type: Full or Neutered
	CoinType  addresspath.CoinType `json:"coin_type"` // Coin type: 21888 for Mainnet, 21777 for Testnet
	Encrypter encrypter.Encrypter  `json:"encrypter"` // Encryption algorithm
	KeyStore  string               `json:"key_store"` // KeyStore that stores the secrets and encrypts using Encrypter
	Purposes  Purposes             `json:"purposes"`  // Contains Purposes of the vault
}

type keyStore struct {
	MasterNode   masterNode `json:"master_node"`   // HD Root Tree (Master node)
	ImportedKeys []string   `json:"imported_keys"` // Imported private keys
}

type masterNode struct {
	Mnemonic string `json:"seed,omitempty"` // Seed phrase or mnemonic (encrypted)
}

type Purposes struct {
	PurposeBLS   purposeBLS   `json:"purpose_bls"`   // BLS Purpose: m/12381'/21888'/1' or 2'
	PurposeBIP44 purposeBIP44 `json:"purpose_bip44"` // BIP44 Purpose: m/44'/21888'/3' or 4'
}

type purposeBLS struct {
	XPubValidator      string `json:"xpub_account"`         // Extended public key for account: m/12381'/21888'/1'/0
	XPubAccount        string `json:"xpub_validator"`       // Extended public key for validator: m/12381'/21888'/2'/0
	NextAccountIndex   uint32 `json:"next_account_index"`   // Index of next derived account
	NextValidatorIndex uint32 `json:"next_validator_index"` // Index of next derived validator
}

type purposeBIP44 struct {
	NextEd25519Index   uint32 `json:"next_ed25519_index"`   // Index of next Ed25519 derived account: m/44'/21888'/3'/0'
	NextSexp256k1Index uint32 `json:"next_secp256k1_index"` // Index of next Secp256k1 derived account: m/44'/21888'/4'/0
	XPubSecp256K1      string `json:"xpub_secp256k1"`       // Extended public key for account: m/44'/21888'/4'/0
}

func CreateVaultFromMnemonic(mnemonic string, coinType addresspath.CoinType) (*Vault, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	blsMasterKey, err := blshdkeychain.NewMaster(seed, false)
	if err != nil {
		return nil, err
	}
	enc := encrypter.NopeEncrypter()

	xPubBLSValidator, err := blsMasterKey.DerivePath([]uint32{
		addresspath.Harden(addresspath.PurposeBLS12381),
		addresspath.Harden(coinType),
		addresspath.Harden(crypto.AddressTypeValidator),
	})
	if err != nil {
		return nil, err
	}

	xPubBLSAccount, err := blsMasterKey.DerivePath([]uint32{
		addresspath.Harden(addresspath.PurposeBLS12381),
		addresspath.Harden(coinType),
		addresspath.Harden(crypto.AddressTypeBLSAccount),
	})
	if err != nil {
		return nil, err
	}

	secp256k1MasterKey, err := secp256k1hdkeychain.NewMaster(seed)
	if err != nil {
		return nil, err
	}
	xPubSecp256k1, err := secp256k1MasterKey.DerivePath([]uint32{
		addresspath.Harden(addresspath.PurposeBIP44),
		addresspath.Harden(coinType),
		addresspath.Harden(crypto.AddressTypeSecp256k1Account),
	})
	if err != nil {
		return nil, err
	}

	store := keyStore{
		MasterNode: masterNode{
			Mnemonic: mnemonic,
		},
		ImportedKeys: make([]string, 0),
	}

	storeDate, err := json.Marshal(store)
	if err != nil {
		return nil, err
	}

	return &Vault{
		Type:      TypeFull,
		CoinType:  coinType,
		Encrypter: enc,
		KeyStore:  string(storeDate),
		Purposes: Purposes{
			PurposeBLS: purposeBLS{
				XPubValidator: xPubBLSValidator.Neuter().String(),
				XPubAccount:   xPubBLSAccount.Neuter().String(),
			},
			PurposeBIP44: purposeBIP44{
				XPubSecp256K1: xPubSecp256k1.Neuter().String(),
			},
		},
	}, nil
}

func (v *Vault) Neuter() {
	v.Type = TypeNeutered
	v.Encrypter = encrypter.NopeEncrypter()
	v.KeyStore = ""
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

func (v *Vault) IsEncrypted() bool {
	return v.Encrypter.IsEncrypted()
}

func (v *Vault) ImportBLSPrivateKey(password string, prv *bls.PrivateKey) (
	valInfo *types.AddressInfo, accInfo *types.AddressInfo, err error,
) {
	if v.IsNeutered() {
		return nil, nil, ErrNeutered
	}

	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return nil, nil, err
	}

	pub := prv.PublicKeyNative()
	addressIndex := len(keyStore.ImportedKeys)

	blsAccPathStr := addresspath.NewPath(
		addresspath.Harden(addresspath.PurposeImportPrivateKey),
		addresspath.Harden(v.CoinType),
		addresspath.Harden(crypto.AddressTypeBLSAccount),
		addresspath.Harden(addressIndex)).String()

	blsValidatorPathStr := addresspath.NewPath(
		addresspath.Harden(addresspath.PurposeImportPrivateKey),
		addresspath.Harden(v.CoinType),
		addresspath.Harden(crypto.AddressTypeValidator),
		addresspath.Harden(addressIndex)).String()

	accInfo = &types.AddressInfo{
		Address:   pub.AccountAddress().String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Address %d", addressIndex+1),
		Path:      blsAccPathStr,
	}

	valInfo = &types.AddressInfo{
		Address:   pub.ValidatorAddress().String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Validator Address %d", addressIndex+1),
		Path:      blsValidatorPathStr,
	}

	keyStore.ImportedKeys = append(keyStore.ImportedKeys, prv.String())

	err = v.encryptKeyStore(keyStore, password)
	if err != nil {
		return nil, nil, err
	}

	return valInfo, accInfo, nil
}

func (v *Vault) ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) (*types.AddressInfo, error) {
	if v.IsNeutered() {
		return nil, ErrNeutered
	}

	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return nil, err
	}

	addressIndex := len(keyStore.ImportedKeys)
	pub := prv.PublicKeyNative()

	accPathStr := addresspath.NewPath(
		addresspath.Harden(addresspath.PurposeImportPrivateKey),
		addresspath.Harden(v.CoinType),
		addresspath.Harden(crypto.AddressTypeEd25519Account),
		addresspath.Harden(addressIndex)).String()

	accInfo := &types.AddressInfo{
		Address:   pub.AccountAddress().String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Address %d", addressIndex+1),
		Path:      accPathStr,
	}

	keyStore.ImportedKeys = append(keyStore.ImportedKeys, prv.String())

	err = v.encryptKeyStore(keyStore, password)
	if err != nil {
		return nil, err
	}

	return accInfo, nil
}

func (v *Vault) ImportSecp256k1PrivateKey(password string, prv *secp256k1.PrivateKey) (*types.AddressInfo, error) {
	if v.IsNeutered() {
		return nil, ErrNeutered
	}

	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return nil, err
	}

	addressIndex := len(keyStore.ImportedKeys)
	pub := prv.PublicKeyNative()

	accPathStr := addresspath.NewPath(
		addresspath.Harden(addresspath.PurposeImportPrivateKey),
		addresspath.Harden(v.CoinType),
		addresspath.Harden(crypto.AddressTypeSecp256k1Account),
		addresspath.Harden(addressIndex)).String()

	accInfo := &types.AddressInfo{
		Address:   pub.AccountAddress().String(),
		PublicKey: pub.String(),
		Label:     fmt.Sprintf("Imported Address %d", addressIndex+1),
		Path:      accPathStr,
	}

	keyStore.ImportedKeys = append(keyStore.ImportedKeys, prv.String())

	err = v.encryptKeyStore(keyStore, password)
	if err != nil {
		return nil, err
	}

	return accInfo, nil
}

// PrivateKeys retrieves the private keys for the given addresses using the provided password.
func (v *Vault) PrivateKeys(password string, paths []addresspath.Path) ([]crypto.PrivateKey, error) {
	if v.IsNeutered() {
		return nil, ErrNeutered
	}

	// Decrypt the key store once to avoid decrypting for each key.
	keyStore, err := v.decryptKeyStore(password)
	if err != nil {
		return nil, err
	}
	seed := bip39.NewSeed(keyStore.MasterNode.Mnemonic, "")

	keys := make([]crypto.PrivateKey, len(paths))
	for i, path := range paths {
		switch path.Purpose() {
		case addresspath.PurposeBLS12381:
			prvKey, err := v.deriveBLSPrivateKey(seed, path)
			if err != nil {
				return nil, err
			}
			keys[i] = prvKey
		case addresspath.PurposeBIP44:
			switch path.AddressType() {
			case crypto.AddressTypeEd25519Account:
				prvKey, err := v.deriveEd25519PrivateKey(seed, path)
				if err != nil {
					return nil, err
				}
				keys[i] = prvKey

			case crypto.AddressTypeSecp256k1Account:
				prvKey, err := v.deriveSecp256k1PrivateKey(seed, path)
				if err != nil {
					return nil, err
				}
				keys[i] = prvKey

			case crypto.AddressTypeTreasury,
				crypto.AddressTypeValidator,
				crypto.AddressTypeBLSAccount:
				// Unsupported Path
			}
		case addresspath.PurposeImportPrivateKey:
			index := addresspath.UnHarden(path.AddressIndex())
			str := keyStore.ImportedKeys[index]

			var prv crypto.PrivateKey
			switch uint32(path.AddressType()) {
			case uint32(crypto.AddressTypeValidator),
				uint32(crypto.AddressTypeBLSAccount):
				prv, err = bls.PrivateKeyFromString(str)
				if err != nil {
					return nil, err
				}

			case uint32(crypto.AddressTypeEd25519Account):
				prv, err = ed25519.PrivateKeyFromString(str)
				if err != nil {
					return nil, err
				}

			case uint32(crypto.AddressTypeSecp256k1Account):
				prv, err = secp256k1.PrivateKeyFromString(str)
				if err != nil {
					return nil, err
				}
			}

			keys[i] = prv
		default:
			return nil, ErrUnsupportedPurpose
		}
	}

	return keys, nil
}

func (v *Vault) NewValidatorAddress(label string) (*types.AddressInfo, error) {
	ext, err := blshdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubValidator)
	if err != nil {
		return nil, err
	}
	index := v.Purposes.PurposeBLS.NextValidatorIndex
	ext, err = ext.DerivePath([]uint32{index})
	if err != nil {
		return nil, err
	}

	pubKey, err := bls.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := pubKey.ValidatorAddress().String()
	info := types.AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: pubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}

	v.Purposes.PurposeBLS.NextValidatorIndex++

	return &info, nil
}

func (v *Vault) NewBLSAccountAddress(label string) (*types.AddressInfo, error) {
	ext, err := blshdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubAccount)
	if err != nil {
		return nil, err
	}
	index := v.Purposes.PurposeBLS.NextAccountIndex
	info, err := v.deriveBLSAccountAddressAt(ext, index, label)
	if err != nil {
		return nil, err
	}

	v.Purposes.PurposeBLS.NextAccountIndex++

	return info, nil
}

func (*Vault) deriveBLSAccountAddressAt(ext *blshdkeychain.ExtendedKey,
	index uint32, label string,
) (*types.AddressInfo, error) {
	ext, err := ext.DerivePath([]uint32{index})
	if err != nil {
		return nil, err
	}

	pubKey, err := bls.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := pubKey.AccountAddress().String()
	info := types.AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: pubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}

	return &info, nil
}

func (v *Vault) NewEd25519AccountAddress(label, password string) (*types.AddressInfo, error) {
	seed, err := v.MnemonicSeed(password)
	if err != nil {
		return nil, err
	}

	masterKey, err := ed25519hdkeychain.NewMaster(seed)
	if err != nil {
		return nil, err
	}

	index := v.Purposes.PurposeBIP44.NextEd25519Index
	info, err := v.deriveEd25519AccountAddressAt(masterKey, index, label)
	if err != nil {
		return nil, err
	}
	v.Purposes.PurposeBIP44.NextEd25519Index++

	return info, nil
}

func (v *Vault) deriveEd25519AccountAddressAt(masterKey *ed25519hdkeychain.ExtendedKey,
	index uint32, label string,
) (*types.AddressInfo, error) {
	ext, err := masterKey.DerivePath([]uint32{
		addresspath.Harden(addresspath.PurposeBIP44),
		addresspath.Harden(v.CoinType),
		addresspath.Harden(crypto.AddressTypeEd25519Account),
		addresspath.Harden(index),
	})
	if err != nil {
		return nil, err
	}

	pubKey, err := ed25519.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := pubKey.AccountAddress().String()
	info := types.AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: pubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}

	return &info, nil
}

func (v *Vault) NewSecp256k1AccountAddress(label, password string) (*types.AddressInfo, error) {
	seed, err := v.MnemonicSeed(password)
	if err != nil {
		return nil, err
	}

	masterKey, err := secp256k1hdkeychain.NewMaster(seed)
	if err != nil {
		return nil, err
	}

	index := v.Purposes.PurposeBIP44.NextSexp256k1Index
	info, err := v.deriveSecp256k1AccountAddressAt(masterKey, index, label)
	if err != nil {
		return nil, err
	}
	v.Purposes.PurposeBIP44.NextSexp256k1Index++

	return info, nil
}

func (v *Vault) deriveSecp256k1AccountAddressAt(masterKey *secp256k1hdkeychain.ExtendedKey,
	index uint32, label string,
) (*types.AddressInfo, error) {
	ext, err := masterKey.DerivePath([]uint32{
		addresspath.Harden(addresspath.PurposeBIP44),
		addresspath.Harden(v.CoinType),
		addresspath.Harden(crypto.AddressTypeSecp256k1Account),
		addresspath.Harden(index),
	})
	if err != nil {
		return nil, err
	}

	pubKey, err := secp256k1.PublicKeyFromBytes(ext.RawPublicKey())
	if err != nil {
		return nil, err
	}

	addr := pubKey.AccountAddress().String()
	info := types.AddressInfo{
		Address:   addr,
		Label:     label,
		PublicKey: pubKey.String(),
		Path:      addresspath.NewPath(ext.Path()...).String(),
	}

	return &info, nil
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

func (*Vault) deriveSecp256k1PrivateKey(mnemonicSeed []byte, path []uint32) (*secp256k1.PrivateKey, error) {
	masterKey, err := secp256k1hdkeychain.NewMaster(mnemonicSeed)
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

	return secp256k1.PrivateKeyFromBytes(prvBytes)
}

// RecoverAddresses automatically recovers used addresses when restoring a wallet from a mnemonic phrase.
// This implementation follows PIP-41 specification for address recovery.
//
// The function recovers both BLS and Ed25519 account addresses, with Ed25519 being the default
// address type for recovery when the wallet is empty.
//
// An address is considered active if its public key is stored in the blockchain database.
// The hasActivity function should return true if the address has been used before.
//
// Limitation: Users cannot automatically recover a used address if it is separated by more than 8
// inactive or empty addresses. In this case, manual address creation is required.
//
// Reference: https://pips.pactus.org/PIPs/pip-41
func (v *Vault) RecoverAddresses(ctx context.Context, password string,
	hasActivity func(addr string) (bool, error),
) ([]types.AddressInfo, error) {
	recovered1, err := v.recoverBLSAccountAddresses(ctx, hasActivity)
	if err != nil {
		return nil, err
	}

	recovered2, err := v.recoverEd25519AccountAddresses(ctx, password, hasActivity)
	if err != nil {
		return nil, err
	}

	recovered1 = append(recovered1, recovered2...)

	return recovered1, nil
}

// scanRecoveredCount scans derived addresses until the gap limit is exceeded and
// returns how many addresses should be recovered according to PIP-41.
func (*Vault) scanRecoveredCount(
	ctx context.Context,
	startIndex uint32,
	deriveAt func(index uint32) (*types.AddressInfo, error),
	hasActivity func(addr string) (bool, error),
) (int, error) {
	recoveredCount := 0
	inactiveCount := 1
	currentIndex := startIndex

	info, err := deriveAt(currentIndex)
	if err != nil {
		return 0, err
	}

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		isActive, err := hasActivity(info.Address)
		if err != nil {
			return 0, err
		}

		if isActive {
			recoveredCount += inactiveCount
			inactiveCount = 1
		} else {
			inactiveCount++
			if inactiveCount > AddressGapLimit {
				break
			}
		}

		currentIndex++
		info, err = deriveAt(currentIndex)
		if err != nil {
			return 0, err
		}
	}

	return recoveredCount, nil
}

// recoverBLSAccountAddresses recovers BLS account addresses following the PIP-41 specification.
func (v *Vault) recoverBLSAccountAddresses(ctx context.Context,
	hasActivity func(addrs string) (bool, error),
) ([]types.AddressInfo, error) {
	ext, err := blshdkeychain.NewKeyFromString(v.Purposes.PurposeBLS.XPubAccount)
	if err != nil {
		return nil, err
	}

	recoveredCount, err := v.scanRecoveredCount(
		ctx,
		v.Purposes.PurposeBLS.NextAccountIndex,
		func(index uint32) (*types.AddressInfo, error) {
			return v.deriveBLSAccountAddressAt(ext, index, "")
		},
		hasActivity,
	)
	if err != nil {
		return nil, err
	}

	recovered := make([]types.AddressInfo, 0, recoveredCount)
	for i := 0; i < recoveredCount; i++ {
		info, _ := v.NewBLSAccountAddress(fmt.Sprintf("BLS Account Address %d", i))
		recovered = append(recovered, *info)
	}

	return recovered, nil
}

// recoverEd25519AccountAddresses recovers Ed25519 account addresses following the PIP-41 specification.
func (v *Vault) recoverEd25519AccountAddresses(ctx context.Context, password string,
	hasActivity func(addrs string) (bool, error),
) ([]types.AddressInfo, error) {
	seed, err := v.MnemonicSeed(password)
	if err != nil {
		return nil, err
	}

	masterKey, err := ed25519hdkeychain.NewMaster(seed)
	if err != nil {
		return nil, err
	}

	recoveredCount, err := v.scanRecoveredCount(
		ctx,
		v.Purposes.PurposeBIP44.NextEd25519Index,
		func(index uint32) (*types.AddressInfo, error) {
			return v.deriveEd25519AccountAddressAt(masterKey, index, "")
		},
		hasActivity,
	)
	if err != nil {
		return nil, err
	}

	recovered := make([]types.AddressInfo, 0, recoveredCount)
	for i := 0; i < recoveredCount; i++ {
		info, _ := v.NewEd25519AccountAddress(fmt.Sprintf("Ed25519 Account Address %d", i), password)
		recovered = append(recovered, *info)
	}

	return recovered, nil
}
