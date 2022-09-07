package crypto

type PrivateKey interface {
	Bytes() []byte
	String() string
	Sign(msg []byte) Signature
	PublicKey() PublicKey
	EqualsTo(right PrivateKey) bool
}
