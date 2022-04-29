package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// For testing we set the parameters to the minimum
	iterations = uint32(1)
	memory = uint32(2)
	parallelism = uint8(3)
}

func TestEncrypt(t *testing.T) {
	e := newArgon2Encrypter("super_secret_passsword")
	msg1 := "hello_world"
	ct := e.encrypt(msg1)
	assert.Equal(t, ct.Method, "ARGON2ID-AES_256_CTR-SHA256")
	msg2, err := e.decrypt(ct)
	assert.NoError(t, err)
	assert.Equal(t, msg1, msg2)
	assert.Equal(t, ct.Params.GetUint32("iterations"), iterations)
	assert.Equal(t, ct.Params.GetUint32("memory"), memory)
	assert.Equal(t, ct.Params.GetUint8("parallelism"), parallelism)

	e2 := newArgon2Encrypter("invalid_password")
	_, err = e2.decrypt(ct)
	assert.Error(t, err)

	ct.Method = "unknown-method"
	_, err = e2.decrypt(ct)
	assert.Error(t, err)
}
