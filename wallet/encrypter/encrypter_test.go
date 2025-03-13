package encrypter

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestNopeEncrypterParams(t *testing.T) {
	enc := NopeEncrypter()
	assert.Equal(t, "", enc.Method)
	assert.Nil(t, enc.Params)
	assert.False(t, enc.IsEncrypted())
}

func TestNopeEncrypter(t *testing.T) {
	enc := NopeEncrypter()

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

func TestDefaultEncrypterParams(t *testing.T) {
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

func TestDefaultEncrypter(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	enc := &Encrypter{
		Method: "ARGON2ID-AES_256_CTR-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "8",
			nameParamParallelism: "1",
			nameParamKeyLen:      "48",
		},
	}

	msg := ts.RandString(ts.RandIntNonZero(100))
	password := ts.RandString(ts.RandIntNonZero(100))

	_, err := enc.Encrypt(msg, "")
	assert.ErrorIs(t, err, ErrInvalidPassword)

	cipher, err := enc.Encrypt(msg, password)
	assert.NoError(t, err)

	dec, err := enc.Decrypt(cipher, password)
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)

	_, err = enc.Decrypt(cipher, "invalid-password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
}

func TestInvalidMethod(t *testing.T) {
	tests := []struct {
		method string
	}{
		{"XXX-AES_256_CTR-MACV1"},
		{"ARGON2ID-XXX-MACV1"},
		{"ARGON2ID-AES_256_CTR-XXX"},
		{"XXX"},
	}

	for _, tt := range tests {
		enc := &Encrypter{
			Method: tt.method,
			Params: params{
				nameParamIterations:  "1",
				nameParamMemory:      "8",
				nameParamParallelism: "1",
				nameParamKeyLen:      "48",
			},
		}

		_, err := enc.Encrypt("foo", "password")
		assert.ErrorIs(t, err, ErrMethodNotSupported)

		_, err = enc.Decrypt("AJFPsGu6bDMJ5iuMWDJS/87xVs7r", "password")
		assert.ErrorIs(t, err, ErrMethodNotSupported)
	}
}

func TestInvalidDecrypt(t *testing.T) {
	enc := &Encrypter{
		Method: "ARGON2ID-AES_256_CTR-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "8",
			nameParamParallelism: "1",
			nameParamKeyLen:      "48",
		},
	}

	_, err := enc.Decrypt("", "password")
	assert.ErrorIs(t, err, ErrInvalidCipher)

	_, err = enc.Decrypt("invalid-base64", "password")
	assert.ErrorIs(t, err, ErrInvalidCipher)

	enc.Params.SetUint32(nameParamKeyLen, 64)
	_, err = enc.Decrypt("AJFPsGu6bDMJ5iuMWDJS/87xVs7r", "password")
	assert.ErrorIs(t, err, ErrInvalidParam)
}

func TestEncrypterV2(t *testing.T) {
	enc := &Encrypter{
		Method: "ARGON2ID-AES_256_CTR-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "8",
			nameParamParallelism: "1",
			nameParamKeyLen:      "32",
		},
	}

	msg := "foo"
	password := "cowboy"
	cipher := "hU/nlRNmHhKXB1tv32Ekt4ctoP7GRLw="

	dec, err := enc.Decrypt(cipher, password)
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)
}

func TestAES256CBC(t *testing.T) {
	enc := &Encrypter{
		Method: "ARGON2ID-AES_256_CBC-MACV1",
		Params: params{
			nameParamIterations:  "1",
			nameParamMemory:      "8",
			nameParamParallelism: "1",
			nameParamKeyLen:      "48",
		},
	}

	msg := "foo"
	password := "cowboy"

	_, err := enc.Encrypt(msg, "")
	assert.ErrorIs(t, err, ErrInvalidPassword)

	cipher, err := enc.Encrypt(msg, password)
	assert.NoError(t, err)

	dec, err := enc.Decrypt(cipher, password)
	assert.NoError(t, err)
	assert.Equal(t, msg, dec)

	_, err = enc.Decrypt(cipher, "invalid-password")
	assert.ErrorIs(t, err, ErrInvalidPassword)
}
