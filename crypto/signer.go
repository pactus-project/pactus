package crypto

type SignableMsg interface {
	SignBytes() []byte
	SetSignature(sig Signature)
	SetPublicKey(pub PublicKey)
}

type Signer interface {
	Address() Address
	PublicKey() PublicKey
	SignData(data []byte) Signature
	SignMsg(msg SignableMsg)
}

type signer struct {
	address    Address
	publicKey  PublicKey
	privateKey PrivateKey
}

func NewSigner(pv PrivateKey) Signer {
	return &signer{
		privateKey: pv,
		publicKey:  pv.PublicKey(),
		address:    pv.PublicKey().Address(),
	}
}

func (s *signer) Address() Address {
	return s.address
}

func (s *signer) PublicKey() PublicKey {
	return s.publicKey
}

func (s *signer) SignMsg(msg SignableMsg) {
	bz := msg.SignBytes()
	sig := s.privateKey.Sign(bz)
	msg.SetSignature(sig)
	msg.SetPublicKey(s.publicKey)
}

func (s *signer) SignData(data []byte) Signature {
	return s.privateKey.Sign(data)
}
