package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
)

type Type uint8

const (
	PayloadTypeSend      = Type(1)
	PayloadTypeBond      = Type(2)
	PayloadTypeSortition = Type(3)
	PayloadTypeUnbond    = Type(4)
	PayloadTypeWithdraw  = Type(5)
)

func (t Type) String() string {
	switch t {
	case PayloadTypeSend:
		return "send"
	case PayloadTypeBond:
		return "bond"
	case PayloadTypeUnbond:
		return "unbond"
	case PayloadTypeWithdraw:
		return "withdraw"
	case PayloadTypeSortition:
		return "sortition"
	}
	return fmt.Sprintf("%d", t)
}

type Payload interface {
	Signer() crypto.Address
	Value() int64
	Type() Type
	SerializeSize() int
	Encode(io.Writer) error
	Decode(io.Reader) error
	SanityCheck() error
	Fingerprint() string
}
