package encrypter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopeEncrypter(t *testing.T) {
	e := NopeEncrypter()
	assert.Equal(t, e.Method, "")
	assert.Nil(t, e.Params)
	assert.False(t, e.IsEncrypted())

	msg := "foo"
	_, err := e.Encrypt(msg, "password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	enc, err := e.Encrypt(msg, "")
	assert.NoError(t, err)
	assert.Equal(t, enc, msg)

	_, err = e.Decrypt(enc, "password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	dec, err := e.Decrypt(enc, "")
	assert.NoError(t, err)
	assert.Equal(t, dec, msg)
}

func TestDefaultEncrypter(t *testing.T) {
	opts := []Option{
		OptionIteration(3),
		OptionMemory(4),
		OptionParallelism(5),
	}
	e := DefaultEncrypter(opts...)
	assert.Equal(t, e.Method, "ARGON2ID-AES_256_CTR-MACV1")
	assert.Equal(t, e.Params["iterations"], "3")
	assert.Equal(t, e.Params["memory"], "4")
	assert.Equal(t, e.Params["parallelism"], "5")
	assert.True(t, e.IsEncrypted())
}

func TestEncrypter(t *testing.T) {
	e := &Encrypter{
		Method: "ARGON2ID-AES_256_CTR-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "1",
			nameParamParallelism: "1",
		},
	}

	password := "cowboy"
	msg := "foo"
	enc, err := e.Encrypt(msg, password)
	assert.NoError(t, err)

	dec, err := e.Decrypt(enc, password)
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)

	_, err = e.Decrypt(enc, "invalid-password")
	assert.ErrorIs(t, err, ErrInvalidPassword)

	_, err = e.Encrypt(enc, "")
	assert.ErrorIs(t, err, ErrInvalidPassword)
}
