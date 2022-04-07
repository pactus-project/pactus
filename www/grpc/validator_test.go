package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetValidator(t *testing.T) {
	conn, client := callServer(t)

	val1 := tMockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid address", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: "",
		})
		assert.Error(t, err, "Error should be returned")
		assert.Nil(t, res, "Response should be empty")
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: crypto.GenerateTestAddress().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator, and the public keys should match", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: val1.Address().String(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().Bytes(), res.GetValidator().PublicKey)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidatorByNumber(t *testing.T) {
	conn, client := callServer(t)

	val1 := tMockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: -1,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: val1.Number() + 1,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator matching with public key and number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: val1.Number(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().Bytes(), res.GetValidator().PublicKey)
		assert.Equal(t, val1.Number(), res.GetValidator().GetNumber())

	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidators(t *testing.T) {
	conn, client := callServer(t)

	t.Run("should return list of validators", func(t *testing.T) {
		res, err := client.GetValidators(tCtx, &zarb.ValidatorsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 21, len(res.GetValidators()))
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
