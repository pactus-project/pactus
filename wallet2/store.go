package wallet2

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet2/db"
	"github.com/pactus-project/pactus/wallet2/vault"
)

const (
	VersionKey   = "version"
	UUIDKey      = "uuid"
	CreatedAtKey = "createdat"
	NetworkKey   = "network"
)

type store struct {
	Version   int               `json:"version"`
	UUID      uuid.UUID         `json:"uuid"`
	CreatedAt time.Time         `json:"created_at"`
	Network   genesis.ChainType `json:"network"`

	db      db.DB
	Vault   *vault.Vault `json:"vault"`
	History history      `json:"history"`
}

func newStore(database db.DB) *store {
	return &store{
		db:      database,
		History: *newHistory(database),
	}
}

func (s *store) Load() error {
	versionValue, err := s.db.GetValue(VersionKey)
	if err != nil {
		return err
	}
	s.Version, _ = strconv.Atoi(versionValue)

	uuidValue, err := s.db.GetValue(UUIDKey)
	if err != nil {
		return err
	}
	s.UUID = uuid.MustParse(uuidValue)

	createdAtValue, err := s.db.GetValue(CreatedAtKey)
	if err != nil {
		return err
	}

	timeFormat := "2006-01-02 15:04:05 -0700 MST"
	s.CreatedAt, err = time.Parse(timeFormat, createdAtValue)
	if err != nil {
		return err
	}

	networkValue, err := s.db.GetValue(NetworkKey)
	if err != nil {
		return err
	}
	pNetwork, _ := strconv.Atoi(networkValue)
	s.Network = genesis.ChainType(pNetwork)

	return nil
}

func (s *store) Save(version int, id uuid.UUID, createdAt time.Time, network genesis.ChainType) error {
	if err := s.db.SetValue(VersionKey, strconv.Itoa(version)); err != nil {
		return err
	}
	s.Version = version

	if err := s.db.SetValue(UUIDKey, id.String()); err != nil {
		return err
	}
	s.UUID = id

	if err := s.db.SetValue(CreatedAtKey, createdAt.String()); err != nil {
		return err
	}
	s.CreatedAt = createdAt

	if err := s.db.SetValue(NetworkKey, strconv.Itoa(int(network))); err != nil {
		return err
	}
	s.Network = network

	return nil
}
