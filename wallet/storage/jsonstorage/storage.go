package jsonstorage

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

type Storage struct {
	path  string
	store store
}

func Create(path string, network genesis.ChainType, vault *vault.Vault) (*Storage, error) {
	store := store{
		Version:    VersionLatest,
		UUID:       uuid.New(),
		CreatedAt:  util.RoundNow(1),
		Network:    network,
		DefaultFee: amount.Amount(10_000_000),
		Vault:      *vault,
		Addresses:  make(map[string]types.AddressInfo),
	}

	if err := store.Save(path); err != nil {
		return nil, err
	}

	return Open(path)
}

func Open(path string) (*Storage, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var store store
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

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) save() error {
	return s.store.Save(s.path)
}

func (s *Storage) WalletInfo() *types.WalletInfo {
	return &types.WalletInfo{
		Path:       s.path,
		Driver:     "JSON (legacy)",
		Version:    s.store.Version,
		Network:    s.store.Network,
		DefaultFee: s.store.DefaultFee,
		UUID:       s.store.UUID.String(),
		Encrypted:  s.store.Vault.IsEncrypted(),
		Neutered:   s.store.Vault.IsNeutered(),
		CreatedAt:  s.store.CreatedAt,
	}
}

func (s *Storage) Vault() *vault.Vault {
	return &s.store.Vault
}

func (s *Storage) UpdateVault(vault *vault.Vault) error {
	s.store.Vault = *vault

	return s.save()
}

func (s *Storage) SetDefaultFee(fee amount.Amount) error {
	s.store.DefaultFee = fee

	return s.save()
}

func (s *Storage) AllAddresses() []types.AddressInfo {
	addrs := make([]types.AddressInfo, 0, len(s.store.Addresses))
	for _, info := range s.store.Addresses {
		addrs = append(addrs, info)
	}

	return addrs
}

func (s *Storage) AddressInfo(address string) (*types.AddressInfo, error) {
	info, exists := s.store.Addresses[address]
	if !exists {
		return nil, storage.ErrNotFound
	}

	return &info, nil
}

func (s *Storage) InsertAddress(info *types.AddressInfo) error {
	s.store.Addresses[info.Address] = *info

	return s.save()
}

func (s *Storage) HasAddress(address string) bool {
	_, exists := s.store.Addresses[address]

	return exists
}

func (s *Storage) AddressCount() int {
	return len(s.store.Addresses)
}

func (s *Storage) UpdateAddress(info *types.AddressInfo) error {
	s.store.Addresses[info.Address] = *info

	return s.save()
}

func (s *Storage) InsertTransaction(info *types.TransactionInfo) error {
	// Unsupported operation
	return nil
}

func (s *Storage) UpdateTransactionStatus(id string, status types.TransactionStatus) error {
	// Unsupported operation
	return nil
}

func (s *Storage) HasTransaction(id string) bool {
	// Unsupported operation
	return false
}

func (s *Storage) GetTransaction(id string) (*types.TransactionInfo, error) {
	// Unsupported operation
	return nil, nil
}

func (s *Storage) ListTransactions(receiver string, count int, skip int) ([]types.TransactionInfo, error) {
	// Unsupported operation
	return nil, nil
}

func (s *Storage) Clone(path string) (storage.IStorage, error) {
	cloned := store{
		Version:    s.store.Version,
		UUID:       uuid.New(),
		CreatedAt:  util.RoundNow(1),
		Network:    s.store.Network,
		DefaultFee: s.store.DefaultFee,
		Vault:      s.store.Vault,
		Addresses:  s.store.Addresses,
	}

	strg := &Storage{
		path:  path,
		store: cloned,
	}

	err := strg.save()
	if err != nil {
		return nil, err
	}

	return strg, nil
}
