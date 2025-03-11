package encrypter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopeEncrypter(t *testing.T) {
	enc := NopeEncrypter()
	assert.Equal(t, "", enc.Method)
	assert.Nil(t, enc.Params)
	assert.False(t, enc.IsEncrypted())

	msg := "foo"
	_, err := enc.Encrypt(msg, "password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	cipher, err := enc.Encrypt(msg, "")
	assert.NoError(t, err)
	assert.Equal(t, msg, cipher)

	_, err = enc.Decrypt(cipher, "password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	decipher, err := enc.Decrypt(cipher, "")
	assert.NoError(t, err)
	assert.Equal(t, msg, decipher)
}

func TestDefaultEncrypter(t *testing.T) {
	opts := []Option{
		OptionIteration(3),
		OptionMemory(4),
		OptionParallelism(5),
	}
	enc := DefaultEncrypter(opts...)
	assert.Equal(t, "ARGON2ID-AES_256_CTR-MACV1", enc.Method)
	assert.Equal(t, "3", enc.Params["iterations"])
	assert.Equal(t, "4", enc.Params["memory"])
	assert.Equal(t, "5", enc.Params["parallelism"])
	assert.Equal(t, "48", enc.Params["keylen"])
	assert.True(t, enc.IsEncrypted())
}

func TestEncrypterV2(t *testing.T) {
	enc := &Encrypter{
		Method: "ARGON2ID-AES_256_CTR-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "1",
			nameParamParallelism: "1",
		},
	}

	msg := "foo"

	_, err := enc.Encrypt(msg, "")
	assert.ErrorIs(t, err, ErrInvalidPassword)

	password := "cowboy"
	cipher, err := enc.Encrypt(msg, password)
	assert.NoError(t, err)

	dec, err := enc.Decrypt(cipher, password)
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)

	_, err = enc.Decrypt(cipher, "invalid-password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
}

func TestEncrypterV3(t *testing.T) {
	enc := &Encrypter{
		Method: "ARGON2ID-AES_256_CTR-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "1",
			nameParamParallelism: "1",
			nameParamKeyLen:      "48",
		},
	}

	msg := "foo"

	_, err := enc.Encrypt(msg, "")
	assert.ErrorIs(t, err, ErrInvalidPassword)

	password := "cowboy"
	cipher, err := enc.Encrypt(msg, password)
	assert.NoError(t, err)

	dec, err := enc.Decrypt(cipher, password)
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)

	_, err = enc.Decrypt(cipher, "invalid-password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
}
