package crypto

// type SignableMsg interface {
// 	SignBytes() []byte
// 	SetSignature(sig Signature)
// 	SetPublicKey(pub PublicKey)
// }

// type Signer interface {
// 	PublicKey() PublicKey
// 	SignData(data []byte) Signature
// 	SignMsg(msg SignableMsg)
// }

// type signer struct {
// 	publicKey  PublicKey
// 	privateKey PrivateKey
// }

// func NewSigner(pv PrivateKey) Signer {
// 	return &signer{
// 		privateKey: pv,
// 		// publicKey:  pv.PublicKey(),
// 	}
// }

// func (s *signer) PublicKey() PublicKey {
// 	return s.publicKey
// }

// func (s *signer) SignMsg(msg SignableMsg) {
// 	sig := s.SignData(msg.SignBytes())
// 	msg.SetSignature(sig)
// 	msg.SetPublicKey(s.publicKey)
// }

// func (s *signer) SignData(data []byte) Signature {
// 	return s.privateKey.Sign(data)
// }
