package wallet

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

type argon2Encrypter struct {
	// TODO: memory safety
	passphrase string
}

func newArgon2Encrypter(passphrase string) encrypter {
	return &argon2Encrypter{
		passphrase: passphrase,
	}
}
func (e *argon2Encrypter) encrypt(message string) encrypted {
	// Random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	exitOnErr(err)

	// Parameters are set based on the spec recommendation
	// Read more here https://datatracker.ietf.org/doc/html/rfc9106#section-4
	iterations := uint32(1)
	memory := uint32(2 ^ 21)
	parallelism := uint8(4)

	cipherKey := e.cipherKey(e.passphrase, salt, iterations, memory, parallelism)

	// Using salt for Initialization Vector (IV)
	iv := salt
	d := aesCrypt([]byte(message), iv, cipherKey)

	// Generate the checksum
	checksum := sha256Checksum(cipherKey[16:32], d)

	params := newParams()
	params.setUint32("iterations", iterations)
	params.setUint32("iterations", iterations)
	params.setUint32("memory", memory)
	params.setUint8("parallelism", parallelism)
	params.setBytes("salt", salt)
	params.setBytes("checksum", checksum)

	cipherText := base64.StdEncoding.EncodeToString(d)

	return encrypted{
		Method:     "ARGON2ID_AES-256-CTR_SHA256",
		Params:     params,
		CipherText: cipherText,
	}
}

func (e *argon2Encrypter) decrypt(ct encrypted) (string, error) {
	salt := ct.Params.getBytes("salt")
	checksum := ct.Params.getBytes("checksum")
	iterations := ct.Params.getUint32("iterations")
	memory := ct.Params.getUint32("memory")
	parallelism := ct.Params.getUint8("parallelism")

	cipherKey := e.cipherKey(e.passphrase, salt, iterations, memory, parallelism)
	d, err := base64.StdEncoding.DecodeString(ct.CipherText)
	exitOnErr(err)

	if !safeCmp(checksum, sha256Checksum(cipherKey[16:32], d)) {
		return "", errors.New("invalid checksum")
	}

	text := aesCrypt(d, salt, cipherKey)
	return string(text), nil
}

func (e *argon2Encrypter) cipherKey(passphrase string, salt []byte, iterations, memory uint32, parallelism uint8) []byte {
	// Argon2 currently has three modes: data-dependent Argon2d, data-independent Argon2i, and a mix of the two, Argon2id.
	return argon2.IDKey([]byte(passphrase), salt, iterations, memory, parallelism, 32)
}
