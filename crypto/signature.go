package crypto

type Signature interface {
	RawBytes() []byte
	String() string
	MarshalText() ([]byte, error)
	UnmarshalText(text []byte) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(bz []byte) error
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR(bs []byte) error
	SanityCheck() error
	EqualsTo(right Signature) bool
}
