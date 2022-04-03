package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	e := newArgon2Encrypter("super_secret_passsword")
	msg1 := "hello_world"
	ct := e.encrypt(msg1)
	assert.Equal(t, ct.Method, "ARGON2ID_AES-256-CTR_SHA256")
	msg2, err := e.decrypt(ct)
	assert.NoError(t, err)
	assert.Equal(t, msg1, msg2)
	assert.Equal(t, ct.Params.GetUint32("iterations"), uint32(1))
	assert.Equal(t, ct.Params.GetUint32("memory"), uint32(0x200000))
	assert.Equal(t, ct.Params.GetUint8("parallelism"), uint8(4))

	e2 := newArgon2Encrypter("invalid_password")
	_, err = e2.decrypt(ct)
	assert.Error(t, err)
}
