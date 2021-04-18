package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/validator"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetValidator(t *testing.T) {
	conn, client := callServer(t)

	k1, k2, k3 := key.GenerateRandomKey(), key.GenerateRandomKey(), key.GenerateRandomKey()
	val1 := validator.NewValidator(k1.PublicKey(), 0, 0)
	val2 := validator.NewValidator(k2.PublicKey(), 1, 0)
	tMockState.Store.UpdateValidator(val1)
	tMockState.Store.UpdateValidator(val2)

	t.Run("Should return nil value due to invalid address", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: "Non existence address",
		})
		assert.Error(t, err, "Error should be returned")
		assert.Nil(t, res, "Response should be empty")
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: k3.Address().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator, and the public keys should match", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{Address: k1.Address().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
	})

	err := conn.Close()
	assert.Nil(t, err, "Error closing connection")
}

func TestGetValidatorByNumber(t *testing.T) {
	conn, client := callServer(t)

	k1, k2 := key.GenerateRandomKey(), key.GenerateRandomKey()
	val1 := validator.NewValidator(k1.PublicKey(), 5, 0)
	val2 := validator.NewValidator(k2.PublicKey(), 6, 0)
	tMockState.Store.UpdateValidator(val1)
	tMockState.Store.UpdateValidator(val2)

	t.Run("Should return nil value due to invalid number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: -3,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: 3,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator matching with public key and number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: 6,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
		assert.Equal(t, int32(val1.Number()), res.GetValidator().GetNumber())

	})

	err := conn.Close()
	assert.Nil(t, err, "Error closing connection")
}

func TestGetValidators(t *testing.T) {
	conn, client := callServer(t)

	t.Run("should setup commiters", func(t *testing.T) {
		res, err := client.GetValidators(tCtx, &zarb.ValidatorsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 4, len(res.GetValidators()))
	})

	err := conn.Close()
	assert.Nil(t, err, "Error closing connection")
}
