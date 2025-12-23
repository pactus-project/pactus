package jsonstorage

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Storage struct {
	path  string
	store store
}

func Create(path string, network genesis.ChainType, vault vault.Vault) (*Storage, error) {
	store := store{
		Version:   VersionLatest,
		UUID:      uuid.New(),
		CreatedAt: time.Now().Round(time.Second).UTC(),
		Network:   network,
		Vault:     vault,
		Addresses: make(map[string]types.AddressInfo),
	}

	strg := &Storage{
		path:  path,
		store: store,
	}

	err := strg.save()
	if err != nil {
		return nil, err
	}

	return strg, nil
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

func (s *Storage) save() error {
	s.store.VaultCRC = s.store.calcVaultCRC()

	data, err := json.MarshalIndent(s.store, "  ", "  ")
	if err != nil {
		return err
	}

	return util.WriteFile(s.path, data)
}

func (s *Storage) WalletInfo() *types.WalletInfo {
	return &types.WalletInfo{
		Path:       s.path,
		Version:    s.store.Version,
		Network:    s.store.Network,
		DefaultFee: s.store.DefaultFee,
		UUID:       s.store.UUID.String(),
		Encrypted:  s.store.Vault.IsEncrypted(),
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

func (s *Storage) AllAddresses() ([]types.AddressInfo, error) {
	var addrs []types.AddressInfo
	for _, info := range s.store.Addresses {
		addrs = append(addrs, info)
	}

	return addrs, nil
}

func (s *Storage) InsertAddress(info *types.AddressInfo) error {
	s.store.Addresses[info.Address] = *info

	return s.save()
}

func (s *Storage) UpdateAddress(info *types.AddressInfo) error {
	s.store.Addresses[info.Address] = *info

	return nil
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
