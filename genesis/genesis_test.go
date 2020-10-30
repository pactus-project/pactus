package genesis

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func TestMarshaling(t *testing.T) {
	addr, pb, _ := crypto.RandomKeyPair()
	acc := account.NewAccount(addr)
	val := validator.NewValidator(pb, 1)
	gen1 := MakeGenesis("test", time.Now().Truncate(0), []*account.Account{acc}, []*validator.Validator{val})
	gen2 := new(Genesis)

	bz, err := json.Marshal(gen1)
	require.NoError(t, err)
	err = json.Unmarshal(bz, gen2)
	require.NoError(t, err)
	require.Equal(t, gen1.Hash(), gen2.Hash())
}
