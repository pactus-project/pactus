package crypto

type Signable interface {
	SignBytes() []byte
	SetSignature(sig *Signature)
}
type Signer struct {
	address    Address
	publicKey  PublicKey
	privateKey PrivateKey
}

func NewSigner(pv PrivateKey) Signer {
	return Signer{
		privateKey: pv,
		publicKey:  pv.PublicKey(),
		address:    pv.PublicKey().Address(),
	}
}

func (s Signer) Address() Address {
	return s.address
}

func (s *Signer) PublicKey() PublicKey {
	return s.publicKey
}

func (s *Signer) SignMsg(msg Signable) {
	bz := msg.SignBytes()
	sig := s.privateKey.Sign(bz)
	msg.SetSignature(sig)
}

func (s *Signer) Sign(data []byte) *Signature {
	return s.privateKey.Sign(data)
}
