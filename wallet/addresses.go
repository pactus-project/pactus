package wallet

import (
	"cmp"
	"slices"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
)

type addresses struct {
	storage storage.IStorage
}

func newAddresses(storage storage.IStorage) addresses {
	return addresses{
		storage: storage,
	}
}

func (a *addresses) AddressInfo(addr string) (*types.AddressInfo, error) {
	return a.storage.AddressInfo(addr)
}

// listAddressConfig contains options for filtering addresses.
type listAddressConfig struct {
	addressTypes []crypto.AddressType
}

var defaultListAddressConfig = listAddressConfig{
	addressTypes: []crypto.AddressType{},
}

// ListAddressOption is a functional option for ListAddresses.
type ListAddressOption func(*listAddressConfig)

// WithAddressTypes filters addresses by the specified type.
func WithAddressTypes(addressTypes []crypto.AddressType) ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = addressTypes
	}
}

// WithAddressType filters addresses by the specified type.
func WithAddressType(addressType crypto.AddressType) ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = []crypto.AddressType{addressType}
	}
}

// OnlyValidatorAddresses filters to show only validator addresses.
func OnlyValidatorAddresses() ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = []crypto.AddressType{crypto.AddressTypeValidator}
	}
}

// OnlyAccountAddresses filters to show only account addresses (BLS and Ed25519).
func OnlyAccountAddresses() ListAddressOption {
	return func(cfg *listAddressConfig) {
		cfg.addressTypes = []crypto.AddressType{
			crypto.AddressTypeBLSAccount,
			crypto.AddressTypeEd25519Account,
		}
	}
}

func (a *addresses) ListAddresses(opts ...ListAddressOption) []types.AddressInfo {
	cfg := defaultListAddressConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	infos := make([]types.AddressInfo, 0)
	for _, info := range a.storage.AllAddresses() {
		if len(cfg.addressTypes) == 0 {
			infos = append(infos, info)

			continue
		}

		addr, err := crypto.AddressFromString(info.Address)
		if err != nil {
			return nil
		}

		for _, addrType := range cfg.addressTypes {
			if addr.Type() == addrType {
				infos = append(infos, info)

				break
			}
		}
	}

	a.sortAddressesByAddressIndex(infos...)
	a.sortAddressesByAddressType(infos...)
	a.sortAddressesByPurpose(infos...)

	return infos
}

func (*addresses) sortAddressesByPurpose(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.Purpose(), pathB.Purpose())
	})
}

func (*addresses) sortAddressesByAddressType(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressType(), pathB.AddressType())
	})
}

func (*addresses) sortAddressesByAddressIndex(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressIndex(), pathB.AddressIndex())
	})
}

// AddressCount returns the number of addresses inside the wallet.
func (a *addresses) AddressCount() int {
	return a.storage.AddressCount()
}

func (a *addresses) ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error {
	pub := prv.PublicKeyNative()
	accAddr := pub.AccountAddress()
	if a.HasAddress(accAddr.String()) {
		return ErrAddressExists
	}

	vault := a.storage.Vault()
	accInfo, valInfo, err := vault.ImportBLSPrivateKey(password, prv)
	if err != nil {
		return err
	}

	err = a.storage.InsertAddress(accInfo)
	if err != nil {
		return err
	}

	err = a.storage.InsertAddress(valInfo)
	if err != nil {
		return err
	}

	return a.storage.UpdateVault(vault)
}

func (a *addresses) ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) error {
	pub := prv.PublicKeyNative()

	accAddr := pub.AccountAddress()
	if a.HasAddress(accAddr.String()) {
		return ErrAddressExists
	}

	vault := a.storage.Vault()
	accInfo, err := vault.ImportEd25519PrivateKey(password, prv)
	if err != nil {
		return err
	}

	err = a.storage.InsertAddress(accInfo)
	if err != nil {
		return err
	}

	return a.storage.UpdateVault(vault)
}

func (a *addresses) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	keys, err := a.PrivateKeys(password, []string{addr})
	if err != nil {
		return nil, err
	}

	return keys[0], nil
}

func (a *addresses) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	paths := make([]addresspath.Path, len(addrs))
	for i, addr := range addrs {
		info, err := a.AddressInfo(addr)
		if err != nil {
			return nil, err
		}

		hdPath, err := addresspath.FromString(info.Path)
		if err != nil {
			return nil, err
		}

		paths[i] = hdPath
	}

	return a.storage.Vault().PrivateKeys(password, paths)
}

// newAddressConfig contains options for creating new addresses.
type newAddressConfig struct {
	password string
}

var defaultNewAddressConfig = newAddressConfig{
	password: "",
}

// NewAddressOption is a functional option for NewAddresa.
type NewAddressOption func(*newAddressConfig)

// WithPassword sets the password for address creation required for Ed25519 accounta.
func WithPassword(password string) NewAddressOption {
	return func(cfg *newAddressConfig) {
		cfg.password = password
	}
}

func (a *addresses) NewAddress(addressType crypto.AddressType, label string, opts ...NewAddressOption,
) (*types.AddressInfo, error) {
	cfg := defaultNewAddressConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	vault := a.storage.Vault()
	var info *types.AddressInfo
	var err error
	switch addressType {
	case crypto.AddressTypeValidator:
		info, err = vault.NewValidatorAddress(label)
	case crypto.AddressTypeBLSAccount:
		info, err = vault.NewBLSAccountAddress(label)
	case crypto.AddressTypeEd25519Account:
		info, err = vault.NewEd25519AccountAddress(label, cfg.password)
	case crypto.AddressTypeTreasury:
		return nil, ErrInvalidAddressType

	default:
		return nil, ErrInvalidAddressType
	}

	if err != nil {
		return nil, err
	}

	err = a.storage.InsertAddress(info)
	if err != nil {
		return nil, err
	}

	err = a.storage.UpdateVault(vault)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// NewBLSAccountAddress create a new BLS-based account address and
// associates it with the given label.
func (a *addresses) NewBLSAccountAddress(label string) (*types.AddressInfo, error) {
	return a.NewAddress(crypto.AddressTypeBLSAccount, label)
}

// NewEd25519AccountAddress create a new Ed25519-based account address and
// associates it with the given label.
// The password is required to access the master private key needed for address generation.
func (a *addresses) NewEd25519AccountAddress(label, password string) (*types.AddressInfo, error) {
	return a.NewAddress(crypto.AddressTypeEd25519Account, label, WithPassword(password))
}

// NewValidatorAddress creates a new BLS validator address and
// associates it with the given label.
func (a *addresses) NewValidatorAddress(label string) (*types.AddressInfo, error) {
	return a.NewAddress(crypto.AddressTypeValidator, label)
}

func (a *addresses) HasAddress(addr string) bool {
	return a.storage.HasAddress(addr)
}

// AddressLabel returns label of the given addresa.
func (a *addresses) AddressLabel(addr string) string {
	info, err := a.AddressInfo(addr)
	if err != nil {
		return ""
	}

	return info.Label
}

// SetAddressLabel updates the label of the given addresa.
func (a *addresses) SetAddressLabel(addr, label string) error {
	info, err := a.AddressInfo(addr)
	if err != nil {
		return err
	}

	info.Label = label

	return a.storage.UpdateAddress(info)
}
