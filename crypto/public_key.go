package crypto

type PublicKey interface {
	RawBytes() []byte
	String() string
	MarshalJSON() ([]byte, error)
	MarshalCBOR() ([]byte, error)
	SanityCheck() error
	Verify(msg []byte, sig Signature) bool
	Address() Address
	EqualsTo(right PublicKey) bool
}
