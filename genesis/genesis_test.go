package genesis_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, addr := ts.GenerateTestAccount(
		testsuite.AccountWithNumber(0),
		testsuite.AccountWithBalance(100000))
	val := ts.GenerateTestValidator(
		testsuite.ValidatorWithNumber(0),
	)
	gen1 := genesis.MakeGenesis(util.RoundNow(10),
		map[crypto.Address]*account.Account{addr: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())
	gen2 := new(genesis.Genesis)

	assert.Equal(t, 10, gen1.Params().BlockIntervalInSecond)

	bz, err := json.MarshalIndent(gen1, " ", " ")
	require.NoError(t, err)
	err = json.Unmarshal(bz, gen2)
	require.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen2.Hash())

	// Test saving and loading
	f := util.TempFilePath()
	assert.NoError(t, gen1.SaveToFile(f))
	gen3, err := genesis.LoadFromFile(f)
	assert.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen3.Hash())

	_, err = genesis.LoadFromFile(util.TempFilePath())
	assert.Error(t, err, "file not found")
}

func TestGenesisMainnet(t *testing.T) {
	gen := genesis.MainnetGenesis()
	assert.Equal(t, len(gen.Validators()), 4)
	assert.Equal(t, len(gen.Accounts()), 5)

	genTime, _ := time.Parse("02 Jan 2006, 15:04 MST", "24 Jan 2024, 20:24 UTC")
	expected, _ := hash.FromString("e4d59e3145c9d718caf178edb33bc2ca7fe43e5b30990c9d57d53a60c4741432")
	assert.Equal(t, expected, gen.Hash())
	assert.Equal(t, genTime, gen.GenesisTime())
	assert.Equal(t, uint32(8640/24), gen.Params().BondInterval)
	assert.Equal(t, uint32(8640*21), gen.Params().UnbondInterval)
	assert.Equal(t, genesis.Mainnet, gen.ChainType())
	assert.Equal(t, amount.Amount(42e15), gen.TotalSupply())
	assert.True(t, gen.ChainType().IsMainnet())
}

func TestCheckGenesisAccountAndValidator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	accs := map[crypto.Address]*account.Account{}
	vals := []*validator.Validator{}
	for i := int32(0); i < 10; i++ {
		pub, _ := ts.RandBLSKeyPair()
		acc := account.NewAccount(i)
		val := validator.NewValidator(pub, i)

		accs[pub.AccountAddress()] = acc
		vals = append(vals, val)
	}
	gen := genesis.MakeGenesis(time.Now(), accs, vals, genesis.DefaultGenesisParams())

	for addr, acc := range gen.Accounts() {
		assert.Equal(t, accs[addr], acc)
	}

	for i, val := range gen.Validators() {
		assert.Equal(t, vals[i].Hash(), val.Hash())
	}
}

func TestGenesisTestnet(t *testing.T) {
	crypto.AddressHRP = "tpc"

	gen := genesis.TestnetGenesis()
	assert.Equal(t, 4, len(gen.Validators()))
	assert.Equal(t, 5, len(gen.Accounts()))

	genTime, _ := time.Parse("2006-01-02", "2024-03-16")
	expected, _ := hash.FromString("13f96e6fbc9e0de0d53537ac5e894fc8e66be1600436db2df1511dc30696e822")
	assert.Equal(t, expected, gen.Hash())
	assert.Equal(t, genTime, gen.GenesisTime())
	assert.Equal(t, uint32(360), gen.Params().BondInterval)
	assert.Equal(t, genesis.Testnet, gen.ChainType())
	assert.Equal(t, amount.Amount(42e15), gen.TotalSupply())
	assert.True(t, gen.ChainType().IsTestnet())

	// reset address HRP global variable to miannet to prevent next tests failing.
	crypto.AddressHRP = "pc"
}
