package crypto

type PrivateKey interface {
	String() string
	SanityCheck() error
	Sign(msg []byte) Signature
	PublicKey() PublicKey
	EqualsTo(right PrivateKey) bool
}
