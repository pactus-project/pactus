package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"os"
)

/// aesCrypt encrypts/decrypts a message using AES-256-CTR and
/// returns the encoded/decoded bytes.
func aesCrypt(message []byte, iv, cipherKey []byte) []byte {
	// Generate the cipher message
	cipherMsg := make([]byte, len(message))
	aesCipher, err := aes.NewCipher(cipherKey)
	exitOnErr(err)

	stream := cipher.NewCTR(aesCipher, iv)
	stream.XORKeyStream(cipherMsg, message)

	return cipherMsg
}

/// sha256MAC calculates the MAC of the given slices base on SHA-256
func sha256MAC(data ...[]byte) []byte {
	h := sha256.New()
	for _, d := range data {
		_, err := h.Write(d)
		exitOnErr(err)
	}

	return h.Sum(nil)[:4]
}

/// safeCmp compares two slices with constant time.
/// Note that we are using the subtle.ConstantTimeCompare() function for this
/// to help prevent timing attacks
func safeCmp(s1, s2 []byte) bool {
	return subtle.ConstantTimeCompare(s1, s2) == 1
}

/// exitOnErr exit the software immediately if an error happens.
/// Panics are not safe because panics print a stack trace,
/// which may not be relevant to the error at all.
func exitOnErr(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}
