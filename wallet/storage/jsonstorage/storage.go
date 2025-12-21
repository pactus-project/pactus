package jsonstorage

import (
	"cmp"
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Storage struct {
	path  string
	store Store
}

func Create(path string, version int, network genesis.ChainType, vault vault.Vault) (*Storage, error) {
	store := Store{
		Version:   version,
		UUID:      uuid.New(),
		CreatedAt: time.Now().Round(time.Second).UTC(),
		Network:   network,
		Vault:     vault,
	}

	return &Storage{
		path:  path,
		store: store,
	}, nil
}

func Open(path string) (*Storage, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	if err := store.ValidateCRC(); err != nil {
		return nil, err
	}

	return &Storage{
		path:  path,
		store: store,
	}, nil
}

func (s *Storage) Save() error {
	s.store.VaultCRC = s.store.calcVaultCRC()

	data, err := json.MarshalIndent(s, "  ", "  ")
	if err != nil {
		return err
	}

	return util.WriteFile(s.path, data)
}

func (s *Storage) WalletInfo() *types.WalletInfo {
	return &types.WalletInfo{
		Version:    s.store.Version,
		Network:    s.store.Network.String(),
		DefaultFee: s.store.DefaultFee,
		UUID:       s.store.UUID.String(),
		Encrypted:  s.IsEncrypted(),
		CreatedAt:  s.store.CreatedAt,
	}
}

func (s *Storage) Version() int {
	return s.store.Version
}

func (s *Storage) Network() genesis.ChainType {
	return s.store.Network
}

func (s *Storage) Path() string {
	return s.path
}

func (s *Storage) IsEncrypted() bool {
	return s.store.Vault.IsEncrypted()
}

func (s *Storage) Upgrade() error {
	return s.store.Upgrade()
}

func (s *Storage) Mnemonic(password string) (string, error) {
	return s.store.Vault.Mnemonic(password)
}

func (s *Storage) UpdatePassword(oldPassword, newPassword string, opts ...encrypter.Option) error {
	return s.UpdatePassword(oldPassword, newPassword, opts...)
}

func (s *Storage) RecoverAddresses(ctx context.Context, password string,
	eventFunc func(addr string) (bool, error),
) error {
	return s.store.Vault.RecoverAddresses(ctx, password, eventFunc)
}

func (s *Storage) AddressCount() int {
	return len(s.store.Addresses)
}

func (s *Storage) HasAddress(addr string) bool {
	_, ok := s.store.Addresses[addr]

	return ok
}

func (s *Storage) AddressInfo(addr string) *types.AddressInfo {
	info, ok := s.store.Addresses[addr]
	if !ok {
		return nil
	}

	return &info
}

func (s *Storage) AddressLabel(addr string) string {
	info, ok := s.store.Addresses[addr]
	if !ok {
		return ""
	}

	return info.Label
}

func (s *Storage) SetAddressLabel(addr, label string) error {
	info, ok := s.store.Addresses[addr]
	if !ok {
		return NewErrAddressNotFound(addr)
	}

	info.Label = label
	s.store.Addresses[addr] = info

	return nil
}

func (s *Storage) ListAddresses(opts ...types.ListAddressOption) []types.AddressInfo {
	infos := make([]types.AddressInfo, 0, s.AddressCount())
	for _, addrInfo := range s.store.Addresses {
		infos = append(infos, addrInfo)
	}

	s.sortAddressesByAddressIndex(infos...)
	s.sortAddressesByAddressType(infos...)
	s.sortAddressesByPurpose(infos...)

	return infos
}

func (*Storage) sortAddressesByPurpose(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.Purpose(), pathB.Purpose())
	})
}

func (*Storage) sortAddressesByAddressType(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressType(), pathB.AddressType())
	})
}

func (*Storage) sortAddressesByAddressIndex(addrs ...types.AddressInfo) {
	slices.SortStableFunc(addrs, func(a, b types.AddressInfo) int {
		pathA, _ := addresspath.FromString(a.Path)
		pathB, _ := addresspath.FromString(b.Path)

		return cmp.Compare(pathA.AddressIndex(), pathB.AddressIndex())
	})
}

func (s *Storage) Neuter(path string) error {
	// Clone the store
	cloned := Storage{
		path:  path,
		store: s.store,
	}

	cloned.store.Vault.Neuter()

	err := cloned.Save()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DefaultFee() amount.Amount {
	return s.store.DefaultFee
}

func (s *Storage) SetDefaultFee(fee amount.Amount) {
	s.store.DefaultFee = fee
}

func (s *Storage) ImportBLSPrivateKey(password string, prv *bls.PrivateKey) error {
	pub := prv.PublicKeyNative()

	accAddr := pub.AccountAddress()
	if s.HasAddress(accAddr.String()) {
		return ErrAddressExists
	}

	valAddr := pub.ValidatorAddress()
	if s.HasAddress(valAddr.String()) {
		return ErrAddressExists
	}

	accInfo, valInfo, err := s.store.Vault.ImportBLSPrivateKey(password, prv)
	if err != nil {
		return err
	}

	s.store.Addresses[accInfo.Address] = *accInfo
	s.store.Addresses[valInfo.Address] = *valInfo

	return nil
}

func (s *Storage) ImportEd25519PrivateKey(password string, prv *ed25519.PrivateKey) error {
	pub := prv.PublicKeyNative()

	accAddr := pub.AccountAddress()
	if s.HasAddress(accAddr.String()) {
		return ErrAddressExists
	}

	accInfo, err := s.store.Vault.ImportEd25519PrivateKey(password, prv)
	if err != nil {
		return err
	}

	s.store.Addresses[accInfo.Address] = *accInfo

	return nil
}

func (s *Storage) PrivateKeys(password string, addrs []string) ([]crypto.PrivateKey, error) {
	paths := make([]addresspath.Path, len(addrs))
	for i, addr := range addrs {
		info := s.AddressInfo(addr)
		if info == nil {
			return nil, NewErrAddressNotFound(addr)
		}

		hdPath, err := addresspath.FromString(info.Path)
		if err != nil {
			return nil, err
		}

		paths[i] = hdPath
	}

	return s.store.Vault.PrivateKeys(password, paths)
}

func (s *Storage) NewValidatorAddress(label string) (*types.AddressInfo, error) {
	info, err := s.store.Vault.NewValidatorAddress(label)
	if err != nil {
		return nil, err
	}

	s.store.Addresses[info.Address] = *info

	return info, nil
}

func (s *Storage) NewBLSAccountAddress(label string) (*types.AddressInfo, error) {
	info, err := s.store.Vault.NewBLSAccountAddress(label)
	if err != nil {
		return nil, err
	}

	s.store.Addresses[info.Address] = *info

	return info, nil
}

func (s *Storage) NewEd25519AccountAddress(label, password string) (*types.AddressInfo, error) {
	info, err := s.store.Vault.NewEd25519AccountAddress(label, password)
	if err != nil {
		return nil, err
	}

	s.store.Addresses[info.Address] = *info

	return info, nil
}

func (s *Storage) AddPending(addr string, amt amount.Amount, txID tx.ID, data []byte) {
	s.store.History.addPending(addr, amt, txID, data)
}

func (s *Storage) HasTransaction(id string) bool {
	return s.store.History.hasTransaction(id)
}

func (s *Storage) GetAddrHistory(addr string) []types.HistoryInfo {
	return s.store.History.getAddrHistory(addr)
}

func (s *Storage) AddActivity(addr string, amt amount.Amount, trx *pactus.GetTransactionResponse) {
	s.store.History.addActivity(addr, amt, trx)
}
