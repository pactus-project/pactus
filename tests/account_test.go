package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
)

func getAccount(t *testing.T, addr crypto.Address) *account.Account {
	url := fmt.Sprintf("http://%s/account/address/%s", tCurlAddress, addr.String())
	for i := 0; i < 10; i++ {
		res, err := http.Get(url)
		if err == nil {
			if res.StatusCode == 200 {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(res.Body)
				assert.NoError(t, err)
				var acc account.Account
				err = acc.UnmarshalJSON(buf.Bytes())
				assert.NoError(t, err)
				return &acc
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
	assert.NoError(t, fmt.Errorf("timeout"))
	return nil
}

func TestTreasuryAccount(t *testing.T) {

	res := getAccount(t, crypto.TreasuryAddress)
	require.NotNil(t, res)
	assert.Equal(t, tGenDoc.Accounts()[0].Balance(), res.Balance()+int64(res.Sequence()*500000000))
}
