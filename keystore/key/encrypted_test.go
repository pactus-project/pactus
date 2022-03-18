package key

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

func TestEncryption(t *testing.T) {
	auth1 := "secret1"
	auth2 := "secret2"
	//Generates Private Key
	_, prv := bls.GenerateTestKeyPair()
	key := NewKey(prv)
	filePath := fmt.Sprintf("/tmp/%s.key", key.Address().String())
	//Encrypts the key json blob
	err := EncryptKeyToFile(key, filePath, auth1, "")
	assert.NoError(t, err)
	// Existing file
	err = EncryptKeyToFile(key, filePath, auth2, "")
	assert.NoError(t, err)
	// Invalid auth
	_, err = DecryptKeyFile(filePath, auth1)
	assert.Error(t, err)
	// Decrypts Json Object
	k2, err := DecryptKeyFile(filePath, auth2)
	assert.NoError(t, err)
	assert.Equal(t, key, k2)
	// wrong password: should fails
	k3, err := DecryptKeyFile(filePath, "Secret")
	assert.Error(t, err)
	assert.Nil(t, k3)
	// invalid file path, should fails
	filePath1 := fmt.Sprintf("/tmp/%s_invalid_path.key", key.Address().String())
	k4, err := DecryptKeyFile(filePath1, auth2)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Nil(t, k4)
}

func TestEncryptionData(t *testing.T) {
	auth := "secret"
	//Generates
	_, prv := bls.GenerateTestKeyPair()
	key := NewKey(prv)
	//Encrypts the key json blob
	ek, err := EncryptKey(key, auth, "")
	f := util.TempFilePath()
	assert.NoError(t, ek.Save(f))

	assert.NoError(t, err)
	//Decrypts Json Object
	k2, err := ek.Decrypt(auth)
	assert.NoError(t, err)
	assert.Equal(t, key, k2)
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
	_, prv := bls.GenerateTestKeyPair()
	key := NewKey(prv)
	ek, _ := EncryptKey(key, "", "")
	k2, _ := ek.Decrypt("")
	assert.Equal(t, key, k2)
}

func TestCheckLabel(t *testing.T) {
	_, prv := bls.GenerateTestKeyPair()
	key := NewKey(prv)
	label := "zarb"
	f := util.TempFilePath()
	assert.NoError(t, EncryptKeyToFile(key, f, "secret", label))
	ek, err := NewEncryptedKey(f)
	assert.NoError(t, err)
	assert.Equal(t, ek.Label, label)
}

func TestInvalidFile(t *testing.T) {
	f := util.TempFilePath()
	_, err := DecryptKeyFile(f, "")
	assert.Error(t, err)
}
