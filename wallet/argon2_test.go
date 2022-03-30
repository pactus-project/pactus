package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	e := argon2Encrypter{}
	pass := "super_secret_passsword"
	msg1 := "hello_world"
	ct := e.encrypt(msg1, pass)
	assert.Equal(t, ct.Method, "ARGON2ID_AES-256-CTR_SHA256")
	msg2, err := e.decrypt(ct, pass)
	assert.NoError(t, err)
	assert.Equal(t, msg1, msg2)

	_, err = e.decrypt(ct, "invalid_password")
	assert.Error(t, err)
}
