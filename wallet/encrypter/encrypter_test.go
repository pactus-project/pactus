package encrypter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopeEncrypter(t *testing.T) {
	e := NopeEncrypter()
	assert.Equal(t, "", e.Method)
	assert.Nil(t, e.Params)
	assert.False(t, e.IsEncrypted())

	msg := "foo"
	_, err := e.Encrypt(msg, "password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	enc, err := e.Encrypt(msg, "")
	assert.NoError(t, err)
	assert.Equal(t, msg, enc)

	_, err = e.Decrypt(enc, "password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
	dec, err := e.Decrypt(enc, "")
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)
}

func TestDefaultEncrypter(t *testing.T) {
	opts := []Option{
		OptionIteration(3),
		OptionMemory(4),
		OptionParallelism(5),
	}
	e := DefaultEncrypter(opts...)
	assert.Equal(t, "ARGON2ID-AES_256_CTR-MACV1", e.Method)
	assert.Equal(t, "3", e.Params["iterations"])
	assert.Equal(t, "4", e.Params["memory"])
	assert.Equal(t, "5", e.Params["parallelism"])
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
