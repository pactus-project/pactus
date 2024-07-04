package crypto

import "github.com/pactus-project/pactus/crypto/bls"

func SignMessageWithPrivateKey(prv string, msg string) (string, error) {
	prvKey, err := bls.PrivateKeyFromString(prv)
	if err != nil {
		return "", err
	}

	return prvKey.Sign([]byte(msg)).String(), nil
}

func VerifyMessage(sigStr string, pubStr string, msg string) (bool, error) {
	sig, err := bls.SignatureFromString(sigStr)
	if err != nil {
		return false, err
	}

	pub, err := bls.PublicKeyFromString(pubStr)
	if err != nil {
		return false, err
	}

	if err = pub.Verify([]byte(msg), sig); err != nil {
		return false, err
	}

	return true, nil
}
