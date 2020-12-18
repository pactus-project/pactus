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
	assert.Equal(t, g.Accounts()[0].Balance(), int64(0x4A9B6384488000))

	expected, _ := crypto.HashFromString("cabed9b2a3ac28b79cbe02cb521f462d6a2a37bba201618b377fe3a671f99fd2")
	assert.Equal(t, g.Hash(), expected)
}
