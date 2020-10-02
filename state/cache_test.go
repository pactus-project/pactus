package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/zarb-chain/zarb-go/account"
	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/genesis"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/store"
	"gitlab.com/zarb-chain/zarb-go/utils"
	"gitlab.com/zarb-chain/zarb-go/validator"
)

func testState(t *testing.T) *State {
	pb, _ := crypto.GenerateRandomKey()
	addr := pb.Address()
	acc := account.NewAccount(addr)
	val := validator.NewValidator(pb, 1)
	gen := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
	conf := config.DefaultConfig()
	conf.Store.Path = utils.TempPath()
	store, err := store.NewStore(conf)
	logger.InitLogger(conf)
	st, err := LoadStateOrNewState(conf, gen, store, nil)
	require.NoError(t, err)
	return st
}

func TestAccountChange(t *testing.T) {
	st := testState(t)
	cache := newCache(st.store)
	pb1, _ := crypto.GenerateRandomKey()
	pb2, _ := crypto.GenerateRandomKey()
	//pb3, _ := crypto.GenerateRandomKey()
	addr1 := pb1.Address()
	addr2 := pb2.Address()
	//addr3 := pb3.Address()

	acc := cache.Account(addr1)
	assert.Nil(t, acc, nil)

	// update cache
	acc1 := account.NewAccount(addr1)
	acc1.AddToBalance(10)
	cache.UpdateAccount(acc1)
	acc11 := cache.Account(addr1)
	assert.Equal(t, acc1, acc11)

	// update state
	acc2 := account.NewAccount(addr2)
	cache.UpdateAccount(acc2)
	acc22 := cache.Account(addr2)
	assert.Equal(t, acc2, acc22)

	// /// update storages
	// val, err := cache.Storage(addr1, binary.Uint64ToWord256(1))
	// assert.NoError(t, err)
	// assert.Equal(t, val, binary.Uint64ToWord256(0))
	// cache.SetStorage(addr1, binary.Uint64ToWord256(1), binary.Uint64ToWord256(2))
	// val, err = cache.Storage(addr1, binary.Uint64ToWord256(2))
	// assert.NoError(t, err)
	// assert.Equal(t, val, binary.Uint64ToWord256(0)) // wrong storage key
	// val, err = cache.Storage(addr1, binary.Uint64ToWord256(1))
	// assert.NoError(t, err)
	// assert.Equal(t, val, binary.Uint64ToWord256(2))

	// /// Update storage then account
	// acc3 := account.NewAccount(addr3)
	// st.updateAccount(acc3)
	// cache.SetStorage(addr3, binary.Uint64ToWord256(1), binary.Uint64ToWord256(2))
	// acc3.AddToBalance(10)
	// cache.UpdateAccount(acc3)
	// acc33, err := cache.Account(addr3)
	// assert.NoError(t, err)
	// assert.Equal(t, acc3, acc33)
	// val, err = cache.Storage(addr3, binary.Uint64ToWord256(1))
	// assert.NoError(t, err)
	// assert.Equal(t, val, binary.Uint64ToWord256(2))

	/// accounts should be untouched while changing storages
	acc11 = cache.Account(addr1)
	assert.Equal(t, acc1, acc11)

	acc22 = cache.Account(addr2)
	assert.Equal(t, acc2, acc22)

	cache.reset()
	assert.Equal(t, cache.accChanges.Len(), 0)
	assert.Equal(t, cache.valChanges.Len(), 0)
}
