package wallet

import "encoding/json"

type encrypter interface {
	encrypt(message, passphrase string) encrypted
	decrypt(ct encrypted, passphrase string) (string, error)
}

type nopeEncrypter struct{}

func newNopeEncrypter() encrypter {
	return &nopeEncrypter{}
}

func (e *nopeEncrypter) encrypt(message, passphrase string) encrypted {
	return encrypted{
		CipherText: message,
	}
}

func (e *nopeEncrypter) decrypt(ct encrypted, passphrase string) (string, error) {
	return ct.CipherText, nil
}

func encryptInterface(e encrypter, i interface{}, passphrase string) encrypted {
	d, err := json.Marshal(i)
	exitOnErr(err)
	return e.encrypt(string(d), passphrase)
}
func decryptInterface(e encrypter, ct encrypted, passphrase string, i interface{}) error {
	s, err := e.decrypt(ct, passphrase)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(s), i)
	exitOnErr(err)
	return nil
}
