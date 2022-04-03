package wallet

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

func newEncrypter(passphrase string, net int) encrypter {
	if len(passphrase) == 0 {
		return newNopeEncrypter()
	}
	if net == 2 {
		return newTestEncrypter()
	}
	return newArgon2Encrypter(passphrase)
}

type testEncrypter struct{}

func newTestEncrypter() encrypter {
	return &testEncrypter{}
}

func (e *testEncrypter) encrypt(message string) encrypted {
	return encrypted{
		CipherText: message,
	}
}

func (e *testEncrypter) decrypt(ct encrypted) (string, error) {
	return ct.CipherText, nil
}
