package crypto

type PrivateKey interface {
	RawBytes() []byte
	String() string
	SanityCheck() error
	Sign(msg []byte) Signature
	PublicKey() PublicKey
	EqualsTo(right PrivateKey) bool
}
