package vault

import (
	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/wallet/encrypter"
	"github.com/zarbchain/zarb-go/wallet/hdkeychain"
)

//
// Deterministic Account Hierarchy
//
// Specification
//
// We define the following 4 levels in BIP32 path:
//
// m / purpose' / coin_type' / account / use
//
// Where:
//   Apostrophe in the path indicates that BIP32 hardened derivation is used.
//   `m` Denotes the master node (or root) of the tree
//   `/` Separates the tree into depths, thus i / j signifies that j is a child of i
//   `purpose` is set to 12381 which is the name of the new curve (BLS12-381).
//   `coin_type` is set 21888 for Mainnet, 21777 for Testnet
//   `account` is a field that provides the ability for a user to have distinct sets of keys.
//   `use` is set to zero for now.
//
// References:
// BIP-44: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
// EIP-2334: https://eips.ethereum.org/EIPS/eip-2334

type AddressInfo struct {
	Address  string
	Pub      string
	Label    string
	Path     string
	Imported bool
	// for ;oca; use only
	importedIndex int
}

const PurposeBLS12381 = uint32(12381)

type Vault struct {
	Encrypter    encrypter.Encrypter `json:"encrypter"` //
	Keystore     keystore            `json:"keystore"`  //
	ImportedKeys []imported          `json:"imported"`  // Imported private keys
	Labels       map[string]string   `json:"labels"`    //
}

type imported struct {
	Addr string `json:"addr"` // Address
	Pub  string `json:"pub"`  // Public key
	Prv  string `json:"prv"`  // Private key (encrypted)
}

type keystore struct {
	CoinType uint32             `json:"coin_type"`      // Coin type: 21888 for Mainnet, 21777 for Testnet
	Mnemonic string             `json:"seed,omitempty"` // Seed phrase or mnemonic (encrypted)
	Purposes map[uint32]purpose `json:"purpose"`        // Purposes: 12381 for BLS signature
}

type purpose struct {
	XPub      string   `json:"xpub"`      // Extended public key
	Addresses []string `json:"addresses"` // Derived addresses
}

func CreateVaultFromMnemonic(mnemonic string, coinType uint32) (*Vault, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed)
	if err != nil {
		return nil, err
	}
	encrypter := encrypter.NopeEncrypter()

	purposeKey, err := masterKey.DerivePath([]uint32{
		12381 + hdkeychain.HardenedKeyStart,
		coinType + hdkeychain.HardenedKeyStart})

	blsPurpose := purpose{
		XPub:      purposeKey.BLSPublicKey().String(),
		Addresses: make([]string, 20),
	}
	for i, _ := range blsPurpose.Addresses {
		ext, err := purposeKey.DerivePath([]uint32{uint32(i), 0})
		if err != nil {
			return nil, err
		}

		blsPurpose.Addresses[i] = ext.Address().String()
	}

	return &Vault{
		Encrypter: encrypter,
		Keystore: keystore{
			CoinType: coinType,
			Mnemonic: mnemonic,
			Purposes: map[uint32]purpose{
				PurposeBLS12381: blsPurpose,
			},
		},
		Labels:       map[string]string{},
		ImportedKeys: []imported{},
	}, nil
}

func (v *Vault) Neuter() *Vault {
	blsPurpose := v.Keystore.Purposes[PurposeBLS12381]
	blsPurposeClone := purpose{
		XPub:      blsPurpose.XPub,
		Addresses: make([]string, len(blsPurpose.Addresses)),
	}
	for i, addr := range blsPurpose.Addresses {
		blsPurposeClone.Addresses[i] = addr
	}

	neutered := &Vault{
		Encrypter: encrypter.NopeEncrypter(),
		Keystore: keystore{
			CoinType: v.Keystore.CoinType,
			Purposes: map[uint32]purpose{
				PurposeBLS12381: blsPurposeClone,
			},
		},
		Labels:       map[string]string{},
		ImportedKeys: []imported{},
	}
	return neutered
}

func (v *Vault) IsNeutered() bool {
	return v.Keystore.Mnemonic == ""
}

func (v *Vault) UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error {
	if v.IsNeutered() {
		return ErrNeutered
	}

	oldEncrypter := v.Encrypter
	newEncrypter := encrypter.NopeEncrypter()
	if newPassword != "" {
		newEncrypter = encrypter.DefaultEncrypter(opts...)
	}

	// Updating mnemonic
	mnemonic, err := oldEncrypter.Decrypt(v.Keystore.Mnemonic, oldPassword)
	if err != nil {
		return err
	}
	v.Keystore.Mnemonic, err = newEncrypter.Encrypt(mnemonic, newPassword)
	util.ExitOnErr(err)

	// Updating imported private keys
	for i, key := range v.ImportedKeys {
		prv, err := oldEncrypter.Decrypt(key.Prv, oldPassword)
		util.ExitOnErr(err)

		v.ImportedKeys[i].Prv, err = newEncrypter.Encrypt(prv, newPassword)
		util.ExitOnErr(err)
	}
	v.Encrypter = newEncrypter
	return nil
}

