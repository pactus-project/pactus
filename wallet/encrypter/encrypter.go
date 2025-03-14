package encrypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pactus-project/pactus/util"
	"golang.org/x/crypto/argon2"
)

// Parameters are set based on the spec recommendation
// Read more here https://datatracker.ietf.org/doc/html/rfc9106#section-4
type argon2dParameters struct {
	iterations  uint32
	memory      uint32
	parallelism uint8
	keyLen      uint32
}

type Option func(p *argon2dParameters)

func OptionIteration(iterations uint32) func(p *argon2dParameters) {
	return func(p *argon2dParameters) {
		p.iterations = iterations
	}
}

func OptionMemory(memory uint32) func(p *argon2dParameters) {
	return func(p *argon2dParameters) {
		p.memory = memory
	}
}

func OptionParallelism(parallelism uint8) func(p *argon2dParameters) {
	return func(p *argon2dParameters) {
		p.parallelism = parallelism
	}
}

const (
	nameParamIterations  = "iterations"
	nameParamMemory      = "memory"
	nameParamParallelism = "parallelism"
	nameParamKeyLen      = "keylen"

	nameFuncNope      = ""
	nameFuncArgon2ID  = "ARGON2ID"
	nameFuncAES256CTR = "AES_256_CTR"
	nameFuncAES256CBC = "AES_256_CBC"
	nameFuncMACv1     = "MACV1"

	// Parameter Choice
	// https://www.rfc-editor.org/rfc/rfc9106.html#section-4
	defaultIterations  = 3
	defaultMemory      = 65536 // 2 ^ 16
	defaultParallelism = 4
	defaultKeyLen      = 48
)

// Encrypter keeps the method and parameters for the cipher algorithm.
type Encrypter struct {
	Method string `json:"method,omitempty"`
	Params params `json:"params,omitempty"`
}

// NopeEncrypter creates a nope encrypter instance.
//
// The nope encrypter doesn't encrypt the message and the cipher message is same
// as original message.
func NopeEncrypter() Encrypter {
	return Encrypter{
		Method: nameFuncNope,
		Params: nil,
	}
}

// DefaultEncrypter creates a new encrypter instance.
// If no option sets it uses the default parameters.
//
// The default encrypter uses Argon2ID as password hasher and AES_256_CTR as
// encryption algorithm.
func DefaultEncrypter(opts ...Option) Encrypter {
	argon2dParameters := &argon2dParameters{
		iterations:  defaultIterations,
		memory:      defaultMemory,
		parallelism: defaultParallelism,
		keyLen:      defaultKeyLen,
	}
	for _, opt := range opts {
		opt(argon2dParameters)
	}

	method := fmt.Sprintf("%s-%s-%s",
		nameFuncArgon2ID, nameFuncAES256CTR, nameFuncMACv1)

	encParams := newParams()
	encParams.SetUint32(nameParamIterations, argon2dParameters.iterations)
	encParams.SetUint32(nameParamMemory, argon2dParameters.memory)
	encParams.SetUint8(nameParamParallelism, argon2dParameters.parallelism)
	encParams.SetUint32(nameParamKeyLen, argon2dParameters.keyLen)

	return Encrypter{
		Method: method,
		Params: encParams,
	}
}

func (e *Encrypter) IsEncrypted() bool {
	return e.Method != nameFuncNope
}

// Encrypt encrypts the `message` using give `password` and returns the cipher message.
func (e *Encrypter) Encrypt(message, password string) (string, error) {
	if e.Method == nameFuncNope {
		if password != "" {
			return "", ErrInvalidPassword
		}

		return message, nil
	}

	// Check if password is empty
	if password == "" {
		return "", ErrInvalidPassword
	}

	funcs := strings.Split(e.Method, "-")
	if len(funcs) != 3 {
		return "", ErrMethodNotSupported
	}

	// Password hasher method
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	iterations := e.Params.GetUint32(nameParamIterations)
	memory := e.Params.GetUint32(nameParamMemory)
	parallelism := e.Params.GetUint8(nameParamParallelism)
	keyLen := e.Params.GetUint32(nameParamKeyLen)

	if keyLen == 32 {
		// Legacy encryption methods where the same salt is reused as the IV.
		// Update to 48 byes key length.
		keyLen = 48
		e.Params.SetUint32(nameParamKeyLen, 48)
	}

	// Password hasher method
	var passwordHash []byte
	switch funcs[0] {
	case nameFuncArgon2ID:
		// Argon2 currently has three modes:
		// - data-dependent Argon2d,
		// - data-independent Argon2i,
		// - a mix of the two, Argon2id.
		passwordHash = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)
	default:
		return "", ErrMethodNotSupported
	}

	// Encrypter method
	// The first 32 bytes are used as the encryption key, and the last 16 bytes are used as the IV.
	cipherKey := passwordHash[:32]
	initVec := passwordHash[32:]
	var cipher []byte
	switch funcs[1] {
	case nameFuncAES256CTR:
		cipher = aes256CTRCrypt([]byte(message), initVec, cipherKey)

	case nameFuncAES256CBC:
		cipher = aes256CBCEncrypt([]byte(message), initVec, cipherKey)

	default:
		return "", ErrMethodNotSupported
	}

	// MAC method
	data := make([]byte, 0)
	switch funcs[2] {
	case nameFuncMACv1:
		// Calculate the MAC
		// We use the MAC to check if the password is correct.
		// Use Cipher key as the key for the MAC.
		// https: //en.wikipedia.org/wiki/Authenticated_encryption#Encrypt-then-MAC_(EtM)
		mac := calcMACv1(cipherKey[16:32], cipher)

		data = append(data, salt...)
		data = append(data, cipher...)
		data = append(data, mac...)

	default:
		return "", ErrMethodNotSupported
	}

	cipherText := base64.StdEncoding.EncodeToString(data)

	return cipherText, nil
}

