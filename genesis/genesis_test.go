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
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, prv := ts.GenerateTestAccount(0)
	acc.AddToBalance(100000)
	val, _ := ts.GenerateTestValidator(0)
	gen1 := genesis.MakeGenesis(util.Now(),
		map[crypto.Address]*account.Account{prv: acc},
		[]*validator.Validator{val}, param.DefaultParams())
	gen2 := new(genesis.Genesis)

	assert.Equal(t, gen1.Params().BlockIntervalInSecond, 10)

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

func TestGenesisTestnet(t *testing.T) {
	crypto.AddressHRP = "tpc"

	gen := genesis.TestnetGenesis()
	assert.Equal(t, len(gen.Validators()), 4)
	assert.Equal(t, len(gen.Accounts()), 5)

	genTime, _ := time.Parse("2006-01-02", "2024-03-16")
	expected, _ := hash.FromString("13f96e6fbc9e0de0d53537ac5e894fc8e66be1600436db2df1511dc30696e822")
	assert.Equal(t, gen.Hash(), expected)
	assert.Equal(t, gen.GenesisTime(), genTime)
	assert.Equal(t, gen.Params().BondInterval, uint32(360))
	assert.Equal(t, gen.ChainType(), genesis.Testnet)
	assert.Equal(t, gen.TotalSupply(), amount.Amount(42e15))

	crypto.AddressHRP = "pc"
}

func TestGenesisMainnet(t *testing.T) {
	gen := genesis.MainnetGenesis()
	assert.Equal(t, len(gen.Validators()), 4)
	assert.Equal(t, len(gen.Accounts()), 5)

	genTime, _ := time.Parse("02 Jan 2006, 15:04 MST", "24 Jan 2024, 20:24 UTC")
	expected, _ := hash.FromString("e4d59e3145c9d718caf178edb33bc2ca7fe43e5b30990c9d57d53a60c4741432")
	assert.Equal(t, gen.Hash(), expected)
	assert.Equal(t, gen.GenesisTime(), genTime)
	assert.Equal(t, gen.Params().BondInterval, uint32(8640/24))
	assert.Equal(t, gen.Params().UnbondInterval, uint32(8640*21))
	assert.Equal(t, gen.ChainType(), genesis.Mainnet)
	assert.Equal(t, gen.TotalSupply(), amount.Amount(42e15))
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
	gen := genesis.MakeGenesis(util.Now(), accs, vals, param.DefaultParams())

	for addr, acc := range gen.Accounts() {
		assert.Equal(t, acc, accs[addr])
	}

	for i, val := range gen.Validators() {
		assert.Equal(t, val.Hash(), vals[i].Hash())
	}
}
