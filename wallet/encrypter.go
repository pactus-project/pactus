package wallet

import "encoding/json"

/// cipher text
type encrypted struct {
	Method     string `json:"method,omitempty"`
	Params     params `json:"params,omitempty"`
	CipherText string `json:"ct"`
}

type encrypter interface {
	encrypt(message string) encrypted
	decrypt(ct encrypted) (string, error)
}

type nopeEncrypter struct{}

func newNopeEncrypter() encrypter {
	return &nopeEncrypter{}
}

func (e *nopeEncrypter) encrypt(message string) encrypted {
	return encrypted{
		CipherText: message,
	}
}

func (e *nopeEncrypter) decrypt(ct encrypted) (string, error) {
	return ct.CipherText, nil
}

func encryptInterface(e encrypter, i interface{}) encrypted {
	d, err := json.Marshal(i)
	exitOnErr(err)
	return e.encrypt(string(d))
}
func decryptInterface(e encrypter, ct encrypted, i interface{}) error {
	s, err := e.decrypt(ct)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(s), i)
	exitOnErr(err)
	return nil
}

func newEncrypter(passphrase string) encrypter {
	if len(passphrase) == 0 {
		return newNopeEncrypter()
	}
	return newArgon2Encrypter(passphrase)
}
