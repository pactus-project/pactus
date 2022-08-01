package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/zarbchain/zarb-go/util"
	"golang.org/x/crypto/argon2"
)

const (
	nameParamIterations  = "iterations"
	nameParamMemory      = "memory"
	nameParamParallelism = "parallelism"

	nameFuncNope      = "NOPE"
	nameFuncArgon2ID  = "ARGON2ID"
	nameFuncAES256CTR = "AES_256_CTR"
	nameFuncMACv1     = "MACV1"
)

var (
	// ErrNotSupported describes an error in which the encrypted method is no
	// known or supported.
	ErrNotSupported = errors.New("encrypted method is not supported")
)

// encrypter keeps the the method and parameters for the cipher algorithm.
type Encrypter struct {
	Method string `json:"method,omitempty"`
	Params params `json:"params"`
}

// NewNopeEncrypter creates new instance of `Encrypter` that has no encryptions method.
// The cipher message is same as original message.
func NewNopeEncrypter() *Encrypter {
	return &Encrypter{
		Method: nameFuncNope,
		Params: nil,
	}
}

// NewDefaultEncrypter creates new instance of `Encrypter` that has use Argon2ID as
// password hasher and AES_256_CTR as encryption algorithm.
func NewDefaultEncrypter() *Encrypter {
	method := fmt.Sprintf("%s-%s-%s",
		nameFuncArgon2ID, nameFuncAES256CTR, nameFuncMACv1)

	// Parameters are set based on the spec recommendation
	// Read more here https://datatracker.ietf.org/doc/html/rfc9106#section-4
	var (
		iterations  = uint32(1)
		memory      = uint32(2 * 1024 * 1024)
		parallelism = uint8(4)
	)

	params := newParams()
	params.SetUint32(nameParamIterations, iterations)
	params.SetUint32(nameParamMemory, memory)
	params.SetUint8(nameParamParallelism, parallelism)

	return &Encrypter{
		Method: method,
		Params: params,
	}
}

// Encrypt encrypts the `message` using give `password` and returns the cipher message
func (e *Encrypter) Encrypt(message string, password string) (string, error) {
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

		iterations := e.Params.GetUint32(nameParamIterations)
		memory := e.Params.GetUint32(nameParamMemory)
		parallelism := e.Params.GetUint8(nameParamParallelism)

		// Argon2 currently has three modes:
		// - data-dependent Argon2d,
		// - data-independent Argon2i,
		// - a mix of the two, Argon2id.
		cipherKey := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, 32)

		// Encrypter method
		switch funcs[1] {
		case nameFuncAES256CTR:
			{
				// Using salt for Initialization Vector (IV)
				iv := salt
				ct := aesCrypt([]byte(message), iv, cipherKey)

				// MAC method
				switch funcs[2] {
				case nameFuncMACv1:
					{
						// Calculate the MAC
						// We use the MAC to check if the password is correct
						// https: //en.wikipedia.org/wiki/Authenticated_encryption#Encrypt-then-MAC_(EtM)
						mac := calcMACv1(cipherKey[16:32], ct)

						data = append(data, salt...)
						data = append(data, ct...)
						data = append(data, mac...)
					}
				default:
					return "", ErrMethodNotSupported
				}

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

// Decrypt decrypts the `cipher` using give `password` and returns the original message
func (e *Encrypter) Decrypt(cipher string, password string) (string, error) {
	if e.Method == nameFuncNope {
		if password != "" {
			return "", ErrInvalidPassword
		}
		return cipher, nil
	}

	funcs := strings.Split(e.Method, "-")
	if len(funcs) != 3 {
		return "", ErrMethodNotSupported
	}

	data, err := base64.StdEncoding.DecodeString(cipher)
	util.ExitOnErr(err)

	var text = ""
	// Minimum length of data should be 20 (16 salt + 4 bytes mac)
	if len(data) < 20 {
		return "", ErrInvalidCipher
	}

	// Password hasher method
	switch funcs[0] {
	case nameFuncArgon2ID:
		salt := data[0:16]

		iterations := e.Params.GetUint32(nameParamIterations)
		memory := e.Params.GetUint32(nameParamMemory)
		parallelism := e.Params.GetUint8(nameParamParallelism)

		cipherKey := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, 32)

		// Encrypter method
		switch funcs[1] {
		case nameFuncAES256CTR:
			{
				iv := salt
				enc := data[16 : len(data)-4]
				text = string(aesCrypt(enc, iv, cipherKey))

				// MAC method
				switch funcs[2] {
				case nameFuncMACv1:
					{
						mac := data[len(data)-4:]
						if !util.SafeCmp(mac, calcMACv1(cipherKey[16:32], enc)) {
							return "", ErrInvalidPassword
						}
					}
				default:
					return "", ErrMethodNotSupported
				}

			}
		default:
			return "", ErrMethodNotSupported
		}

	default:
		return "", ErrMethodNotSupported
	}

	return string(text), nil
}

// aesCrypt encrypts/decrypts a message using AES-256-CTR and
// returns the encoded/decoded bytes.
func aesCrypt(message []byte, iv, cipherKey []byte) []byte {
	// Generate the cipher message
	cipherMsg := make([]byte, len(message))
	aesCipher, err := aes.NewCipher(cipherKey)
	util.ExitOnErr(err)

	stream := cipher.NewCTR(aesCipher, iv)
	stream.XORKeyStream(cipherMsg, message)

	return cipherMsg
}

// calcMACv1 calculates the 4 bytes MAC of the given slices base on SHA-256.
func calcMACv1(data ...[]byte) []byte {
	h := sha256.New()
	for _, d := range data {
		_, err := h.Write(d)
		util.ExitOnErr(err)
	}

	return h.Sum(nil)[:4]
}
