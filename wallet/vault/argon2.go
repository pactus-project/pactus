package vault

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Parameters are set based on the spec recommendation
// Read more here https://datatracker.ietf.org/doc/html/rfc9106#section-4
var (
	iterations  = uint32(1)
	memory      = uint32(2 * 1024 * 1024)
	parallelism = uint8(4)
)

type argon2Encrypter struct {
	method   string
	password string
}

const (
	nameParamIterations  = "iterations"
	nameParamMemory      = "memory"
	nameParamParallelism = "parallelism"
	nameParamSalt        = "salt"
	nameParamMAC         = "mac"

	nameFuncArgon2ID  = "ARGON2ID"
	nameFuncAES256CTR = "AES_256_CTR"
	nameFuncMACv1     = "MACV1"
)

func newArgon2Encrypter(password string) *argon2Encrypter {
	method := fmt.Sprintf("%s-%s-%s",
		nameFuncArgon2ID, nameFuncAES256CTR, nameFuncMACv1)

	return &argon2Encrypter{
		method:   method,
		password: password,
	}
}
func (e *argon2Encrypter) encrypt(message string) encrypted {
	// Random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	exitOnErr(err)

	cipherKey := e.cipherKey(e.password, salt, iterations, memory, parallelism)

	// Using salt for Initialization Vector (IV)
	iv := salt
	d := aesCrypt([]byte(message), iv, cipherKey)

	// Calculate the MAC
	mac := calcMACv1(cipherKey[16:32], d)

	params := newParams()
	params.SetUint32(nameParamIterations, iterations)
	params.SetUint32(nameParamMemory, memory)
	params.SetUint8(nameParamParallelism, parallelism)
	params.SetBytes(nameParamSalt, salt)
	params.SetBytes(nameParamMAC, mac)

	cipherText := base64.StdEncoding.EncodeToString(d)

	return encrypted{
		Method:     e.method,
		Params:     params,
		CipherText: cipherText,
	}
}

func (e *argon2Encrypter) decrypt(ct encrypted) (string, error) {
	if ct.Method != e.method {
		return "", NewErrUnknownMethod(ct.Method)
	}

	iterations := ct.Params.GetUint32(nameParamIterations)
	memory := ct.Params.GetUint32(nameParamMemory)
	parallelism := ct.Params.GetUint8(nameParamParallelism)
	salt := ct.Params.GetBytes(nameParamSalt)
	mac := ct.Params.GetBytes(nameParamMAC)

	cipherKey := e.cipherKey(e.password, salt, iterations, memory, parallelism)
	d, err := base64.StdEncoding.DecodeString(ct.CipherText)
	exitOnErr(err)

	// Using MAC to check if the password is correct
	// https: //en.wikipedia.org/wiki/Authenticated_encryption#Encrypt-then-MAC_(EtM)
	if !safeCmp(mac, calcMACv1(cipherKey[16:32], d)) {
		return "", ErrInvalidPassword
	}

	text := aesCrypt(d, salt, cipherKey)
	return string(text), nil
}

func (e *argon2Encrypter) cipherKey(password string, salt []byte, iterations, memory uint32, parallelism uint8) []byte {
	// Argon2 currently has three modes:
	// - data-dependent Argon2d,
	// - data-independent Argon2i,
	// - a mix of the two, Argon2id.
	return argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, 32)
}
