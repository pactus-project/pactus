package key

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const (
	//Scrypt parameters
	keyHeaderKDF = "scrypt"
	scryptDKLen  = 32

	scryptN = 2
	scryptP = 1
	scryptR = 8

	version = 3
)

type EncryptedKey struct {
	Address    crypto.Address     `json:"address"`
	Crypto     *cryptoJSON        `json:"crypto,omitempty"`
	PrivateKey *crypto.PrivateKey `json:"privatekey,omitempty"`
	Label      string             `json:"label,omitempty"`
	Version    int                `json:"version"`
}

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

// NewEncryptedKey reads the encrypted file and returns an instance of EncryptedKey
func NewEncryptedKey(path string) (*EncryptedKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ek := new(EncryptedKey)
	if err := json.Unmarshal(data, ek); err != nil {
		return nil, err
	}
	return ek, nil
}

// DecryptKeyFile decrypts the file and returns Key
func DecryptKeyFile(path, auth string) (*Key, error) {
	ek, err := NewEncryptedKey(path)
	if err != nil {
		return nil, err
	}
	return ek.Decrypt(auth)
}

// Save saves the encrypted key into file
func (ek *EncryptedKey) Save(p string) error {
	j, err := json.Marshal(ek)
	if err != nil {
		return err
	}
	return util.WriteFile(p, j)
}

// Decrypt decrypts the Key from a json blob and returns the plaintext of the private key
func (ek *EncryptedKey) Decrypt(auth string) (*Key, error) {
	if ek.PrivateKey != nil {
		return NewKey(ek.Address, *ek.PrivateKey)
	}
	if ek.Crypto.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("cipher not supported: %v", ek.Crypto.Cipher)
	}
	mac, err := hex.DecodeString(ek.Crypto.MAC)
	if err != nil {
		return nil, err
	}
	iv, err := hex.DecodeString(ek.Crypto.CipherParams.IV)
	if err != nil {
		return nil, err
	}
	cipherText, err := hex.DecodeString(ek.Crypto.CipherText)
	if err != nil {
		return nil, err
	}
	derivedKey, err := getKDFKey(ek.Crypto, auth)
	if err != nil {
		return nil, err
	}
	calculatedMAC := hash.Hash256(append(derivedKey[16:32], cipherText...))
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, fmt.Errorf("could not decrypt key with given passphrase")
	}
	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	pv, err := bls.PrivateKeyFromRawBytes(plainText)
	if err != nil {
		return nil, err
	}
	return NewKey(ek.Address, pv)

}

func getKDFKey(cryptoJSON *cryptoJSON, auth string) ([]byte, error) {

	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)

	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("unsupported KDF: %s", cryptoJSON.KDF)
}

// EncryptKeyToFile encrypts the key and saves it to file
func EncryptKeyToFile(key *Key, filePath, auth, label string) error {
	ek, err := EncryptKey(key, auth, label)
	if err != nil {
		return err
	}
	return ek.Save(filePath)
}

// EncryptKey encrypts a key and returns the encrypted byte array
func EncryptKey(key *Key, auth, label string) (*EncryptedKey, error) {
	if auth == "" {
		pv := key.PrivateKey()
		return &EncryptedKey{
			Address:    key.data.Address,
			PrivateKey: &pv,
			Label:      label,
			Version:    version,
		}, nil
	}

	authArray := []byte(auth)
	salt := getEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}

	encryptKey := derivedKey[:16]
	keyBytes := key.PrivateKey().RawBytes()

	iv := getEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return nil, err
	}
	mac := hash.Hash256(append(derivedKey[16:32], cipherText...))

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}
	cryptoStruct := &cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}

	return &EncryptedKey{
		Address: key.data.Address,
		Crypto:  cryptoStruct,
		Label:   label,
		Version: version,
	}, nil
}

func getEntropyCSPRNG(n int) []byte {
	mainBuff := make([]byte, n)
	_, err := io.ReadFull(crand.Reader, mainBuff)
	if err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return mainBuff
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}
