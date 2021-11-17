package crypto

type Signature interface {
	RawBytes() []byte
	String() string
	MarshalJSON() ([]byte, error)
	MarshalCBOR() ([]byte, error)
	SanityCheck() error
	EqualsTo(right Signature) bool
}
