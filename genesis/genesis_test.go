package genesis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
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

func TestGenesisTestNet(t *testing.T) {
	g := Testnet()
	assert.Equal(t, len(g.Validators()), 4)
	assert.Equal(t, len(g.Accounts()), 1)

	for _, v := range g.Validators() {
		assert.Equal(t, v.Address(), v.PublicKey().Address())
	}

	assert.Equal(t, g.Accounts()[0].Address(), crypto.TreasuryAddress)
	assert.Equal(t, g.Accounts()[0].Balance(), int64(4200000000000000))

	expected, _ := crypto.HashFromString("e5df62fe9fb49cd0471633b3d605fe178a9434e4788c4a57b98dd1c6c6c5bed1")
	assert.Equal(t, g.Hash(), expected)
}

func TestCheckGenesisAccountAndValidator(t *testing.T) {
	accs := []*account.Account{}
	vals := []*validator.Validator{}
	for i := 0; i < 10; i++ {
		a, pub, _ := crypto.GenerateTestKeyPair()
		acc := account.NewAccount(a, i)
		val := validator.NewValidator(pub, i, 0)

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
