package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func getValidator(t *testing.T, addr crypto.Address) *validator.Validator {
	res := tCapnpServer.GetValidator(tCtx, func(p capnp.ZarbServer_getValidator_Params) error {
		assert.NoError(t, p.SetAddress([]byte(addr.String())))
		return nil
	}).Result()

	st, err := res.Struct()
	if err != nil {
		return nil
	}

	d, _ := st.Data()
	val := new(validator.Validator)
	assert.NoError(t, val.Decode(d))
	return val
}

func TestGetValidator(t *testing.T) {
	val := getValidator(t, tSigners[tNodeIdx2].Address())
	require.NotNil(t, val)
	assert.Equal(t, val.Number(), 1)
}
