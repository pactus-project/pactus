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

func TestDecrypt(t *testing.T) {
	msg := "foo"
	password := "cowboy"

	tests := []struct {
		name   string
		enc    Encrypter
		cipher string
	}{
		{
			name: "Legacy cipher with keylen 32",
			enc: Encrypter{
				Method: "ARGON2ID-AES_256_CTR-MACV1",
				Params: params{
					nameParamIterations:  "1",
					nameParamMemory:      "8",
					nameParamParallelism: "1",
					nameParamKeyLen:      "32",
				},
			},
			// 854fe79513661e1297075b6fdf6124b7 872da0 fec644bc
			cipher: "hU/nlRNmHhKXB1tv32Ekt4ctoP7GRLw=",
		},
		{
			name: "Stream Cipher with keylen 48",
			enc: Encrypter{
				Method: "ARGON2ID-AES_256_CTR-MACV1",
				Params: params{
					nameParamIterations:  "1",
					nameParamMemory:      "8",
					nameParamParallelism: "1",
					nameParamKeyLen:      "48",
				},
			},
			// f000c3271de35c14162b0ddceeac0492 b4da66 b04cae34
			cipher: "8ADDJx3jXBQWKw3c7qwEkrTaZrBMrjQ=",
		},
		{
			name: "Block Cipher with keylen 48",
			enc: Encrypter{
				Method: "ARGON2ID-AES_256_CBC-MACV1",
				Params: params{
					nameParamIterations:  "1",
					nameParamMemory:      "8",
					nameParamParallelism: "1",
					nameParamKeyLen:      "48",
				},
			},
			// b14b280a9e5f7907671a405b8b8c918a 639f113cef2b6d37e7d590678b5d5a45 b422f34e
			cipher: "sUsoCp5feQdnGkBbi4yRimOfETzvK20359WQZ4tdWkW0IvNO",
		},
	}

	for _, tt := range tests {
		dec, err := tt.enc.Decrypt(tt.cipher, password)
		assert.NoError(t, err)
		assert.Equal(t, msg, dec)
	}
}
