package state

import (
	"bytes"
	"testing"
	"time"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// realStoreData is a version-5 chain on a LevelDB store, configured the way the
// node configures it: the recent transaction cache spans exactly the
// transaction time to live.
type realStoreData struct {
	*testData

	genDoc    *genesis.Genesis
	valKeys   []*bls.ValidatorKey
	storeConf *store.Config
}

func setupWithRealStore(t *testing.T, ttl uint32) *realStoreData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	genValKeys := make([]*bls.ValidatorKey, 0, 4)
	genVals := make([]*validator.Validator, 0, 4)
	for i := int32(0); i < 4; i++ {
		valKey := ts.RandValKey()
		genValKeys = append(genValKeys, valKey)
		genVals = append(genVals, validator.NewValidator(valKey.PublicKey(), i))
	}

	genParams := genesis.DefaultGenesisParams()
	genParams.CommitteeSize = 7
	genParams.BondInterval = 10
	genParams.BlockVersion = protocol.ProtocolVersion5
	genParams.TransactionToLiveInterval = ttl

	treasury := account.NewAccount(0)
	treasury.AddToBalance(21 * 1e15)
	holder := account.NewAccount(1)
	holder.AddToBalance(21 * 1e15)
	genAccPubKey, genAccPrvKey := ts.RandEd25519KeyPair()
	genDoc := genesis.MakeGenesis(util.RoundNow(10).Add(-8640*time.Second),
		map[crypto.Address]*account.Account{
			crypto.TreasuryAddress:        treasury,
			genAccPubKey.AccountAddress(): holder,
		}, genVals, genParams)

	// Same private settings as cmd.MakeConfig.
	storeConf := store.DefaultConfig()
	storeConf.Path = util.TempDirPath()
	storeConf.TxCacheWindow = genParams.TransactionToLiveInterval
	storeConf.SeedCacheWindow = genParams.SortitionInterval

	fakeTxPool := txpool.NewFakeTxPool(ts)
	fakeTxPool.EXPECT().SetNewSandboxAndRecheck(gomock.Any()).Return().AnyTimes()
	fakeTxPool.EXPECT().HandleCommittedBlock(gomock.Any()).Return().AnyTimes()
	fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()

	chain := &realStoreData{
		testData: &testData{
			TestSuite:  ts,
			fakeTxPool: fakeTxPool,
			genValKeys: genValKeys,
			genAccKey:  genAccPrvKey,
		},
		genDoc:    genDoc,
		valKeys:   []*bls.ValidatorKey{genValKeys[0], ts.RandValKey()},
		storeConf: storeConf,
	}
	chain.open(t)
	t.Cleanup(func() { chain.state.Close() })

	return chain
}

func (td *realStoreData) open(t *testing.T) {
	t.Helper()

	str, err := store.NewStore(td.storeConf)
	require.NoError(t, err)
	loaded, err := LoadOrNewState(t.Context(), td.genDoc, td.valKeys, str, td.fakeTxPool,
		pipeline.New[any](t.Context()))
	require.NoError(t, err)
	td.state = loaded.(*state)
}

// restart closes the store and loads the state again from disk, which rebuilds
// the recent transaction cache from the stored blocks.
func (td *realStoreData) restart(t *testing.T) {
	t.Helper()

	td.state.Close()
	td.open(t)
}

// PIP-50, Replay: on a real store, the node remembers an executed transaction
// for at least as long as its lock time is valid, also after a restart. So a
// Delete replayed at any later height never runs again, up to and past the
// edge of its lock time, and the re-created slot survives.
func TestAnchorReplayWindowOnRealStore(t *testing.T) {
	const ttl = 6
	chain := setupWithRealStore(t, ttl)
	sender := chain.sender()

	chain.commit(t, chain.propose(t, block.Txs{
		chain.anchor(chain.nextHeight(), payload.AnchorActionSet, bytes.Repeat([]byte{0x71}, 32), "first",
			executor.MinAnchorDeposit, 1),
	}))
	deleteHeight := chain.nextHeight()
	deleted := chain.anchor(deleteHeight, payload.AnchorActionDelete, nil, "", 0, 1)
	chain.commit(t, chain.propose(t, block.Txs{deleted}))
	chain.commit(t, chain.propose(t, block.Txs{
		chain.anchor(chain.nextHeight(), payload.AnchorActionSet, bytes.Repeat([]byte{0x72}, 32), "second",
			executor.MinAnchorDeposit, 1),
	}))
	recreated := chain.mustAccount(t, sender).Hash()

	// The lock time of `deleted` (its block height) is valid up to lastValid.
	lastValid := deleteHeight + ttl
	for height := chain.nextHeight(); height <= lastValid+2; height = chain.nextHeight() {
		if height == lastValid {
			chain.restart(t)
		}

		next := chain.propose(t, block.Txs{deleted})
		require.False(t, blockHasTx(next, deleted.ID()))
		forced := rebuildBlock(next, next.Header().Version(), appendTxs(next.Transactions(), deleted))
		err := chain.state.ValidateBlock(forced, 0)
		switch {
		case height <= lastValid+1:
			// Still in the cache: one block longer than the lock time lasts.
			require.ErrorIs(t, err, execution.TransactionCommittedError{ID: deleted.ID()}, "height %d", height)
		default:
			require.ErrorIs(t, err, execution.LockTimeExpiredError{LockTime: deleteHeight}, "height %d", height)
		}

		chain.commit(t, next)
		require.Equal(t, recreated, chain.mustAccount(t, sender).Hash(), "height %d", height)
	}
	require.Equal(t, lastValid+3, chain.nextHeight())
	require.Equal(t, types.Height(ttl), lastValid-deleteHeight)
}