func (v *Vault) Label(addr string) string {
	lbl, ok := v.Labels[addr]
	if !ok {
		return ""
	}
	return lbl
}

func (v *Vault) SetLabel(addr, label string) error {
	_, ok := v.Labels[addr]
	if !ok {
		return NewErrAddressNotFound(addr)
	}

	v.Labels[addr] = label
	return nil
}

func (v *Vault) AddressLabels() []AddressInfo {
	addrs := make([]AddressInfo, 0, v.AddressCount())

	for _, p := range v.Keystore.Purposes {
		for _, a := range p.Addresses {
			addrs = append(addrs, AddressInfo{
				Address:  a,
				Label:    v.Label(a),
				Imported: false,
			})
		}
	}

	for _, i := range v.ImportedKeys {
		addrs = append(addrs, AddressInfo{
			Address:  i.Addr,
			Label:    v.Label(i.Addr),
			Imported: true,
		})
	}
	return addrs
}

func (v *Vault) IsEncrypted() bool {
	return v.Encrypter.IsEncrypted()
}

func (v *Vault) AddressCount() int {
	count := len(v.ImportedKeys)
	for _, p := range v.Keystore.Purposes {
		count += len(p.Addresses)
	}
	return count
}

func (v *Vault) ImportPrivateKey(password string, prvStr string) error {
	if v.IsNeutered() {
		return ErrNeutered
	}

	prv, err := bls.PrivateKeyFromString(prvStr)
	if err != nil {
		return err
	}

	addr := prv.PublicKey().Address().String()
	if v.Contains(addr) {
		return ErrAddressExists
	}
	// Decrypt seed to make sure the password is correct
	_, err = v.Mnemonic(password)
	if err != nil {
		return err
	}

	encPrv, err := v.Encrypter.Encrypt(prv.String(), password)
	if err != nil {
		return err
	}
	v.ImportedKeys = append(v.ImportedKeys, imported{
		Prv:  encPrv,
		Pub:  prv.PublicKey().String(),
		Addr: prv.PublicKey().Address().String(),
	})

	return nil
}

func (v *Vault) PrivateKey(password, addr string) (string, error) {
	info := v.GetAddressInfo(addr)
	if info == nil {
		return "", NewErrAddressNotFound(addr)
	}

	if info.Imported {
		ct := v.ImportedKeys[info.importedIndex].Prv
		prv, err := v.Encrypter.Decrypt(ct, password)
		if err != nil {
			return "", err
		}
		return prv, nil
	} else {
		mnemonic, err := v.Mnemonic(password)
		if err != nil {
			return "", err
		}
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
		if err != nil {
			return "", err
		}
		masterKey, err := hdkeychain.NewMaster(seed)
		if err != nil {
			return "", err
		}
		path, _ := stringToDerivePath(info.Path)
		ext, err := masterKey.DerivePath(path)
		if err != nil {
			return "", err
		}
		prv, err := ext.BLSPrivateKey()
		if err != nil {
			return "", err
		}
		return prv.String(), nil
	}
}

func (v *Vault) DeriveNewAddress(password, label string, purpose uint32) (string, error) {
	p, ok := v.Keystore.Purposes[purpose]
	if ok {
		ext, err := hdkeychain.NewKeyFromString(p.XPub)
		if err != nil {
			return "", err
		}
		index := uint32(len(p.Addresses))
		ext, err = ext.DerivePath([]uint32{index, 0})
		if err != nil {
			return "", err
		}
		return ext.Address().String(), nil
	}

	return "", nil
}

func (v *Vault) GetAddressInfo(addr string) *AddressInfo {
	for _, p := range v.Keystore.Purposes {
		for i, a := range p.Addresses {
			if a == addr {
				pubKey, err := hdkeychain.NewKeyFromString(p.XPub)
				if err != nil {
					return nil
				}
				ext, err := pubKey.DerivePath([]uint32{uint32(i), 0})
				if err != nil {
					return nil
				}

				return &AddressInfo{
					Address: addr,
					Label:   v.Label(addr),
					Pub:     ext.BLSPublicKey().String(),
					Path:    derivePathToString(ext.Path()),
				}
			}
		}
	}

	for i, k := range v.ImportedKeys {
		if k.Addr == addr {
			return &AddressInfo{
				Address:       addr,
				Label:         v.Label(addr),
				Pub:           k.Pub,
				Path:          "",
				Imported:      true,
				importedIndex: i,
			}
		}
	}

	return nil
}

func (v *Vault) Contains(addr string) bool {
	return v.GetAddressInfo(addr) != nil
}

func (v *Vault) Mnemonic(password string) (string, error) {
	dec, err := v.Encrypter.Decrypt(v.Keystore.Mnemonic, password)
	if err != nil {
		return "", err
	}
	return dec, nil
}
