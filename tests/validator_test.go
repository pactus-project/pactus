package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func getValidator(t *testing.T, addr crypto.Address) *validator.Validator {
	url := fmt.Sprintf("http://%s/validator/address/%s", tCurlAddress, addr.String())
	for i := 0; i < 10; i++ {
		res, err := http.Get(url)
		if err == nil {
			if res.StatusCode == 200 {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(res.Body)
				assert.NoError(t, err)
				var val validator.Validator
				err = val.UnmarshalJSON(buf.Bytes())
				assert.NoError(t, err)
				return &val
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	assert.NoError(t, fmt.Errorf("timeout"))
	return nil
}

func TestValidator(t *testing.T) {

	res := getValidator(t, tSigners["node_2"].Address())
	require.NotNil(t, res)
	assert.Zero(t, res.Stake())
	assert.Equal(t, res.Number(), 1)
}
