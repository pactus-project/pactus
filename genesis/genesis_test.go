package genesis

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func TestMarshaling(t *testing.T) {
	acc, _ := account.GenerateTestAccount(0)
	acc.AddToBalance(100000)
	val, _ := validator.GenerateTestValidator(0)
	gen1 := MakeGenesis("test", time.Now().Truncate(0), []*account.Account{acc}, []*validator.Validator{val}, 5)
	gen2 := new(Genesis)

	bz, err := json.MarshalIndent(gen1, " ", " ")
	require.NoError(t, err)
	err = json.Unmarshal(bz, gen2)
	require.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen2.Hash())
}

func TestGenesisTestNet(t *testing.T) {
	g := Testnet()
	assert.Equal(t, len(g.Validators()), 4)
	assert.Equal(t, len(g.Accounts()), 1)

	for _, v := range g.Validators() {
		assert.Equal(t, v.Address(), v.PublicKey().Address())
	}

	assert.Equal(t, g.Accounts()[0].Address(), crypto.TreasuryAddress)
	assert.Equal(t, g.Accounts()[0].Balance(), int64(0x775F05A074000))

	expected, _ := crypto.HashFromString("2dc57c69f70d74e0d1c5dba7b30dcf0903402c37e523efee3b910bdca73a2234")
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
	gen := MakeGenesis("test", time.Now().Truncate(0), accs, vals, 5)

	genAccs := gen.Accounts()
	genVals := gen.Validators()
	for i := 0; i < 10; i++ {
		assert.Equal(t, genAccs[i], accs[i])
		assert.Equal(t, genVals[i].Hash(), vals[i].Hash())
	}
}
