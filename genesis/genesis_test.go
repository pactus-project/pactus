package genesis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestMarshaling(t *testing.T) {
	acc, _ := account.GenerateTestAccount(0)
	acc.AddToBalance(100000)
	val, _ := validator.GenerateTestValidator(0)
	gen1 := MakeGenesis(util.Now(), []*account.Account{acc}, []*validator.Validator{val}, param.DefaultParams())
	gen2 := new(Genesis)

	assert.Equal(t, gen1.Params().BlockTimeInSecond, 10)

	bz, err := json.MarshalIndent(gen1, " ", " ")
	require.NoError(t, err)
	err = json.Unmarshal(bz, gen2)
	require.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen2.Hash())

	// Test saving and loading
	f := util.TempFilePath()
	assert.NoError(t, gen1.SaveToFile(f))
	gen3, err := LoadFromFile(f)
	assert.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen3.Hash())
}

// func TestGenesisTestNet(t *testing.T) {
// 	g := Testnet()
// 	assert.Equal(t, len(g.Validators()), 4)
// 	assert.Equal(t, len(g.Accounts()), 1)

// 	for _, v := range g.Validators() {
// 		assert.Equal(t, v.Address(), v.PublicKey().Address())
// 	}

// 	assert.Equal(t, g.Accounts()[0].Address(), crypto.TreasuryAddress)
// 	assert.Equal(t, g.Accounts()[0].Balance(), int64(2100000000000000))

// 	expected, _ := hash.FromString("4d22446ce560c591575b8205de52da0bd99757c9254bb248d1c9853208733c30")
// 	assert.Equal(t, g.Hash(), expected)
// }

func TestCheckGenesisAccountAndValidator(t *testing.T) {
	accs := []*account.Account{}
	vals := []*validator.Validator{}
	for i := 0; i < 10; i++ {
		pub, _ := bls.GenerateTestKeyPair()
		acc := account.NewAccount(pub.Address(), i)
		val := validator.NewValidator(pub, i)

		accs = append(accs, acc)
		vals = append(vals, val)
	}
	gen := MakeGenesis(util.Now(), accs, vals, param.DefaultParams())

	genAccs := gen.Accounts()
	genVals := gen.Validators()
	for i := 0; i < 10; i++ {
		assert.Equal(t, genAccs[i], accs[i])
		assert.Equal(t, genVals[i].Hash(), vals[i].Hash())
	}
}
