package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/validator"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetValidator(t *testing.T) {
	conn, client := callServer(t)

	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	pub3, _ := bls.GenerateTestKeyPair()

	val1 := validator.NewValidator(pub1, 0)
	val2 := validator.NewValidator(pub2, 1)
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
			Address: pub3.Address().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator, and the public keys should match", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &zarb.ValidatorRequest{Address: pub1.Address().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
	})

	err := conn.Close()
	assert.Nil(t, err, "Error closing connection")
}

func TestGetValidatorByNumber(t *testing.T) {
	conn, client := callServer(t)

	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()

	val1 := validator.NewValidator(pub1, 5)
	val2 := validator.NewValidator(pub2, 6)
	val2.UpdateLastBondingHeight(100)
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
			Number: 5,
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

	t.Run("should return list of validators", func(t *testing.T) {
		res, err := client.GetValidators(tCtx, &zarb.ValidatorsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 21, len(res.GetValidators()))
	})

	err := conn.Close()
	assert.Nil(t, err, "Error closing connection")
}
