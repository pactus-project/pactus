package genesis_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	acc, signer := ts.GenerateTestAccount(0)
	acc.AddToBalance(100000)
	val, _ := ts.GenerateTestValidator(0)
	gen1 := genesis.MakeGenesis(util.Now(),
		map[crypto.Address]*account.Account{signer.Address(): acc},
		[]*validator.Validator{val}, param.DefaultParams())
	gen2 := new(genesis.Genesis)

	assert.Equal(t, gen1.Params().BlockTimeInSecond, 10)

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

func TestGenesisTestNet(t *testing.T) {
	g := genesis.TestnetGenesis()
	assert.Equal(t, len(g.Validators()), 4)
	assert.Equal(t, len(g.Accounts()), 1)

	assert.Equal(t, g.Accounts()[crypto.TreasuryAddress].Balance(), int64(21e15))

	genTime, _ := time.Parse("2006-01-02", "2023-05-08")
	assert.Equal(t, g.GenesisTime(), genTime)
	assert.Equal(t, g.Params().BondInterval, uint32(120))

	expected, _ := hash.FromString("d5b41d8c2d45a951329512915ae7b2f4a445aeda380a1d40c50b5a5743dfb3fa")
	assert.Equal(t, g.Hash(), expected)
}

func TestCheckGenesisAccountAndValidator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	accs := map[crypto.Address]*account.Account{}
	vals := []*validator.Validator{}
	for i := int32(0); i < 10; i++ {
		pub, _ := ts.RandomBLSKeyPair()
		acc := account.NewAccount(i)
		val := validator.NewValidator(pub, i)

		accs[pub.Address()] = acc
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
