package crypto

type PublicKey interface {
	RawBytes() []byte
	String() string
	MarshalText() ([]byte, error)
	UnmarshalText(text []byte) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(bz []byte) error
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR(bs []byte) error
	Verify(msg []byte, sig Signature) bool
	Address() Address
	SanityCheck() error
	EqualsTo(right PublicKey) bool
}
