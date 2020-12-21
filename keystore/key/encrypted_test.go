package key

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestEncryption(t *testing.T) {
	auth := "secret"
	//Generates Private Key
	k1 := GenKey()
	filePath := fmt.Sprintf("/tmp/%s.key", k1.Address().String())
	//Encrypts the key json blob
	err := EncryptKeyToFile(k1, filePath, auth, "")
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
	ek, err := EncryptKey(k1, auth, "")

	assert.NoError(t, err)
	//Decrypts Json Object
	k2, err := ek.Decrypt(auth)
	assert.NoError(t, err)
	assert.Equal(t, k1, k2)
	// wrong password: should fails
	k3, err := ek.Decrypt("Secret")
	assert.Error(t, err)
	assert.Nil(t, k3)

	ek.Crypto.Cipher = "invalid" /// manipulated key data
	k4, err := ek.Decrypt(auth)
	assert.Error(t, err)
	assert.Nil(t, k4)

	ek.Crypto.CipherParams.IV = "invalid" /// manipulated key data
	k4, err = ek.Decrypt(auth)
	assert.Error(t, err)
	assert.Nil(t, k4)

	ek.Crypto.CipherText = "invalid" /// manipulated key data
	k4, err = ek.Decrypt(auth)
	assert.Error(t, err)
	assert.Nil(t, k4)

	ek.Crypto.KDF = "invalid" /// manipulated key data
	k4, err = ek.Decrypt(auth)
	assert.Error(t, err)
	assert.Nil(t, k4)

	ek.Crypto.MAC = "invalid" /// manipulated key data
	k4, err = ek.Decrypt(auth)
	assert.Error(t, err)
	assert.Nil(t, k4)

	ek.Crypto.KDFParams = nil /// manipulated key data
	k4, err = ek.Decrypt(auth)
	assert.Error(t, err)
	assert.Nil(t, k4)
}

func TestNonEncryptied(t *testing.T) {
	k1 := GenKey()
	ek, _ := EncryptKey(k1, "", "")
	k2, _ := ek.Decrypt("")
	assert.Equal(t, k1, k2)
}

func TestCheckLabel(t *testing.T) {
	k1 := GenKey()
	label := "zarb"
	f := util.TempFilePath()
	assert.NoError(t, EncryptKeyToFile(k1, f, "secret", label))
	ek, err := NewEncryptedKey(f)
	assert.NoError(t, err)
	assert.Equal(t, ek.Label, label)
}

func TestInvalidFile(t *testing.T) {
	f := util.TempFilePath()
	_, err := DecryptKeyFile(f, "")
	assert.Error(t, err)
}
