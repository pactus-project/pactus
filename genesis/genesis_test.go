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

	addr, acc := ts.GenerateTestAccount(
		testsuite.AccountWithNumber(0),
		testsuite.AccountWithBalance(100000),
	)
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
	require.NoError(t, gen1.SaveToFile(f))
	gen3, err := genesis.LoadFromFile(f)
	require.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen3.Hash())

	_, err = genesis.LoadFromFile(util.TempFilePath())
	require.Error(t, err, "file not found")
}

func TestGenesisMainnet(t *testing.T) {
	gen := genesis.MainnetGenesis()
	assert.Len(t, gen.Validators(), 4)
	assert.Len(t, gen.Accounts(), 5)

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
	assert.Len(t, gen.Validators(), 4)
	assert.Len(t, gen.Accounts(), 5)

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

func TestValidate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pubKey, _ := ts.RandBLSKeyPair()
	pubKey2, _ := ts.RandBLSKeyPair()
	addr := ts.RandAccAddress().String()
	addr2 := ts.RandAccAddress().String()

	makeJSON := func(params, accounts, validators string) []byte {
		return []byte(`{
  "genesis_time": "2024-01-24T20:24:00Z",
  "params": ` + params + `,
  "accounts": [` + accounts + `],
  "validators": [` + validators + `]
}`)
	}
	validParams := `{"block_version": 1, "block_interval_in_second": 10, "committee_size": 4}`
	acc := func(address string) string {
		return `{"address": "` + address + `", "balance": 1}`
	}
	val := func(pub string) string {
		return `{"public_key": "` + pub + `"}`
	}

	tests := []struct {
		name  string
		data  []byte
		valid bool
	}{
		{
			"valid genesis",
			makeJSON(validParams, acc(addr)+","+acc(addr2), val(pubKey.String())+","+val(pubKey2.String())),
			true,
		},
		{
			"no accounts and validators",
			makeJSON(validParams, "", ""),
			true,
		},
		{
			"missing params",
			[]byte(`{"genesis_time": "2024-01-24T20:24:00Z", "accounts": [], "validators": []}`),
			false,
		},
		{
			"zero block interval",
			makeJSON(`{"block_version": 1, "block_interval_in_second": 0, "committee_size": 4}`, "", ""),
			false,
		},
		{
			"zero committee size",
			makeJSON(`{"block_version": 1, "block_interval_in_second": 10, "committee_size": 0}`, "", ""),
			false,
		},
		{
			"invalid account address",
			makeJSON(validParams, acc("not-an-address"), ""),
			false,
		},
		{
			"duplicated account address",
			makeJSON(validParams, acc(addr)+","+acc(addr), ""),
			false,
		},
		{
			"invalid validator public key",
			makeJSON(validParams, "", val("not-a-public-key")),
			false,
		},
		{
			"duplicated validator public key",
			makeJSON(validParams, "", val(pubKey.String())+","+val(pubKey.String())),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := util.TempFilePath()
			require.NoError(t, util.WriteFile(file, tt.data))

			gen, err := genesis.LoadFromFile(file)
			require.NoError(t, err)

			err = gen.Validate()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, &genesis.InvalidGenesisError{})
			}
		})
	}
}

func TestValidateEmbeddedGenesis(t *testing.T) {
	require.NoError(t, genesis.MainnetGenesis().Validate())

	crypto.ToTestnetHRP()
	t.Cleanup(crypto.ToMainnetHRP)
	require.NoError(t, genesis.TestnetGenesis().Validate())
}
