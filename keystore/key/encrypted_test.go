package key

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	auth := "secret"
	//Generates Private Key
	k1 := GenKey()
	filePath := fmt.Sprintf("/tmp/%s.key", k1.Address().String())
	//Encrypts the key json blob
	err := EncryptKeyFile(k1, filePath, auth, "")
	assert.NoError(t, err)
	//Decrypts Json Object
	k2, err := DecryptKeyFile(filePath, auth)
	assert.NoError(t, err)
	assert.Equal(t, k1, k2)
	// wrong password: should fails
	k3, err := DecryptKeyFile(filePath, "Secret")
	assert.Error(t, err)
	assert.Nil(t, k3)
	// invalid file path, should fails
	filePath1 := fmt.Sprintf("/tmp/%s_invalid_path.key", k1.Address().String())
	k4, err := DecryptKeyFile(filePath1, auth)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Nil(t, k4)
}

func TestEncryptionData(t *testing.T) {
	auth := "secret"
	//Generates
	k1 := GenKey()
	//Encrypts the key json blob
	bs, err := EncryptKey(k1, auth, "")

	assert.NoError(t, err)
	//Decrypts Json Object
	k2, err := DecryptKey(bs, auth)
	assert.NoError(t, err)
	assert.Equal(t, k1, k2)
	// wrong password: should fails
	k3, err := DecryptKey(bs, "Secret")
	assert.Error(t, err)
	assert.Nil(t, k3)
	//Decrypts Json Object, should fails
	bs[0] = 0 /// manipulated byte stream
	k4, err := DecryptKey(bs, auth)
	assert.Error(t, err)
	assert.Nil(t, k4)
}

func TestNonEncryptied(t *testing.T) {
	k1 := GenKey()
	bs, _ := EncryptKey(k1, "", "")
	fmt.Println(string(bs))
	k2, _ := DecryptKey(bs, "")
	assert.Equal(t, k1, k2)
}
