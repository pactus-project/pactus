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

	k, k1, k2 := key.GenerateRandomKey(), key.GenerateRandomKey(), key.GenerateRandomKey()
	val := validator.NewValidator(k.PublicKey(), 0, 0)
	val1 := validator.NewValidator(k1.PublicKey(), 1, 0)
	tMockState.Store.Validators[k.Address()] = val
	tMockState.Store.Validators[k1.Address()] = val1

	t.Run("Should return nil value due to invalid address", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: "Non existence address",
		})
		// assert.Error(t, err)
		assert.NotNil(t, err, "Error should be returned")
		assert.Nil(t, res, "Response should be empty")
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{
			Address: k2.Address().String(),
		})

		assert.Nil(t, res)
		assert.NotEmpty(t, err)
	})

	t.Run("Should return validator, verbosity 0", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{Address: k1.Address().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.Json)
	})

	t.Run("Should return transaction, verbosity 1", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{Address: k.Address().String(), Verbosity: 1})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Json)
	})

	err := conn.Close()

	assert.Nil(t, err, "Error closing connection")
}

func TestGetValidatorByNumber(t *testing.T) {
	conn, client := callServer(t)

	k, k1 := key.GenerateRandomKey(), key.GenerateRandomKey()
	val := validator.NewValidator(k.PublicKey(), 0, 0)
	val1 := validator.NewValidator(k1.PublicKey(), 1, 0)
	tMockState.Store.Validators[k.Address()] = val
	tMockState.Store.Validators[k1.Address()] = val1

	t.Run("Should return nil value due to invalid number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: -3,
		})
		// assert.Error(t, err)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: 3,
		})
		// assert.Error(t, err)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator json, verbosity 0", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number: 1,
		})
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Empty(t, res.Json)
	})

	t.Run("Should return transaction json, verbosity 1", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &zarb.ValidatorByNumberRequest{
			Number:    0,
			Verbosity: 1,
		})
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.NotEmpty(t, res.Json)
	})

	err := conn.Close()

	assert.Nil(t, err, "Error closing connection")
}