// Decrypt decrypts the `cipher` using give `password` and returns the original message.
func (e *Encrypter) Decrypt(cipherText, password string) (string, error) {
	if e.Method == nameFuncNope {
		if password != "" {
			return "", ErrInvalidPassword
		}

		return cipherText, nil
	}

	funcs := strings.Split(e.Method, "-")
	if len(funcs) != 3 {
		return "", ErrMethodNotSupported
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", ErrInvalidCipher
	}

	// Minimum length of data should be 20 (16 salt + 4 bytes mac)
	if len(data) < 20 {
		return "", ErrInvalidCipher
	}

	iterations := e.Params.GetUint32(nameParamIterations)
	memory := e.Params.GetUint32(nameParamMemory)
	parallelism := e.Params.GetUint8(nameParamParallelism)
	keyLen := e.Params.GetUint32(nameParamKeyLen)

	// Password hasher method
	salt := data[0:16]
	var passwordHash []byte

	switch funcs[0] {
	case nameFuncArgon2ID:
		passwordHash = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

	default:
		return "", ErrMethodNotSupported
	}

	// Encrypter method
	var initVec, cipherKey []byte

	switch keyLen {
	case 32:
		// This case supports legacy encryption methods where the same salt is reused as the IV.
		cipherKey = passwordHash
		initVec = salt

	case 48:
		// The first 32 bytes are used as the encryption key, and the last 16 bytes are used as the IV.
		cipherKey = passwordHash[:32]
		initVec = passwordHash[32:]

	default:
		return "", ErrInvalidParam
	}

	cipher := data[16 : len(data)-4]
	var msg []byte

	switch funcs[1] {
	case nameFuncAES256CTR:
		msg = aes256CTRCrypt(cipher, initVec, cipherKey)

	case nameFuncAES256CBC:
		msg = aes256CBCDecrypt(cipher, initVec, cipherKey)

	default:
		return "", ErrMethodNotSupported
	}

	// MAC method
	switch funcs[2] {
	case nameFuncMACv1:
		mac := data[len(data)-4:]
		if !util.SafeCmp(mac, calcMACv1(cipherKey[16:32], cipher)) {
			return "", ErrInvalidPassword
		}

	default:
		return "", ErrMethodNotSupported
	}

	return string(msg), nil
}

// aes256CTRCrypt encrypts or decrypts a message using AES-256-CTR mode.
// It requires a 32-byte (256-bit) cipher key and a 16-byte (128-bit) initialization vector (IV).
// Returns the encrypted or decrypted output.
func aes256CTRCrypt(input, initVec, cipherKey []byte) []byte {
	aesCipher, _ := aes.NewCipher(cipherKey)

	output := make([]byte, len(input))
	stream := cipher.NewCTR(aesCipher, initVec)
	stream.XORKeyStream(output, input)

	return output
}

// aes256CBCEncrypt encrypts a plain message using AES-256-CBC mode.
// It requires a 32-byte (256-bit) cipher key and a 16-byte (128-bit) initialization vector (IV).
// Returns the encrypted cipher message.
func aes256CBCEncrypt(plainMsg, initVec, cipherKey []byte) []byte {
	aesCipher, _ := aes.NewCipher(cipherKey)

	plainMsg = pkcs7Padding(plainMsg, aes.BlockSize)
	cipherMsg := make([]byte, len(plainMsg))
	enc := cipher.NewCBCEncrypter(aesCipher, initVec)
	enc.CryptBlocks(cipherMsg, plainMsg)

	return cipherMsg
}

// aes256CBCDecrypt decrypts a cipher message using AES-256-CBC mode.
// It requires a 32-byte (256-bit) cipher key and a 16-byte (128-bit) initialization vector (IV).
// Returns the decrypted plain message.
func aes256CBCDecrypt(cipherMsg, initVec, cipherKey []byte) []byte {
	aesCipher, _ := aes.NewCipher(cipherKey)

	plainMsg := make([]byte, len(cipherMsg))
	dec := cipher.NewCBCDecrypter(aesCipher, initVec)
	dec.CryptBlocks(plainMsg, cipherMsg)
	plainMsg = pkcs7UnPadding(plainMsg)

	return plainMsg
}

func pkcs7Padding(cipherMsg []byte, blockSize int) []byte {
	padding := blockSize - len(cipherMsg)%blockSize
	padMsg := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(cipherMsg, padMsg...)
}

func pkcs7UnPadding(plainMsg []byte) []byte {
	length := len(plainMsg)
	unpadding := int(plainMsg[length-1])
	if length-unpadding <= 0 {
		return plainMsg
	}

	return plainMsg[:(length - unpadding)]
}

// calcMACv1 calculates the 4 bytes MAC of the given slices base on SHA-256.
func calcMACv1(data ...[]byte) []byte {
	hasher := sha256.New()
	for _, d := range data {
		_, _ = hasher.Write(d)
	}

	return hasher.Sum(nil)[:4]
}
