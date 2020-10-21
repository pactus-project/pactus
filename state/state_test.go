package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/utils"
	"github.com/zarbchain/zarb-go/validator"
)

func mockState(t *testing.T) (*State, crypto.Address) {
	pb, _ := crypto.GenerateRandomKey()
	addr := pb.Address()
	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)
	val := validator.NewValidator(pb, 1)
	gen := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
	loggerConfig := logger.DefaultConfig()
	loggerConfig.Levels["default"] = "error"
	logger.InitLogger(loggerConfig)
	storeConfig := store.DefaultConfig()
	storeConfig.Path = utils.TempDirName()
	store, err := store.NewStore(storeConfig)
	require.NoError(t, err)
	txPoolConfig := txpool.DefaultConfig()
	txPool, err := txpool.NewTxPool(txPoolConfig, nil)
	require.NoError(t, err)
	st, err := LoadOrNewState(gen, store, txPool, nil)
	require.NoError(t, err)
	return st, addr
}

func TestBlockValidate(t *testing.T) {
	st, addr := mockState(t)
	block, _ := st.ProposeBlock(1, addr, nil)
	err := st.ApplyBlock(&block, 0)
	require.NoError(t, err)

}
