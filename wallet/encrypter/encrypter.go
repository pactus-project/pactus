package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
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

// ErrNotSupported describes an error in which the encrypted method is no
// known or supported.
var ErrNotSupported = errors.New("encrypted method is not supported")

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

	data := make([]byte, 0)

	// Password hasher method
	switch funcs[0] {
	case nameFuncArgon2ID:
		salt := make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			return "", err
		}

		iterations := e.Params.GetUint32(nameParamIterations, defaultIterations)
		memory := e.Params.GetUint32(nameParamMemory, defaultMemory)
		parallelism := e.Params.GetUint8(nameParamParallelism, defaultParallelism)
		keyLen := e.Params.GetUint32(nameParamKeyLen, defaultKeyLen)

		// Argon2 currently has three modes:
		// - data-dependent Argon2d,
		// - data-independent Argon2i,
		// - a mix of the two, Argon2id.
		derivedBytes := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

		// Encrypter method
		switch funcs[1] {
		case nameFuncAES256CTR:
			cipherKey := derivedBytes[:32]
			iv := derivedBytes[32:]
			cipher := aesCTRCrypt([]byte(message), iv, cipherKey)

			// MAC method
			switch funcs[2] {
			case nameFuncMACv1:
				// Calculate the MAC
				// We use the MAC to check if the password is correct
				// https: //en.wikipedia.org/wiki/Authenticated_encryption#Encrypt-then-MAC_(EtM)
				mac := calcMACv1(cipherKey[16:32], cipher)

				data = append(data, salt...)
				data = append(data, cipher...)
				data = append(data, mac...)

			default:
				return "", ErrMethodNotSupported
			}

		default:
			return "", ErrMethodNotSupported
		}

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
	exitOnErr(err)

	var text string
	// Minimum length of data should be 20 (16 salt + 4 bytes mac)
	if len(data) < 20 {
		return "", ErrInvalidCipher
	}

	// Password hasher method
	switch funcs[0] {
	case nameFuncArgon2ID:
		salt := data[0:16]

		iterations := e.Params.GetUint32(nameParamIterations, defaultIterations)
		memory := e.Params.GetUint32(nameParamMemory, defaultMemory)
		parallelism := e.Params.GetUint8(nameParamParallelism, defaultParallelism)
		keyLen := e.Params.GetUint32(nameParamKeyLen, defaultKeyLen)

		derivedByte := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

		// Encrypter method
		switch funcs[1] {
		case nameFuncAES256CTR:
			var initVec, cipherKey []byte

			switch keyLen {
			case 0:
				// This case supports legacy encryption methods where the same salt is reused as the IV.
				cipherKey = derivedByte
				initVec = salt

			case 48:
				// The first 32 bytes are used as the encryption key, and the last 16 bytes are used as the IV.
				cipherKey = derivedByte[:32]
				initVec = derivedByte[32:]

			default:
				return "", ErrInvalidParam
			}

			enc := data[16 : len(data)-4]
			text = string(aesCTRCrypt(enc, initVec, cipherKey))

			// MAC method
			switch funcs[2] {
			case nameFuncMACv1:
				mac := data[len(data)-4:]
				if !util.SafeCmp(mac, calcMACv1(cipherKey[16:32], enc)) {
					return "", ErrInvalidPassword
				}
			default:
				return "", ErrMethodNotSupported
			}
		default:
			return "", ErrMethodNotSupported
		}

	default:
		return "", ErrMethodNotSupported
	}

	return text, nil
}

// aesCTRCrypt encrypts/decrypts a message using AES-256-CTR and
// returns the encoded/decoded bytes.
func aesCTRCrypt(message, initVec, cipherKey []byte) []byte {
	// Generate the cipher message
	cipherMsg := make([]byte, len(message))
	aesCipher, err := aes.NewCipher(cipherKey)
	exitOnErr(err)

	stream := cipher.NewCTR(aesCipher, initVec)
	stream.XORKeyStream(cipherMsg, message)

	return cipherMsg
}

// calcMACv1 calculates the 4 bytes MAC of the given slices base on SHA-256.
func calcMACv1(data ...[]byte) []byte {
	hasher := sha256.New()
	for _, d := range data {
		_, err := hasher.Write(d)
		exitOnErr(err)
	}

	return hasher.Sum(nil)[:4]
}

// exitOnErr exit the software immediately if an error happens.
// Panics are not safe because panics print a stack trace,
// which may not be relevant to the error at all.
func exitOnErr(e error) {
	if e != nil {
		os.Exit(1)
	}
}
