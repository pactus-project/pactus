package wallet

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

type argon2Encrypter struct{}

func newArgon2Encrypter() encrypter {
	return &argon2Encrypter{}
}
func (e *argon2Encrypter) encrypt(message, passphrase string) encrypted {
	// Random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	exitOnErr(err)

	cipherKey := e.cipherKey(passphrase, salt)

	// Using salt for Initialization Vector (IV)
	iv := salt
	cipherMsg := aesCrypt([]byte(message), iv, cipherKey)

	// Generate the checksum
	checksum := sha256Checksum(cipherKey[16:32], cipherMsg)

	params := params{}
	params.setBytes("salt", salt)
	params.setBytes("checksum", checksum)

	return encrypted{
		Method:  "ARGON2ID_AES-256-CTR_SHA256",
		Params:  params,
		Message: base64.StdEncoding.EncodeToString(cipherMsg),
	}
}

func (e *argon2Encrypter) decrypt(ct encrypted, passphrase string) (string, error) {
	salt := ct.Params.getBytes("salt")
	checksum := ct.Params.getBytes("checksum")

	cipherKey := e.cipherKey(passphrase, salt)
	cipherMsg, err := base64.StdEncoding.DecodeString(ct.Message)
	exitOnErr(err)

	if !safeCmp(checksum, sha256Checksum(cipherKey[16:32], cipherMsg)) {
		return "", errors.New("invalid checksum")
	}

	msg := aesCrypt(cipherMsg, salt, cipherKey)
	return string(msg), nil
}

func (e *argon2Encrypter) cipherKey(passphrase string, salt []byte) []byte {
	iterations := uint32(1)
	parallelism := uint8(4)
	memory := uint32(2 ^ 21)

	// Argon2 currently has three modes: data-dependent Argon2d, data-independent Argon2i, and a mix of the two, Argon2id.
	// Parameters are set based on the spec recommendation
	// Read more here https://datatracker.ietf.org/doc/html/rfc9106#section-4
	return argon2.IDKey([]byte(passphrase), salt, iterations, memory, parallelism, 32)
}
