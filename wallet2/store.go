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
	versionPair, err := s.db.GetPairByKey(VersionKey)
	if err != nil {
		return err
	}
	s.Version, _ = strconv.Atoi(versionPair.Value)

	uuidPair, err := s.db.GetPairByKey(UUIDKey)
	if err != nil {
		return err
	}
	s.UUID = uuid.MustParse(uuidPair.Value)

	createdAtPair, err := s.db.GetPairByKey(CreatedAtKey)
	if err != nil {
		return err
	}
	s.CreatedAt = createdAtPair.CreatedAt

	networkPair, err := s.db.GetPairByKey(NetworkKey)
	if err != nil {
		return err
	}
	pNetwork, _ := strconv.Atoi(networkPair.Value)
	s.Network = genesis.ChainType(pNetwork)

	return nil
}

func (s *store) Save(version int, id uuid.UUID, createdAt time.Time, network genesis.ChainType) error {
	if _, err := s.db.InsertIntoPair(VersionKey, strconv.Itoa(version)); err != nil {
		return err
	}

	if _, err := s.db.InsertIntoPair(UUIDKey, id.String()); err != nil {
		return err
	}

	if _, err := s.db.InsertIntoPair(CreatedAtKey, createdAt.String()); err != nil {
		return err
	}

	if _, err := s.db.InsertIntoPair(NetworkKey, strconv.Itoa(int(network))); err != nil {
		return err
	}

	return nil
}
