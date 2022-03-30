package wallet

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
		Message: message,
	}
}

func (e *nopeEncrypter) decrypt(ct encrypted, passphrase string) (string, error) {
	return ct.Message, nil
}
