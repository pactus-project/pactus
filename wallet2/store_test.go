package wallet2

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet2/db"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	someDB, _ := db.NewDB(context.Background(), ":memory:")
	_ = someDB.CreateTables()

	s := newStore(someDB)

	version := 1
	id := uuid.New()
	createdAt := time.Now().Round(time.Second).UTC()
	network := genesis.Mainnet

	err := s.Save(version, id, createdAt, network)

	assert.NoError(t, err)

	s.Version = 0
	s.UUID = uuid.Nil
	s.CreatedAt = time.Now()
	s.Network = genesis.Localnet

	err = s.Load()

	assert.NoError(t, err)

	assert.Equal(t, version, s.Version)
	assert.Equal(t, id, s.UUID)
	assert.Equal(t, createdAt, s.CreatedAt)
	assert.Equal(t, network, s.Network)
}
