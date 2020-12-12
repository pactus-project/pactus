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
	addr, pb, _ := crypto.RandomKeyPair()
	acc := account.NewAccount(addr, 0)
	acc.AddToBalance(100000)
	val := validator.NewValidator(pb, 0, 0)
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

	assert.Equal(t, g.Accounts()[0].Address(), crypto.MintbaseAddress)
	assert.Equal(t, g.Accounts()[0].Balance(), int64(0x4A9B6384488000))

	expected, _ := crypto.HashFromString("804c372bd7c0327a7c84389983f159b43510c3596097f2bb36f8856136dd829a")
	assert.Equal(t, g.Hash(), expected)
}
