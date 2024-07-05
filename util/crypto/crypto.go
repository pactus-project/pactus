package crypto

import "github.com/pactus-project/pactus/crypto/bls"

func SignMessageWithPrivateKey(prv, msg string) (string, error) {
	prvKey, err := bls.PrivateKeyFromString(prv)
	if err != nil {
		return "", err
	}

	return prvKey.Sign([]byte(msg)).String(), nil
}

func VerifyMessage(sigStr, pubStr, msg string) bool {
	sig, err := bls.SignatureFromString(sigStr)
	if err != nil {
		return false
	}

	pub, err := bls.PublicKeyFromString(pubStr)
	if err != nil {
		return false
	}

	if err := pub.Verify([]byte(msg), sig); err != nil {
		return false
	}

	return true
}
