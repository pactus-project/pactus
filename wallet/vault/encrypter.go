package vault

// encrypted keeps the cipher text and the method parameters for the
// chiper algorithm.
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
		Method:     "",
		Params:     newParams(),
		CipherText: message,
	}
}

func (e *nopeEncrypter) decrypt(ct encrypted) (string, error) {
	if ct.Method != "" {
		return "", ErrInvalidPassword
	}
	return ct.CipherText, nil
}

func newEncrypter(password string) encrypter {
	if len(password) == 0 {
		return newNopeEncrypter()
	}
	return newArgon2Encrypter(password)
}
